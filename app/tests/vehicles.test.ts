import { expect } from '@playwright/test';
import { testColorScheme } from './utils';

testColorScheme('Display bus details', async ({ page }) => {
  await page.goto('/map/bus/bus-3');
  await expect(page.getByRole('heading', { name: '62 Russee, Schiefe Horn' })).toBeVisible();
  await page.waitForSelector('#map[data-idle="true"]'); // wait for the map to have loaded
  await expect(page).toHaveScreenshot();
});
