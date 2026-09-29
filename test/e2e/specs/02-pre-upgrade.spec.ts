import { test } from '../helpers/fixtures'
import { shoot } from '../helpers/screenshot'
import { env } from '../helpers/env'
import { login } from '../helpers/jellyfin'

const user = env('PLAYWRIGHT_DEVICE_USER')
const password = env('PLAYWRIGHT_DEVICE_PASSWORD')

test.describe('jellyfin before the upgrade', () => {
  test('sign in to the released version', async ({ page }, testInfo) => {
    await login(page, user, password)
    await shoot(page, testInfo, 'before-upgrade')
  })
})
