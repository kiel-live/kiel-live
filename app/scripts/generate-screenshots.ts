/**
 * Regenerates the README screenshots (screenshots.png / screenshots_dark.png) using Playwright.
 *
 * Captures three panels against the dummy API - a stop popup, the plain map, and a vehicle
 * popup - for both color schemes, and stitches each trio side by side into the two collages
 * that live at the repo root.
 *
 * Usage: pnpm run generate:screenshots
 */
import type { Browser } from '@playwright/test';
import type { Buffer } from 'node:buffer';

import { spawn } from 'node:child_process';
import { writeFile } from 'node:fs/promises';
import process from 'node:process';
import { setTimeout as delay } from 'node:timers/promises';

import { chromium } from '@playwright/test';

type ColorScheme = 'light' | 'dark';

const BASE_URL = 'http://localhost:4173';
const PANEL_WIDTH = 390;
const PANEL_HEIGHT = 842;
const FIXED_TIME = new Date('2024-02-02T10:30:00');
const OUTPUT = {
  light: new URL('../../screenshots.png', import.meta.url),
  dark: new URL('../../screenshots_dark.png', import.meta.url),
} satisfies Record<ColorScheme, URL>;

async function isServerUp(): Promise<boolean> {
  try {
    const res = await fetch(BASE_URL);
    return res.ok;
  } catch {
    return false;
  }
}

async function startServer() {
  if (await isServerUp()) {
    return null;
  }

  const server = spawn('sh', ['-c', 'pnpm run build:test && pnpm run preview'], {
    env: { ...process.env, VITE_USE_DUMMY_API: 'true' },
    stdio: 'inherit',
    detached: true,
  });

  for (let attempt = 0; attempt < 60; attempt++) {
    if (await isServerUp()) {
      return server;
    }
    await delay(1000);
  }

  server.kill();
  throw new Error(`Server did not become ready at ${BASE_URL}`);
}

function stopServer(server: ReturnType<typeof spawn> | null) {
  if (server?.pid) {
    process.kill(-server.pid, 'SIGTERM');
  }
}

async function capturePanel(
  browser: Browser,
  colorScheme: ColorScheme,
  goto: string,
  ready: (page: import('@playwright/test').Page) => Promise<void>,
): Promise<Buffer> {
  const context = await browser.newContext({
    baseURL: BASE_URL,
    viewport: { width: PANEL_WIDTH, height: PANEL_HEIGHT },
    deviceScaleFactor: 1,
    colorScheme,
  });
  const page = await context.newPage();
  await page.clock.setFixedTime(FIXED_TIME);
  await page.goto(goto);
  await ready(page);
  const buffer = await page.screenshot();
  await context.close();
  return buffer;
}

async function waitForMapToLoad(page: import('@playwright/test').Page) {
  await page.waitForSelector('#map[data-idle="true"]');
}

// dragging the map marks it as "moved manually", which shrinks the popup from
// 3/4 to 1/2 height (see Home.vue's popupSize) - leaves more map visible
async function panMapSlightly(page: import('@playwright/test').Page) {
  const mapBox = await page.getByRole('region', { name: 'Map' }).boundingBox();
  if (!mapBox) {
    return;
  }
  const x = mapBox.x + mapBox.width / 2;
  const y = mapBox.y + Math.min(120, mapBox.height - 10);
  await page.mouse.move(x, y);
  await page.mouse.down();
  await page.mouse.move(x - 40, y - 15, { steps: 10 });
  await page.mouse.up();
  await page.waitForTimeout(400);
}

async function captureStopPopup(browser: Browser, colorScheme: ColorScheme) {
  return capturePanel(browser, colorScheme, '/map/bus-stop/kvg-2387', async (page) => {
    await page.getByRole('heading', { name: 'Hauptbahnhof' }).waitFor();
    await waitForMapToLoad(page);
    await panMapSlightly(page);
  });
}

async function capturePlainMap(browser: Browser, colorScheme: ColorScheme) {
  return capturePanel(browser, colorScheme, '/', async (page) => {
    await waitForMapToLoad(page);
    // the root route lands on the favorites popup - clicking the map dismisses it
    await page.getByRole('region', { name: 'Map' }).click({ position: { x: 100, y: 300 } });
    await page.getByRole('heading', { name: 'Favorites' }).waitFor({ state: 'hidden' });
  });
}

async function captureVehiclePopup(browser: Browser, colorScheme: ColorScheme) {
  return capturePanel(browser, colorScheme, '/map/bus/bus-3', async (page) => {
    await page.getByRole('heading', { name: '62 Russee, Schiefe Horn' }).waitFor();
    await waitForMapToLoad(page);
    await panMapSlightly(page);
  });
}

async function stitch(browser: Browser, panels: Buffer[]): Promise<Buffer> {
  const context = await browser.newContext({
    viewport: { width: PANEL_WIDTH * panels.length, height: PANEL_HEIGHT },
    deviceScaleFactor: 1,
  });
  const page = await context.newPage();

  const images = panels
    .map(
      (panel) =>
        `<img src="data:image/png;base64,${panel.toString('base64')}" width="${PANEL_WIDTH}" height="${PANEL_HEIGHT}">`,
    )
    .join('');

  await page.setContent(
    `<!doctype html><html><head><style>
      body { margin: 0; }
      div { display: flex; }
      img { display: block; }
    </style></head><body><div>${images}</div></body></html>`,
  );
  await page.waitForFunction(() => [...document.images].every((img) => img.complete));

  const buffer = await page.screenshot();
  await context.close();
  return buffer;
}

async function main() {
  const server = await startServer();
  const browser = await chromium.launch();

  try {
    for (const colorScheme of ['light', 'dark'] as const) {
      const panels = await Promise.all([
        captureStopPopup(browser, colorScheme),
        capturePlainMap(browser, colorScheme),
        captureVehiclePopup(browser, colorScheme),
      ]);
      const collage = await stitch(browser, panels);
      await writeFile(OUTPUT[colorScheme], collage);
      console.log(`Wrote ${OUTPUT[colorScheme].pathname}`);
    }
  } finally {
    await browser.close();
    stopServer(server);
  }
}

await main();
