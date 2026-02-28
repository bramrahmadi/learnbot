const nextJest = require("next/jest");

const createJestConfig = nextJest({ dir: "./" });

const config = {
  coverageProvider: "v8",
  testEnvironment: "jsdom",
  setupFilesAfterFramework: ["<rootDir>/jest.setup.ts"],
  moduleNameMapper: { "^@/(.*)$": "<rootDir>/src/$1" },
};

module.exports = createJestConfig(config);
