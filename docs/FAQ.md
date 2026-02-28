# LearnBot — Frequently Asked Questions

## Table of Contents

- [General Questions](#general-questions)
- [Account & Profile](#account--profile)
- [Resume Upload](#resume-upload)
- [Job Matching](#job-matching)
- [Skill Gap Analysis](#skill-gap-analysis)
- [Learning Recommendations](#learning-recommendations)
- [Privacy & Data](#privacy--data)
- [Technical Issues](#technical-issues)

---

## General Questions

### What is LearnBot?

LearnBot is an AI-powered career development platform that helps you:
- Analyze your professional skills from your resume
- Find jobs that match your profile with an acceptance likelihood score
- Identify skill gaps between your current abilities and target roles
- Get personalized learning recommendations to close those gaps

### Is LearnBot free to use?

LearnBot's core features (resume analysis, job matching, skill gap analysis, and learning recommendations) are free. Some recommended learning resources may have their own costs (e.g., paid Udemy courses), but LearnBot always shows free alternatives when available.

### What makes LearnBot different from other job platforms?

Unlike traditional job boards, LearnBot:
- **Quantifies your fit** with an Acceptance Likelihood Score (not just keyword matching)
- **Explains the gaps** — tells you exactly what skills you're missing and why
- **Provides a learning path** — not just "you need Docker" but "here's how to learn Docker in 8 hours"
- **Tracks your progress** — your recommendations update as you develop new skills

### What career stages does LearnBot support?

LearnBot supports all career stages:
- **Entry level** — Recent graduates and career changers
- **Mid-level** — Professionals looking to advance
- **Senior level** — Experienced professionals targeting leadership roles
- **Executive** — C-suite and VP-level transitions

---

## Account & Profile

### How do I create an account?

1. Go to the LearnBot homepage
2. Click **Get Started**
3. Enter your name, email, and password
4. Click **Create Account**

You'll be automatically logged in and guided through the onboarding process.

### I forgot my password. How do I reset it?

Password reset via email is coming in a future update. For now, please contact support at support@learnbot.example.com with your registered email address.

### How long does my login session last?

Your session lasts **24 hours**. After that, you'll need to log in again. We recommend logging in at the start of each work session.

### Can I have multiple accounts?

Each email address can only be associated with one account. If you need to change your email, contact support.

### What is the Profile Completeness Score?

The Profile Completeness Score (0–100%) measures how much information you've provided. A higher score means:
- More accurate job matching
- Better skill gap analysis
- More relevant learning recommendations

To reach 100%, add: headline, summary, location, skills, work experience, education, and upload a resume.

### Can I manually add skills without uploading a resume?

Yes! Go to **Profile** → **Skills** → **Add Skill**. You can add any skill with a proficiency level (beginner, intermediate, advanced, expert).

---

## Resume Upload

### What file formats are supported?

LearnBot supports:
- **PDF** (recommended)
- **DOCX** (Microsoft Word)

Maximum file size: **10 MB**

### My resume didn't parse correctly. What should I do?

Resume parsing accuracy depends on the formatting. For best results:
- Use a clean, single-column layout
- Avoid tables, text boxes, and complex formatting
- Use standard section headings (Experience, Education, Skills)
- Ensure text is selectable (not a scanned image)

After upload, review the extracted data and manually correct any errors.

### Can I upload multiple resumes?

Yes, you can upload a new resume at any time. Each upload creates a new version. Your profile is updated with the latest resume's data, but manually added information is preserved.

### Will uploading a new resume delete my existing profile data?

No. Uploading a new resume:
- **Updates** skills and experience extracted from the new resume
- **Preserves** skills and experience you added manually
- **Creates** a new resume version (old versions are kept for reference)

### How long does resume parsing take?

Typically 10–30 seconds depending on the file size and complexity.

### What confidence score is acceptable?

- **90–100%** — High confidence, data is likely accurate
- **70–89%** — Good confidence, review flagged items
- **Below 70%** — Low confidence, manually verify this data

---

## Job Matching

### Where do the job listings come from?

LearnBot aggregates jobs from:
- **LinkedIn Jobs** — Updated daily
- **Indeed** — Updated daily
- **Company career pages** — Configurable list of direct employer pages

Jobs are deduplicated and refreshed daily.

### How often are job listings updated?

Job listings are refreshed **daily at 2am UTC**. Jobs not seen for 7 days are automatically marked as expired.

### What does the Acceptance Likelihood Score mean?

The Acceptance Likelihood Score (0–100%) estimates how well your profile matches a job's requirements. It's calculated from:

| Factor | Weight |
|--------|--------|
| Skill Match | 35% |
| Experience Match | 25% |
| Education Match | 15% |
| Location Fit | 10% |
| Industry Relevance | 15% |

**Score guide:**
- 80–100%: Ready to Apply 🟢
- 60–79%: Close Match 🟡
- 40–59%: Stretch Goal 🟠
- 0–39%: Not Ready 🔴

### Why is my score low for a job I'm qualified for?

A few reasons this might happen:
1. **Skills not in your profile** — Add skills you have that weren't extracted from your resume
2. **Proficiency levels** — Your proficiency may be set lower than required
3. **Experience years** — The job may require more years than your profile shows
4. **Skill naming** — "Node.js" and "NodeJS" are treated as the same skill, but "JavaScript" and "JS" may not be

### Can I apply to jobs directly through LearnBot?

LearnBot shows you the job and provides a direct link to the employer's application page. You apply directly on the employer's website.

### Why don't I see salary information for all jobs?

Salary information is only shown when the employer includes it in the job posting. Many employers don't disclose salary ranges publicly.

---

## Skill Gap Analysis

### What is a skill gap?

A skill gap is a skill required (or preferred) by a target job that you either don't have or have at a lower proficiency level than required.

### What do the gap priority levels mean?

| Priority | Meaning |
|----------|---------|
| 🔴 Critical | Must-have skill — your application will likely be rejected without it |
| 🟡 Important | Strongly preferred — having it significantly improves your chances |
| 🟢 Nice to Have | Optional — adds value but won't disqualify you |

### How is the Readiness Score calculated?

The Readiness Score (0–100%) is a weighted measure of how many required and preferred skills you have at the required proficiency level. Critical skills have higher weight than nice-to-have skills.

### Can I run a gap analysis without a specific job?

Yes! On the Analysis page, you can enter a job title or paste a job description to analyze gaps without selecting a specific job listing.

### How do I mark a skill gap as addressed?

When you complete a learning resource for a skill, update your skill proficiency in your profile. The gap analysis will automatically reflect the improvement.

---

## Learning Recommendations

### How are learning resources selected?

LearnBot maintains a curated catalog of 60+ high-quality learning resources. Resources are selected based on:
- Relevance to the skill gap
- Quality and rating
- Difficulty level match
- Cost (free resources are always shown)
- Certificate availability

### Can I filter for free resources only?

Yes! When viewing your learning plan, click **Customize Plan** and toggle **Prefer Free Resources**.

### How accurate are the estimated learning hours?

Estimated hours are based on the resource provider's stated duration. Actual time may vary based on your learning pace and prior knowledge.

### What if I've already completed a recommended resource?

Update your skill proficiency in your profile to reflect your new knowledge. The learning plan will automatically update to remove completed gaps.

### Can I save resources for later?

Resource bookmarking is coming in a future update. For now, you can access your learning plan at any time from the **Learning** page.

---

## Privacy & Data

### What data does LearnBot collect?

LearnBot collects:
- Account information (name, email, password hash)
- Profile data (skills, experience, education)
- Resume files (stored securely in encrypted cloud storage)
- Usage data (pages visited, features used — anonymized)

### How is my resume stored?

Resume files are stored in encrypted cloud storage (AWS S3 with AES-256 encryption). Files are only accessible to you and the LearnBot system.

### Can I delete my data?

Yes. To delete your account and all associated data, contact support@learnbot.example.com. We will process deletion requests within 30 days in compliance with GDPR.

### Is my data shared with employers?

No. LearnBot does not share your profile data with employers. You control when and where you apply.

### Is LearnBot GDPR compliant?

Yes. LearnBot is designed with GDPR compliance in mind:
- You can request a copy of your data
- You can request deletion of your data
- Data is processed only for the purposes described in our Privacy Policy
- We use soft deletes to preserve referential integrity while honoring deletion requests

### How long is my data retained?

- **Active accounts:** Data is retained as long as your account is active
- **Deleted accounts:** Data is purged within 30 days of deletion request
- **Resume files:** Retained for the lifetime of your account

---

## Technical Issues

### The page isn't loading. What should I do?

1. Refresh the page (Ctrl+R / Cmd+R)
2. Clear your browser cache
3. Try a different browser
4. Check your internet connection
5. If the issue persists, contact support

### I'm getting a "Session expired" error.

Your JWT token has expired (sessions last 24 hours). Simply log in again to continue.

### Resume upload is failing. What are the common causes?

- **File too large** — Maximum size is 10 MB
- **Wrong format** — Only PDF and DOCX are supported
- **Corrupted file** — Try re-saving the file and uploading again
- **Scanned PDF** — LearnBot cannot parse image-based PDFs; use a text-based PDF

### The job search isn't returning results.

Try:
- Removing some filters (especially skill filters)
- Using a broader search term
- Checking if the job aggregator has run recently (jobs refresh daily)

### My skills aren't being recognized in job matching.

Ensure your skill names match common industry terminology:
- Use "React" not "ReactJS" or "React.js"
- Use "Node.js" not "NodeJS"
- Use "PostgreSQL" not "Postgres" or "PSQL"

You can also add skills manually with the exact name used in job descriptions.

### I found a bug. How do I report it?

Please report bugs on [GitHub Issues](https://github.com/learnbot/learnbot/issues) with:
- Steps to reproduce
- Expected behavior
- Actual behavior
- Browser and OS information
- Screenshots if applicable

---

## Still Have Questions?

- **User Guide:** [docs/USER_GUIDE.md](USER_GUIDE.md)
- **Help Center:** [docs/help/](help/)
- **Email Support:** support@learnbot.example.com
- **GitHub Issues:** [github.com/learnbot/learnbot/issues](https://github.com/learnbot/learnbot/issues)
