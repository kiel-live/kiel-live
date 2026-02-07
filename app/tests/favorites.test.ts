import { expect, test } from '@playwright/test';
import { setLiteMode, testColorScheme, waitForMapToLoad } from './utils';

test('Can redirect to favorites', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
});

test('Can add stop to favorites', async ({ page }) => {
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Lehmberg' })).toBeVisible();
});

test('Can remove stop from favorites', async ({ page }) => {
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.getByRole('button', { name: 'Delete favorite' }).click();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Lehmberg' })).not.toBeVisible();
});

test('Clicking a favorite stop opens the stop page', async ({ page }) => {
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await page.getByRole('link', { name: 'Lehmberg' }).click();
  await expect(page.getByRole('heading', { name: 'Lehmberg' })).toBeVisible();
});

test('Favorites are stored in local storage', async ({ page }) => {
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.reload();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Lehmberg' })).toBeVisible();
});

testColorScheme('Favorites are correctly displayed', async ({ page }) => {
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await waitForMapToLoad(page);
  await expect(page).toHaveScreenshot();
});

testColorScheme('Favorites are correctly displayed in lite mode', async ({ page }) => {
  await setLiteMode(page);
  await page.goto('/map/bus-stop/stop-3');
  await page.getByRole('button', { name: 'Add favorite' }).click();
  await page.getByRole('link', { name: 'Favorites' }).click();
  await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Lehmberg' })).toBeVisible();
  await expect(page).toHaveScreenshot();
});
