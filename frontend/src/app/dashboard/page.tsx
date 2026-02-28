"use client";
import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useAuthStore, useToken, useUser } from "@/store/authStore";
import { jobsAPI, analysisAPI, type Job, type GapAnalysisResult } from "@/lib/api";
import { Card, CardHeader } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Badge, ScoreBadge, GapCategoryBadge } from "@/components/ui/Badge";
import { InlineLoader } from "@/components/ui/LoadingSpinner";
import { RadarChart, PolarGrid, PolarAngleAxis, Radar, ResponsiveContainer, Tooltip } from "recharts";

export default function DashboardPage() {
  const token = useToken();
  const user = useUser();
  const { profile, loadProfile } = useAuthStore();
  const [recommendedJobs, setRecommendedJobs] = useState<Job[]>([]);
  const [gapResult, setGapResult] = useState<GapAnalysisResult | null>(null);
  const [loadingJobs, setLoadingJobs] = useState(true);
  const [loadingGaps, setLoadingGaps] = useState(true);

  useEffect(() => {
    if (token) { loadProfile(); loadRecommendedJobs(); loadGapAnalysis(); }
  }, [token]); // eslint-disable-line react-hooks/exhaustive-deps

  const loadRecommendedJobs = async () => {
    if (!token) return;
    setLoadingJobs(true);
    try { const jobs = await jobsAPI.getRecommendations(token); setRecommendedJobs(jobs.slice(0, 3)); }
    catch { /* silently fail */ } finally { setLoadingJobs(false); }
  };

  const loadGapAnalysis = async () => {
    if (!token) return;
    setLoadingGaps(true);
    try { const result = await analysisAPI.analyzeGaps(token, "job-001"); setGapResult(result); }
    catch { /* silently fail */ } finally { setLoadingGaps(false); }
  };

  const skills = profile?.skills || [];
  const readinessScore = gapResult?.readiness_score ?? 0;
  const radarData = gapResult?.visual_data?.radar_chart
    ? gapResult.visual_data.radar_chart.labels.map((label, i) => ({
        subject: label,
        candidate: Math.round((gapResult.visual_data.radar_chart.candidate_scores[i] || 0) * 100),
        required: 100,
      }))
    : [];

  return (
    <div className="page-container section animate-fade-in">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Welcome back, {user?.full_name?.split(" ")[0] || "there"}! 👋</h1>
        <p className="text-gray-600 mt-1">Here&apos;s your career development overview.</p>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-8">
        {[{ value: skills.length, label: "Skills" }, { value: `${Math.round(readinessScore)}%`, label: "Readiness" }, { value: gapResult?.critical_gap_count ?? "—", label: "Critical Gaps" }, { value: `${gapResult?.total_estimated_learning_hours ?? "—"}h`, label: "Learning Hours" }].map((s) => (
          <Card key={s.label} padding="sm" className="text-center">
            <div className="text-2xl font-bold text-primary-600">{s.value}</div>
            <div className="text-xs text-gray-500 mt-1">{s.label}</div>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader title="Recommended Jobs" subtitle="Ranked by your acceptance likelihood" action={<Link href="/jobs"><Button variant="ghost" size="sm">View all →</Button></Link>} />
            {loadingJobs ? <InlineLoader label="Loading recommendations..." /> : recommendedJobs.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <p className="text-sm">No job recommendations yet.</p>
                <Link href="/jobs" className="text-primary-600 text-sm hover:underline mt-1 block">Browse jobs →</Link>
              </div>
            ) : (
              <div className="space-y-3">
                {recommendedJobs.map((job) => (
                  <Link key={job.id} href={`/jobs/${job.id}`} className="block p-4 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50/30 transition-colors">
                    <div className="flex items-start justify-between gap-3">
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-gray-900 truncate">{job.title}</p>
                        <p className="text-sm text-gray-500">{job.company}</p>
                        <div className="flex flex-wrap gap-1.5 mt-2">
                          <Badge variant="gray">{job.location_type}</Badge>
                          <Badge variant="gray">{job.experience_level}</Badge>
                        </div>
                      </div>
                      {job.match_score !== undefined && <ScoreBadge score={job.match_score} />}
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </Card>

          <Card>
            <CardHeader title="Skill Gap Overview" subtitle="For: Senior Backend Engineer (Go)" action={<Link href="/analysis"><Button variant="ghost" size="sm">Full analysis →</Button></Link>} />
            {loadingGaps ? <InlineLoader label="Analyzing gaps..." /> : gapResult ? (
              <div>
                <div className="mb-4">
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-gray-600">Readiness Score</span>
                    <span className="font-medium text-gray-900">{Math.round(readinessScore)}%</span>
                  </div>
                  <div className="progress-bar">
                    <div className={`progress-fill ${readinessScore >= 75 ? "bg-success-500" : readinessScore >= 50 ? "bg-warning-500" : "bg-danger-500"}`} style={{ width: `${readinessScore}%` }} />
                  </div>
                </div>
                {gapResult.top_priority_gaps?.slice(0, 3).map((gap) => (
                  <div key={gap.skill_name} className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                    <div className="flex items-center gap-2">
                      <GapCategoryBadge category={gap.category} />
                      <span className="text-sm font-medium text-gray-800">{gap.skill_name}</span>
                    </div>
                    <span className="text-xs text-gray-500">~{gap.estimated_learning_hours}h</span>
                  </div>
                ))}
                {gapResult.total_gaps > 3 && <p className="text-xs text-gray-500 mt-2 text-center">+{gapResult.total_gaps - 3} more gaps</p>}
              </div>
            ) : (
              <div className="text-center py-6 text-gray-500 text-sm">
                <p>Add skills to your profile to see gap analysis.</p>
                <Link href="/onboarding" className="text-primary-600 hover:underline mt-1 block">Complete your profile →</Link>
              </div>
            )}
          </Card>
        </div>

        <div className="space-y-6">
          <Card>
            <CardHeader title="Your Profile" />
            <div className="space-y-3">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 font-bold text-lg">{user?.full_name?.charAt(0) || "U"}</div>
                <div><p className="font-medium text-gray-900">{user?.full_name}</p><p className="text-sm text-gray-500">{user?.email}</p></div>
              </div>
              {profile?.headline && <p className="text-sm text-gray-600">{profile.headline}</p>}
              {skills.length > 0 && (
                <div>
                  <p className="text-xs text-gray-500 mb-2">Top skills</p>
                  <div className="flex flex-wrap gap-1">
                    {skills.slice(0, 6).map((skill) => <span key={skill.name} className="skill-tag">{skill.name}</span>)}
                    {skills.length > 6 && <span className="skill-tag">+{skills.length - 6}</span>}
                  </div>
                </div>
              )}
              <Link href="/onboarding"><Button variant="secondary" size="sm" className="w-full">Edit Profile</Button></Link>
            </div>
          </Card>

          {radarData.length > 0 && (
            <Card>
              <CardHeader title="Skills Coverage" subtitle="vs. job requirements" />
              <ResponsiveContainer width="100%" height={200}>
                <RadarChart data={radarData}>
                  <PolarGrid />
                  <PolarAngleAxis dataKey="subject" tick={{ fontSize: 10, fill: "#6b7280" }} />
                  <Radar name="Required" dataKey="required" stroke="#e5e7eb" fill="#e5e7eb" fillOpacity={0.3} />
                  <Radar name="You" dataKey="candidate" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.4} />
                  <Tooltip formatter={(value: number) => [`${value}%`]} />
                </RadarChart>
              </ResponsiveContainer>
            </Card>
          )}

          <Card>
            <CardHeader title="Quick Actions" />
            <div className="space-y-2">
              {[{ href: "/jobs", label: "🔍 Search Jobs" }, { href: "/analysis", label: "📊 Analyze Gaps" }, { href: "/learning", label: "📚 Learning Plan" }].map((a) => (
                <Link key={a.href} href={a.href} className="block">
                  <Button variant="secondary" size="sm" className="w-full justify-start">{a.label}</Button>
                </Link>
              ))}
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
