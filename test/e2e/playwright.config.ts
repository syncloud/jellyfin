import { defineConfig, devices } from '@playwright/test'
import { env } from './helpers/env'

const appDomain = env('PLAYWRIGHT_APP_DOMAIN')
const artifactDir = env('PLAYWRIGHT_ARTIFACT_DIR')
const baseURL = `https://${appDomain}`

export default defineConfig({
  testDir: './specs',
  fullyParallel: false,
  workers: 1,
  retries: 1,
  maxFailures: 1,
  reporter: [['list'], ['html', { outputFolder: `${artifactDir}/playwright/report`, open: 'never' }]],
  outputDir: `${artifactDir}/playwright/test-results`,
  globalTeardown: './globalTeardown.ts',
  timeout: 420_000,
  expect: { timeout: 60_000 },
  use: {
    baseURL,
    ignoreHTTPSErrors: true,
    trace: 'retain-on-failure',
    video: 'on',
  },
  projects: [
    {
      name: 'desktop',
      use: { ...devices['Desktop Chrome'], baseURL, ignoreHTTPSErrors: true, viewport: { width: 1440, height: 960 } },
    },
    {
      name: 'mobile',
      use: { ...devices['Pixel 7'], baseURL, ignoreHTTPSErrors: true },
    },
  ],
  metadata: {
    appDomain,
    artifactDir,
  },
})
