import { expect, test } from '@playwright/test';

test('Can render map', async ({ page }) => {
  await page.goto('/');

  // wait for the map to have loaded
  await page.waitForSelector('#map[data-idle="true"]');

  await expect(page).toHaveScreenshot();
});
