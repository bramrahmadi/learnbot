import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { Navbar } from "../Navbar";

// Mock Next.js navigation hooks
jest.mock("next/navigation", () => ({
  usePathname: jest.fn(() => "/dashboard"),
}));

// Mock Next.js Link component
jest.mock("next/link", () => {
  const MockLink = ({ children, href, ...props }: { children: React.ReactNode; href: string; [key: string]: unknown }) => (
    <a href={href} {...props}>{children}</a>
  );
  MockLink.displayName = "MockLink";
  return MockLink;
});

// Mock the auth store
const mockLogout = jest.fn();
const mockUseAuthStore = jest.fn();
const mockUseIsAuthenticated = jest.fn();
const mockUseUser = jest.fn();

jest.mock("@/store/authStore", () => ({
  useAuthStore: (selector: (state: { logout: () => void }) => unknown) => mockUseAuthStore(selector),
  useIsAuthenticated: () => mockUseIsAuthenticated(),
  useUser: () => mockUseUser(),
}));

// Setup mock implementations
beforeEach(() => {
  mockUseAuthStore.mockImplementation((selector: (state: { logout: () => void }) => unknown) => {
    return selector({ logout: mockLogout });
  });
  mockLogout.mockClear();
});

describe("Navbar - Unauthenticated", () => {
  beforeEach(() => {
    mockUseIsAuthenticated.mockReturnValue(false);
    mockUseUser.mockReturnValue(null);
  });

  it("renders the LearnBot brand", () => {
    render(<Navbar />);
    expect(screen.getByText("LearnBot")).toBeInTheDocument();
  });

  it("shows Sign in link when not authenticated", () => {
    render(<Navbar />);
    expect(screen.getByRole("link", { name: "Sign in" })).toBeInTheDocument();
  });

  it("shows Get started link when not authenticated", () => {
    render(<Navbar />);
    expect(screen.getByRole("link", { name: "Get started" })).toBeInTheDocument();
  });

  it("does not show nav links when not authenticated", () => {
    render(<Navbar />);
    expect(screen.queryByRole("link", { name: "Dashboard" })).not.toBeInTheDocument();
    expect(screen.queryByRole("link", { name: "Jobs" })).not.toBeInTheDocument();
  });

  it("links to home page when not authenticated", () => {
    render(<Navbar />);
    const brandLink = screen.getByRole("link", { name: /LearnBot/i });
    expect(brandLink).toHaveAttribute("href", "/");
  });

  it("has navigation role", () => {
    render(<Navbar />);
    expect(screen.getByRole("navigation")).toBeInTheDocument();
  });
});

describe("Navbar - Authenticated", () => {
  beforeEach(() => {
    mockUseIsAuthenticated.mockReturnValue(true);
    mockUseUser.mockReturnValue({ id: "1", email: "test@example.com", full_name: "Test User" });
  });

  it("shows navigation links when authenticated", () => {
    render(<Navbar />);
    expect(screen.getAllByRole("link", { name: "Dashboard" }).length).toBeGreaterThan(0);
    expect(screen.getAllByRole("link", { name: "Jobs" }).length).toBeGreaterThan(0);
    expect(screen.getAllByRole("link", { name: "Gap Analysis" }).length).toBeGreaterThan(0);
    expect(screen.getAllByRole("link", { name: "Learning Plan" }).length).toBeGreaterThan(0);
  });

  it("shows user name when authenticated", () => {
    render(<Navbar />);
    expect(screen.getByText("Test User")).toBeInTheDocument();
  });

  it("shows Sign out button when authenticated", () => {
    render(<Navbar />);
    expect(screen.getByRole("button", { name: "Sign out" })).toBeInTheDocument();
  });

  it("calls logout when Sign out is clicked", () => {
    render(<Navbar />);
    fireEvent.click(screen.getByRole("button", { name: "Sign out" }));
    expect(mockLogout).toHaveBeenCalledTimes(1);
  });

  it("links to dashboard when authenticated", () => {
    render(<Navbar />);
    const brandLink = screen.getByRole("link", { name: /LearnBot/i });
    expect(brandLink).toHaveAttribute("href", "/dashboard");
  });

  it("shows mobile menu toggle button", () => {
    render(<Navbar />);
    expect(screen.getByRole("button", { name: "Toggle menu" })).toBeInTheDocument();
  });

  it("toggles mobile menu when toggle button is clicked", () => {
    render(<Navbar />);
    const toggleBtn = screen.getByRole("button", { name: "Toggle menu" });
    
    // Initially closed
    expect(toggleBtn).toHaveAttribute("aria-expanded", "false");
    
    // Open menu
    fireEvent.click(toggleBtn);
    expect(toggleBtn).toHaveAttribute("aria-expanded", "true");
    
    // Close menu
    fireEvent.click(toggleBtn);
    expect(toggleBtn).toHaveAttribute("aria-expanded", "false");
  });

  it("shows email in mobile menu when open", () => {
    render(<Navbar />);
    const toggleBtn = screen.getByRole("button", { name: "Toggle menu" });
    fireEvent.click(toggleBtn);
    expect(screen.getByText("test@example.com")).toBeInTheDocument();
  });
});

describe("Navbar - Active link highlighting", () => {
  beforeEach(() => {
    mockUseIsAuthenticated.mockReturnValue(true);
    mockUseUser.mockReturnValue({ id: "1", email: "test@example.com", full_name: "Test User" });
  });

  it("marks current page link with aria-current=page", () => {
    const { usePathname } = require("next/navigation");
    usePathname.mockReturnValue("/dashboard");
    
    render(<Navbar />);
    const dashboardLinks = screen.getAllByRole("link", { name: "Dashboard" });
    const activeLink = dashboardLinks.find(link => link.getAttribute("aria-current") === "page");
    expect(activeLink).toBeTruthy();
  });
});
