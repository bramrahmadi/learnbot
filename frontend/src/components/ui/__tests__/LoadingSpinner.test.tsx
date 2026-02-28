import React from "react";
import { render, screen } from "@testing-library/react";
import { LoadingSpinner, PageLoader, InlineLoader } from "../LoadingSpinner";

describe("LoadingSpinner", () => {
  it("renders with default props", () => {
    render(<LoadingSpinner />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("has default aria-label of Loading...", () => {
    render(<LoadingSpinner />);
    expect(screen.getByRole("status")).toHaveAttribute("aria-label", "Loading...");
  });

  it("uses custom label", () => {
    render(<LoadingSpinner label="Fetching data..." />);
    expect(screen.getByRole("status")).toHaveAttribute("aria-label", "Fetching data...");
  });

  it("shows label text for medium size", () => {
    render(<LoadingSpinner size="md" label="Loading..." />);
    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("shows label text for large size", () => {
    render(<LoadingSpinner size="lg" label="Loading..." />);
    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("does not show label text for small size", () => {
    render(<LoadingSpinner size="sm" label="Loading..." />);
    // The label text should not be visible for sm size
    expect(screen.queryByText("Loading...")).not.toBeInTheDocument();
  });

  it("applies small size class", () => {
    const { container } = render(<LoadingSpinner size="sm" />);
    const spinner = container.querySelector(".spinner");
    expect(spinner).toHaveClass("w-4", "h-4");
  });

  it("applies medium size class", () => {
    const { container } = render(<LoadingSpinner size="md" />);
    const spinner = container.querySelector(".spinner");
    expect(spinner).toHaveClass("w-8", "h-8");
  });

  it("applies large size class", () => {
    const { container } = render(<LoadingSpinner size="lg" />);
    const spinner = container.querySelector(".spinner");
    expect(spinner).toHaveClass("w-12", "h-12");
  });

  it("applies custom className", () => {
    const { container } = render(<LoadingSpinner className="custom-class" />);
    expect(container.firstChild).toHaveClass("custom-class");
  });
});

describe("PageLoader", () => {
  it("renders a full-page loader", () => {
    render(<PageLoader />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("shows LearnBot loading text", () => {
    render(<PageLoader />);
    expect(screen.getByText("Loading LearnBot...")).toBeInTheDocument();
  });

  it("renders with min-h-screen class", () => {
    const { container } = render(<PageLoader />);
    expect(container.firstChild).toHaveClass("min-h-screen");
  });
});

describe("InlineLoader", () => {
  it("renders with default label", () => {
    render(<InlineLoader />);
    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders with custom label", () => {
    render(<InlineLoader label="Fetching jobs..." />);
    expect(screen.getByText("Fetching jobs...")).toBeInTheDocument();
  });

  it("renders a spinner", () => {
    render(<InlineLoader />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });
});
