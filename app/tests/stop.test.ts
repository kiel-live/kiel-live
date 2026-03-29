import { expect } from '@playwright/test';
import { setLiteMode, testColorScheme, waitForMapToLoad } from './utils';

testColorScheme('Display stop details', async ({ page }) => {
  await page.goto('/map/stop/kvg-2387');
  await expect(page.getByRole('heading', { name: 'Hauptbahnhof' })).toBeVisible();
  await waitForMapToLoad(page);
  await expect(page).toHaveScreenshot();
});

testColorScheme('Display stop details in lite mode', async ({ page }) => {
  await setLiteMode(page);
  await page.goto('/map/stop/kvg-2387');
  await expect(page.getByRole('heading', { name: 'Hauptbahnhof' })).toBeVisible();
  await expect(page).toHaveScreenshot();
});
