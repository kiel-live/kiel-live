import { expect, test } from '@playwright/test';
import { waitForMapToLoad } from './utils';

const ANIMATION_DURATION = 400;

test.describe(
  'Popup behavior',
  {
    tag: ['@mobile-only'],
  },
  () => {
    test('Clicking the map closes the popup', async ({ page }) => {
      await page.goto('/search');
      await page.getByRole('region', { name: 'Map' }).click({ position: { x: 100, y: 100 } });
      await expect(page.getByRole('heading', { name: 'Search result' })).not.toBeVisible();
    });

    test('moving the popup down closes it', async ({ page }) => {
      await page.goto('/favorites');
      await page
        .getByRole('button', { name: 'Drag to resize popup' })
        .dragTo(page.getByRole('link', { name: 'Settings' }));

      await expect(page.getByRole('heading', { name: 'Favorites' })).not.toBeVisible();
    });

    test('moving the popup up makes it larger', async ({ page }) => {
      await page.goto('/favorites');
      await page
        .getByRole('button', { name: 'Drag to resize popup' })
        .dragTo(page.getByRole('textbox', { name: 'Search' }));

      await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
      await waitForMapToLoad(page);
      await expect(page).toHaveScreenshot();
    });

    test('moving the popup up and then to the middle keeps it half size', async ({ page }) => {
      await page.goto('/favorites');
      const dragHandle = page.getByRole('button', { name: 'Drag to resize popup' });
      await dragHandle.dragTo(page.getByRole('textbox', { name: 'Search' }));
      await page.waitForTimeout(ANIMATION_DURATION);
      await dragHandle.hover();
      await page.mouse.down();
      await page.mouse.move(100, page.viewportSize()!.height / 2);
      await page.mouse.up();

      await expect(page.getByRole('heading', { name: 'Favorites' })).toBeVisible();
      await waitForMapToLoad(page);
      await expect(page).toHaveScreenshot();
    });

    test('moving the popup up and then down closes it', async ({ page }) => {
      await page.goto('/favorites');
      const dragHandle = page.getByRole('button', { name: 'Drag to resize popup' });
      await dragHandle.dragTo(page.getByRole('textbox', { name: 'Search' }));
      await page.waitForTimeout(ANIMATION_DURATION);
      await dragHandle.dragTo(page.getByRole('link', { name: 'Settings' }));

      await expect(page.getByRole('heading', { name: 'Favorites' })).not.toBeVisible();
    });

    test('popup stays maximized when switching content', async ({ page }) => {
      await page.goto('/favorites');
      const dragHandle = page.getByRole('button', { name: 'Drag to resize popup' });
      await dragHandle.dragTo(page.getByRole('textbox', { name: 'Search' }));
      await page.getByRole('textbox', { name: 'Search' }).fill('stop');

      await expect(page.getByRole('heading', { name: 'Search result' })).toBeVisible();
      await waitForMapToLoad(page);
      await expect(page).toHaveScreenshot();
    });
  },
);
