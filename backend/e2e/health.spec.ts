import { test, expect } from '@playwright/test';
import { PortalPage, AdminPage } from './pages';

test.describe('Health Check API', () => {
  test('GET /health returns ok status with version info', async ({ request }) => {
    const response = await request.get('/health');

    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);

    const body = await response.json();
    expect(body.status).toBe('ok');
    expect(body.version).toBeDefined();
    expect(body.version.commit).toBeDefined();
    expect(body.version.short).toBeDefined();
    expect(body.version.branch).toBeDefined();
    expect(body.version.build_time).toBeDefined();
  });
});

test.describe('Version Display in UI', () => {
  test('Portal page displays version in footer', async ({ page }) => {
    const portalPage = new PortalPage(page);
    await portalPage.goto();
    await portalPage.expectPageLoaded();
    await portalPage.expectVersionVisible();
  });

  test('Admin page displays version in sidebar', async ({ page }) => {
    const adminPage = new AdminPage(page);
    await adminPage.goto();
    await adminPage.expectPageLoaded();
    await adminPage.expectVersionVisible();
  });
});

test.describe('CSRF Protection', () => {
  test('Same-origin form submission works', async ({ page }) => {
    const portalPage = new PortalPage(page);
    await portalPage.goto();

    const response = await portalPage.fillAndSubmitForm('test value');
    expect(response?.ok()).toBeTruthy();
  });

  test('Cross-origin POST request is blocked', async ({ request }) => {
    const response = await request.post('/test/submit', {
      headers: {
        'Origin': 'https://evil-site.com',
        'Sec-Fetch-Site': 'cross-site',
      },
      data: { test: 'value' },
    });

    expect(response.status()).toBe(403);
  });

  test('Request without Sec-Fetch-Site header from same origin works', async ({ request }) => {
    const response = await request.post('/test/submit', {
      data: { test: 'value' },
    });

    expect(response.ok()).toBeTruthy();
  });
});

test.describe('Navigation', () => {
  test('Portal navigation links are present', async ({ page }) => {
    const portalPage = new PortalPage(page);
    await portalPage.goto();
    await portalPage.expectNavigationVisible();
  });

  test('Admin navigation links are present', async ({ page }) => {
    const adminPage = new AdminPage(page);
    await adminPage.goto();
    await adminPage.expectNavigationVisible();
  });
});
