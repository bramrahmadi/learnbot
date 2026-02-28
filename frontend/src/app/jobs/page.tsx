"use client";
import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useToken } from "@/store/authStore";
import { jobsAPI, type Job, type JobSearchParams } from "@/lib/api";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Input, Select } from "@/components/ui/Input";
import { Badge, ScoreBadge } from "@/components/ui/Badge";
import { InlineLoader } from "@/components/ui/LoadingSpinner";

export default function JobsPage() {
  const token = useToken();
  const [jobs, setJobs] = useState<Job[]>([]);
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [savedJobs, setSavedJobs] = useState<Set<string>>(new Set());
  const [query, setQuery] = useState("");
  const [locationFilter, setLocationFilter] = useState("");
  const [levelFilter, setLevelFilter] = useState("");
  const [offset, setOffset] = useState(0);
  const LIMIT = 10;

  useEffect(() => { if (token) searchJobs(); }, [token, locationFilter, levelFilter, offset]); // eslint-disable-line react-hooks/exhaustive-deps

  const searchJobs = async () => {
    if (!token) return;
    setIsLoading(true);
    try {
      const params: JobSearchParams = { query: query || undefined, location_type: locationFilter || undefined, experience_level: levelFilter || undefined, limit: LIMIT, offset };
      const result = await jobsAPI.search(token, params);
      setJobs(result.jobs || []); setTotal(result.meta?.total || 0);
    } catch { setJobs([]); } finally { setIsLoading(false); }
  };

  const handleSearch = (e: React.FormEvent) => { e.preventDefault(); setOffset(0); searchJobs(); };

  return (
    <div className="page-container section animate-fade-in">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Job Search</h1>
        <p className="text-gray-600 mt-1">Find jobs that match your skills and get your acceptance likelihood score.</p>
      </div>

      <Card className="mb-6">
        <form onSubmit={handleSearch} className="flex flex-col sm:flex-row gap-3">
          <Input value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Search jobs, companies, skills..." className="flex-1" />
          <Select options={[{ value: "", label: "Any location" }, { value: "remote", label: "Remote" }, { value: "hybrid", label: "Hybrid" }, { value: "on_site", label: "On-site" }]} value={locationFilter} onChange={(e) => { setLocationFilter(e.target.value); setOffset(0); }} className="sm:w-40" />
          <Select options={[{ value: "", label: "Any level" }, { value: "entry", label: "Entry" }, { value: "mid", label: "Mid" }, { value: "senior", label: "Senior" }, { value: "lead", label: "Lead" }]} value={levelFilter} onChange={(e) => { setLevelFilter(e.target.value); setOffset(0); }} className="sm:w-36" />
          <Button type="submit" loading={isLoading}>Search</Button>
        </form>
      </Card>

      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-gray-600">{isLoading ? "Searching..." : `${total} job${total !== 1 ? "s" : ""} found`}</p>
      </div>

      {isLoading ? <InlineLoader label="Searching jobs..." /> : jobs.length === 0 ? (
        <Card className="text-center py-12"><p className="text-gray-500 text-lg mb-2">No jobs found</p><p className="text-gray-400 text-sm">Try adjusting your search filters</p></Card>
      ) : (
        <div className="space-y-4">
          {jobs.map((job) => (
            <Card key={job.id} hover>
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-start gap-3">
                    <div className="flex-1">
                      <Link href={`/jobs/${job.id}`} className="font-semibold text-gray-900 hover:text-primary-600 transition-colors">{job.title}</Link>
                      <p className="text-sm text-gray-600 mt-0.5">{job.company}</p>
                    </div>
                    {job.match_score !== undefined && <ScoreBadge score={job.match_score} />}
                  </div>
                  <div className="flex flex-wrap gap-2 mt-3">
                    <Badge variant="gray">{job.location_type}</Badge>
                    <Badge variant="gray">{job.experience_level}</Badge>
                    {job.required_skills?.slice(0, 3).map((skill) => <span key={skill} className="skill-tag">{skill}</span>)}
                    {(job.required_skills?.length || 0) > 3 && <span className="skill-tag">+{job.required_skills.length - 3}</span>}
                  </div>
                </div>
                <div className="flex flex-col items-end gap-2 flex-shrink-0">
                  <button onClick={() => setSavedJobs((prev) => { const n = new Set(prev); n.has(job.id) ? n.delete(job.id) : n.add(job.id); return n; })}
                    className={`p-1.5 rounded-lg transition-colors ${savedJobs.has(job.id) ? "text-warning-500 bg-warning-50" : "text-gray-400 hover:text-warning-500 hover:bg-warning-50"}`}
                    aria-label={savedJobs.has(job.id) ? "Unsave job" : "Save job"}>
                    {savedJobs.has(job.id) ? "★" : "☆"}
                  </button>
                  <Link href={`/jobs/${job.id}`}><Button size="sm" variant="secondary">View →</Button></Link>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      {total > LIMIT && (
        <div className="flex items-center justify-center gap-3 mt-6">
          <Button variant="secondary" size="sm" disabled={offset === 0} onClick={() => setOffset(Math.max(0, offset - LIMIT))}>← Previous</Button>
          <span className="text-sm text-gray-600">Page {Math.floor(offset / LIMIT) + 1} of {Math.ceil(total / LIMIT)}</span>
          <Button variant="secondary" size="sm" disabled={offset + LIMIT >= total} onClick={() => setOffset(offset + LIMIT)}>Next →</Button>
        </div>
      )}
    </div>
  );
}
