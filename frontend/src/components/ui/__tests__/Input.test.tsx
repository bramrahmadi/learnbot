import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { Input, Textarea, Select } from "../Input";

describe("Input", () => {
  it("renders without label", () => {
    render(<Input placeholder="Enter text" />);
    expect(screen.getByPlaceholderText("Enter text")).toBeInTheDocument();
  });

  it("renders with label", () => {
    render(<Input label="Email address" />);
    expect(screen.getByLabelText("Email address")).toBeInTheDocument();
  });

  it("shows required asterisk when required", () => {
    render(<Input label="Email" required />);
    expect(screen.getByText("*")).toBeInTheDocument();
  });

  it("does not show asterisk when not required", () => {
    render(<Input label="Email" />);
    expect(screen.queryByText("*")).not.toBeInTheDocument();
  });

  it("shows error message when error prop is provided", () => {
    render(<Input label="Email" error="Email is required" />);
    expect(screen.getByRole("alert")).toHaveTextContent("Email is required");
  });

  it("applies input-error class when error is provided", () => {
    render(<Input label="Email" error="Error" />);
    const input = screen.getByLabelText("Email");
    expect(input).toHaveClass("input-error");
  });

  it("sets aria-invalid when error is provided", () => {
    render(<Input label="Email" error="Error" />);
    const input = screen.getByLabelText("Email");
    expect(input).toHaveAttribute("aria-invalid", "true");
  });

  it("does not set aria-invalid when no error", () => {
    render(<Input label="Email" />);
    const input = screen.getByLabelText("Email");
    expect(input).toHaveAttribute("aria-invalid", "false");
  });

  it("shows hint text when hint is provided and no error", () => {
    render(<Input label="Email" hint="Enter your email address" />);
    expect(screen.getByText("Enter your email address")).toBeInTheDocument();
  });

  it("does not show hint when error is present", () => {
    render(<Input label="Email" error="Error" hint="Hint text" />);
    expect(screen.queryByText("Hint text")).not.toBeInTheDocument();
  });

  it("calls onChange when value changes", () => {
    const onChange = jest.fn();
    render(<Input label="Email" onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Email"), { target: { value: "test@example.com" } });
    expect(onChange).toHaveBeenCalledTimes(1);
  });

  it("uses provided id", () => {
    render(<Input id="custom-id" label="Email" />);
    expect(screen.getByLabelText("Email")).toHaveAttribute("id", "custom-id");
  });

  it("generates id from label when no id provided", () => {
    render(<Input label="Email address" />);
    expect(screen.getByLabelText("Email address")).toHaveAttribute("id", "email-address");
  });

  it("applies custom className", () => {
    render(<Input label="Email" className="custom-class" />);
    expect(screen.getByLabelText("Email")).toHaveClass("custom-class");
  });
});

describe("Textarea", () => {
  it("renders with label", () => {
    render(<Textarea label="Description" />);
    expect(screen.getByLabelText("Description")).toBeInTheDocument();
  });

  it("shows error message", () => {
    render(<Textarea label="Description" error="Description is required" />);
    expect(screen.getByText("Description is required")).toBeInTheDocument();
  });

  it("shows hint text when no error", () => {
    render(<Textarea label="Description" hint="Max 500 characters" />);
    expect(screen.getByText("Max 500 characters")).toBeInTheDocument();
  });

  it("applies resize-none class", () => {
    render(<Textarea label="Description" />);
    const textarea = screen.getByLabelText("Description");
    expect(textarea).toHaveClass("resize-none");
  });
});

describe("Select", () => {
  const options = [
    { value: "option1", label: "Option 1" },
    { value: "option2", label: "Option 2" },
    { value: "option3", label: "Option 3" },
  ];

  it("renders with label", () => {
    render(<Select label="Choose option" options={options} />);
    expect(screen.getByLabelText("Choose option")).toBeInTheDocument();
  });

  it("renders all options", () => {
    render(<Select label="Choose option" options={options} />);
    expect(screen.getByText("Option 1")).toBeInTheDocument();
    expect(screen.getByText("Option 2")).toBeInTheDocument();
    expect(screen.getByText("Option 3")).toBeInTheDocument();
  });

  it("shows error message", () => {
    render(<Select label="Choose option" options={options} error="Selection required" />);
    expect(screen.getByText("Selection required")).toBeInTheDocument();
  });

  it("calls onChange when selection changes", () => {
    const onChange = jest.fn();
    render(<Select label="Choose option" options={options} onChange={onChange} />);
    fireEvent.change(screen.getByLabelText("Choose option"), { target: { value: "option2" } });
    expect(onChange).toHaveBeenCalledTimes(1);
  });
});
