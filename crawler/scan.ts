import puppeteer from 'puppeteer';

export interface ScanResult {
  url: string;
  trackers: string[];
  cookies: any[];
  requests: string[];
  timestamp: string;
};

export async function scanForTrackers(url: string): Promise<ScanResult> {
  console.log(`Starting scan for: ${url}`);
  
  let browser;
  try {
    browser = await puppeteer.launch({ 
      headless: true,
      args: [
        '--no-sandbox', 
        '--disable-setuid-sandbox',
        '--disable-dev-shm-usage',
        '--disable-accelerated-2d-canvas',
        '--no-first-run',
        '--no-zygote',
        '--disable-gpu'
      ]
    });
    
    console.log('Browser launched successfully');
    
    const page = await browser.newPage();
    const requests: string[] = [];
    const trackers: string[] = [];
    
    // Set a user agent to avoid being blocked
    await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36');
    
    // Common tracker domains to detect
    const trackerDomains = [
      'google-analytics.com',
      'googletagmanager.com',
      'facebook.com',
      'doubleclick.net',
      'amazon-adsystem.com',
      'googlesyndication.com',
      'adsystem.amazon.com',
      'scorecardresearch.com',
      'quantserve.com',
      'hotjar.com'
    ];

    // Intercept network requests
    page.on('request', (request) => {
      const requestUrl = request.url();
      requests.push(requestUrl);
      
      // Check if request is to a known tracker domain
      for (const trackerDomain of trackerDomains) {
        if (requestUrl.includes(trackerDomain)) {
          trackers.push(requestUrl);
          break;
        }
      }
    });

    
    page.on('error', (error) => {
      console.error('Page error:', error);
    });

    page.on('pageerror', (error) => {
      console.error('Page script error:', error);
    });

    console.log('Navigating to page...');
    
    // Navigate to the page with timeout
    // await page.goto(url, { 
    //   waitUntil: 'networkidle2',
    //   timeout: 30000 
    // });
    
    console.log('Page loaded, getting cookies...');
    
   
    const cookies = await browser.cookies();
    
    console.log(`Scan completed. Found ${trackers.length} trackers, ${cookies.length} cookies, ${requests.length} requests`);
    
    // Remove duplicate trackers and requests
    return {
      url,
      trackers: [...new Set(trackers)], 
      cookies,
      requests: [...new Set(requests)], 
      timestamp: new Date().toISOString()
    };
    
  } catch (error: unknown) {
    console.error('Error during scan:', error);
    const errorMessage = error instanceof Error ? error.message : String(error);
    throw new Error(`Failed to scan ${url}: ${errorMessage}`);
  } finally {
    if (browser) {
      console.log('Closing browser...');
      await browser.close();
    }
  };
}