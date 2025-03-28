import { expect, test } from '@playwright/test';

test('Settings', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('link', { name: 'Settings' }).click();
  await expect(page.getByRole('link', { name: 'Follow @kiel.live' })).toBeVisible();
  await expect(page.getByRole('main').getByRole('link', { name: 'Settings' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Contact us' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Analytics' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Changelog' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Develop on Github' })).toBeVisible();

  await expect(page).toHaveScreenshot();
});

test('Lite mode', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('link', { name: 'Settings' }).click();
  await page.getByRole('main').getByRole('link', { name: 'Settings' }).click();
  await page.getByRole('checkbox', { name: 'Lite mode If you enable the' }).check();
  await page.getByRole('link', { name: 'Search' }).click();

  // ensure map is not visible
  await expect(page.getByRole('region', { name: 'Map' })).not.toBeVisible();

  // ensure search is visible
  await page.getByRole('textbox', { name: 'Search' }).click();
  await page.getByRole('textbox', { name: 'Search' }).fill('haupt');
  await expect(page.getByRole('heading', { name: 'Search result' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Hauptbahnhof' })).toBeVisible();

  await expect(page).toHaveScreenshot();
});

test('Dark Theme', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('link', { name: 'Settings' }).click();
  await page.getByRole('main').getByRole('link', { name: 'Settings' }).click();
  await page.getByLabel('ThemeChoose a theme that fits').selectOption('dark');

  await expect(page).toHaveScreenshot();
});

test('Contact Us', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('link', { name: 'Settings' }).click();
  await page.getByRole('link', { name: 'Contact us' }).click();

  await expect(page.getByRole('textbox')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Send email' })).toBeVisible();
  await page.getByRole('textbox').fill('GaLiGrü');

  await expect(page).toHaveScreenshot();
});

test('Changelog', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('link', { name: 'Settings' }).click();
  await page.getByRole('link', { name: 'Changelog' }).click();
  await expect(page.getByRole('heading', { name: 'Changelog' })).toBeVisible();

  await expect(page).toHaveScreenshot();
});
