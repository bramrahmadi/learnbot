import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { Button } from "../Button";

describe("Button", () => {
  it("renders children", () => { render(<Button>Click me</Button>); expect(screen.getByText("Click me")).toBeInTheDocument(); });
  it("applies primary variant by default", () => { render(<Button>Primary</Button>); expect(screen.getByRole("button").className).toContain("btn-primary"); });
  it("applies secondary variant", () => { render(<Button variant="secondary">Secondary</Button>); expect(screen.getByRole("button").className).toContain("btn-secondary"); });
  it("shows loading spinner when loading=true", () => { render(<Button loading>Loading</Button>); const btn = screen.getByRole("button"); expect(btn).toBeDisabled(); expect(btn.querySelector(".spinner")).toBeInTheDocument(); });
  it("is disabled when disabled prop is set", () => { render(<Button disabled>Disabled</Button>); expect(screen.getByRole("button")).toBeDisabled(); });
  it("calls onClick when clicked", () => { const onClick = jest.fn(); render(<Button onClick={onClick}>Click</Button>); fireEvent.click(screen.getByRole("button")); expect(onClick).toHaveBeenCalledTimes(1); });
  it("does not call onClick when disabled", () => { const onClick = jest.fn(); render(<Button disabled onClick={onClick}>Click</Button>); fireEvent.click(screen.getByRole("button")); expect(onClick).not.toHaveBeenCalled(); });
  it("applies large size class", () => { render(<Button size="lg">Large</Button>); expect(screen.getByRole("button").className).toContain("btn-lg"); });
  it("applies small size class", () => { render(<Button size="sm">Small</Button>); expect(screen.getByRole("button").className).toContain("btn-sm"); });
});
