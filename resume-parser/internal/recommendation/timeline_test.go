package recommendation

import (
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// Timeline tests
// ─────────────────────────────────────────────────────────────────────────────

func TestBuildTimeline_EmptyPhases(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: 10}
	timeline := buildTimeline(nil, prefs, "Test Job")

	if timeline.TotalWeeks != 0 {
		t.Errorf("expected 0 weeks for empty phases, got %d", timeline.TotalWeeks)
	}
	if timeline.TotalHours != 0 {
		t.Errorf("expected 0 hours for empty phases, got %.1f", timeline.TotalHours)
	}
}

func TestBuildTimeline_WeeklyHoursSet(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: 15}
	timeline := buildTimeline(nil, prefs, "Test Job")

	if timeline.WeeklyHours != 15 {
		t.Errorf("expected WeeklyHours=15, got %.1f", timeline.WeeklyHours)
	}
}

func TestBuildTimeline_DefaultWeeklyHours(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: 0}
	timeline := buildTimeline(nil, prefs, "Test Job")

	if timeline.WeeklyHours != 10 {
		t.Errorf("expected default WeeklyHours=10, got %.1f", timeline.WeeklyHours)
	}
}

func TestBuildTimeline_WeeksAreSequential(t *testing.T) {
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 20,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 20},
						EstimatedCompletionHours: 20,
					},
				},
			},
			TotalHours:     20,
			EstimatedWeeks: 2,
		},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	timeline := buildTimeline(phases, prefs, "Test Job")

	for i, week := range timeline.Weeks {
		if week.WeekNumber != i+1 {
			t.Errorf("week[%d] has WeekNumber=%d, expected %d", i, week.WeekNumber, i+1)
		}
	}
}

func TestBuildTimeline_CumulativeHoursIncreasing(t *testing.T) {
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 30,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 30},
						EstimatedCompletionHours: 30,
					},
				},
			},
			TotalHours:     30,
			EstimatedWeeks: 3,
		},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	timeline := buildTimeline(phases, prefs, "Test Job")

	for i := 1; i < len(timeline.Weeks); i++ {
		if timeline.Weeks[i].CumulativeHours < timeline.Weeks[i-1].CumulativeHours {
			t.Errorf("cumulative hours decreased at week %d: %.1f < %.1f",
				i+1, timeline.Weeks[i].CumulativeHours, timeline.Weeks[i-1].CumulativeHours)
		}
	}
}

func TestBuildTimeline_PhaseNumberAssigned(t *testing.T) {
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 10,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 10},
						EstimatedCompletionHours: 10,
					},
				},
			},
			TotalHours:     10,
			EstimatedWeeks: 1,
		},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	timeline := buildTimeline(phases, prefs, "Test Job")

	for _, week := range timeline.Weeks {
		if week.PhaseNumber != 1 {
			t.Errorf("expected PhaseNumber=1, got %d", week.PhaseNumber)
		}
	}
}

func TestBuildTimeline_TargetDateSet(t *testing.T) {
	prefs := UserPreferences{
		WeeklyHoursAvailable: 10,
		TargetDate:           "2025-12-31",
	}
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 10,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 10},
						EstimatedCompletionHours: 10,
					},
				},
			},
			TotalHours:     10,
			EstimatedWeeks: 1,
		},
	}

	timeline := buildTimeline(phases, prefs, "Test Job")

	if timeline.TargetCompletionDate != "2025-12-31" {
		t.Errorf("expected TargetCompletionDate='2025-12-31', got %s", timeline.TargetCompletionDate)
	}
}

func TestBuildTimeline_EstimatedDateWhenNoTarget(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: 10}
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 10,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 10},
						EstimatedCompletionHours: 10,
					},
				},
			},
			TotalHours:     10,
			EstimatedWeeks: 1,
		},
	}

	timeline := buildTimeline(phases, prefs, "Test Job")

	// Should have an estimated date.
	if timeline.TargetCompletionDate == "" {
		t.Error("expected non-empty TargetCompletionDate when phases exist")
	}
}

