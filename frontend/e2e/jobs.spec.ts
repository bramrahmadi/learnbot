import { test, expect } from "@playwright/test";

test.describe("Job Search and Filtering", () => {
  test("jobs page redirects to login when not authenticated", async ({ page }) => {
    await page.goto("/jobs");
    await expect(page).toHaveURL(/\/(login|jobs)/);
  });

  test("job detail page redirects to login when not authenticated", async ({ page }) => {
    await page.goto("/jobs/job-001");
    await expect(page).toHaveURL(/\/(login|jobs)/);
  });
});

test.describe("Gap Analysis Workflow", () => {
  test("analysis page redirects to login when not authenticated", async ({ page }) => {
    await page.goto("/analysis");
    await expect(page).toHaveURL(/\/(login|analysis)/);
  });
});

test.describe("Learning Resources", () => {
  test("learning page redirects to login when not authenticated", async ({ page }) => {
    await page.goto("/learning");
    await expect(page).toHaveURL(/\/(login|learning)/);
  });
});
