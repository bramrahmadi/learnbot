import React from "react";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "danger" | "ghost";
  size?: "sm" | "md" | "lg";
  loading?: boolean;
  children: React.ReactNode;
}

export function Button({ variant = "primary", size = "md", loading = false, children, className = "", disabled, ...props }: ButtonProps) {
  const v = { primary: "btn-primary", secondary: "btn-secondary", danger: "btn-danger", ghost: "btn-ghost" };
  const s = { sm: "btn-sm", md: "", lg: "btn-lg" };
  return (
    <button className={`btn ${v[variant]} ${s[size]} ${className}`} disabled={disabled || loading} {...props}>
      {loading && <span className="spinner w-4 h-4" aria-hidden="true" />}
      {children}
    </button>
  );
}
