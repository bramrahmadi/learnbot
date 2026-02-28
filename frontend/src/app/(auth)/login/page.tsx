"use client";
import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card } from "@/components/ui/Card";

export default function LoginPage() {
  const router = useRouter();
  const { login, isLoading, error, clearError } = useAuthStore();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const validate = () => {
    const errors: Record<string, string> = {};
    if (!email) errors.email = "Email is required";
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) errors.email = "Enter a valid email address";
    if (!password) errors.password = "Password is required";
    return errors;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();
    const errors = validate();
    if (Object.keys(errors).length > 0) { setFieldErrors(errors); return; }
    setFieldErrors({});
    try { await login(email, password); router.push("/dashboard"); } catch { /* handled by store */ }
  };

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center py-12 px-4">
      <div className="w-full max-w-md animate-slide-up">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">Welcome back</h1>
          <p className="text-gray-600 mt-2">Sign in to continue your career journey</p>
        </div>
        <Card>
          <form onSubmit={handleSubmit} noValidate className="space-y-4">
            {error && <div className="bg-danger-50 border border-danger-200 text-danger-700 rounded-lg px-4 py-3 text-sm" role="alert">{error}</div>}
            <Input label="Email address" type="email" value={email} onChange={(e) => setEmail(e.target.value)} error={fieldErrors.email} placeholder="you@example.com" autoComplete="email" required />
            <Input label="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} error={fieldErrors.password} placeholder="••••••••" autoComplete="current-password" required />
            <Button type="submit" loading={isLoading} className="w-full" size="lg">Sign in</Button>
          </form>
          <div className="mt-6 text-center text-sm text-gray-600">
            Don&apos;t have an account?{" "}
            <Link href="/register" className="text-primary-600 font-medium hover:text-primary-700">Create one free</Link>
          </div>
        </Card>
      </div>
    </div>
  );
}
