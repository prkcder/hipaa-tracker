// import puppeteer from 'puppeteer';
// const knownTrackers: Record<string, { name: string; risk: string }> = {
//   'google-analytics.com': { name: 'Google Analytics', risk: 'High' },
//   'facebook.net': { name: 'Facebook Pixel', risk: 'High' },
//   'newrelic.com': { name: 'New Relic', risk: 'Medium' },
// };
// type Tracker = {
//   tracker_name: string;
//   tracker_domain: string;
//   risk_level: string;
//   page_url: string;
// };
// (async () => {
//   const targetUrl = process.argv[2];
//   if (!targetUrl) {
//     console.error('Usage: node scan.js <url>');
//     process.exit(1);
//   }
//   const browser = await puppeteer.launch();
//   const page = await browser.newPage();
//   const trackers = new Set<string>();
//   page.on('request', (req) => {
//     const url = new URL(req.url());
//     const matchedDomain = Object.keys(knownTrackers).find((domain) =>
//       url.hostname.includes(domain)
//     );
//     if (matchedDomain) {
//       trackers.add(
//         JSON.stringify({
//           tracker_name: knownTrackers[matchedDomain].name,
//           tracker_domain: matchedDomain,
//           risk_level: knownTrackers[matchedDomain].risk,
//           page_url: page.url(),
//         })
//       );
//     }
//   });
//   try {
//     await page.goto(targetUrl, { waitUntil: 'networkidle2' });
//     await new Promise((res) => setTimeout(res, 3000)); // wait for trackers
//   } catch (err) {
//     console.error('Failed to load page:', err);
//     process.exit(1);
//   } finally {
//     await browser.close();
//   }
//   const results: Tracker[] = [...trackers].map((t) => JSON.parse(t));
//   console.log(JSON.stringify(results, null, 2));
// })();
// import puppeteer from 'puppeteer';
// export async function runScan(url: string) {
//     const browser = await puppeteer.launch({ headless: 'new' });
//     const page = await browser.newPage();
//     await page.goto(url, { waitUntil: 'networkidle2' });
//     const scripts = await page.evaluate(() => {
//         return Array.from(document.scripts).map((s) => s.src).filter(Boolean);
//     });
//     await browser.close();
//     // You can expand the logic to match domains and assign risk levels
//     const results = scripts.map((src) => ({
//         src,
//         risk: src.includes('google') ? 'high' : 'medium', // demo logic
//     }));
//     return {
//         url,
//         scannedAt: new Date().toISOString(),
//         trackerCount: results.length,
//         trackers: results,
//     };
// }
// scan.ts
import puppeteer from 'puppeteer';
export async function scanForTrackers(url) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto(url, { waitUntil: 'networkidle2' });
    const requests = new Set();
    page.on('request', (req) => {
        const reqUrl = new URL(req.url());
        requests.add(reqUrl.hostname);
    });
    // Wait for page + scripts to settle
    await page.waitForTimeout(5000);
    await browser.close();
    const results = Array.from(requests).map((domain) => ({
        url,
        domain,
    }));
    return results;
}
//# sourceMappingURL=scan.js.map