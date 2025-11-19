import { test, expect } from '@playwright/test';

const BASE_URL = 'https://endymuhardin.github.io/website-stmik';

test.describe('Deployed Website Check', () => {
  test.describe('Indonesian Homepage', () => {
    test('should load homepage without errors', async ({ page }) => {
      const response = await page.goto(BASE_URL);
      expect(response?.status()).toBe(200);
    });

    test('should display site name and logo', async ({ page }) => {
      await page.goto(BASE_URL);

      // Check title
      await expect(page).toHaveTitle(/STMIK Tazkia/);

      // Check logo
      const logo = page.locator('img[alt*="STMIK"]').first();
      await expect(logo).toBeVisible();

      // Check h1 heading
      const heading = page.locator('h1:has-text("STMIK Tazkia")');
      await expect(heading).toBeVisible();
    });

    test('should have no console errors', async ({ page }) => {
      const errors: string[] = [];
      page.on('console', msg => {
        if (msg.type() === 'error') {
          errors.push(msg.text());
        }
      });

      await page.goto(BASE_URL);
      await page.waitForLoadState('networkidle');

      expect(errors).toEqual([]);
    });

    test('should load all images without 404', async ({ page }) => {
      const failed404s: string[] = [];

      page.on('response', response => {
        if (response.status() === 404 && response.request().resourceType() === 'image') {
          failed404s.push(response.url());
        }
      });

      await page.goto(BASE_URL);
      await page.waitForLoadState('networkidle');

      if (failed404s.length > 0) {
        console.log('Failed images:', failed404s);
      }
      expect(failed404s).toEqual([]);
    });

    test('should load all stylesheets without errors', async ({ page }) => {
      const failedCSS: string[] = [];

      page.on('response', response => {
        if (response.status() !== 200 && response.request().resourceType() === 'stylesheet') {
          failedCSS.push(`${response.url()} - ${response.status()}`);
        }
      });

      await page.goto(BASE_URL);
      await page.waitForLoadState('networkidle');

      if (failedCSS.length > 0) {
        console.log('Failed CSS:', failedCSS);
      }
      expect(failedCSS).toEqual([]);
    });

    test('should have navigation menu', async ({ page }) => {
      await page.goto(BASE_URL);

      // Check for navigation links
      const homeLink = page.locator('a:has-text("Beranda"), nav a[href*="index"]').first();
      const programsLink = page.locator('a:has-text("Program")').first();

      // At least some navigation should exist
      const navLinks = await page.locator('nav a, header a').count();
      expect(navLinks).toBeGreaterThan(0);
    });

    test('should have footer with contact info', async ({ page }) => {
      await page.goto(BASE_URL);

      const footer = page.locator('footer');
      await expect(footer).toBeVisible();

      // Check for contact information
      const hasPhone = await footer.locator('text=/\\+62|phone|telepon/i').count() > 0;
      const hasAddress = await footer.locator('text=/bogor|dramaga|alamat|address/i').count() > 0;

      expect(hasPhone || hasAddress).toBeTruthy();
    });

    test('should have working CTA buttons', async ({ page }) => {
      await page.goto(BASE_URL);

      // Check for "Daftar Sekarang" button
      const ctaButton = page.locator('a:has-text("Daftar Sekarang")').first();
      await expect(ctaButton).toBeVisible();

      // Button should have href
      const href = await ctaButton.getAttribute('href');
      expect(href).toBeTruthy();
    });

    test('should display features section', async ({ page }) => {
      await page.goto(BASE_URL);

      // Check for features/keunggulan section
      const featuresHeading = page.locator('h2:has-text("Keunggulan")');
      await expect(featuresHeading).toBeVisible();

      // Should have multiple feature items
      const featureCount = await page.locator('h3').count();
      expect(featureCount).toBeGreaterThan(2);
    });

    test('should display programs section', async ({ page }) => {
      await page.goto(BASE_URL);

      // Check for programs section
      const programsHeading = page.locator('h2:has-text("Program Studi")');
      await expect(programsHeading).toBeVisible();

      // Check for specific programs
      const systemsInfo = page.locator('text=/Sistem Informasi/i');
      const computerEng = page.locator('text=/Teknik Informatika/i');

      await expect(systemsInfo).toBeVisible();
      await expect(computerEng).toBeVisible();
    });
  });

  test.describe('English Homepage', () => {
    test('should load English homepage', async ({ page }) => {
      const response = await page.goto(`${BASE_URL}/en/`);
      expect(response?.status()).toBe(200);
    });

    test('should display content in English', async ({ page }) => {
      await page.goto(`${BASE_URL}/en/`);

      // Check for English content
      const hasEnglishContent = await page.locator('text=/Learn more|Apply now|Information Systems|Computer Engineering/i').count() > 0;
      expect(hasEnglishContent).toBeGreaterThan(0);
    });

    test('should load logo on English page', async ({ page }) => {
      const failed404s: string[] = [];

      page.on('response', response => {
        if (response.status() === 404 && response.request().resourceType() === 'image') {
          failed404s.push(response.url());
        }
      });

      await page.goto(`${BASE_URL}/en/`);
      await page.waitForLoadState('networkidle');

      if (failed404s.length > 0) {
        console.log('Failed images on EN page:', failed404s);
      }
      expect(failed404s).toEqual([]);
    });
  });

  test.describe('Language Switching', () => {
    test('should have language switcher', async ({ page }) => {
      await page.goto(BASE_URL);

      // Look for language switcher (EN/ID links)
      const langSwitcher = await page.locator('a[href*="/en"], button:has-text("EN"), a:has-text("English")').count();
      expect(langSwitcher).toBeGreaterThan(0);
    });
  });

  test.describe('Responsive Design', () => {
    test('should be mobile responsive', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto(BASE_URL);

      // Check if logo is visible on mobile
      const logo = page.locator('img[alt*="STMIK"]').first();
      await expect(logo).toBeVisible();

      // Check if main heading is visible
      const heading = page.locator('h1');
      await expect(heading).toBeVisible();
    });

    test('should be tablet responsive', async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 });
      await page.goto(BASE_URL);

      const heading = page.locator('h1');
      await expect(heading).toBeVisible();
    });

    test('should be desktop responsive', async ({ page }) => {
      await page.setViewportSize({ width: 1920, height: 1080 });
      await page.goto(BASE_URL);

      const heading = page.locator('h1');
      await expect(heading).toBeVisible();
    });
  });

  test.describe('SEO and Meta Tags', () => {
    test('should have proper meta tags', async ({ page }) => {
      await page.goto(BASE_URL);

      // Check meta description
      const metaDescription = page.locator('meta[name="description"]');
      await expect(metaDescription).toHaveAttribute('content', /.+/);

      // Check Open Graph tags
      const ogTitle = page.locator('meta[property="og:title"]');
      await expect(ogTitle).toHaveCount(1);
    });

    test('should have sitemap', async ({ page }) => {
      const response = await page.goto(`${BASE_URL}/sitemap-index.xml`);
      expect(response?.status()).toBe(200);

      const content = await page.content();
      expect(content).toContain('<?xml');
      expect(content).toContain('sitemap');
    });

    test('should have robots.txt', async ({ page }) => {
      const response = await page.goto(`${BASE_URL}/robots.txt`);
      expect(response?.status()).toBe(200);
    });
  });

  test.describe('Performance', () => {
    test('should load within reasonable time', async ({ page }) => {
      const startTime = Date.now();
      await page.goto(BASE_URL);
      await page.waitForLoadState('networkidle');
      const loadTime = Date.now() - startTime;

      console.log(`Page load time: ${loadTime}ms`);
      expect(loadTime).toBeLessThan(10000); // 10 seconds max
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper heading hierarchy', async ({ page }) => {
      await page.goto(BASE_URL);

      // Should have h1
      const h1Count = await page.locator('h1').count();
      expect(h1Count).toBeGreaterThanOrEqual(1);

      // Should not have more than one h1
      expect(h1Count).toBeLessThanOrEqual(1);
    });

    test('should have alt text on images', async ({ page }) => {
      await page.goto(BASE_URL);

      const images = page.locator('img');
      const count = await images.count();

      for (let i = 0; i < count; i++) {
        const img = images.nth(i);
        const alt = await img.getAttribute('alt');
        expect(alt).toBeTruthy();
      }
    });

    test('should have proper lang attribute', async ({ page }) => {
      await page.goto(BASE_URL);

      const htmlLang = await page.locator('html').getAttribute('lang');
      expect(htmlLang).toBe('id');
    });

    test('should have proper lang attribute on EN page', async ({ page }) => {
      await page.goto(`${BASE_URL}/en/`);

      const htmlLang = await page.locator('html').getAttribute('lang');
      expect(htmlLang).toBe('en');
    });
  });
});
