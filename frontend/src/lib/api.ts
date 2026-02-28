const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8090";

export interface APIResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: { code: string; message: string; details?: Array<{ field: string; message: string }> };
  meta?: { total: number; limit: number; offset: number };
}

export interface AuthResponse {
  token: string;
  expires_at: string;
  user: { id: string; email: string; full_name: string };
}

export interface UserProfile {
  user_id: string;
  email: string;
  full_name: string;
  headline?: string;
  summary?: string;
  location_city?: string;
  location_country?: string;
  linkedin_url?: string;
  github_url?: string;
  website_url?: string;
  years_of_experience?: number;
  is_open_to_work?: boolean;
  skills?: Skill[];
  updated_at?: string;
}

export interface Skill {
  name: string;
  proficiency: string;
  years_of_experience?: number;
  is_primary?: boolean;
}

export interface Job {
  id: string;
  title: string;
  company: string;
  location_type: string;
  experience_level: string;
  required_skills: string[];
  preferred_skills?: string[];
  posted_at?: string;
  match_score?: number;
  description?: string;
  min_years_experience?: number;
  industry?: string;
  location_city?: string;
  location_country?: string;
  salary_min?: number;
  salary_max?: number;
  salary_currency?: string;
  apply_url?: string;
}

export interface JobMatch {
  job_id: string;
  overall_score: number;
  skill_match: number;
  experience_match: number;
  education_match: number;
  location_fit: number;
  industry_match: number;
  matched_skills: string[];
  missing_skills: string[];
  recommendation: string;
}

export interface SkillGap {
  skill_name: string;
  category: string;
  priority_score: number;
  importance_score: number;
  estimated_learning_hours: number;
  transferability_score: number;
  current_level?: string;
  target_level: string;
  semantic_similarity_score: number;
  closest_existing_skill?: string;
  recommendations: Array<{ title: string; description: string; resource_type: string; estimated_hours: number; priority: number }>;
  difficulty: string;
}

export interface GapAnalysisResult {
  critical_gaps: SkillGap[];
  important_gaps: SkillGap[];
  nice_to_have_gaps: SkillGap[];
  total_gaps: number;
  critical_gap_count: number;
  important_gap_count: number;
  nice_to_have_gap_count: number;
  total_estimated_learning_hours: number;
  readiness_score: number;
  top_priority_gaps: SkillGap[];
  matched_skills: string[];
  visual_data: {
    radar_chart: { labels: string[]; candidate_scores: number[]; required_scores: number[] };
    gaps_by_category: Array<{ category: string; gap_count: number; total_learning_hours: number; average_priority: number }>;
    learning_timeline: Array<{ order: number; skill_name: string; category: string; estimated_hours: number; cumulative_hours: number; rationale: string }>;
  };
}

export interface Resource {
  id: string; title: string; description: string; url: string; provider: string;
  resource_type: string; difficulty: string; cost_type: string; cost_usd?: number;
  duration_hours?: number; duration_label?: string; skills: string[]; primary_skill: string;
  rating?: number; rating_count?: number; has_certificate: boolean; has_hands_on: boolean; is_verified: boolean;
}

export interface ResourceRecommendation {
  resource: Resource; relevance_score: number; is_alternative: boolean;
  recommendation_reason: string; estimated_completion_hours: number;
}

export interface SkillRecommendation {
  skill_name: string; gap_category: string; priority_score: number;
  primary_resource?: ResourceRecommendation; alternative_resources?: ResourceRecommendation[];
  estimated_hours_to_job_ready: number; current_level?: string; target_level: string;
}

export interface LearningPhase {
  phase_number: number; phase_name: string; phase_description: string;
  skills: SkillRecommendation[]; total_hours: number; estimated_weeks: number; milestone: string;
}

export interface WeeklySchedule {
  week_number: number; phase_number: number; skill_focus: string; resource_title: string;
  hours_planned: number; cumulative_hours: number; activities: string[];
  is_checkpoint: boolean; checkpoint_description?: string;
}

export interface LearningPlan {
  job_title: string; readiness_score: number; total_gaps: number; total_estimated_hours: number;
  phases: LearningPhase[];
  timeline: { total_weeks: number; total_hours: number; weekly_hours: number; weeks: WeeklySchedule[]; target_completion_date?: string };
  matched_skills: string[];
  summary: { headline: string; critical_gap_count: number; important_gap_count: number; free_resource_count: number; paid_resource_count: number; estimated_total_cost_usd: number; top_skills_to_learn: string[]; quick_wins: string[] };
}

export class APIError extends Error {
  constructor(public code: string, message: string, public details?: Array<{ field: string; message: string }>) {
    super(message); this.name = "APIError";
  }
}

