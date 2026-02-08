import { expect, test } from '@playwright/test';
import { setLiteMode, waitForMapToLoad } from './utils';

async function injectSafeAreaVars(page: any) {
  await page.addStyleTag({
    content: `
      :root {
        --safe-area-top: 48px;
        --safe-area-bottom: 32px;
      }
    `,
  });
}

test('Safe area on map', async ({ page }) => {
  await page.goto('/');
  await injectSafeAreaVars(page);
  await waitForMapToLoad(page);

  await expect(page).toHaveScreenshot();
});

test('Safe area on settings', async ({ page }) => {
  await page.goto('/');
  await injectSafeAreaVars(page);
  await page.getByRole('link', { name: 'Settings' }).click();
  await expect(page.getByRole('link', { name: 'Follow @kiel.live' })).toBeVisible();

  await expect(page).toHaveScreenshot();
});

test('Safe area in lite mode', async ({ page }) => {
  await setLiteMode(page);
  await injectSafeAreaVars(page);
  await page.getByRole('link', { name: 'Search' }).click();

  // ensure map is not visible
  await expect(page.getByRole('region', { name: 'Map' })).not.toBeVisible();

  // ensure search is visible
  await page.getByRole('textbox', { name: 'Search' }).click();
  await page.getByRole('textbox', { name: 'Search' }).fill('haupt');
  await expect(page.getByRole('heading', { name: 'Search result' })).toBeVisible();

  await expect(page).toHaveScreenshot();
});