func TestBuildTimeline_ActivitiesPopulated(t *testing.T) {
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 10,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 10},
						EstimatedCompletionHours: 10,
					},
				},
			},
			TotalHours:     10,
			EstimatedWeeks: 1,
		},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	timeline := buildTimeline(phases, prefs, "Test Job")

	for _, week := range timeline.Weeks {
		if len(week.Activities) == 0 {
			t.Errorf("week %d has no activities", week.WeekNumber)
		}
	}
}

func TestBuildTimeline_CheckpointForCriticalSkill(t *testing.T) {
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 10,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 10},
						EstimatedCompletionHours: 10,
					},
				},
			},
			TotalHours:     10,
			EstimatedWeeks: 1,
		},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	timeline := buildTimeline(phases, prefs, "Test Job")

	// The last week of a critical skill should be a checkpoint.
	hasCheckpoint := false
	for _, week := range timeline.Weeks {
		if week.IsCheckpoint {
			hasCheckpoint = true
			break
		}
	}
	if !hasCheckpoint {
		t.Error("expected at least one checkpoint week for critical skill")
	}
}

func TestBuildTimeline_MultipleSkillsMultipleWeeks(t *testing.T) {
	phases := []LearningPhase{
		{
			PhaseNumber: 1,
			PhaseName:   "Critical Skills",
			Skills: []SkillRecommendation{
				{
					SkillName:                "Python",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 20,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Python Course", DurationHours: 20},
						EstimatedCompletionHours: 20,
					},
				},
				{
					SkillName:                "Go",
					GapCategory:              "critical",
					EstimatedHoursToJobReady: 10,
					PrimaryResource: &RecommendedResource{
						Resource:                 ResourceEntry{Title: "Go Course", DurationHours: 10},
						EstimatedCompletionHours: 10,
					},
				},
			},
			TotalHours:     30,
			EstimatedWeeks: 3,
		},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	timeline := buildTimeline(phases, prefs, "Test Job")

	// Should have at least 3 weeks (2 for Python + 1 for Go) + 1 review week.
	if timeline.TotalWeeks < 3 {
		t.Errorf("expected at least 3 weeks, got %d", timeline.TotalWeeks)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Handler tests
// ─────────────────────────────────────────────────────────────────────────────

func TestBuildWeekActivities_FirstWeek(t *testing.T) {
	activities := buildWeekActivities("Python", "Python Course", 0, 3, nil)

	if len(activities) == 0 {
		t.Error("expected non-empty activities for first week")
	}
	// First week should mention starting the resource.
	found := false
	for _, a := range activities {
		if containsSubstring(a, "Start") || containsSubstring(a, "start") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected first week to mention starting the resource, got: %v", activities)
	}
}

func TestBuildWeekActivities_LastWeek(t *testing.T) {
	activities := buildWeekActivities("Python", "Python Course", 2, 3, nil)

	if len(activities) == 0 {
		t.Error("expected non-empty activities for last week")
	}
	// Last week should mention completing the resource.
	found := false
	for _, a := range activities {
		if containsSubstring(a, "Complete") || containsSubstring(a, "complete") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected last week to mention completing the resource, got: %v", activities)
	}
}

func TestJoinSkills_LessThanMax(t *testing.T) {
	skills := []string{"Python", "Go"}
	result := joinSkills(skills, 3)
	if result != "Python, Go" {
		t.Errorf("expected 'Python, Go', got %q", result)
	}
}

func TestJoinSkills_MoreThanMax(t *testing.T) {
	skills := []string{"Python", "Go", "Rust", "Java"}
	result := joinSkills(skills, 3)
	if !containsSubstring(result, "and 1 more") {
		t.Errorf("expected 'and 1 more' in result, got %q", result)
	}
}
