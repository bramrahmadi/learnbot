"use client";
import React, { useState, useEffect, Suspense } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useToken } from "@/store/authStore";
import { analysisAPI, type GapAnalysisResult, type SkillGap } from "@/lib/api";
import { Card, CardHeader } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { GapCategoryBadge, DifficultyBadge } from "@/components/ui/Badge";
import { InlineLoader } from "@/components/ui/LoadingSpinner";
import { RadarChart, PolarGrid, PolarAngleAxis, Radar, ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, Cell } from "recharts";

const SAMPLE_JOBS = [
  { id: "job-001", title: "Senior Backend Engineer (Go)" }, { id: "job-002", title: "Machine Learning Engineer" },
  { id: "job-003", title: "Full Stack JavaScript Developer" }, { id: "job-004", title: "DevOps Engineer" }, { id: "job-005", title: "Data Engineer" },
];

function GapCard({ gap }: { gap: SkillGap }) {
  const [expanded, setExpanded] = useState(false);
  return (
    <div className="border border-gray-200 rounded-lg overflow-hidden">
      <button className="w-full flex items-center justify-between p-4 text-left hover:bg-gray-50 transition-colors" onClick={() => setExpanded(!expanded)} aria-expanded={expanded}>
        <div className="flex items-center gap-3">
          <GapCategoryBadge category={gap.category} />
          <span className="font-medium text-gray-900">{gap.skill_name}</span>
          <DifficultyBadge difficulty={gap.difficulty} />
        </div>
        <div className="flex items-center gap-3 flex-shrink-0">
          <span className="text-xs text-gray-500">~{gap.estimated_learning_hours}h</span>
          <span className="text-gray-400">{expanded ? "▲" : "▼"}</span>
        </div>
      </button>
      {expanded && (
        <div className="px-4 pb-4 border-t border-gray-100 bg-gray-50">
          <div className="grid grid-cols-2 gap-4 mt-3 mb-4 text-sm">
            <div><span className="text-gray-500">Priority Score</span><div className="font-medium text-gray-900 mt-0.5">{(gap.priority_score * 100).toFixed(0)}%</div></div>
            <div><span className="text-gray-500">Target Level</span><div className="font-medium text-gray-900 mt-0.5 capitalize">{gap.target_level}</div></div>
            {gap.closest_existing_skill && <div><span className="text-gray-500">Related Skill</span><div className="font-medium text-gray-900 mt-0.5">{gap.closest_existing_skill}</div></div>}
          </div>
          {gap.recommendations?.length > 0 && (
            <div>
              <p className="text-xs font-medium text-gray-500 mb-2">Recommendations</p>
              <div className="space-y-2">
                {gap.recommendations.slice(0, 2).map((rec, i) => (
                  <div key={i} className="bg-white rounded-lg p-3 border border-gray-200">
                    <p className="text-sm font-medium text-gray-800">{rec.title}</p>
                    <p className="text-xs text-gray-500 mt-0.5">{rec.description}</p>
                    <div className="flex items-center gap-2 mt-1">
                      <span className="text-xs text-gray-400 capitalize">{rec.resource_type}</span>
                      <span className="text-xs text-gray-400">·</span>
                      <span className="text-xs text-gray-400">~{rec.estimated_hours}h</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function AnalysisPageInner() {
  const token = useToken();
  const searchParams = useSearchParams();
  const [selectedJobId, setSelectedJobId] = useState(searchParams.get("job_id") || "job-001");
  const [result, setResult] = useState<GapAnalysisResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => { if (token) runAnalysis(); }, [token, selectedJobId]); // eslint-disable-line react-hooks/exhaustive-deps

  const runAnalysis = async () => {
    if (!token) return;
    setIsLoading(true); setError(null);
    try { const data = await analysisAPI.analyzeGaps(token, selectedJobId); setResult(data); }
    catch (err: unknown) { setError(err instanceof Error ? err.message : "Analysis failed"); }
    finally { setIsLoading(false); }
  };

  const radarData = result?.visual_data?.radar_chart
    ? result.visual_data.radar_chart.labels.map((label, i) => ({ subject: label, you: Math.round((result.visual_data.radar_chart.candidate_scores[i] || 0) * 100), required: 100 }))
    : [];
  const barData = result?.visual_data?.gaps_by_category?.map((cat) => ({ name: cat.category.charAt(0).toUpperCase() + cat.category.slice(1), gaps: cat.gap_count })) || [];
  const COLORS = ["#ef4444", "#f59e0b", "#6b7280"];

  return (
    <div className="page-container section animate-fade-in">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Skill Gap Analysis</h1>
        <p className="text-gray-600 mt-1">See exactly which skills you need for your target role.</p>
      </div>

      <Card className="mb-6">
        <div className="flex flex-col sm:flex-row gap-3 items-end">
          <div className="flex-1">
            <label className="label">Target Job</label>
            <select value={selectedJobId} onChange={(e) => setSelectedJobId(e.target.value)} className="input">
              {SAMPLE_JOBS.map((job) => <option key={job.id} value={job.id}>{job.title}</option>)}
            </select>
          </div>
          <Button onClick={runAnalysis} loading={isLoading}>Analyze</Button>
        </div>
      </Card>

      {error && <div className="mb-4 bg-danger-50 border border-danger-200 text-danger-700 rounded-lg px-4 py-3 text-sm" role="alert">{error}</div>}

      {isLoading ? <InlineLoader label="Analyzing your skill gaps..." /> : result ? (
        <div className="space-y-6">
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            {[{ value: `${Math.round(result.readiness_score)}%`, label: "Readiness", color: "text-primary-600" }, { value: result.critical_gap_count, label: "Critical Gaps", color: "text-danger-600" }, { value: result.important_gap_count, label: "Important Gaps", color: "text-warning-600" }, { value: `${result.total_estimated_learning_hours}h`, label: "Learning Hours", color: "text-gray-600" }].map((s) => (
              <Card key={s.label} padding="sm" className="text-center">
                <div className={`text-2xl font-bold ${s.color}`}>{s.value}</div>
                <div className="text-xs text-gray-500 mt-1">{s.label}</div>
              </Card>
            ))}
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {radarData.length > 0 && (
              <Card>
                <CardHeader title="Skills Coverage" subtitle="Your skills vs. requirements" />
                <ResponsiveContainer width="100%" height={250}>
                  <RadarChart data={radarData}>
                    <PolarGrid /><PolarAngleAxis dataKey="subject" tick={{ fontSize: 11, fill: "#6b7280" }} />
                    <Radar name="Required" dataKey="required" stroke="#e5e7eb" fill="#e5e7eb" fillOpacity={0.3} />
                    <Radar name="You" dataKey="you" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.5} />
                    <Tooltip formatter={(v: number) => [`${v}%`]} />
                  </RadarChart>
                </ResponsiveContainer>
              </Card>
            )}
            {barData.length > 0 && (
              <Card>
                <CardHeader title="Gaps by Category" subtitle="Number of gaps per category" />
                <ResponsiveContainer width="100%" height={250}>
                  <BarChart data={barData} margin={{ top: 5, right: 10, left: -20, bottom: 5 }}>
                    <XAxis dataKey="name" tick={{ fontSize: 12 }} /><YAxis tick={{ fontSize: 12 }} /><Tooltip />
                    <Bar dataKey="gaps" radius={[4, 4, 0, 0]}>{barData.map((_, index) => <Cell key={index} fill={COLORS[index % COLORS.length]} />)}</Bar>
                  </BarChart>
                </ResponsiveContainer>
              </Card>
            )}
          </div>

          {result.matched_skills?.length > 0 && (
            <Card><CardHeader title="✅ Skills You Already Have" subtitle={`${result.matched_skills.length} skills matched`} />
              <div className="flex flex-wrap gap-2">{result.matched_skills.map((skill) => <span key={skill} className="skill-tag">{skill}</span>)}</div>
            </Card>
          )}

          {result.critical_gaps?.length > 0 && (
            <Card><CardHeader title="🔴 Critical Gaps" subtitle="Must-have skills you're missing" />
              <div className="space-y-2">{result.critical_gaps.map((gap) => <GapCard key={gap.skill_name} gap={gap} />)}</div>
            </Card>
          )}

          {result.important_gaps?.length > 0 && (
            <Card><CardHeader title="🟡 Important Gaps" subtitle="Preferred skills that strengthen your application" />
              <div className="space-y-2">{result.important_gaps.map((gap) => <GapCard key={gap.skill_name} gap={gap} />)}</div>
            </Card>
          )}

          <Card className="bg-blue-50 border-blue-200">
            <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
              <div><p className="font-semibold text-blue-900">Ready to close these gaps?</p><p className="text-sm text-blue-700 mt-0.5">Get a personalized week-by-week learning plan.</p></div>
              <Link href={`/learning?job_id=${selectedJobId}`}><Button>Get Learning Plan →</Button></Link>
            </div>
          </Card>
        </div>
      ) : (
        <Card className="text-center py-12"><p className="text-gray-500">Select a job and click Analyze to see your skill gaps.</p></Card>
      )}
    </div>
  );
}

export default function AnalysisPage() {
  return (
    <Suspense fallback={<div className="page-container section"><p className="text-gray-500">Loading...</p></div>}>
      <AnalysisPageInner />
    </Suspense>
  );
}
