import fs from "fs";
import open from "open";
import puppeteer from "puppeteer";
import { startFlow } from "lighthouse/lighthouse-core/fraggle-rock/api.js";

async function captureReport() {
  const browser = await puppeteer.launch({ headless: false });
  const page = await browser.newPage();

  const flow = await startFlow(page, { name: "Kiel-Live snapshots" });
  await flow.navigate('http://127.0.0.1:5173/');

  const aboutSelector = 'a[href="/settings/about"]';
  await page.waitForSelector(aboutSelector);
  await flow.snapshot({ stepName: "Page loaded" });
  await page.click(aboutSelector);

  const settingsSelector = 'a[href="/settings/settings"]';
  await page.waitForSelector(settingsSelector);
  await flow.snapshot({ stepName: "Settings loaded" });

  await page.click(settingsSelector);
  await flow.snapshot({ stepName: "Settings opened" });

  await page.click('a[href="/settings/about"]');
  await page.click('a[href="/settings/changelog"]');
  await flow.snapshot({ stepName: "Changelog opened" });

  await page.click('a[href="/"]');
  await page.click("input");
  await flow.snapshot({ stepName: "Search opened" });

  await page.type("input", "test");
  await flow.snapshot({ stepName: "Search typed" });

  await page.click('a[href="/map/bus-stop/kvg-2932"]');
  await flow.snapshot({ stepName: "Stop opened" });

  await page.click("div.overflow-y-auto a:first-child");
  await flow.snapshot({ stepName: "Route opened" });

  browser.close();

  const report = await flow.generateReport();
  fs.writeFileSync("flow.report.html", report);
  open("flow.report.html", { wait: false });
}

captureReport();
