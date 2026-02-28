"use client";
import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card } from "@/components/ui/Card";

export default function RegisterPage() {
  const router = useRouter();
  const { register, isLoading, error, clearError } = useAuthStore();
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const validate = () => {
    const errors: Record<string, string> = {};
    if (!fullName.trim()) errors.full_name = "Full name is required";
    if (!email) errors.email = "Email is required";
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) errors.email = "Enter a valid email address";
    if (!password) errors.password = "Password is required";
    else if (password.length < 8) errors.password = "Password must be at least 8 characters";
    return errors;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();
    const errors = validate();
    if (Object.keys(errors).length > 0) { setFieldErrors(errors); return; }
    setFieldErrors({});
    try { await register(email, password, fullName); router.push("/onboarding"); } catch { /* handled by store */ }
  };

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center py-12 px-4">
      <div className="w-full max-w-md animate-slide-up">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">Start your career journey</h1>
          <p className="text-gray-600 mt-2">Create your free account in seconds</p>
        </div>
        <Card>
          <form onSubmit={handleSubmit} noValidate className="space-y-4">
            {error && <div className="bg-danger-50 border border-danger-200 text-danger-700 rounded-lg px-4 py-3 text-sm" role="alert">{error}</div>}
            <Input label="Full name" type="text" value={fullName} onChange={(e) => setFullName(e.target.value)} error={fieldErrors.full_name} placeholder="Jane Doe" autoComplete="name" required />
            <Input label="Email address" type="email" value={email} onChange={(e) => setEmail(e.target.value)} error={fieldErrors.email} placeholder="you@example.com" autoComplete="email" required />
            <Input label="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} error={fieldErrors.password} placeholder="Min. 8 characters" autoComplete="new-password" required hint="At least 8 characters" />
            <Button type="submit" loading={isLoading} className="w-full" size="lg">Create free account</Button>
            <p className="text-xs text-gray-500 text-center">By creating an account, you agree to our Terms of Service and Privacy Policy.</p>
          </form>
          <div className="mt-6 text-center text-sm text-gray-600">
            Already have an account?{" "}
            <Link href="/login" className="text-primary-600 font-medium hover:text-primary-700">Sign in</Link>
          </div>
        </Card>
      </div>
    </div>
  );
}
