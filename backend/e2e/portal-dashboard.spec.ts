import { test, expect } from '@playwright/test';
import { RegistrationPage, LoginPage, PortalPage } from './pages';

// Generate unique identifiers for test data with prefix to avoid collisions
function generateUniqueEmail(prefix: string): string {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 10000);
  return `${prefix}${timestamp}${random}@example.com`;
}

test.describe('Candidate Portal Dashboard', () => {
  let portalPage: PortalPage;
  let testEmail: string;
  const testPassword = 'testpassword123';
  const testName = 'Dashboard Test User';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate through full registration
    const page = await browser.newPage();
    const registrationPage = new RegistrationPage(page);
    testEmail = generateUniqueEmail('dash');

    await registrationPage.goto();
    await registrationPage.expectPageLoaded();

    // Step 1: Account
    await registrationPage.fillStep1WithEmail(testEmail, testPassword);
    await registrationPage.expectStep2Visible();

    // Step 2: Personal Info
    await registrationPage.fillStep2(
      testName,
      'Jl. Dashboard Test No. 123',
      'Jakarta',
      'DKI Jakarta'
    );
    await registrationPage.expectStep3Visible();

    // Step 3: Education
    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    const radioCount = await prodiRadios.count();

    if (radioCount > 0) {
      await registrationPage.inputHighSchool.fill('SMA Negeri 1 Jakarta');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();
      await registrationPage.expectStep4Visible();

      // Step 4: Source Tracking
      await registrationPage.fillStep4('google');
    }

    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    // Login before each test
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    portalPage = new PortalPage(page);
  });

  test('should display dashboard after login', async () => {
    await portalPage.expectPageLoaded();
  });

  test('should display candidate name in welcome banner', async () => {
    await portalPage.expectWelcomeMessage(testName);
  });

  test('should display registration status', async () => {
    // After completing all 4 steps, status is "Dalam Proses" (prospecting)
    await portalPage.expectStatus('Dalam Proses');
  });

  test('should display checklist items', async () => {
    await portalPage.expectChecklistVisible();
    const itemCount = await portalPage.getChecklistItemCount();
    expect(itemCount).toBeGreaterThan(0);
  });

  test('should display assigned consultant', async () => {
    // Consultant should be assigned during registration
    // Check if consultant info is visible
    const consultantSection = portalPage.consultantSection;
    await expect(consultantSection).toBeVisible();
  });

  test('should display announcements section', async () => {
    await portalPage.expectAnnouncementsVisible();
  });

  test('should navigate to documents page', async ({ page }) => {
    await portalPage.clickDocuments();
    await expect(page).toHaveURL('/portal/documents');
  });

  test('should navigate to payments page', async ({ page }) => {
    await portalPage.clickPayments();
    await expect(page).toHaveURL('/portal/payments');
  });

  test('should logout successfully', async ({ page }) => {
    await portalPage.logout();
    // After logout, should redirect to login page
    await expect(page).toHaveURL('/login');
  });

  test('should redirect to login when accessing portal without session', async ({ page }) => {
    // Clear cookies
    await page.context().clearCookies();
    // Try to access portal directly
    await page.goto('/portal');
    // Should redirect to login
    await expect(page).toHaveURL('/login');
  });
});

test.describe('Portal Documents Page', () => {
  let testEmail: string;
  const testPassword = 'testpassword123';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate
    const page = await browser.newPage();
    const registrationPage = new RegistrationPage(page);
    testEmail = generateUniqueEmail('docs');

    await registrationPage.goto();
    await registrationPage.fillStep1WithEmail(testEmail, testPassword);
    await registrationPage.expectStep2Visible();
    await registrationPage.fillStep2(
      'Documents Test User',
      'Jl. Documents Test',
      'Bandung',
      'Jawa Barat'
    );
    await registrationPage.expectStep3Visible();

    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    if ((await prodiRadios.count()) > 0) {
      await registrationPage.inputHighSchool.fill('SMA Test');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();
      await registrationPage.expectStep4Visible();
      await registrationPage.fillStep4('instagram');
    }

    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();
  });

  test('should display documents page', async ({ page }) => {
    await page.goto('/portal/documents');
    // Check for document types using headings
    await expect(page.getByRole('heading', { name: 'KTP' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Pas Foto' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Ijazah' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Transkrip Nilai' })).toBeVisible();
  });

  test('should show document upload requirements', async ({ page }) => {
    await page.goto('/portal/documents');
    // Check for format info - use first() to avoid ambiguity
    await expect(page.locator('text=JPG, PNG, PDF').first()).toBeVisible();
  });
});

test.describe('Portal Payments Page', () => {
  let testEmail: string;
  const testPassword = 'testpassword123';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate
    const page = await browser.newPage();
    const registrationPage = new RegistrationPage(page);
    testEmail = generateUniqueEmail('pay');

    await registrationPage.goto();
    await registrationPage.fillStep1WithEmail(testEmail, testPassword);
    await registrationPage.expectStep2Visible();
    await registrationPage.fillStep2(
      'Payments Test User',
      'Jl. Payments Test',
      'Surabaya',
      'Jawa Timur'
    );
    await registrationPage.expectStep3Visible();

    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    if ((await prodiRadios.count()) > 0) {
      await registrationPage.inputHighSchool.fill('SMA Test');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();
      await registrationPage.expectStep4Visible();
      await registrationPage.fillStep4('youtube');
    }

    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();
  });

  test('should display payments page', async ({ page }) => {
    await page.goto('/portal/payments');
    // Page should load without error - use exact heading role
    await expect(page.getByRole('heading', { name: 'Pembayaran', exact: true })).toBeVisible();
  });

  test('should show payment summary', async ({ page }) => {
    await page.goto('/portal/payments');
    // Check for summary section (even if empty) - use first() to avoid ambiguity
    await expect(page.locator('text=Total').first()).toBeVisible();
  });
});

test.describe('Portal Announcements Page', () => {
  let testEmail: string;
  const testPassword = 'testpassword123';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate
    const page = await browser.newPage();
    const registrationPage = new RegistrationPage(page);
    testEmail = generateUniqueEmail('ann');

    await registrationPage.goto();
    await registrationPage.fillStep1WithEmail(testEmail, testPassword);
    await registrationPage.expectStep2Visible();
    await registrationPage.fillStep2(
      'Announcements Test User',
      'Jl. Announcements Test',
      'Yogyakarta',
      'DI Yogyakarta'
    );
    await registrationPage.expectStep3Visible();

    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    if ((await prodiRadios.count()) > 0) {
      await registrationPage.inputHighSchool.fill('SMA Test');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();
      await registrationPage.expectStep4Visible();
      await registrationPage.fillStep4('tiktok');
    }

    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();
  });

  test('should display announcements page', async ({ page }) => {
    await page.goto('/portal/announcements');
    // Page should load without error - use heading role
    await expect(page.getByRole('heading', { name: 'Pengumuman' })).toBeVisible();
  });
});
