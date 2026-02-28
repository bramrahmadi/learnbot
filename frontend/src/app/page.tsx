"use client";
import React from "react";
import Link from "next/link";
import { useIsAuthenticated } from "@/store/authStore";
import { Button } from "@/components/ui/Button";

const features = [
  { icon: "📄", title: "Resume Analysis", description: "Upload your resume and our AI extracts your skills, experience, and education automatically." },
  { icon: "🎯", title: "Job Matching", description: "Get an acceptance likelihood score for any job. Know exactly where you stand before applying." },
  { icon: "🔍", title: "Skill Gap Analysis", description: "Identify critical, important, and nice-to-have skills you're missing for your target role." },
  { icon: "📚", title: "Personalized Learning", description: "Get a week-by-week learning plan with curated courses, certifications, and resources." },
  { icon: "📈", title: "Progress Tracking", description: "Track your learning progress and watch your readiness score improve over time." },
  { icon: "🤖", title: "AI-Powered", description: "Powered by advanced algorithms that understand skill relationships and career paths." },
];

export default function LandingPage() {
  const isAuthenticated = useIsAuthenticated();
  return (
    <div className="animate-fade-in">
      {/* Hero */}
      <section className="relative overflow-hidden bg-gradient-to-br from-primary-600 via-primary-700 to-primary-900 text-white">
        <div className="relative page-container py-20 sm:py-28">
          <div className="max-w-3xl mx-auto text-center">
            <div className="inline-flex items-center gap-2 bg-white/10 rounded-full px-4 py-1.5 text-sm font-medium mb-6">
              <span className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
              AI-Powered Career Development
            </div>
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight mb-6">
              Land Your Dream Job <span className="text-primary-200">Faster</span>
            </h1>
            <p className="text-lg sm:text-xl text-primary-100 mb-8 max-w-2xl mx-auto">
              LearnBot analyzes your skills, identifies gaps, and creates a personalized learning plan to get you job-ready — faster than ever.
            </p>
            <div className="flex flex-col sm:flex-row gap-3 justify-center">
              {isAuthenticated ? (
                <Link href="/dashboard"><Button size="lg" className="w-full sm:w-auto bg-white text-primary-700 hover:bg-primary-50">Go to Dashboard →</Button></Link>
              ) : (
                <>
                  <Link href="/register"><Button size="lg" className="w-full sm:w-auto bg-white text-primary-700 hover:bg-primary-50 focus:ring-white">Get started free →</Button></Link>
                  <Link href="/login"><Button size="lg" variant="ghost" className="w-full sm:w-auto text-white border border-white/30 hover:bg-white/10">Sign in</Button></Link>
                </>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* Stats */}
      <section className="bg-white border-b border-gray-200">
        <div className="page-container py-8">
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-6">
            {[{ value: "60+", label: "Curated Resources" }, { value: "5", label: "Sample Jobs" }, { value: "100%", label: "Free to Start" }, { value: "∞", label: "Career Paths" }].map((s) => (
              <div key={s.label} className="text-center">
                <div className="text-3xl font-bold text-primary-600">{s.value}</div>
                <div className="text-sm text-gray-500 mt-1">{s.label}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="section bg-gray-50">
        <div className="page-container">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">Everything you need to accelerate your career</h2>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">From resume parsing to personalized learning plans, LearnBot covers every step of your career development journey.</p>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((f) => (
              <div key={f.title} className="card p-6 hover:shadow-md transition-shadow duration-200">
                <div className="text-3xl mb-3">{f.icon}</div>
                <h3 className="text-base font-semibold text-gray-900 mb-2">{f.title}</h3>
                <p className="text-sm text-gray-600">{f.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="section bg-primary-600 text-white">
        <div className="page-container text-center">
          <h2 className="text-3xl font-bold mb-4">Ready to accelerate your career?</h2>
          <p className="text-primary-100 text-lg mb-8 max-w-xl mx-auto">Join thousands of professionals using LearnBot to close skill gaps and land their dream jobs.</p>
          {!isAuthenticated && <Link href="/register"><Button size="lg" className="bg-white text-primary-700 hover:bg-primary-50 focus:ring-white">Start for free →</Button></Link>}
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-400 py-8">
        <div className="page-container">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <div className="w-6 h-6 bg-primary-600 rounded flex items-center justify-center"><span className="text-white font-bold text-xs">LB</span></div>
              <span className="text-white font-semibold">LearnBot</span>
            </div>
            <p className="text-sm">© {new Date().getFullYear()} LearnBot. AI-powered career development.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
