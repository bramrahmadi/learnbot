"use client";
import React, { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore, useToken } from "@/store/authStore";
import { resumeAPI, profileAPI, type Skill } from "@/lib/api";
import { Button } from "@/components/ui/Button";
import { Input, Textarea } from "@/components/ui/Input";
import { Card } from "@/components/ui/Card";

type Step = "upload" | "profile" | "skills" | "preferences";
const STEPS: { id: Step; label: string }[] = [
  { id: "upload", label: "Resume" }, { id: "profile", label: "Profile" },
  { id: "skills", label: "Skills" }, { id: "preferences", label: "Preferences" },
];

export default function OnboardingPage() {
  const router = useRouter();
  const token = useToken();
  const { updateProfile } = useAuthStore();
  const [currentStep, setCurrentStep] = useState<Step>("upload");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [headline, setHeadline] = useState("");
  const [summary, setSummary] = useState("");
  const [locationCity, setLocationCity] = useState("");
  const [yearsExp, setYearsExp] = useState("");
  const [isOpenToWork, setIsOpenToWork] = useState(true);
  const [skills, setSkills] = useState<Skill[]>([]);
  const [newSkillName, setNewSkillName] = useState("");
  const [newSkillProficiency, setNewSkillProficiency] = useState("intermediate");
  const [targetRole, setTargetRole] = useState("");

  const currentStepIndex = STEPS.findIndex((s) => s.id === currentStep);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault(); setIsDragging(false);
    const f = e.dataTransfer.files[0]; if (f) setFile(f);
  }, []);

  const handleUploadResume = async () => {
    if (!file || !token) return;
    setIsLoading(true); setError(null);
    try {
      const data = await resumeAPI.upload(token, file);
      if (data.personal?.location) setLocationCity(data.personal.location);
      if (data.skills?.length > 0) setSkills(data.skills.slice(0, 20).map((s: { name: string }) => ({ name: s.name, proficiency: "intermediate", is_primary: false })));
      setCurrentStep("profile");
    } catch (err: unknown) { setError(err instanceof Error ? err.message : "Failed to parse resume"); }
    finally { setIsLoading(false); }
  };

  const handleSaveProfile = async () => {
    if (!token) return;
    setIsLoading(true); setError(null);
    try {
      await updateProfile({ headline, summary, location_city: locationCity, years_of_experience: yearsExp ? parseFloat(yearsExp) : undefined, is_open_to_work: isOpenToWork });
      setCurrentStep("skills");
    } catch (err: unknown) { setError(err instanceof Error ? err.message : "Failed to save profile"); }
    finally { setIsLoading(false); }
  };

  const addSkill = () => {
    if (!newSkillName.trim()) return;
    setSkills((prev) => [...prev, { name: newSkillName.trim(), proficiency: newSkillProficiency, is_primary: false }]);
    setNewSkillName("");
  };

  const handleSaveSkills = async () => {
    if (!token) return;
    setIsLoading(true); setError(null);
    try { await profileAPI.updateSkills(token, skills); setCurrentStep("preferences"); }
    catch (err: unknown) { setError(err instanceof Error ? err.message : "Failed to save skills"); }
    finally { setIsLoading(false); }
  };

  const handleFinish = async () => {
    if (!token) return;
    setIsLoading(true);
    try { if (targetRole) await updateProfile({ headline: headline || `Targeting: ${targetRole}` }); router.push("/dashboard"); }
    catch { router.push("/dashboard"); }
  };

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-gray-50 py-8 px-4">
      <div className="max-w-2xl mx-auto">
        {/* Progress */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-2">
            {STEPS.map((step, index) => (
              <React.Fragment key={step.id}>
                <div className="flex flex-col items-center">
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors ${index < currentStepIndex ? "bg-primary-600 text-white" : index === currentStepIndex ? "bg-primary-600 text-white ring-4 ring-primary-100" : "bg-gray-200 text-gray-500"}`}>
                    {index < currentStepIndex ? "✓" : index + 1}
                  </div>
                  <span className="text-xs text-gray-500 mt-1 hidden sm:block">{step.label}</span>
                </div>
                {index < STEPS.length - 1 && <div className={`flex-1 h-0.5 mx-2 transition-colors ${index < currentStepIndex ? "bg-primary-600" : "bg-gray-200"}`} />}
              </React.Fragment>
            ))}
          </div>
        </div>

        {error && <div className="mb-4 bg-danger-50 border border-danger-200 text-danger-700 rounded-lg px-4 py-3 text-sm" role="alert">{error}</div>}

        {currentStep === "upload" && (
          <Card>
            <h2 className="text-xl font-bold text-gray-900 mb-2">Upload your resume</h2>
            <p className="text-gray-600 text-sm mb-6">We&apos;ll automatically extract your skills, experience, and education. Supports PDF and DOCX files.</p>
            <div className={`border-2 border-dashed rounded-xl p-8 text-center transition-colors ${isDragging ? "border-primary-400 bg-primary-50" : file ? "border-success-500 bg-success-50" : "border-gray-300 hover:border-gray-400"}`}
              onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }} onDragLeave={() => setIsDragging(false)} onDrop={handleDrop}>
              {file ? (
                <div>
                  <div className="text-3xl mb-2">📄</div>
                  <p className="font-medium text-gray-900">{file.name}</p>
                  <p className="text-sm text-gray-500 mt-1">{(file.size / 1024).toFixed(0)} KB</p>
                  <button onClick={() => setFile(null)} className="text-sm text-danger-600 hover:text-danger-700 mt-2">Remove</button>
                </div>
              ) : (
                <div>
                  <div className="text-4xl mb-3">📤</div>
                  <p className="font-medium text-gray-700 mb-1">Drag & drop your resume here</p>
                  <p className="text-sm text-gray-500 mb-4">or</p>
                  <label className="btn-secondary btn cursor-pointer">Browse files<input type="file" accept=".pdf,.docx" onChange={(e) => { const f = e.target.files?.[0]; if (f) setFile(f); }} className="sr-only" /></label>
                  <p className="text-xs text-gray-400 mt-3">PDF or DOCX, max 10MB</p>
                </div>
              )}
            </div>
            <div className="flex gap-3 mt-6">
              <Button onClick={handleUploadResume} disabled={!file} loading={isLoading} className="flex-1">Parse Resume</Button>
              <Button variant="ghost" onClick={() => setCurrentStep("profile")}>Skip for now</Button>
            </div>
          </Card>
        )}

        {currentStep === "profile" && (
          <Card>
            <h2 className="text-xl font-bold text-gray-900 mb-2">Set up your profile</h2>
            <p className="text-gray-600 text-sm mb-6">Tell us about yourself to get personalized recommendations.</p>
            <div className="space-y-4">
              <Input label="Professional headline" value={headline} onChange={(e) => setHeadline(e.target.value)} placeholder="e.g. Senior Software Engineer" />
              <Textarea label="Professional summary" value={summary} onChange={(e) => setSummary(e.target.value)} placeholder="Brief overview of your experience and goals..." rows={3} />
              <div className="grid grid-cols-2 gap-4">
                <Input label="City" value={locationCity} onChange={(e) => setLocationCity(e.target.value)} placeholder="San Francisco" />
                <Input label="Years of experience" type="number" value={yearsExp} onChange={(e) => setYearsExp(e.target.value)} placeholder="5" min="0" max="50" />
              </div>
              <div className="flex items-center gap-3">
                <input type="checkbox" id="open-to-work" checked={isOpenToWork} onChange={(e) => setIsOpenToWork(e.target.checked)} className="w-4 h-4 text-primary-600 rounded border-gray-300 focus:ring-primary-500" />
                <label htmlFor="open-to-work" className="text-sm text-gray-700">I&apos;m open to new opportunities</label>
              </div>
            </div>
            <div className="flex gap-3 mt-6">
              <Button onClick={handleSaveProfile} loading={isLoading} className="flex-1">Continue</Button>
              <Button variant="ghost" onClick={() => setCurrentStep("skills")}>Skip</Button>
            </div>
          </Card>
        )}

        {currentStep === "skills" && (
          <Card>
            <h2 className="text-xl font-bold text-gray-900 mb-2">Your skills</h2>
            <p className="text-gray-600 text-sm mb-6">{skills.length > 0 ? `We found ${skills.length} skills. Add or remove as needed.` : "Add your technical and professional skills."}</p>
            <div className="flex gap-2 mb-4">
              <Input value={newSkillName} onChange={(e) => setNewSkillName(e.target.value)} placeholder="Add a skill (e.g. Python)" onKeyDown={(e) => e.key === "Enter" && addSkill()} className="flex-1" />
              <select value={newSkillProficiency} onChange={(e) => setNewSkillProficiency(e.target.value)} className="input w-36">
                <option value="beginner">Beginner</option><option value="intermediate">Intermediate</option>
                <option value="advanced">Advanced</option><option value="expert">Expert</option>
              </select>
              <Button onClick={addSkill} variant="secondary">Add</Button>
            </div>
            <div className="flex flex-wrap gap-2 min-h-[80px] p-3 bg-gray-50 rounded-lg border border-gray-200">
              {skills.length === 0 ? <p className="text-sm text-gray-400 self-center w-full text-center">No skills added yet</p> :
                skills.map((skill, index) => (
                  <div key={index} className="inline-flex items-center gap-1.5 bg-white border border-gray-200 rounded-full px-3 py-1 text-sm">
                    <span className="font-medium text-gray-800">{skill.name}</span>
                    <span className="text-gray-400 text-xs">·</span>
                    <span className="text-gray-500 text-xs capitalize">{skill.proficiency}</span>
                    <button onClick={() => setSkills((prev) => prev.filter((_, i) => i !== index))} className="ml-1 text-gray-400 hover:text-danger-500" aria-label={`Remove ${skill.name}`}>×</button>
                  </div>
                ))}
            </div>
            <div className="flex gap-3 mt-6">
              <Button onClick={handleSaveSkills} loading={isLoading} className="flex-1">Continue</Button>
              <Button variant="ghost" onClick={() => setCurrentStep("preferences")}>Skip</Button>
            </div>
          </Card>
        )}

        {currentStep === "preferences" && (
          <Card>
            <h2 className="text-xl font-bold text-gray-900 mb-2">Career preferences</h2>
            <p className="text-gray-600 text-sm mb-6">Tell us about your career goals to get better job recommendations.</p>
            <div className="space-y-4">
              <Input label="Target role" value={targetRole} onChange={(e) => setTargetRole(e.target.value)} placeholder="e.g. Senior Backend Engineer" hint="What role are you aiming for?" />
              <div className="bg-primary-50 rounded-lg p-4 text-sm text-primary-700">
                <p className="font-medium mb-1">🎉 You&apos;re all set!</p>
                <p>Your profile is ready. Head to the dashboard to see job recommendations and your personalized learning plan.</p>
              </div>
            </div>
            <div className="flex gap-3 mt-6">
              <Button onClick={handleFinish} loading={isLoading} className="flex-1">Go to Dashboard →</Button>
            </div>
          </Card>
        )}
      </div>
    </div>
  );
}
