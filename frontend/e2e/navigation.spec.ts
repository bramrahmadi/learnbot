import { test, expect } from "@playwright/test";

test.describe("Navigation", () => {
  test("navbar shows sign in and get started when not authenticated", async ({ page }) => {
    await page.goto("/");
    const nav = page.getByRole("navigation");
    await expect(nav.getByRole("link", { name: "Sign in" })).toBeVisible();
    await expect(nav.getByRole("link", { name: "Get started" })).toBeVisible();
  });

  test("landing page features section is visible", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByText("Resume Analysis")).toBeVisible();
    await expect(page.getByText("Job Matching")).toBeVisible();
    await expect(page.getByText("Skill Gap Analysis")).toBeVisible();
    await expect(page.getByText("Personalized Learning")).toBeVisible();
  });

  test("clicking Get started navigates to register", async ({ page }) => {
    await page.goto("/");
    await page.getByRole("link", { name: "Get started free" }).first().click();
    await expect(page).toHaveURL("/register");
  });

  test("clicking Sign in navigates to login", async ({ page }) => {
    await page.goto("/");
    await page.getByRole("link", { name: "Sign in" }).first().click();
    await expect(page).toHaveURL("/login");
  });

  test("logo links to home when not authenticated", async ({ page }) => {
    await page.goto("/login");
    await page.getByRole("link", { name: /LearnBot/ }).click();
    await expect(page).toHaveURL("/");
  });

  test("page has correct title", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveTitle(/LearnBot/);
  });
});

test.describe("Accessibility", () => {
  test("landing page has main landmark", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByRole("main")).toBeVisible();
  });

  test("navigation has aria-label", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByRole("navigation", { name: "Main navigation" })).toBeVisible();
  });

  test("register form has proper labels", async ({ page }) => {
    await page.goto("/register");
    await expect(page.getByLabel("Full name")).toBeVisible();
    await expect(page.getByLabel("Email address")).toBeVisible();
    await expect(page.getByLabel("Password")).toBeVisible();
  });

  test("login form has proper labels", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByLabel("Email address")).toBeVisible();
    await expect(page.getByLabel("Password")).toBeVisible();
  });

  test("error messages are visible after form submission", async ({ page }) => {
    await page.goto("/login");
    await page.getByRole("button", { name: "Sign in" }).click();
    await expect(page.getByText("Email is required")).toBeVisible();
  });
});