async function request<T>(path: string, options: RequestInit = {}, token?: string): Promise<T> {
  const headers: Record<string, string> = { "Content-Type": "application/json", ...(options.headers as Record<string, string>) };
  if (token) headers["Authorization"] = `Bearer ${token}`;
  const response = await fetch(`${BASE_URL}${path}`, { ...options, headers });
  const data: APIResponse<T> = await response.json();
  if (!data.success || !response.ok) throw new APIError(data.error?.code || "UNKNOWN_ERROR", data.error?.message || "An unexpected error occurred", data.error?.details);
  return data.data as T;
}

export const authAPI = {
  register: (email: string, password: string, fullName: string) =>
    request<AuthResponse>("/api/auth/register", { method: "POST", body: JSON.stringify({ email, password, full_name: fullName }) }),
  login: (email: string, password: string) =>
    request<AuthResponse>("/api/auth/login", { method: "POST", body: JSON.stringify({ email, password }) }),
};

export const profileAPI = {
  getProfile: (token: string) => request<UserProfile>("/api/users/profile", {}, token),
  updateProfile: (token: string, data: Partial<UserProfile>) =>
    request<UserProfile>("/api/users/profile", { method: "PUT", body: JSON.stringify(data) }, token),
  getSkills: (token: string) => request<{ user_id: string; skills: Skill[]; count: number }>("/api/profile/skills", {}, token),
  updateSkills: (token: string, skills: Skill[]) =>
    request<{ user_id: string; skills: Skill[]; count: number }>("/api/profile/skills", { method: "PUT", body: JSON.stringify({ skills }) }, token),
};

export const resumeAPI = {
  upload: async (token: string, file: File) => {
    const formData = new FormData();
    formData.append("resume", file);
    const response = await fetch(`${BASE_URL}/api/resume/upload`, { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: formData });
    const data = await response.json();
    if (!data.success) throw new APIError(data.error?.code || "PARSE_ERROR", data.error?.message || "Failed to parse resume");
    return data.data;
  },
};

export interface JobSearchParams { query?: string; skills?: string[]; location_type?: string; experience_level?: string; industry?: string; limit?: number; offset?: number; }

export const jobsAPI = {
  search: (token: string, params: JobSearchParams = {}) =>
    fetch(`${BASE_URL}/api/jobs/search`, { method: "POST", headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` }, body: JSON.stringify(params) })
      .then(async (r) => { const d = await r.json(); if (!d.success) throw new APIError(d.error?.code, d.error?.message); return { jobs: d.data as Job[], meta: d.meta }; }),
  getRecommendations: (token: string) => request<Job[]>("/api/jobs/recommendations", {}, token),
  getJob: (jobId: string) => request<Job>(`/api/jobs/${jobId}`),
  getJobMatch: (token: string, jobId: string) => request<JobMatch>(`/api/jobs/${jobId}/match`, {}, token),
};

export interface JobRequirementsInput { title: string; required_skills: string[]; preferred_skills?: string[]; min_years_experience?: number; experience_level?: string; location_type?: string; industry?: string; }

export const analysisAPI = {
  analyzeGaps: (token: string, jobId?: string, job?: JobRequirementsInput) =>
    request<GapAnalysisResult>("/api/analysis/gaps", { method: "POST", body: JSON.stringify({ job_id: jobId, job }) }, token),
  getTrainingPlan: (token: string, jobId?: string, job?: JobRequirementsInput, preferences?: { prefer_free?: boolean; max_budget_usd?: number; weekly_hours_available?: number; prefer_hands_on?: boolean; prefer_certificates?: boolean; target_date?: string }) =>
    request<LearningPlan>("/api/training/recommendations", { method: "POST", body: JSON.stringify({ job_id: jobId, job, preferences }) }, token),
};

export interface ResourceSearchParams { skill?: string; type?: string; difficulty?: string; free?: boolean; has_certificate?: boolean; has_hands_on?: boolean; min_rating?: number; q?: string; limit?: number; offset?: number; }

export const resourcesAPI = {
  search: (params: ResourceSearchParams = {}) => {
    const query = new URLSearchParams();
    if (params.skill) query.set("skill", params.skill);
    if (params.type) query.set("type", params.type);
    if (params.difficulty) query.set("difficulty", params.difficulty);
    if (params.free) query.set("free", "true");
    if (params.has_certificate) query.set("has_certificate", "true");
    if (params.has_hands_on) query.set("has_hands_on", "true");
    if (params.min_rating) query.set("min_rating", String(params.min_rating));
    if (params.q) query.set("q", params.q);
    if (params.limit) query.set("limit", String(params.limit));
    if (params.offset) query.set("offset", String(params.offset));
    return fetch(`${BASE_URL}/api/resources/search?${query}`)
      .then(async (r) => { const d = await r.json(); if (!d.success) throw new APIError(d.error?.code, d.error?.message); return { resources: d.data as Resource[], meta: d.meta }; });
  },
};
