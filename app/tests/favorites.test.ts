import { expect, test } from '@playwright/test';

test('Can add stop to favorites', async ({ page }) => {
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Lehmberg' })).toBeVisible();
  await page.getByRole('link', { name: 'Lehmberg' }).click();
  await expect(page.getByRole('heading', { name: 'Lehmberg' })).toBeVisible();
});

test('Can remove stop from favorites', async ({ page }) => {
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.getByRole('button', { name: 'Delete favorite' }).click();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Lehmberg' })).not.toBeVisible();
});
