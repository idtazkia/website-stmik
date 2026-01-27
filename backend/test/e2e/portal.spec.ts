import { test, expect } from '@playwright/test';
import { PortalPage, RegistrationPage, LoginPage } from './pages';

// Generate unique identifiers for test data
function generateUniqueEmail(): string {
  const timestamp = Date.now();
  return `test${timestamp}@example.com`;
}

test.describe('Portal Dashboard - Authentication', () => {
  test('should redirect to login when accessing portal without session', async ({ page }) => {
    // Clear any existing cookies
    await page.context().clearCookies();

    // Try to access portal directly
    await page.goto('/portal');

    // Should redirect to login
    await expect(page).toHaveURL('/login');
  });

  test('should redirect to login when session is invalid', async ({ page }) => {
    // Set an invalid session cookie
    await page.context().addCookies([{
      name: 'session',
      value: 'invalid-token',
      domain: 'localhost',
      path: '/',
    }]);

    await page.goto('/portal');

    // Should redirect to login
    await expect(page).toHaveURL('/login');
  });
});

test.describe('Portal Dashboard - Candidate View', () => {
  let testEmail: string;
  const testPassword = 'testpassword123';
  const testName = 'Portal Test User';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate with complete registration
    const page = await browser.newPage();
    const registrationPage = new RegistrationPage(page);
    testEmail = generateUniqueEmail();

    await registrationPage.goto();
    await registrationPage.expectPageLoaded();
    await registrationPage.fillStep1WithEmail(testEmail, testPassword);
    await registrationPage.expectStep2Visible();
    await registrationPage.fillStep2(testName, 'Jl. Test No. 789', 'Bogor', 'Jawa Barat');
    await registrationPage.expectStep3Visible();

    // Check if there are programs available
    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    const radioCount = await prodiRadios.count();

    if (radioCount > 0) {
      await registrationPage.inputHighSchool.fill('SMA Negeri 1 Bogor');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();
      await registrationPage.expectStep4Visible();
      await registrationPage.fillStep4('google');
    }

    await page.close();
  });

  test('should display dashboard with candidate information after login', async ({ page }) => {
    const loginPage = new LoginPage(page);
    const portalPage = new PortalPage(page);

    // Login
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    // Verify dashboard is loaded
    await portalPage.expectPageLoaded();

    // Verify candidate name is displayed
    await portalPage.expectWelcomeMessage(testName);
  });

  test('should display checklist items', async ({ page }) => {
    const loginPage = new LoginPage(page);
    const portalPage = new PortalPage(page);

    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    await portalPage.expectChecklistVisible();

    // Should have at least 4 checklist items (email verification, personal info, education, documents, payment)
    const itemCount = await portalPage.getChecklistItemCount();
    expect(itemCount).toBeGreaterThanOrEqual(4);
  });

  test('should display announcements section', async ({ page }) => {
    const loginPage = new LoginPage(page);
    const portalPage = new PortalPage(page);

    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    await portalPage.expectAnnouncementsVisible();
  });

  test('should display status badge', async ({ page }) => {
    const loginPage = new LoginPage(page);
    const portalPage = new PortalPage(page);

    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    // Should display status (registered or another valid status)
    await expect(portalPage.statusBadge).toBeVisible();
  });

  test('should navigate to documents page', async ({ page }) => {
    const loginPage = new LoginPage(page);
    const portalPage = new PortalPage(page);

    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    await portalPage.clickDocuments();
    await expect(page).toHaveURL('/portal/documents');
  });

  test('should navigate to payments page', async ({ page }) => {
    const loginPage = new LoginPage(page);
    const portalPage = new PortalPage(page);

    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    await portalPage.clickPayments();
    await expect(page).toHaveURL('/portal/payments');
  });

  test('should logout successfully', async ({ page }) => {
    const loginPage = new LoginPage(page);
    const portalPage = new PortalPage(page);

    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    await portalPage.logout();
    await expect(page).toHaveURL('/login');

    // Trying to access portal should redirect to login
    await page.goto('/portal');
    await expect(page).toHaveURL('/login');
  });
});

test.describe('Portal Documents Page', () => {
  test('should redirect to login when not authenticated', async ({ page }) => {
    await page.context().clearCookies();
    await page.goto('/portal/documents');
    await expect(page).toHaveURL('/login');
  });
});

test.describe('Portal Payments Page', () => {
  test('should redirect to login when not authenticated', async ({ page }) => {
    await page.context().clearCookies();
    await page.goto('/portal/payments');
    await expect(page).toHaveURL('/login');
  });
});
