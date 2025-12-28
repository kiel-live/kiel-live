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
