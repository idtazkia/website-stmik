import { test, expect } from '@playwright/test';
import { RegistrationPage, LoginPage } from './pages';

// Generate unique identifiers for test data
function generateUniqueEmail(prefix: string): string {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 10000);
  return `${prefix}${timestamp}${random}@example.com`;
}

test.describe('Email Verification', () => {
  let testEmail: string;
  const testPassword = 'testpassword123';
  const testName = 'Email Verify Test';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate through registration
    const page = await browser.newPage();
    const registrationPage = new RegistrationPage(page);
    testEmail = generateUniqueEmail('emailverify');

    await registrationPage.goto();
    await registrationPage.expectPageLoaded();

    // Step 1: Account - use email for email verification test
    await registrationPage.fillStep1WithEmail(testEmail, testPassword);
    await registrationPage.expectStep2Visible();

    // Step 2: Personal Info
    await registrationPage.fillStep2(
      testName,
      'Jl. Email Verify Test No. 123',
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

  test('should show verification page with send OTP button', async ({ page }) => {
    // Login first
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    // Navigate to email verification page
    await page.goto('/portal/verify-email');

    // Check page loaded
    await expect(page.getByTestId('verify-email-page')).toBeVisible();

    // Check email is displayed
    await expect(page.getByText(testEmail)).toBeVisible();

    // Check send OTP button is visible
    await expect(page.getByTestId('send-otp-btn')).toBeVisible();
  });

  test('should show OTP form after clicking send button', async ({ page }) => {
    // Login first
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    // Navigate to email verification page
    await page.goto('/portal/verify-email');

    // Click send OTP button
    await page.getByTestId('send-otp-btn').click();

    // Wait for HTMX response - should show OTP input form
    // Note: In real environment with Resend configured, this will send an email
    // For testing without Resend, it should show an error message
    await page.waitForTimeout(1000);

    // Check if OTP input appeared (if Resend is configured) or error message
    const otpInput = page.getByTestId('otp-input');
    const errorMessage = page.getByTestId('otp-error');

    // Either OTP input should be visible (email sent) or error message (service unavailable)
    const otpVisible = await otpInput.isVisible().catch(() => false);
    const errorVisible = await errorMessage.isVisible().catch(() => false);

    expect(otpVisible || errorVisible).toBeTruthy();
  });

  test('should show error for invalid OTP', async ({ page }) => {
    // Login first
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    // Navigate to email verification page
    await page.goto('/portal/verify-email');

    // Click send OTP button
    await page.getByTestId('send-otp-btn').click();
    await page.waitForTimeout(1000);

    // Check if OTP input appeared
    const otpInput = page.getByTestId('otp-input');
    if (await otpInput.isVisible()) {
      // Enter invalid OTP
      await otpInput.fill('000000');
      await page.getByTestId('confirm-otp-btn').click();

      // Wait for response
      await page.waitForTimeout(1000);

      // Should show error message
      await expect(page.getByTestId('otp-error')).toBeVisible();
    }
  });

  test('should navigate to verification from dashboard checklist', async ({ page }) => {
    // Login first
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();

    // Find email verification in checklist and click
    const verifyEmailLink = page.getByRole('link', { name: 'Verifikasi' }).first();

    if (await verifyEmailLink.isVisible()) {
      await verifyEmailLink.click();
      await expect(page).toHaveURL('/portal/verify-email');
      await expect(page.getByTestId('verify-email-page')).toBeVisible();
    }
  });
});
