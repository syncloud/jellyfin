import { Page, expect } from '@playwright/test'

export function userMenu(page: Page) {
  return page.locator('#app-user-menu')
}

export async function login(page: Page, user: string, password: string) {
  await page.goto('/')
  await page.locator('#txtManualName').waitFor()
  await page.locator('#txtManualName').fill(user)
  await page.locator('#txtManualPassword').fill(password)
  await page.locator('#txtManualPassword').press('Enter')
  await expect(page.getByRole('heading', { name: 'Nothing here.' })).toBeVisible()
}

export async function openDashboard(page: Page) {
  await page.getByRole('button', { name: 'User Menu' }).click()
  const menu = userMenu(page)
  await expect(menu.getByRole('menuitem', { name: 'Settings', exact: true })).toBeVisible()
  await expect(menu.getByRole('menuitem', { name: 'Sign Out', exact: true })).toBeVisible()
  await menu.getByRole('menuitem', { name: 'Dashboard', exact: true }).click()
  await expect(page.getByRole('button', { name: 'Scan All Libraries' })).toBeVisible()
}

export async function openDashboardSection(page: Page, href: string) {
  const drawer = page.getByRole('button', { name: 'Open Menu' })
  if (await drawer.isVisible()) {
    await drawer.click()
  }
  await page.locator(`a[href="${href}"]`).filter({ visible: true }).first().click()
}

export async function scanAllLibraries(page: Page) {
  await page.getByRole('button', { name: 'Scan All Libraries' }).click()
}

export async function expectScanHasRun(page: Page) {
  await openDashboardSection(page, '#/dashboard/tasks')
  const task = page.getByRole('heading', { name: 'Scan Media Library' })
  await expect(task).toBeVisible()
  await expect(task.locator('xpath=following-sibling::*').filter({ hasText: 'Last ran' }).first()).toBeVisible()
}

export async function expectServerVersion(page: Page, version: string) {
  await expect(page.getByText(version, { exact: true }).first()).toBeVisible()
}
