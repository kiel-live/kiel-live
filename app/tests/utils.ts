import type { Page } from '@playwright/test';
import { test } from '@playwright/test';

export function testColorScheme(name: string, testFunc: ({ page }: { page: Page }) => Promise<void>) {
  (['dark', 'light'] as const).forEach((theme) => {
    test(`${name} [${theme}]`, async ({ page }) => {
      await page.emulateMedia({ colorScheme: theme });
      await testFunc({ page });
    });
  });
}

export async function setLiteMode(page: Page) {
  await page.goto('/');
  await page.getByRole('link', { name: 'Settings' }).click();
  await page.getByRole('main').getByRole('link', { name: 'Settings' }).click();
  await page.getByRole('checkbox', { name: 'Lite mode If you enable the' }).check();
}

export async function waitForMapToLoad(page: Page) {
  await page.waitForSelector('#map[data-idle="true"]');
}
