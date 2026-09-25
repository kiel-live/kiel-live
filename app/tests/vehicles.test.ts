import { expect } from '@playwright/test';
import { setLiteMode, testColorScheme, waitForMapToLoad } from './utils';

testColorScheme('Display bus details', async ({ page }) => {
  await page.clock.setFixedTime(new Date('2024-02-02T10:30:00'));
  await page.goto('/map/bus/bus-3');
  await expect(page.getByRole('heading', { name: '62 Russee, Schiefe Horn' })).toBeVisible();
  await waitForMapToLoad(page);
  await expect(page).toHaveScreenshot();
});

testColorScheme('Display bus details in lite mode', async ({ page }) => {
  await page.clock.setFixedTime(new Date('2024-02-02T10:30:00'));
  await setLiteMode(page);
  await page.goto('/map/bus/bus-3');
  await expect(page.getByRole('heading', { name: '62 Russee, Schiefe Horn' })).toBeVisible();
  await expect(page).toHaveScreenshot();
});
