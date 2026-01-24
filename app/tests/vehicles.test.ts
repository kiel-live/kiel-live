import { expect } from '@playwright/test';
import { setLiteMode, testColorScheme, waitForMapToLoad } from './utils';

testColorScheme('Display bus details', async ({ page }) => {
  await page.goto('/map/bus/bus-3');
  await expect(page.getByRole('heading', { name: '62 Russee, Schiefe Horn' })).toBeVisible();
  await waitForMapToLoad(page);
  await expect(page).toHaveScreenshot();
});

testColorScheme('Display bus details in lite mode', async ({ page }) => {
  await setLiteMode(page);
  await page.goto('/map/bus/bus-3');
  await expect(page.getByRole('heading', { name: '62 Russee, Schiefe Horn' })).toBeVisible();
  await expect(page).toHaveScreenshot();
});
