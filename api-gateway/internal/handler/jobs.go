// Package handler – jobs.go implements job matching endpoints.
package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/learnbot/api-gateway/internal/middleware"
	"github.com/learnbot/api-gateway/internal/types"
	"github.com/learnbot/resume-parser/pkg/scoring"
)

// ─────────────────────────────────────────────────────────────────────────────
// Built-in job catalog (MVP – replace with database in production)
// ─────────────────────────────────────────────────────────────────────────────

// sampleJobs is a small in-memory job catalog for the MVP.
var sampleJobs = []types.JobDetail{
	{
		JobSummary: types.JobSummary{
			ID:              "job-001",
			Title:           "Senior Backend Engineer (Go)",
			Company:         "TechCorp",
			LocationType:    "remote",
			ExperienceLevel: "senior",
			RequiredSkills:  []string{"Go", "PostgreSQL", "Docker", "Kubernetes"},
			PostedAt:        "2025-01-15",
		},
		Description:     "Build high-performance backend services using Go. Work with distributed systems and cloud infrastructure.",
		PreferredSkills: []string{"Terraform", "AWS", "Redis"},
		MinExperience:   5,
		Industry:        "software",
		SalaryMin:       intPtr(130000),
		SalaryMax:       intPtr(180000),
		SalaryCurrency:  "USD",
		ApplyURL:        "https://techcorp.example.com/jobs/001",
	},
	{
		JobSummary: types.JobSummary{
			ID:              "job-002",
			Title:           "Machine Learning Engineer",
			Company:         "AI Startup",
			LocationType:    "hybrid",
			ExperienceLevel: "mid",
			RequiredSkills:  []string{"Python", "TensorFlow", "PyTorch", "Kubernetes"},
			PostedAt:        "2025-01-20",
		},
		Description:     "Design and implement ML models for production. Work on NLP and computer vision projects.",
		PreferredSkills: []string{"Docker", "AWS", "Spark"},
		MinExperience:   3,
		Industry:        "artificial intelligence",
		SalaryMin:       intPtr(120000),
		SalaryMax:       intPtr(160000),
		SalaryCurrency:  "USD",
		ApplyURL:        "https://aistartup.example.com/jobs/002",
	},
	{
		JobSummary: types.JobSummary{
			ID:              "job-003",
			Title:           "Full Stack JavaScript Developer",
			Company:         "WebAgency",
			LocationType:    "remote",
			ExperienceLevel: "mid",
			RequiredSkills:  []string{"JavaScript", "React", "Node.js", "PostgreSQL"},
			PostedAt:        "2025-01-22",
		},
		Description:     "Build modern web applications using React and Node.js. Work with REST APIs and databases.",
		PreferredSkills: []string{"TypeScript", "Docker", "AWS"},
		MinExperience:   3,
		Industry:        "software",
		SalaryMin:       intPtr(90000),
		SalaryMax:       intPtr(130000),
		SalaryCurrency:  "USD",
		ApplyURL:        "https://webagency.example.com/jobs/003",
	},
	{
		JobSummary: types.JobSummary{
			ID:              "job-004",
			Title:           "DevOps Engineer",
			Company:         "CloudFirst",
			LocationType:    "remote",
			ExperienceLevel: "senior",
			RequiredSkills:  []string{"Kubernetes", "Docker", "Terraform", "AWS"},
			PostedAt:        "2025-01-25",
		},
		Description:     "Manage cloud infrastructure and CI/CD pipelines. Implement infrastructure as code.",
		PreferredSkills: []string{"Go", "Python", "Ansible", "Prometheus"},
		MinExperience:   5,
		Industry:        "software",
		SalaryMin:       intPtr(120000),
		SalaryMax:       intPtr(160000),
		SalaryCurrency:  "USD",
		ApplyURL:        "https://cloudfirst.example.com/jobs/004",
	},
	{
		JobSummary: types.JobSummary{
			ID:              "job-005",
			Title:           "Data Engineer",
			Company:         "DataCo",
			LocationType:    "hybrid",
			ExperienceLevel: "mid",
			RequiredSkills:  []string{"Python", "SQL", "Spark", "Kafka"},
			PostedAt:        "2025-01-28",
		},
		Description:     "Build and maintain data pipelines. Work with large-scale data processing systems.",
		PreferredSkills: []string{"Airflow", "dbt", "AWS", "Docker"},
		MinExperience:   3,
		Industry:        "data",
		SalaryMin:       intPtr(100000),
		SalaryMax:       intPtr(140000),
		SalaryCurrency:  "USD",
		ApplyURL:        "https://dataco.example.com/jobs/005",
	},
}

