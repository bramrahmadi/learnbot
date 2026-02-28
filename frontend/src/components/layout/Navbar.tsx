"use client";
import React, { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuthStore, useIsAuthenticated, useUser } from "@/store/authStore";

const navLinks = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/jobs", label: "Jobs" },
  { href: "/analysis", label: "Gap Analysis" },
  { href: "/learning", label: "Learning Plan" },
];

export function Navbar() {
  const pathname = usePathname();
  const isAuthenticated = useIsAuthenticated();
  const user = useUser();
  const logout = useAuthStore((s) => s.logout);
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50" role="navigation" aria-label="Main navigation">
      <div className="page-container">
        <div className="flex items-center justify-between h-16">
          <Link href={isAuthenticated ? "/dashboard" : "/"} className="flex items-center gap-2">
            <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">LB</span>
            </div>
            <span className="font-bold text-gray-900 text-lg">LearnBot</span>
          </Link>

          {isAuthenticated && (
            <div className="hidden md:flex items-center gap-1">
              {navLinks.map((link) => (
                <Link key={link.href} href={link.href}
                  className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${pathname.startsWith(link.href) ? "bg-primary-50 text-primary-700" : "text-gray-600 hover:bg-gray-100"}`}
                  aria-current={pathname.startsWith(link.href) ? "page" : undefined}>
                  {link.label}
                </Link>
              ))}
            </div>
          )}

          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <>
                <div className="hidden md:flex items-center gap-3">
                  <span className="text-sm text-gray-600">{user?.full_name || user?.email}</span>
                  <button onClick={logout} className="btn-ghost btn text-sm" aria-label="Sign out">Sign out</button>
                </div>
                <button className="md:hidden p-2 rounded-lg text-gray-600 hover:bg-gray-100" onClick={() => setMobileOpen(!mobileOpen)} aria-expanded={mobileOpen} aria-label="Toggle menu">
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    {mobileOpen ? <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /> : <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />}
                  </svg>
                </button>
              </>
            ) : (
              <div className="flex items-center gap-2">
                <Link href="/login" className="btn-ghost btn text-sm">Sign in</Link>
                <Link href="/register" className="btn-primary btn text-sm">Get started</Link>
              </div>
            )}
          </div>
        </div>

        {isAuthenticated && mobileOpen && (
          <div className="md:hidden border-t border-gray-200 py-3 space-y-1">
            {navLinks.map((link) => (
              <Link key={link.href} href={link.href}
                className={`block px-3 py-2 rounded-lg text-sm font-medium ${pathname.startsWith(link.href) ? "bg-primary-50 text-primary-700" : "text-gray-600 hover:bg-gray-100"}`}
                onClick={() => setMobileOpen(false)}>
                {link.label}
              </Link>
            ))}
            <div className="pt-2 border-t border-gray-200">
              <p className="px-3 py-1 text-xs text-gray-500">{user?.email}</p>
              <button onClick={() => { logout(); setMobileOpen(false); }} className="w-full text-left px-3 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-lg">Sign out</button>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
