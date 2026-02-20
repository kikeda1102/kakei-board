/** @type {import('jest').Config} */
module.exports = {
  preset: "ts-jest",
  testEnvironment: "jsdom",
  roots: ["<rootDir>/src"],
  moduleNameMapper: {
    "\\.(css|less|scss)$": "<rootDir>/src/__mocks__/styleMock.ts",
    "\\.(jpg|jpeg|png|gif|svg)$": "<rootDir>/src/__mocks__/fileMock.ts",
  },
  setupFiles: ["<rootDir>/src/test-setup.js"],
  setupFilesAfterEnv: ["@testing-library/jest-dom"],
  transform: {
    "^.+\\.tsx?$": [
      "ts-jest",
      {
        tsconfig: "tsconfig.jest.json",
        diagnostics: false,
        astTransformers: {
          before: [
            {
              path: "ts-jest-mock-import-meta",
              options: {
                metaObjectReplacement: {
                  env: {
                    VITE_API_BASE_URL: "http://localhost:8080",
                  },
                },
              },
            },
          ],
        },
      },
    ],
  },
};
