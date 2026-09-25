import { expect, test } from '@playwright/test';
import { setLiteMode, testColorScheme, waitForMapToLoad } from './utils';

test('Clicking the search field opens the search popup', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('textbox', { name: 'Search' }).click();
  await expect(page.getByRole('heading', { name: 'Search result' })).toBeVisible();
  await expect(page.getByText('Search for a stop or a vehicle')).toBeVisible();
});

test('Searching displays stops', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('textbox', { name: 'Search' }).fill('ziegelteich');
  await expect(page.getByRole('heading', { name: 'Search result' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Ziegelteich' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Hauptbahnhof' })).not.toBeVisible();

  await page.getByRole('textbox', { name: 'Search' }).fill('lange reihe');
  await expect(page.getByRole('link', { name: 'Lange Reihe' })).toBeVisible();
});

test('Clicking a stop in the search result opens the stop page', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('textbox', { name: 'Search' }).fill('ziegelteich');
  await page.getByRole('link', { name: 'Ziegelteich' }).click();
  await expect(page.getByRole('heading', { name: 'Ziegelteich' })).toBeVisible();
});

testColorScheme('Search results are correctly displayed', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('textbox', { name: 'Search' }).fill('haupt');
  await waitForMapToLoad(page);
  await expect(page).toHaveScreenshot();
});

testColorScheme('Search results are correctly displayed in lite mode', async ({ page }) => {
  await setLiteMode(page);
  await page.goto('/');
  await page.getByRole('textbox', { name: 'Search' }).fill('haupt');
  await expect(page).toHaveScreenshot();
});
