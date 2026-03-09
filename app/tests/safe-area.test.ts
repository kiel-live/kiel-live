import type { Page } from '@playwright/test';
import { expect, test } from '@playwright/test';
import { setLiteMode, waitForMapToLoad } from './utils';

async function injectSafeAreaVars(page: Page) {
  await page.addStyleTag({
    content: `
      :root {
        --safe-area-top: 48px;
        --safe-area-bottom: 32px;
      }
    `,
  });
}

test.describe(
  'Safe area',
  {
    tag: ['@mobile-only'],
  },
  () => {
    test('Safe area on map', async ({ page }) => {
      await page.goto('/');
      await injectSafeAreaVars(page);
      await waitForMapToLoad(page);

      await expect(page).toHaveScreenshot();
    });

    test('Safe area on map with maximized popup', async ({ page }) => {
      await page.goto('/');
      await injectSafeAreaVars(page);
      await waitForMapToLoad(page);
      await page
        .getByRole('button', { name: 'Drag to resize popup' })
        .dragTo(page.getByRole('textbox', { name: 'Search' }));

      await expect(page).toHaveScreenshot();
    });

    test('Safe area on settings', async ({ page }) => {
      await page.goto('/');
      await injectSafeAreaVars(page);
      await page.getByRole('link', { name: 'Settings' }).click();

      await expect(page).toHaveScreenshot();
    });

    test('Safe area in lite mode', async ({ page }) => {
      await setLiteMode(page);
      await injectSafeAreaVars(page);
      await page.getByRole('link', { name: 'Search' }).click();

      await expect(page).toHaveScreenshot();
    });
  },
);
