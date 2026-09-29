import { test } from '../helpers/fixtures'
import { shoot } from '../helpers/screenshot'
import { env } from '../helpers/env'
import { login, openDashboard, scanAllLibraries, expectScanHasRun } from '../helpers/jellyfin'

const user = env('PLAYWRIGHT_DEVICE_USER')
const password = env('PLAYWRIGHT_DEVICE_PASSWORD')

test.describe('jellyfin smoke', () => {
  test('login, open the dashboard and scan the libraries', async ({ page }, testInfo) => {
    await test.step('login', async () => {
      await login(page, user, password)
      await shoot(page, testInfo, 'main')
    })

    await test.step('open the dashboard', async () => {
      await openDashboard(page)
      await shoot(page, testInfo, 'dashboard')
    })

    await test.step('scan the libraries', async () => {
      await scanAllLibraries(page)
      await shoot(page, testInfo, 'scan')
      await expectScanHasRun(page)
      await shoot(page, testInfo, 'scan-done')
    })
  })
})
