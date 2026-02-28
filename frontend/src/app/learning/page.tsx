"use client";
import React, { useState, useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { useToken } from "@/store/authStore";
import { analysisAPI, type LearningPlan, type LearningPhase, type SkillRecommendation, type WeeklySchedule } from "@/lib/api";
import { Card, CardHeader } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { CostBadge, DifficultyBadge } from "@/components/ui/Badge";
import { InlineLoader } from "@/components/ui/LoadingSpinner";

const SAMPLE_JOBS = [
  { id: "job-001", title: "Senior Backend Engineer (Go)" }, { id: "job-002", title: "Machine Learning Engineer" },
  { id: "job-003", title: "Full Stack JavaScript Developer" }, { id: "job-004", title: "DevOps Engineer" }, { id: "job-005", title: "Data Engineer" },
];

function ResourceCard({ rec }: { rec: SkillRecommendation["primary_resource"] }) {
  if (!rec) return null;
  const { resource } = rec;
  return (
    <div className="bg-white border border-gray-200 rounded-lg p-4">
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 min-w-0">
          <a href={resource.url} target="_blank" rel="noopener noreferrer" className="font-medium text-gray-900 hover:text-primary-600 transition-colors text-sm">{resource.title} ↗</a>
          <p className="text-xs text-gray-500 mt-0.5">{resource.provider}</p>
          <div className="flex flex-wrap gap-1.5 mt-2">
            <DifficultyBadge difficulty={resource.difficulty} />
            <CostBadge costType={resource.cost_type} />
            {resource.has_certificate && <span className="badge badge-blue">🏆 Certificate</span>}
            {resource.has_hands_on && <span className="badge badge-green">🛠 Hands-on</span>}
          </div>
        </div>
        <div className="text-right flex-shrink-0">
          {resource.rating && <div className="text-xs text-gray-500">⭐ {resource.rating.toFixed(1)}</div>}
          {resource.duration_label && <div className="text-xs text-gray-500 mt-0.5">⏱ {resource.duration_label}</div>}
        </div>
      </div>
      {rec.recommendation_reason && <p className="text-xs text-gray-500 mt-2 italic">{rec.recommendation_reason}</p>}
    </div>
  );
}

function PhaseCard({ phase, completedSkills, onToggleSkill }: { phase: LearningPhase; completedSkills: Set<string>; onToggleSkill: (skill: string) => void }) {
  const [expanded, setExpanded] = useState(phase.phase_number === 1);
  const completedCount = phase.skills.filter((s) => completedSkills.has(s.skill_name)).length;
  return (
    <Card className="overflow-hidden">
      <button className="w-full flex items-center justify-between p-6 text-left" onClick={() => setExpanded(!expanded)} aria-expanded={expanded}>
        <div className="flex items-center gap-3">
          <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold text-white ${phase.phase_number === 1 ? "bg-danger-500" : phase.phase_number === 2 ? "bg-warning-500" : "bg-gray-400"}`}>{phase.phase_number}</div>
          <div>
            <h3 className="font-semibold text-gray-900">{phase.phase_name}</h3>
            <p className="text-sm text-gray-500">{completedCount}/{phase.skills.length} skills · {phase.total_hours}h · ~{phase.estimated_weeks} weeks</p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          {completedCount === phase.skills.length && phase.skills.length > 0 && <span className="text-success-600 text-sm font-medium">✓ Complete</span>}
          <span className="text-gray-400">{expanded ? "▲" : "▼"}</span>
        </div>
      </button>
      {expanded && (
        <div className="px-6 pb-6 border-t border-gray-100">
          <p className="text-sm text-gray-600 mt-4 mb-4">{phase.phase_description}</p>
          <div className="space-y-4">
            {phase.skills.map((skillRec: SkillRecommendation) => (
              <div key={skillRec.skill_name} className="border border-gray-200 rounded-lg p-4">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <input type="checkbox" id={`skill-${skillRec.skill_name}`} checked={completedSkills.has(skillRec.skill_name)} onChange={() => onToggleSkill(skillRec.skill_name)} className="w-4 h-4 text-primary-600 rounded border-gray-300 focus:ring-primary-500" />
                    <label htmlFor={`skill-${skillRec.skill_name}`} className={`font-medium cursor-pointer ${completedSkills.has(skillRec.skill_name) ? "line-through text-gray-400" : "text-gray-900"}`}>{skillRec.skill_name}</label>
                  </div>
                  <span className="text-xs text-gray-500">~{skillRec.estimated_hours_to_job_ready}h</span>
                </div>
                {skillRec.primary_resource && <ResourceCard rec={skillRec.primary_resource} />}
                {skillRec.alternative_resources && skillRec.alternative_resources.length > 0 && (
                  <details className="mt-2">
                    <summary className="text-xs text-gray-500 cursor-pointer hover:text-gray-700">{skillRec.alternative_resources.length} alternative resource(s)</summary>
                    <div className="mt-2 space-y-2">{skillRec.alternative_resources.map((alt, i) => <ResourceCard key={i} rec={alt} />)}</div>
                  </details>
                )}
              </div>
            ))}
          </div>
          {phase.milestone && <div className="mt-4 bg-success-50 border border-success-200 rounded-lg p-3 text-sm text-success-700">🎯 Milestone: {phase.milestone}</div>}
        </div>
      )}
    </Card>
  );
}

function TimelineView({ weeks }: { weeks: WeeklySchedule[] }) {
  return (
    <div className="space-y-2">
      {weeks.map((week) => (
        <div key={week.week_number} className={`flex gap-3 p-3 rounded-lg border ${week.is_checkpoint ? "border-primary-200 bg-primary-50" : "border-gray-200 bg-white"}`}>
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center text-xs font-medium text-gray-600">{week.week_number}</div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between gap-2">
              <p className="text-sm font-medium text-gray-900 truncate">{week.skill_focus}</p>
              <span className="text-xs text-gray-500 flex-shrink-0">{week.hours_planned}h</span>
            </div>
            <p className="text-xs text-gray-500 truncate">{week.resource_title}</p>
            {week.is_checkpoint && week.checkpoint_description && <p className="text-xs text-primary-600 mt-1">✓ {week.checkpoint_description}</p>}
          </div>
        </div>
      ))}
    </div>
  );
}

function LearningPageInner() {
  const token = useToken();
  const searchParams = useSearchParams();
  const [selectedJobId, setSelectedJobId] = useState(searchParams.get("job_id") || "job-001");
  const [weeklyHours, setWeeklyHours] = useState(10);
  const [preferFree, setPreferFree] = useState(false);
  const [plan, setPlan] = useState<LearningPlan | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [completedSkills, setCompletedSkills] = useState<Set<string>>(new Set());
  const [activeTab, setActiveTab] = useState<"phases" | "timeline">("phases");

  useEffect(() => { if (token) generatePlan(); }, [token, selectedJobId]); // eslint-disable-line react-hooks/exhaustive-deps

  const generatePlan = async () => {
    if (!token) return;
    setIsLoading(true); setError(null);
    try { const data = await analysisAPI.getTrainingPlan(token, selectedJobId, undefined, { weekly_hours_available: weeklyHours, prefer_free: preferFree }); setPlan(data); }
    catch (err: unknown) { setError(err instanceof Error ? err.message : "Failed to generate plan"); }
    finally { setIsLoading(false); }
  };

  const toggleSkill = (skill: string) => setCompletedSkills((prev) => { const n = new Set(prev); n.has(skill) ? n.delete(skill) : n.add(skill); return n; });
  const totalSkills = plan?.phases.reduce((sum, p) => sum + p.skills.length, 0) || 0;
  const progressPct = totalSkills > 0 ? (completedSkills.size / totalSkills) * 100 : 0;

  return (
    <div className="page-container section animate-fade-in">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Learning Plan</h1>
        <p className="text-gray-600 mt-1">Your personalized week-by-week path to job readiness.</p>
      </div>

      <Card className="mb-6">
        <div className="flex flex-col sm:flex-row gap-4 items-end">
          <div className="flex-1">
            <label className="label">Target Job</label>
            <select value={selectedJobId} onChange={(e) => setSelectedJobId(e.target.value)} className="input">
              {SAMPLE_JOBS.map((job) => <option key={job.id} value={job.id}>{job.title}</option>)}
            </select>
          </div>
          <div className="w-32">
            <label className="label">Hours/week</label>
            <input type="number" value={weeklyHours} onChange={(e) => setWeeklyHours(Number(e.target.value))} min={1} max={40} className="input" />
          </div>
          <div className="flex items-center gap-2 pb-2">
            <input type="checkbox" id="prefer-free" checked={preferFree} onChange={(e) => setPreferFree(e.target.checked)} className="w-4 h-4 text-primary-600 rounded border-gray-300" />
            <label htmlFor="prefer-free" className="text-sm text-gray-700">Free only</label>
          </div>
          <Button onClick={generatePlan} loading={isLoading}>Generate Plan</Button>
        </div>
      </Card>

      {error && <div className="mb-4 bg-danger-50 border border-danger-200 text-danger-700 rounded-lg px-4 py-3 text-sm" role="alert">{error}</div>}

      {isLoading ? <InlineLoader label="Generating your learning plan..." /> : plan ? (
        <div className="space-y-6">
          <Card className="bg-gradient-to-r from-primary-600 to-primary-700 text-white">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <div>
                <p className="text-primary-100 text-sm">{plan.summary.headline}</p>
                <div className="flex flex-wrap gap-3 mt-2 text-sm">
                  <span>📅 {plan.timeline.total_weeks} weeks</span>
                  <span>⏱ {plan.timeline.total_hours}h total</span>
                  <span>💰 ${plan.summary.estimated_total_cost_usd.toFixed(0)} est. cost</span>
                </div>
              </div>
              <div className="text-right"><div className="text-3xl font-bold">{Math.round(plan.readiness_score)}%</div><div className="text-primary-200 text-xs">Current readiness</div></div>
            </div>
          </Card>

          {totalSkills > 0 && (
            <Card padding="sm">
              <div className="flex justify-between text-sm mb-2"><span className="text-gray-600">Your Progress</span><span className="font-medium">{completedSkills.size}/{totalSkills} skills completed</span></div>
              <div className="progress-bar"><div className="progress-fill bg-success-500" style={{ width: `${progressPct}%` }} /></div>
            </Card>
          )}

          {plan.summary.quick_wins?.length > 0 && (
            <Card><CardHeader title="⚡ Quick Wins" subtitle="Skills you can learn in under 20 hours" />
              <div className="flex flex-wrap gap-2">{plan.summary.quick_wins.map((skill) => <span key={skill} className="skill-tag">{skill}</span>)}</div>
            </Card>
          )}

          <div className="flex gap-1 bg-gray-100 rounded-lg p-1 w-fit">
            {(["phases", "timeline"] as const).map((tab) => (
              <button key={tab} className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${activeTab === tab ? "bg-white text-gray-900 shadow-sm" : "text-gray-600 hover:text-gray-900"}`} onClick={() => setActiveTab(tab)}>
                {tab === "phases" ? "Learning Phases" : "Weekly Timeline"}
              </button>
            ))}
          </div>

          {activeTab === "phases" && (
            <div className="space-y-4">{plan.phases.map((phase) => <PhaseCard key={phase.phase_number} phase={phase} completedSkills={completedSkills} onToggleSkill={toggleSkill} />)}</div>
          )}

          {activeTab === "timeline" && (
            <Card>
              <CardHeader title="Week-by-Week Schedule" subtitle={`${plan.timeline.total_weeks} weeks · ${plan.timeline.weekly_hours}h/week`} />
              {plan.timeline.target_completion_date && <p className="text-sm text-gray-600 mb-4">📅 Estimated completion: <strong>{new Date(plan.timeline.target_completion_date).toLocaleDateString("en-US", { month: "long", year: "numeric" })}</strong></p>}
              <TimelineView weeks={plan.timeline.weeks} />
            </Card>
          )}

          {plan.matched_skills?.length > 0 && (
            <Card><CardHeader title="✅ Skills You Already Have" subtitle={`${plan.matched_skills.length} skills matched`} />
              <div className="flex flex-wrap gap-2">{plan.matched_skills.map((skill) => <span key={skill} className="skill-tag">{skill}</span>)}</div>
            </Card>
          )}
        </div>
      ) : (
        <Card className="text-center py-12"><p className="text-gray-500">Select a job and click Generate Plan to get your personalized learning path.</p></Card>
      )}
    </div>
  );
}

export default function LearningPage() {
  return (
    <Suspense fallback={<div className="page-container section"><p className="text-gray-500">Loading...</p></div>}>
      <LearningPageInner />
    </Suspense>
  );
}