func intPtr(n int) *int { return &n }

// ─────────────────────────────────────────────────────────────────────────────
// JobsHandler
// ─────────────────────────────────────────────────────────────────────────────

// JobsHandler handles job matching endpoints.
type JobsHandler struct{}

// NewJobsHandler creates a new JobsHandler.
func NewJobsHandler() *JobsHandler {
	return &JobsHandler{}
}

// RegisterRoutes registers job routes on the mux.
//
//	POST /api/jobs/search          – search jobs with filters
//	GET  /api/jobs/recommendations – get recommended jobs for current user
//	GET  /api/jobs/{id}            – get job details
//	GET  /api/jobs/{id}/match      – get acceptance likelihood for a job
func (h *JobsHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/api/jobs/search",
		authMiddleware(http.HandlerFunc(h.Search)))
	mux.Handle("/api/jobs/recommendations",
		authMiddleware(http.HandlerFunc(h.Recommendations)))
	mux.HandleFunc("/api/jobs/", h.handleJobByID)
}

// Search handles POST /api/jobs/search.
//
// Request body:
//
//	{
//	  "query": "backend engineer",
//	  "skills": ["Go", "PostgreSQL"],
//	  "location_type": "remote",
//	  "experience_level": "senior",
//	  "limit": 20,
//	  "offset": 0
//	}
func (h *JobsHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var req types.JobSearchRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	// Apply defaults.
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Filter jobs.
	var results []types.JobSummary
	for _, job := range sampleJobs {
		if !matchesJobFilter(job, req) {
			continue
		}
		results = append(results, job.JobSummary)
	}

	// Apply pagination.
	total := len(results)
	if req.Offset >= total {
		results = nil
	} else {
		end := req.Offset + req.Limit
		if end > total {
			end = total
		}
		results = results[req.Offset:end]
	}

	WriteSuccessWithMeta(w, http.StatusOK, results, &types.ResponseMeta{
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
}

// Recommendations handles GET /api/jobs/recommendations.
// Returns jobs ranked by acceptance likelihood for the current user.
func (h *JobsHandler) Recommendations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteMethodNotAllowed(w)
		return
	}

	userID := middleware.GetUserID(r)
	profile := buildCandidateProfile(userID)

	// Score all jobs for this user.
	var scored []scoredJob

	for _, job := range sampleJobs {
		jobReqs := jobToRequirements(job)
		breakdown := scoring.Calculate(profile, jobReqs)
		summary := job.JobSummary
		score := breakdown.OverallScore
		summary.MatchScore = &score
		scored = append(scored, scoredJob{job: summary, score: breakdown.OverallScore})
	}

	// Sort by score descending.
	sortScoredJobs(scored)

	results := make([]types.JobSummary, len(scored))
	for i, s := range scored {
		results[i] = s.job
	}

	WriteSuccess(w, http.StatusOK, results)
}

// handleJobByID routes GET /api/jobs/{id} and GET /api/jobs/{id}/match.
func (h *JobsHandler) handleJobByID(w http.ResponseWriter, r *http.Request) {
	// Parse path: /api/jobs/{id}[/match]
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	parts := strings.SplitN(path, "/", 2)
	jobID := parts[0]

	if jobID == "" {
		WriteNotFound(w, "job")
		return
	}

	// Find job.
	var found *types.JobDetail
	for i := range sampleJobs {
		if sampleJobs[i].ID == jobID {
			found = &sampleJobs[i]
			break
		}
	}
	if found == nil {
		WriteNotFound(w, "job")
		return
	}

	// Route to sub-handler.
	if len(parts) == 2 && parts[1] == "match" {
		h.getJobMatch(w, r, found)
	} else {
		h.getJobDetail(w, r, found)
	}
}

// getJobDetail handles GET /api/jobs/{id}.
func (h *JobsHandler) getJobDetail(w http.ResponseWriter, r *http.Request, job *types.JobDetail) {
	if r.Method != http.MethodGet {
		WriteMethodNotAllowed(w)
		return
	}
	WriteSuccess(w, http.StatusOK, job)
}

