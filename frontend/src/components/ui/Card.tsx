import React from "react";

interface CardProps { children: React.ReactNode; className?: string; hover?: boolean; padding?: "none" | "sm" | "md" | "lg"; }
export function Card({ children, className = "", hover = false, padding = "md" }: CardProps) {
  const p = { none: "", sm: "p-4", md: "p-6", lg: "p-8" };
  return <div className={`${hover ? "card-hover" : "card"} ${p[padding]} ${className}`}>{children}</div>;
}

interface CardHeaderProps { title: string; subtitle?: string; action?: React.ReactNode; icon?: React.ReactNode; }
export function CardHeader({ title, subtitle, action, icon }: CardHeaderProps) {
  return (
    <div className="flex items-start justify-between mb-4">
      <div className="flex items-center gap-3">
        {icon && <div className="flex-shrink-0 w-10 h-10 rounded-lg bg-primary-50 flex items-center justify-center text-primary-600">{icon}</div>}
        <div>
          <h3 className="text-base font-semibold text-gray-900">{title}</h3>
          {subtitle && <p className="text-sm text-gray-500 mt-0.5">{subtitle}</p>}
        </div>
      </div>
      {action && <div className="flex-shrink-0 ml-4">{action}</div>}
    </div>
  );
}
