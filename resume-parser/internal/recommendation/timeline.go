// Package recommendation – timeline.go implements the learning timeline generator.
// It converts a set of learning phases into a week-by-week study schedule.
package recommendation

import (
	"fmt"
	"math"
	"time"
)

// buildTimeline generates a week-by-week learning timeline from the phases.
//
// Algorithm:
//  1. Flatten all skill recommendations across phases in order.
//  2. Assign each resource to one or more weeks based on duration.
//  3. Respect the user's weekly hours available.
//  4. Insert checkpoint weeks at phase boundaries.
//  5. Calculate cumulative hours and target completion date.
func buildTimeline(phases []LearningPhase, prefs UserPreferences, jobTitle string) LearningTimeline {
	weeklyHours := prefs.WeeklyHoursAvailable
	if weeklyHours <= 0 {
		weeklyHours = 10
	}

	var weeks []WeeklySchedule
	weekNum := 1
	cumulativeHours := 0.0

	for _, phase := range phases {
		for _, skillRec := range phase.Skills {
			// Determine the resource to schedule.
			var resourceTitle string
			var resourceHours float64

			if skillRec.PrimaryResource != nil {
				resourceTitle = skillRec.PrimaryResource.Resource.Title
				resourceHours = skillRec.PrimaryResource.EstimatedCompletionHours
			} else {
				resourceTitle = "Self-study: " + skillRec.SkillName
				resourceHours = float64(skillRec.EstimatedHoursToJobReady)
			}

			if resourceHours <= 0 {
				resourceHours = float64(skillRec.EstimatedHoursToJobReady)
			}

			// Calculate how many weeks this resource spans.
			weeksNeeded := math.Ceil(resourceHours / weeklyHours)
			if weeksNeeded < 1 {
				weeksNeeded = 1
			}

			for w := 0; w < int(weeksNeeded); w++ {
				hoursThisWeek := weeklyHours
				remaining := resourceHours - float64(w)*weeklyHours
				if remaining < weeklyHours {
					hoursThisWeek = remaining
				}

				cumulativeHours += hoursThisWeek
				isLastWeekOfResource := w == int(weeksNeeded)-1

				activities := buildWeekActivities(skillRec.SkillName, resourceTitle, w, int(weeksNeeded), skillRec.PrimaryResource)

				week := WeeklySchedule{
					WeekNumber:      weekNum,
					PhaseNumber:     phase.PhaseNumber,
					SkillFocus:      skillRec.SkillName,
					ResourceTitle:   resourceTitle,
					HoursPlanned:    roundTo1(hoursThisWeek),
					CumulativeHours: roundTo1(cumulativeHours),
					Activities:      activities,
					IsCheckpoint:    isLastWeekOfResource && skillRec.GapCategory == "critical",
				}

				if week.IsCheckpoint {
					week.CheckpointDescription = fmt.Sprintf(
						"Complete %s and verify %s proficiency through practice exercises.",
						resourceTitle, skillRec.SkillName)
				}

				weeks = append(weeks, week)
				weekNum++
			}
		}

		// Add a phase review week if the phase has multiple skills.
		if len(phase.Skills) > 1 {
			cumulativeHours += weeklyHours * 0.5 // half week for review
			weeks = append(weeks, WeeklySchedule{
				WeekNumber:            weekNum,
				PhaseNumber:           phase.PhaseNumber,
				SkillFocus:            "Phase Review",
				ResourceTitle:         fmt.Sprintf("Review & consolidate %s", phase.PhaseName),
				HoursPlanned:          roundTo1(weeklyHours * 0.5),
				CumulativeHours:       roundTo1(cumulativeHours),
				Activities:            buildReviewActivities(phase),
				IsCheckpoint:          true,
				CheckpointDescription: phase.Milestone,
			})
			weekNum++
		}
	}

	// Calculate target completion date.
	targetDate := ""
	if prefs.TargetDate != "" {
		targetDate = prefs.TargetDate
	} else if len(weeks) > 0 {
		// Estimate based on start date = today.
		completionDate := time.Now().AddDate(0, 0, len(weeks)*7)
		targetDate = completionDate.Format("2006-01-02")
	}

	return LearningTimeline{
		TotalWeeks:           len(weeks),
		TotalHours:           roundTo1(cumulativeHours),
		WeeklyHours:          weeklyHours,
		Weeks:                weeks,
		TargetCompletionDate: targetDate,
	}
}

// buildWeekActivities generates suggested activities for a study week.
func buildWeekActivities(
	skillName, resourceTitle string,
	weekIndex, totalWeeks int,
	resource *RecommendedResource,
) []string {
	var activities []string

	if weekIndex == 0 {
		// First week of a resource.
		activities = append(activities, fmt.Sprintf("Start '%s'", resourceTitle))
		activities = append(activities, fmt.Sprintf("Set up development environment for %s", skillName))
		if resource != nil && resource.Resource.HasHandsOn {
			activities = append(activities, "Complete introductory exercises")
		}
	} else if weekIndex == totalWeeks-1 {
		// Last week of a resource.
		activities = append(activities, fmt.Sprintf("Complete '%s'", resourceTitle))
		activities = append(activities, fmt.Sprintf("Build a small project using %s", skillName))
		activities = append(activities, "Review key concepts and take notes")
	} else {
		// Middle weeks.
		activities = append(activities, fmt.Sprintf("Continue '%s' (week %d of %d)", resourceTitle, weekIndex+1, totalWeeks))
		if resource != nil && resource.Resource.HasHandsOn {
			activities = append(activities, "Complete hands-on exercises")
		}
		activities = append(activities, fmt.Sprintf("Practice %s concepts", skillName))
	}

	// Add practice recommendation for technical skills.
	if resource != nil && resource.Resource.ResourceType != "practice" {
		activities = append(activities, fmt.Sprintf("Supplement with LeetCode/HackerRank problems for %s", skillName))
	}

	return activities
}

// buildReviewActivities generates activities for a phase review week.
func buildReviewActivities(phase LearningPhase) []string {
	skillNames := make([]string, 0, len(phase.Skills))
	for _, s := range phase.Skills {
		skillNames = append(skillNames, s.SkillName)
	}

	activities := []string{
		fmt.Sprintf("Review all %s skills covered in Phase %d", phase.PhaseName, phase.PhaseNumber),
		"Build an integration project combining learned skills",
		"Update your resume/portfolio with new skills",
	}

	if len(skillNames) > 0 {
		activities = append(activities, fmt.Sprintf("Practice interview questions for: %s", joinSkills(skillNames, 3)))
	}

	return activities
}

// joinSkills joins skill names with commas, truncating at maxCount.
func joinSkills(skills []string, maxCount int) string {
	if len(skills) <= maxCount {
		result := ""
		for i, s := range skills {
			if i > 0 {
				result += ", "
			}
			result += s
		}
		return result
	}
	result := ""
	for i := 0; i < maxCount; i++ {
		if i > 0 {
			result += ", "
		}
		result += skills[i]
	}
	result += fmt.Sprintf(" and %d more", len(skills)-maxCount)
	return result
}
