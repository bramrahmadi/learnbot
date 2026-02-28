# LearnBot — User Guide

## Table of Contents

1. [Getting Started](#1-getting-started)
2. [Creating Your Account](#2-creating-your-account)
3. [Building Your Profile](#3-building-your-profile)
4. [Uploading Your Resume](#4-uploading-your-resume)
5. [Exploring Job Matches](#5-exploring-job-matches)
6. [Understanding Your Acceptance Likelihood Score](#6-understanding-your-acceptance-likelihood-score)
7. [Skill Gap Analysis](#7-skill-gap-analysis)
8. [Your Personalized Learning Plan](#8-your-personalized-learning-plan)
9. [Dashboard Overview](#9-dashboard-overview)
10. [Managing Your Career Goals](#10-managing-your-career-goals)
11. [Account Settings](#11-account-settings)

---

## 1. Getting Started

Welcome to **LearnBot** — your AI-powered career development companion. LearnBot helps you:

- **Understand your current skills** through intelligent resume analysis
- **Find the right jobs** with personalized matching and acceptance likelihood scores
- **Close skill gaps** with curated learning recommendations
- **Track your progress** toward your career goals

### What You'll Need

- A resume in PDF or DOCX format (optional but recommended)
- Your career goals and target role in mind
- About 10 minutes to set up your profile

---

## 2. Creating Your Account

### Step 1: Register

1. Go to [LearnBot](http://localhost:3000) (or your organization's URL)
2. Click **Get Started** or **Sign Up**
3. Enter your details:
   - **Full Name** — Your display name
   - **Email Address** — Used for login (must be unique)
   - **Password** — Minimum 8 characters
4. Click **Create Account**

You'll be automatically logged in and redirected to the onboarding flow.

### Step 2: Login (Returning Users)

1. Click **Sign In**
2. Enter your email and password
3. Click **Sign In**

> **Forgot your password?** Password reset is coming in a future update. Contact support if you need access.

---

## 3. Building Your Profile

Your profile is the foundation of LearnBot's recommendations. The more complete your profile, the better your job matches and learning recommendations.

### Profile Sections

| Section | Impact on Recommendations |
|---------|--------------------------|
| Professional headline | Job title matching |
| Summary | Context for AI analysis |
| Location | Location-based job filtering |
| Skills | Core matching algorithm |
| Work experience | Experience level matching |
| Education | Education requirement matching |
| Certifications | Qualification matching |

### Completing Your Profile

1. Navigate to **Dashboard** → **Profile**
2. Click **Edit Profile**
3. Fill in your:
   - **Professional Headline** (e.g., "Senior Software Engineer")
   - **Professional Summary** (2–3 sentences about your background)
   - **Location** (city and country)
   - **Open to Work** toggle (makes you visible to job matching)
4. Click **Save Changes**

### Profile Completeness Score

Your profile has a **completeness score** (0–100%). A higher score means:
- More accurate job matching
- Better skill gap analysis
- More relevant learning recommendations

**To reach 100%:**
- ✅ Add a professional headline
- ✅ Write a summary
- ✅ Set your location
- ✅ Upload a resume (or add skills manually)
- ✅ Add work experience
- ✅ Add education
- ✅ Add at least 3 skills

---

## 4. Uploading Your Resume

Uploading your resume is the fastest way to build a complete profile. LearnBot automatically extracts your skills, experience, and education.

### Supported Formats

| Format | Max Size |
|--------|----------|
| PDF | 10 MB |
| DOCX (Microsoft Word) | 10 MB |

### How to Upload

1. Go to **Dashboard** → **Upload Resume** (or the **Onboarding** page)
2. Click **Choose File** or drag and drop your resume
3. Click **Upload & Analyze**
4. Wait 10–30 seconds while LearnBot parses your resume

### What Gets Extracted

LearnBot extracts:
- **Skills** — Programming languages, frameworks, tools, soft skills
- **Work Experience** — Company names, job titles, dates, responsibilities
- **Education** — Institutions, degrees, fields of study
- **Certifications** — Professional certifications and licenses

### Reviewing Extracted Data

After upload, you'll see a summary of what was extracted. You can:
- **Confirm** the extracted data to add it to your profile
- **Edit** any incorrect information
- **Add** skills or experience that wasn't detected

> **Tip:** If your resume uses unusual formatting or fonts, some data may not be extracted correctly. Review the results and manually add any missing information.

### Confidence Scores

Each extracted item has a **confidence score** (0–100%). Items with low confidence (< 70%) are flagged for your review.

---

## 5. Exploring Job Matches

### Finding Jobs

1. Navigate to **Jobs** in the top navigation
2. Browse recommended jobs (sorted by acceptance likelihood)
3. Use filters to narrow results:
   - **Location Type:** Remote, On-site, Hybrid
   - **Experience Level:** Entry, Mid, Senior, Lead, Executive
   - **Skills:** Filter by specific technologies

### Searching for Jobs

1. Click **Search Jobs**
2. Enter a job title or keyword (e.g., "backend engineer", "data scientist")
3. Optionally add skill filters
4. Click **Search**

### Job Cards

Each job card shows:
- **Job Title** and **Company**
- **Location** and **Location Type** (remote/on-site/hybrid)
- **Acceptance Likelihood Score** (colored badge)
- **Required Skills** (with match indicators)
- **Posted Date**

### Job Detail Page

Click any job to see:
- Full job description
- Complete skill requirements
- Your match breakdown
- Salary range (when available)
- Direct link to apply

---

## 6. Understanding Your Acceptance Likelihood Score

The **Acceptance Likelihood Score** (0–100%) estimates how well your profile matches a job's requirements. It's calculated from five factors:

| Factor | Weight | What It Measures |
|--------|--------|-----------------|
| Skill Match | 35% | How many required skills you have |
| Experience Match | 25% | Years of experience vs. requirements |
| Education Match | 15% | Degree level vs. requirements |
| Location Fit | 10% | Location compatibility |
| Industry Relevance | 15% | Industry experience alignment |

### Score Interpretation

| Score | Category | Meaning |
|-------|----------|---------|
| 80–100% | 🟢 **Ready to Apply** | Strong match — apply now |
| 60–79% | 🟡 **Close Match** | Good fit with minor gaps |
| 40–59% | 🟠 **Stretch Goal** | Significant gaps to address |
| 0–39% | 🔴 **Not Ready** | Major gaps — focus on learning first |

### Improving Your Score

To improve your score for a specific job:
1. Click **View Match Details** on the job card
2. See your breakdown by factor
3. Click **Analyze Skill Gaps** to see what to learn
4. Follow the recommended learning plan

---

## 7. Skill Gap Analysis

Skill Gap Analysis shows you exactly what skills you need to develop for a target role.

### Running a Gap Analysis

**Option 1: From a Job Listing**
1. Open a job detail page
2. Click **Analyze My Gaps**
3. View your personalized gap report

**Option 2: From the Analysis Page**
1. Navigate to **Analysis** in the top navigation
2. Enter a job title or paste a job description
3. Click **Analyze**

### Understanding Your Gap Report

The gap report shows:

**Readiness Score** — Overall readiness percentage for the role

**Skill Gaps by Priority:**
- 🔴 **Critical** — Must-have skills you're missing (highest priority)
- 🟡 **Important** — Strongly preferred skills you're missing
- 🟢 **Nice to Have** — Optional skills that would strengthen your application

**For each gap:**
- Skill name
- Required proficiency level
- Your current proficiency (if any)
- Estimated learning time
- Recommended resources

**Your Strengths** — Skills you already have that match the job

### Gap Analysis Visualization

The radar chart shows your skill coverage across key areas:
- Technical skills
- Tools and frameworks
- Cloud/infrastructure
- Soft skills

---

## 8. Your Personalized Learning Plan

LearnBot creates a customized learning plan based on your skill gaps and preferences.

### Accessing Your Learning Plan

1. Navigate to **Learning** in the top navigation
2. Or click **Get Learning Plan** from a gap analysis

### Learning Plan Structure

Your plan is organized into **phases**:

**Phase 1: Critical Skills** — Address must-have gaps first  
**Phase 2: Important Skills** — Strengthen your profile  
**Phase 3: Nice-to-Have** — Polish and differentiate

### Resource Types

| Type | Description |
|------|-------------|
| 📚 Course | Structured online course (Udemy, Coursera, etc.) |
| 🏆 Certification | Professional certification program |
| 📖 Documentation | Official documentation and guides |
| 🎥 Video | Tutorial videos and series |
| 📕 Book | Technical books |
| 💻 Practice | Hands-on practice platforms (LeetCode, etc.) |
| 📝 Article | Blog posts and articles |
| 🔨 Project | Project-based learning |

### Customizing Your Plan

Click **Customize Plan** to set preferences:
- **Budget** — Maximum spend on paid resources
- **Weekly Hours** — Hours available per week for learning
- **Prefer Free** — Show only free resources
- **Prefer Hands-On** — Prioritize practical exercises
- **Prefer Certificates** — Prioritize resources with certificates

### Resource Details

Each resource shows:
- Title and provider
- Difficulty level
- Estimated hours
- Cost (free or price)
- Rating
- Certificate availability
- Hands-on exercises availability

---

## 9. Dashboard Overview

The Dashboard is your home base in LearnBot.

### Dashboard Sections

**Profile Summary**
- Your profile completeness score
- Quick links to edit profile sections
- Open to Work status

**Job Recommendations**
- Top 5 recommended jobs
- Quick acceptance likelihood scores
- Link to view all jobs

**Skill Overview**
- Your top skills by proficiency
- Skills added from resume vs. manually
- Link to manage skills

**Learning Progress**
- Active learning plan summary
- Resources in progress
- Estimated completion time

**Career Goals**
- Active goals and progress
- Upcoming milestones

---

## 10. Managing Your Career Goals

Career goals help LearnBot tailor recommendations to your specific objectives.

### Setting a Career Goal

1. Go to **Dashboard** → **Career Goals**
2. Click **Add Goal**
3. Enter:
   - **Goal Title** (e.g., "Become a Senior Backend Engineer")
   - **Target Role** — The job title you're aiming for
   - **Target Date** — When you want to achieve this
   - **Priority** — 1 (highest) to 5 (lowest)
4. Click **Save Goal**

### Tracking Progress

- Update your progress percentage as you complete learning milestones
- Mark goals as **Achieved** when you land the role
- **Pause** goals you're temporarily not pursuing
- **Abandon** goals that are no longer relevant

### Goal Statuses

| Status | Description |
|--------|-------------|
| 🟢 Active | Currently pursuing |
| ✅ Achieved | Goal completed |
| ⏸️ Paused | Temporarily on hold |
| ❌ Abandoned | No longer pursuing |

---

## 11. Account Settings

### Updating Your Profile

1. Click your name in the top navigation
2. Select **Profile Settings**
3. Update your information
4. Click **Save**

### Managing Your Skills

**Add a skill manually:**
1. Go to **Profile** → **Skills**
2. Click **Add Skill**
3. Enter the skill name and proficiency level
4. Click **Add**

**Update skill proficiency:**
1. Click the skill you want to update
2. Select the new proficiency level
3. Click **Save**

**Remove a skill:**
1. Click the skill
2. Click **Remove**

### Privacy Settings

- **Open to Work** — Toggle whether you appear in job matching
- **Profile Visibility** — Control who can see your profile (coming soon)

### Uploading a New Resume

You can upload a new resume at any time:
1. Go to **Profile** → **Resume**
2. Click **Upload New Resume**
3. Your profile will be updated with the new data

> **Note:** Uploading a new resume replaces your current resume but doesn't delete manually added skills or experience.

### Logging Out

Click your name in the top navigation → **Sign Out**

Your session expires automatically after 24 hours.

---

## Tips for Best Results

1. **Keep your profile updated** — Add new skills and experience as you gain them
2. **Be specific with skills** — "React" is better than "JavaScript frameworks"
3. **Set realistic proficiency levels** — Honest self-assessment leads to better matches
4. **Review gap analyses regularly** — Job requirements change over time
5. **Complete learning resources** — Mark resources as complete to track progress
6. **Set career goals** — Goals help LearnBot prioritize recommendations

---

## Getting Help

- **FAQ:** [docs/FAQ.md](FAQ.md)
- **Help Center:** [docs/help/](help/)
- **Support:** support@learnbot.example.com
- **GitHub Issues:** [github.com/learnbot/learnbot/issues](https://github.com/learnbot/learnbot/issues)
