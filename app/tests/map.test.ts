import { expect } from '@playwright/test';
import { testColorScheme, waitForMapToLoad } from './utils';

testColorScheme('Can render map', async ({ page }) => {
  await page.goto('/');

  await waitForMapToLoad(page);

  await expect(page).toHaveScreenshot();
});
