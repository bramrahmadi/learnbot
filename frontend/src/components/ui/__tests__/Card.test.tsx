import React from "react";
import { render, screen } from "@testing-library/react";
import { Card, CardHeader } from "../Card";

describe("Card", () => {
  it("renders children", () => {
    render(<Card>Card content</Card>);
    expect(screen.getByText("Card content")).toBeInTheDocument();
  });

  it("applies default card class", () => {
    const { container } = render(<Card>Content</Card>);
    expect(container.firstChild).toHaveClass("card");
  });

  it("applies card-hover class when hover=true", () => {
    const { container } = render(<Card hover>Content</Card>);
    expect(container.firstChild).toHaveClass("card-hover");
  });

  it("applies default medium padding", () => {
    const { container } = render(<Card>Content</Card>);
    expect(container.firstChild).toHaveClass("p-6");
  });

  it("applies small padding when padding=sm", () => {
    const { container } = render(<Card padding="sm">Content</Card>);
    expect(container.firstChild).toHaveClass("p-4");
  });

  it("applies large padding when padding=lg", () => {
    const { container } = render(<Card padding="lg">Content</Card>);
    expect(container.firstChild).toHaveClass("p-8");
  });

  it("applies no padding when padding=none", () => {
    const { container } = render(<Card padding="none">Content</Card>);
    const el = container.firstChild as HTMLElement;
    expect(el.className).not.toContain("p-4");
    expect(el.className).not.toContain("p-6");
    expect(el.className).not.toContain("p-8");
  });

  it("applies custom className", () => {
    const { container } = render(<Card className="custom-class">Content</Card>);
    expect(container.firstChild).toHaveClass("custom-class");
  });
});

describe("CardHeader", () => {
  it("renders title", () => {
    render(<CardHeader title="Test Title" />);
    expect(screen.getByText("Test Title")).toBeInTheDocument();
  });

  it("renders subtitle when provided", () => {
    render(<CardHeader title="Title" subtitle="Subtitle text" />);
    expect(screen.getByText("Subtitle text")).toBeInTheDocument();
  });

  it("does not render subtitle when not provided", () => {
    render(<CardHeader title="Title" />);
    expect(screen.queryByText("Subtitle text")).not.toBeInTheDocument();
  });

  it("renders action when provided", () => {
    render(<CardHeader title="Title" action={<button>Action</button>} />);
    expect(screen.getByRole("button", { name: "Action" })).toBeInTheDocument();
  });

  it("renders icon when provided", () => {
    render(<CardHeader title="Title" icon={<span data-testid="icon">★</span>} />);
    expect(screen.getByTestId("icon")).toBeInTheDocument();
  });

  it("does not render icon container when icon is not provided", () => {
    const { container } = render(<CardHeader title="Title" />);
    // The icon wrapper div should not be present
    const iconWrapper = container.querySelector(".bg-primary-50");
    expect(iconWrapper).not.toBeInTheDocument();
  });
});
