import { test as base } from '@playwright/test'
import { shoot } from './screenshot'

export const test = base.extend({})

test.afterEach(async ({ page }, testInfo) => {
  if (testInfo.status === testInfo.expectedStatus) {
    return
  }
  console.log(`failure url: ${page.url()}`)
  await shoot(page, testInfo, 'failure')
})

export { expect } from '@playwright/test'
