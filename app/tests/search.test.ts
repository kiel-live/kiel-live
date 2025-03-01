import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');

  // wait for the map to have loaded #map data-idle="true"
  await page.waitForSelector('#map[data-idle="true"]');

  // input search query
  await page.fill('input[title="Search"]', 'Stop');

  // assert search results
  await expect(page).toHaveURL('/search');

  // assert screenshot
  await expect(page).toHaveScreenshot('search-results.png');
});
