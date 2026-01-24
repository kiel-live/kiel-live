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
  await page.getByRole('textbox', { name: 'Search' }).fill('stop');
  await expect(page.getByRole('heading', { name: 'Search result' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Dummy Stop 1' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Dummy Stop 2' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Hauptbahnhof' })).not.toBeVisible();
});

test('Clicking a stop in the search result opens the stop page', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('textbox', { name: 'Search' }).fill('stop');
  await page.getByRole('link', { name: 'Dummy Stop 1' }).click();
  await expect(page.getByRole('heading', { name: 'Dummy Stop 1' })).toBeVisible();
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
