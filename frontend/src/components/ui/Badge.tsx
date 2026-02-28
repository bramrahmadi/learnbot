import React from "react";

type BadgeVariant = "blue" | "green" | "yellow" | "red" | "gray";
interface BadgeProps { children: React.ReactNode; variant?: BadgeVariant; className?: string; }

export function Badge({ children, variant = "gray", className = "" }: BadgeProps) {
  const v: Record<BadgeVariant, string> = { blue: "badge-blue", green: "badge-green", yellow: "badge-yellow", red: "badge-red", gray: "badge-gray" };
  return <span className={`${v[variant]} ${className}`}>{children}</span>;
}

export function ScoreBadge({ score }: { score: number }) {
  const variant: BadgeVariant = score >= 75 ? "green" : score >= 50 ? "yellow" : "red";
  return <Badge variant={variant}>{Math.round(score)}% match</Badge>;
}

export function GapCategoryBadge({ category }: { category: string }) {
  const config: Record<string, { variant: BadgeVariant; label: string }> = {
    critical: { variant: "red", label: "Critical" },
    important: { variant: "yellow", label: "Important" },
    nice_to_have: { variant: "gray", label: "Nice to have" },
  };
  const { variant, label } = config[category] || { variant: "gray" as BadgeVariant, label: category };
  return <Badge variant={variant}>{label}</Badge>;
}

export function DifficultyBadge({ difficulty }: { difficulty: string }) {
  const config: Record<string, { variant: BadgeVariant }> = {
    beginner: { variant: "green" }, intermediate: { variant: "blue" },
    advanced: { variant: "yellow" }, expert: { variant: "red" }, all_levels: { variant: "gray" },
  };
  const { variant } = config[difficulty] || { variant: "gray" as BadgeVariant };
  return <Badge variant={variant}>{difficulty.charAt(0).toUpperCase() + difficulty.slice(1).replace("_", " ")}</Badge>;
}

export function CostBadge({ costType }: { costType: string }) {
  const isFree = costType === "free" || costType === "free_audit";
  return <Badge variant={isFree ? "green" : "gray"}>{isFree ? "Free" : costType === "paid" ? "Paid" : costType.replace("_", " ")}</Badge>;
}
