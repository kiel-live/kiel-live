import { devices, expect, test } from '@playwright/test';

test.use({
  ...devices['Pixel 5'],
});

test('Clicking the map closes the popup', async ({ page }) => {
  await page.goto('/search');
  await page.getByRole('region', { name: 'Map' }).click({ position: { x: 100, y: 100 } });
  await expect(page.getByRole('heading', { name: 'Search result' })).not.toBeVisible();
});
