const nextJest = require("next/jest");

// Force test environment for React 19 compatibility
process.env.NODE_ENV = "test";

const createJestConfig = nextJest({ dir: "./" });

const config = {
  coverageProvider: "v8",
  testEnvironment: "jest-environment-jsdom",
  setupFiles: ["<rootDir>/jest.polyfill.js"],
  setupFilesAfterEnv: ["<rootDir>/jest.setup.ts"],
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/src/$1",
  },
  testPathIgnorePatterns: ["/node_modules/", "/e2e/", "/.next/"],
  testMatch: ["**/__tests__/**/*.{ts,tsx}", "**/*.test.{ts,tsx}"],
  collectCoverageFrom: [
    "src/**/*.{ts,tsx}",
    "!src/**/*.d.ts",
    "!src/**/index.ts",
    "!src/app/layout.tsx",
    "!src/app/globals.css",
    "!src/__tests__/**",
    "!src/__mocks__/**",
  ],
  coverageThreshold: {
    global: {
      branches: 50,
      functions: 50,
      lines: 50,
      statements: 50,
    },
  },
};

module.exports = createJestConfig(config);