// getJobMatch handles GET /api/jobs/{id}/match.
// Returns the acceptance likelihood score for the current user vs this job.
func (h *JobsHandler) getJobMatch(w http.ResponseWriter, r *http.Request, job *types.JobDetail) {
	if r.Method != http.MethodGet {
		WriteMethodNotAllowed(w)
		return
	}

	// Get user ID from query param or auth context.
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		// Try auth context.
		userID = middleware.GetUserID(r)
	}

	profile := buildCandidateProfile(userID)
	jobReqs := jobToRequirements(*job)
	breakdown := scoring.Calculate(profile, jobReqs)

	// Build recommendation text.
	recommendation := buildMatchRecommendation(breakdown.OverallScore)

	WriteSuccess(w, http.StatusOK, types.JobMatchResponse{
		JobID:           job.ID,
		OverallScore:    breakdown.OverallScore,
		SkillMatch:      breakdown.SkillMatchScore,
		ExperienceMatch: breakdown.ExperienceMatchScore,
		EducationMatch:  breakdown.EducationMatchScore,
		LocationFit:     breakdown.LocationFitScore,
		IndustryMatch:   breakdown.IndustryRelevanceScore,
		MatchedSkills:   breakdown.MatchedRequiredSkills,
		MissingSkills:   breakdown.MissingRequiredSkills,
		Recommendation:  recommendation,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper functions
// ─────────────────────────────────────────────────────────────────────────────

// buildCandidateProfile builds a scoring.CandidateProfile from the user's stored profile.
func buildCandidateProfile(userID string) scoring.CandidateProfile {
	if userID == "" {
		return scoring.CandidateProfile{}
	}

	profile := globalProfileStore.get(userID)
	skills := make([]scoring.CandidateSkill, len(profile.Skills))
	for i, s := range profile.Skills {
		skills[i] = scoring.CandidateSkill{
			Name:        s.Name,
			Proficiency: s.Proficiency,
		}
	}

	return scoring.CandidateProfile{
		Skills:            skills,
		YearsOfExperience: profile.YearsOfExperience,
		RemotePreference:  "any",
	}
}

// jobToRequirements converts a JobDetail to scoring.JobRequirements.
func jobToRequirements(job types.JobDetail) scoring.JobRequirements {
	return scoring.JobRequirements{
		Title:              job.Title,
		RequiredSkills:     job.RequiredSkills,
		PreferredSkills:    job.PreferredSkills,
		MinYearsExperience: job.MinExperience,
		LocationType:       job.LocationType,
		Industry:           job.Industry,
		ExperienceLevel:    job.ExperienceLevel,
	}
}

// matchesJobFilter returns true if a job matches the search filter.
func matchesJobFilter(job types.JobDetail, filter types.JobSearchRequest) bool {
	// Location type filter.
	if filter.LocationType != "" && !strings.EqualFold(job.LocationType, filter.LocationType) {
		return false
	}

	// Experience level filter.
	if filter.ExperienceLevel != "" && !strings.EqualFold(job.ExperienceLevel, filter.ExperienceLevel) {
		return false
	}

	// Industry filter.
	if filter.Industry != "" && !strings.EqualFold(job.Industry, filter.Industry) {
		return false
	}

	// Skills filter: job must require at least one of the requested skills.
	if len(filter.Skills) > 0 {
		found := false
		for _, filterSkill := range filter.Skills {
			for _, jobSkill := range job.RequiredSkills {
				if strings.EqualFold(filterSkill, jobSkill) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	// Query filter: match against title or company.
	if filter.Query != "" {
		q := strings.ToLower(filter.Query)
		if !strings.Contains(strings.ToLower(job.Title), q) &&
			!strings.Contains(strings.ToLower(job.Company), q) {
			return false
		}
	}

	return true
}

// scoredJob holds a job with its match score.
type scoredJob struct {
	job   types.JobSummary
	score float64
}

// sortScoredJobs sorts jobs by score descending (simple insertion sort for small slices).
func sortScoredJobs(jobs []scoredJob) {
	for i := 1; i < len(jobs); i++ {
		key := jobs[i]
		j := i - 1
		for j >= 0 && jobs[j].score < key.score {
			jobs[j+1] = jobs[j]
			j--
		}
		jobs[j+1] = key
	}
}

// buildMatchRecommendation returns a recommendation string based on the score.
func buildMatchRecommendation(score float64) string {
	switch {
	case score >= 80:
		return "Strong match – you're ready to apply!"
	case score >= 60:
		return "Good match – address a few gaps to strengthen your application."
	case score >= 40:
		return "Moderate match – significant skill gaps to address before applying."
	default:
		return "Low match – consider building more relevant skills first."
	}
}

// queryParamInt parses an integer query parameter with a default value.
func queryParamInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}
