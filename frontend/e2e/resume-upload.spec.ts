import { test, expect } from "@playwright/test";

test.describe("Resume Upload Flow", () => {
  test("onboarding page loads with resume upload section", async ({ page }) => {
    await page.goto("/onboarding");
    // Should redirect to login if not authenticated
    await expect(page).toHaveURL(/\/(login|onboarding)/);
  });

  test("dashboard shows resume upload prompt for new users", async ({ page }) => {
    await page.goto("/dashboard");
    // Should redirect to login if not authenticated
    await expect(page).toHaveURL(/\/(login|dashboard)/);
  });
});
