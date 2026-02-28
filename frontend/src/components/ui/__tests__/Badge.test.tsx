import React from "react";
import { render, screen } from "@testing-library/react";
import { Badge, ScoreBadge, GapCategoryBadge, DifficultyBadge, CostBadge } from "../Badge";

describe("Badge", () => {
  it("renders children", () => { render(<Badge>Test</Badge>); expect(screen.getByText("Test")).toBeInTheDocument(); });
  it("applies blue variant", () => { render(<Badge variant="blue">Blue</Badge>); expect(screen.getByText("Blue").className).toContain("badge-blue"); });
  it("applies green variant", () => { render(<Badge variant="green">Green</Badge>); expect(screen.getByText("Green").className).toContain("badge-green"); });
  it("applies red variant", () => { render(<Badge variant="red">Red</Badge>); expect(screen.getByText("Red").className).toContain("badge-red"); });
});

describe("ScoreBadge", () => {
  it("shows green for high scores", () => { render(<ScoreBadge score={80} />); expect(screen.getByText("80% match").className).toContain("badge-green"); });
  it("shows yellow for medium scores", () => { render(<ScoreBadge score={60} />); expect(screen.getByText("60% match").className).toContain("badge-yellow"); });
  it("shows red for low scores", () => { render(<ScoreBadge score={30} />); expect(screen.getByText("30% match").className).toContain("badge-red"); });
  it("rounds score to nearest integer", () => { render(<ScoreBadge score={75.7} />); expect(screen.getByText("76% match")).toBeInTheDocument(); });
});

describe("GapCategoryBadge", () => {
  it("shows Critical for critical category", () => { render(<GapCategoryBadge category="critical" />); expect(screen.getByText("Critical")).toBeInTheDocument(); });
  it("shows Important for important category", () => { render(<GapCategoryBadge category="important" />); expect(screen.getByText("Important")).toBeInTheDocument(); });
  it("shows Nice to have for nice_to_have category", () => { render(<GapCategoryBadge category="nice_to_have" />); expect(screen.getByText("Nice to have")).toBeInTheDocument(); });
});

describe("DifficultyBadge", () => {
  it("shows Beginner for beginner difficulty", () => { render(<DifficultyBadge difficulty="beginner" />); expect(screen.getByText("Beginner")).toBeInTheDocument(); });
  it("shows Intermediate for intermediate difficulty", () => { render(<DifficultyBadge difficulty="intermediate" />); expect(screen.getByText("Intermediate")).toBeInTheDocument(); });
});

describe("CostBadge", () => {
  it("shows Free for free cost type", () => { render(<CostBadge costType="free" />); expect(screen.getByText("Free")).toBeInTheDocument(); });
  it("shows Free for free_audit cost type", () => { render(<CostBadge costType="free_audit" />); expect(screen.getByText("Free")).toBeInTheDocument(); });
  it("shows Paid for paid cost type", () => { render(<CostBadge costType="paid" />); expect(screen.getByText("Paid")).toBeInTheDocument(); });
});
