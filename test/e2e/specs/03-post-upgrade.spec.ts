import { test } from '../helpers/fixtures'
import { shoot } from '../helpers/screenshot'
import { env } from '../helpers/env'
import { login, openDashboard, scanAllLibraries, expectScanHasRun } from '../helpers/jellyfin'

const user = env('PLAYWRIGHT_DEVICE_USER')
const password = env('PLAYWRIGHT_DEVICE_PASSWORD')

test.describe('jellyfin after the upgrade', () => {
  test('the migrated server is usable', async ({ page }, testInfo) => {
    await test.step('sign in to the upgraded version', async () => {
      await login(page, user, password)
      await shoot(page, testInfo, 'after-upgrade')
    })

    await test.step('open the dashboard', async () => {
      await openDashboard(page)
      await shoot(page, testInfo, 'after-upgrade-dashboard')
    })

    await test.step('scan the migrated libraries', async () => {
      await scanAllLibraries(page)
      await expectScanHasRun(page)
      await shoot(page, testInfo, 'after-upgrade-scan-done')
    })
  })
})
