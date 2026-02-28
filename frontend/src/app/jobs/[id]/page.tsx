"use client";
import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useToken } from "@/store/authStore";
import { jobsAPI, type Job, type JobMatch } from "@/lib/api";
import { Card, CardHeader } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { PageLoader } from "@/components/ui/LoadingSpinner";

function ScoreBar({ label, value }: { label: string; value: number }) {
  const pct = Math.round(value * 100);
  const color = pct >= 75 ? "bg-success-500" : pct >= 50 ? "bg-warning-500" : "bg-danger-500";
  return (
    <div className="mb-3">
      <div className="flex justify-between text-sm mb-1"><span className="text-gray-600">{label}</span><span className="font-medium text-gray-900">{pct}%</span></div>
      <div className="progress-bar"><div className={`progress-fill ${color}`} style={{ width: `${pct}%` }} /></div>
    </div>
  );
}

export default function JobDetailPage() {
  const params = useParams();
  const jobId = params.id as string;
  const token = useToken();
  const [job, setJob] = useState<Job | null>(null);
  const [match, setMatch] = useState<JobMatch | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [loadingMatch, setLoadingMatch] = useState(false);

  useEffect(() => { loadJob(); }, [jobId]); // eslint-disable-line react-hooks/exhaustive-deps
  useEffect(() => { if (token && job) loadMatch(); }, [token, job]); // eslint-disable-line react-hooks/exhaustive-deps

  const loadJob = async () => {
    setIsLoading(true);
    try { const data = await jobsAPI.getJob(jobId); setJob(data); }
    catch { setJob(null); } finally { setIsLoading(false); }
  };

  const loadMatch = async () => {
    if (!token) return;
    setLoadingMatch(true);
    try { const data = await jobsAPI.getJobMatch(token, jobId); setMatch(data); }
    catch { /* silently fail */ } finally { setLoadingMatch(false); }
  };

  if (isLoading) return <PageLoader />;
  if (!job) return (
    <div className="page-container section text-center">
      <p className="text-gray-500">Job not found.</p>
      <Link href="/jobs" className="text-primary-600 hover:underline mt-2 block">← Back to jobs</Link>
    </div>
  );

  const overallScore = match?.overall_score ?? 0;

  return (
    <div className="page-container section animate-fade-in">
      <Link href="/jobs" className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 mb-6">← Back to jobs</Link>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <div className="flex items-start justify-between gap-4">
              <div>
                <h1 className="text-xl font-bold text-gray-900">{job.title}</h1>
                <p className="text-gray-600 mt-1">{job.company}</p>
                <div className="flex flex-wrap gap-2 mt-3">
                  <Badge variant="gray">{job.location_type}</Badge>
                  <Badge variant="gray">{job.experience_level}</Badge>
                  {job.industry && <Badge variant="blue">{job.industry}</Badge>}
                  {job.location_city && <Badge variant="gray">📍 {job.location_city}</Badge>}
                </div>
                {job.salary_min && job.salary_max && <p className="text-sm text-gray-600 mt-2">💰 ${job.salary_min.toLocaleString()} – ${job.salary_max.toLocaleString()} {job.salary_currency}</p>}
              </div>
              {job.apply_url && <a href={job.apply_url} target="_blank" rel="noopener noreferrer" className="btn-primary btn flex-shrink-0">Apply Now ↗</a>}
            </div>
          </Card>

          {job.description && <Card><CardHeader title="Job Description" /><p className="text-sm text-gray-700 leading-relaxed whitespace-pre-line">{job.description}</p></Card>}

          <Card>
            <CardHeader title="Required Skills" />
            <div className="flex flex-wrap gap-2 mb-4">
              {job.required_skills?.map((skill) => {
                const isMissing = match?.missing_skills?.includes(skill);
                return <span key={skill} className={isMissing ? "skill-tag-missing" : "skill-tag"}>{isMissing ? "✗ " : "✓ "}{skill}</span>;
              })}
            </div>
            {job.preferred_skills && job.preferred_skills.length > 0 && (
              <><p className="text-sm font-medium text-gray-700 mb-2">Preferred Skills</p>
              <div className="flex flex-wrap gap-2">{job.preferred_skills.map((skill) => <span key={skill} className="skill-tag">{skill}</span>)}</div></>
            )}
          </Card>
        </div>

        <div className="space-y-6">
          <Card>
            <CardHeader title="Your Match Score" />
            {loadingMatch ? <div className="text-center py-4 text-sm text-gray-500">Calculating...</div> : match ? (
              <div>
                <div className="text-center mb-4">
                  <div className={`inline-flex items-center justify-center w-20 h-20 rounded-full text-2xl font-bold text-white ${overallScore >= 75 ? "bg-success-500" : overallScore >= 50 ? "bg-warning-500" : "bg-danger-500"}`}>{Math.round(overallScore)}%</div>
                  <p className="text-sm text-gray-600 mt-2">{match.recommendation}</p>
                </div>
                <ScoreBar label="Skill Match" value={match.skill_match} />
                <ScoreBar label="Experience" value={match.experience_match} />
                <ScoreBar label="Education" value={match.education_match} />
                <ScoreBar label="Location Fit" value={match.location_fit} />
                <ScoreBar label="Industry" value={match.industry_match} />
                {match.missing_skills?.length > 0 && (
                  <div className="mt-4 pt-4 border-t border-gray-200">
                    <p className="text-xs font-medium text-gray-500 mb-2">Missing skills ({match.missing_skills.length})</p>
                    <div className="flex flex-wrap gap-1">{match.missing_skills.map((skill) => <span key={skill} className="skill-tag-missing text-xs">{skill}</span>)}</div>
                  </div>
                )}
              </div>
            ) : (
              <div className="text-center py-4"><p className="text-sm text-gray-500 mb-3">Sign in to see your match score</p><Link href="/login"><Button size="sm">Sign in</Button></Link></div>
            )}
          </Card>

          <Card>
            <div className="space-y-2">
              <Link href={`/analysis?job_id=${jobId}`} className="block"><Button variant="secondary" size="sm" className="w-full">📊 Analyze Skill Gaps</Button></Link>
              <Link href={`/learning?job_id=${jobId}`} className="block"><Button variant="secondary" size="sm" className="w-full">📚 Get Learning Plan</Button></Link>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
