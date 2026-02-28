import type { Metadata } from "next";
import "./globals.css";
import { Navbar } from "@/components/layout/Navbar";

export const metadata: Metadata = {
  title: { default: "LearnBot – AI Career Development", template: "%s | LearnBot" },
  description: "AI-powered career development platform. Analyze skill gaps, get personalized training recommendations, and accelerate your career growth.",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="h-full">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet" />
      </head>
      <body className="h-full">
        <Navbar />
        <main id="main-content" className="min-h-[calc(100vh-4rem)]">{children}</main>
      </body>
    </html>
  );
}
