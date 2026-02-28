import { test, expect } from "@playwright/test";

test.describe("Authentication", () => {
  test("landing page loads and shows hero section", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveTitle(/LearnBot/);
    await expect(page.getByText("Land Your Dream Job")).toBeVisible();
    await expect(page.getByText("Get started free")).toBeVisible();
  });

  test("landing page has navigation links", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByRole("link", { name: "Sign in" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Get started" })).toBeVisible();
  });

  test("register page loads", async ({ page }) => {
    await page.goto("/register");
    await expect(page.getByText("Start your career journey")).toBeVisible();
    await expect(page.getByLabel("Full name")).toBeVisible();
    await expect(page.getByLabel("Email address")).toBeVisible();
    await expect(page.getByLabel("Password")).toBeVisible();
  });

  test("register form validates required fields", async ({ page }) => {
    await page.goto("/register");
    await page.getByRole("button", { name: "Create free account" }).click();
    await expect(page.getByText("Full name is required")).toBeVisible();
    await expect(page.getByText("Email is required")).toBeVisible();
    await expect(page.getByText("Password is required")).toBeVisible();
  });

  test("register form validates email format", async ({ page }) => {
    await page.goto("/register");
    await page.getByLabel("Full name").fill("Test User");
    await page.getByLabel("Email address").fill("not-an-email");
    await page.getByLabel("Password").fill("password123");
    await page.getByRole("button", { name: "Create free account" }).click();
    await expect(page.getByText("Enter a valid email address")).toBeVisible();
  });

  test("register form validates password length", async ({ page }) => {
    await page.goto("/register");
    await page.getByLabel("Full name").fill("Test User");
    await page.getByLabel("Email address").fill("test@example.com");
    await page.getByLabel("Password").fill("short");
    await page.getByRole("button", { name: "Create free account" }).click();
    await expect(page.getByText("Password must be at least 8 characters")).toBeVisible();
  });

  test("login page loads", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByText("Welcome back")).toBeVisible();
    await expect(page.getByLabel("Email address")).toBeVisible();
    await expect(page.getByLabel("Password")).toBeVisible();
  });

  test("login form validates required fields", async ({ page }) => {
    await page.goto("/login");
    await page.getByRole("button", { name: "Sign in" }).click();
    await expect(page.getByText("Email is required")).toBeVisible();
    await expect(page.getByText("Password is required")).toBeVisible();
  });

  test("login page has link to register", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByRole("link", { name: "Create one free" })).toBeVisible();
  });

  test("register page has link to login", async ({ page }) => {
    await page.goto("/register");
    await expect(page.getByRole("link", { name: "Sign in" })).toBeVisible();
  });
});
