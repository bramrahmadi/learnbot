import React from "react";

interface LoadingSpinnerProps { size?: "sm" | "md" | "lg"; className?: string; label?: string; }
export function LoadingSpinner({ size = "md", className = "", label = "Loading..." }: LoadingSpinnerProps) {
  const s = { sm: "w-4 h-4", md: "w-8 h-8", lg: "w-12 h-12" };
  return (
    <div className={`flex flex-col items-center justify-center gap-3 ${className}`} role="status" aria-label={label}>
      <div className={`spinner ${s[size]}`} />
      {size !== "sm" && <p className="text-sm text-gray-500">{label}</p>}
    </div>
  );
}

export function PageLoader() {
  return <div className="min-h-screen flex items-center justify-center"><LoadingSpinner size="lg" label="Loading LearnBot..." /></div>;
}

export function InlineLoader({ label = "Loading..." }: { label?: string }) {
  return (
    <div className="flex items-center gap-2 text-sm text-gray-500 py-4">
      <LoadingSpinner size="sm" /><span>{label}</span>
    </div>
  );
}
