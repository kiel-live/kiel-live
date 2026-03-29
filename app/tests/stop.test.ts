import { expect } from '@playwright/test';
import { setLiteMode, testColorScheme, waitForMapToLoad } from './utils';

testColorScheme('Display stop details', async ({ page }) => {
  await page.clock.setFixedTime(new Date('2024-02-02T10:30:00'));
  await page.goto('/map/stop/kvg-2387');
  await expect(page.getByRole('heading', { name: 'Hauptbahnhof' })).toBeVisible();
  await waitForMapToLoad(page);
  await expect(page).toHaveScreenshot();
});

testColorScheme('Display stop details in lite mode', async ({ page }) => {
  await page.clock.setFixedTime(new Date('2024-02-02T10:30:00'));
  await setLiteMode(page);
  await page.goto('/map/stop/kvg-2387');
  await expect(page.getByRole('heading', { name: 'Hauptbahnhof' })).toBeVisible();
  await expect(page).toHaveScreenshot();
});
