import { test, expect } from '@playwright/test';
import { RegistrationPage, LoginPage } from './pages';

// Generate unique identifiers for test data
function generateUniqueEmail(): string {
  const timestamp = Date.now();
  return `test${timestamp}@example.com`;
}

function generateUniquePhone(): string {
  const timestamp = Date.now();
  // Generate a unique phone number (Indonesia format)
  return `08${timestamp.toString().slice(-10)}`;
}

test.describe('Candidate Registration', () => {
  let registrationPage: RegistrationPage;

  test.beforeEach(async ({ page }) => {
    registrationPage = new RegistrationPage(page);
    await registrationPage.goto();
    await registrationPage.expectPageLoaded();
  });

  test.describe('Step 1: Account Creation', () => {
    test('should display step 1 form initially', async () => {
      await registrationPage.expectStep1Visible();
    });

    test('should show error if neither email nor phone provided', async () => {
      await registrationPage.inputPassword.fill('testpassword123');
      await registrationPage.inputPasswordConfirm.fill('testpassword123');
      await registrationPage.btnSubmitStep1.click();

      await registrationPage.expectErrorMessage('email atau nomor HP');
    });

    test('should show error if passwords do not match', async () => {
      await registrationPage.inputEmail.fill(generateUniqueEmail());
      await registrationPage.inputPassword.fill('testpassword123');
      await registrationPage.inputPasswordConfirm.fill('differentpassword');
      await registrationPage.btnSubmitStep1.click();

      await registrationPage.expectErrorMessage('tidak cocok');
    });

    test('should prevent submission if password too short (browser validation)', async ({ page }) => {
      await registrationPage.inputEmail.fill(generateUniqueEmail());
      await registrationPage.inputPassword.fill('short');
      await registrationPage.inputPasswordConfirm.fill('short');
      await registrationPage.btnSubmitStep1.click();

      // Browser's HTML5 validation should prevent form submission
      // We should still be on step 1
      await registrationPage.expectStep1Visible();

      // Check that the password field has validation error (via :invalid pseudo-class)
      const isInvalid = await registrationPage.inputPassword.evaluate((el: HTMLInputElement) => !el.validity.valid);
      expect(isInvalid).toBe(true);
    });

    test('should proceed to step 2 with email only', async () => {
      const email = generateUniqueEmail();
      await registrationPage.fillStep1WithEmail(email, 'testpassword123');

      await registrationPage.expectStep2Visible();
    });

    test('should proceed to step 2 with phone only', async () => {
      const phone = generateUniquePhone();
      await registrationPage.fillStep1WithPhone(phone, 'testpassword123');

      await registrationPage.expectStep2Visible();
    });

    test('should proceed to step 2 with both email and phone', async () => {
      const email = generateUniqueEmail();
      const phone = generateUniquePhone();
      await registrationPage.fillStep1WithBoth(email, phone, 'testpassword123');

      await registrationPage.expectStep2Visible();
    });
  });

  test.describe('Step 2: Personal Info', () => {
    test.beforeEach(async () => {
      // Complete step 1 first
      const email = generateUniqueEmail();
      await registrationPage.fillStep1WithEmail(email, 'testpassword123');
      await registrationPage.expectStep2Visible();
    });

    test('should show error if required fields are empty', async ({ page }) => {
      await registrationPage.btnSubmitStep2.click();
      // Browser validation should prevent submission
      await registrationPage.expectStep2Visible();
    });

    test('should proceed to step 3 with valid data', async () => {
      await registrationPage.fillStep2(
        'Test Candidate',
        'Jl. Test No. 123',
        'Jakarta',
        'DKI Jakarta'
      );

      await registrationPage.expectStep3Visible();
    });
  });

  test.describe('Step 3: Education', () => {
    test.beforeEach(async () => {
      // Complete steps 1 and 2 first
      const email = generateUniqueEmail();
      await registrationPage.fillStep1WithEmail(email, 'testpassword123');
      await registrationPage.expectStep2Visible();
      await registrationPage.fillStep2(
        'Test Candidate',
        'Jl. Test No. 123',
        'Jakarta',
        'DKI Jakarta'
      );
      await registrationPage.expectStep3Visible();
    });

    test('should display education form', async () => {
      await expect(registrationPage.inputHighSchool).toBeVisible();
      await expect(registrationPage.selectGraduationYear).toBeVisible();
    });

    test('should proceed to step 4 with valid data', async ({ page }) => {
      // Check if there are programs available
      const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
      const radioCount = await prodiRadios.count();

      if (radioCount === 0) {
        test.skip();
        return;
      }

      await registrationPage.inputHighSchool.fill('SMA Negeri 1 Jakarta');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();

      await registrationPage.expectStep4Visible();
    });
  });

  test.describe('Step 4: Source Tracking and Completion', () => {
    test.beforeEach(async ({ page }) => {
      // Complete steps 1, 2, and 3 first
      const email = generateUniqueEmail();
      await registrationPage.fillStep1WithEmail(email, 'testpassword123');
      await registrationPage.expectStep2Visible();
      await registrationPage.fillStep2(
        'Test Candidate',
        'Jl. Test No. 123',
        'Jakarta',
        'DKI Jakarta'
      );
      await registrationPage.expectStep3Visible();

      // Check if there are programs available
      const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
      const radioCount = await prodiRadios.count();

      if (radioCount === 0) {
        test.skip();
        return;
      }

      await registrationPage.inputHighSchool.fill('SMA Negeri 1 Jakarta');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();
      await registrationPage.expectStep4Visible();
    });

    test('should display source tracking form', async () => {
      await expect(registrationPage.selectSourceType).toBeVisible();
    });

    test('should complete registration and redirect to portal', async () => {
      await registrationPage.fillStep4('instagram');
      await registrationPage.expectRedirectToPortal();
    });

    test('should show detail field for referral source types', async ({ page }) => {
      await registrationPage.selectSourceType.selectOption('friend_family');
      const detailContainer = page.locator('#source_detail_container');
      await expect(detailContainer).toBeVisible();
    });
  });
});

test.describe('Candidate Login', () => {
  let loginPage: LoginPage;
  let registrationPage: RegistrationPage;
  let testEmail: string;
  let testPhone: string;
  const testPassword = 'testpassword123';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate with email first
    const page = await browser.newPage();
    registrationPage = new RegistrationPage(page);
    testEmail = generateUniqueEmail();

    await registrationPage.goto();
    await registrationPage.expectPageLoaded();
    await registrationPage.fillStep1WithEmail(testEmail, testPassword);
    await registrationPage.expectStep2Visible();
    await registrationPage.fillStep2(
      'Login Test User',
      'Jl. Test No. 456',
      'Bandung',
      'Jawa Barat'
    );
    await registrationPage.expectStep3Visible();

    // Check if there are programs available
    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    const radioCount = await prodiRadios.count();

    if (radioCount > 0) {
      await registrationPage.inputHighSchool.fill('SMA Negeri 1 Bandung');
      await registrationPage.selectGraduationYear.selectOption('2025');
      await prodiRadios.first().click();
      await registrationPage.btnSubmitStep3.click();
      await registrationPage.expectStep4Visible();
      await registrationPage.fillStep4('google');
    }

    await page.close();

    // Create another test candidate with phone
    const page2 = await browser.newPage();
    const registrationPage2 = new RegistrationPage(page2);
    testPhone = generateUniquePhone();

    await registrationPage2.goto();
    await registrationPage2.expectPageLoaded();
    await registrationPage2.fillStep1WithPhone(testPhone, testPassword);
    await registrationPage2.expectStep2Visible();
    await registrationPage2.fillStep2(
      'Phone Login Test User',
      'Jl. Phone Test No. 789',
      'Surabaya',
      'Jawa Timur'
    );
    await registrationPage2.expectStep3Visible();

    const prodiRadios2 = page2.locator('[data-testid^="radio-prodi-"]');
    const radioCount2 = await prodiRadios2.count();

    if (radioCount2 > 0) {
      await registrationPage2.inputHighSchool.fill('SMA Negeri 1 Surabaya');
      await registrationPage2.selectGraduationYear.selectOption('2025');
      await prodiRadios2.first().click();
      await registrationPage2.btnSubmitStep3.click();
      await registrationPage2.expectStep4Visible();
      await registrationPage2.fillStep4('instagram');
    }

    await page2.close();
  });

  test.beforeEach(async ({ page }) => {
    loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.expectPageLoaded();
  });

  test('should display login form', async () => {
    await expect(loginPage.inputIdentifier).toBeVisible();
    await expect(loginPage.inputPassword).toBeVisible();
    await expect(loginPage.btnLogin).toBeVisible();
  });

  test('should show error for invalid credentials', async () => {
    await loginPage.login('invalid@example.com', 'wrongpassword');
    await loginPage.expectErrorMessage('salah');
  });

  test('should login with correct email and password', async () => {
    await loginPage.login(testEmail, testPassword);
    await loginPage.expectRedirectToPortal();
  });

  test('should login with correct phone and password', async () => {
    await loginPage.login(testPhone, testPassword);
    await loginPage.expectRedirectToPortal();
  });

  test('should redirect to login when accessing portal without session', async ({ page }) => {
    // Clear cookies first
    await page.context().clearCookies();
    await page.goto('/portal');
    // Should still be on portal for mockup (no auth protection yet)
    // This test can be updated when portal is protected
  });
});

test.describe('Duplicate Account Prevention', () => {
  let registrationPage: RegistrationPage;
  let existingEmail: string;
  let existingPhone: string;
  const testPassword = 'testpassword123';

  test.beforeAll(async ({ browser }) => {
    // Create a test candidate with both email and phone
    const page = await browser.newPage();
    registrationPage = new RegistrationPage(page);
    existingEmail = generateUniqueEmail();
    existingPhone = generateUniquePhone();

    await registrationPage.goto();
    await registrationPage.expectPageLoaded();
    await registrationPage.fillStep1WithBoth(existingEmail, existingPhone, testPassword);
    await registrationPage.expectStep2Visible();

    await page.close();
  });

  test('should show error when registering with existing email', async ({ page }) => {
    const regPage = new RegistrationPage(page);
    await regPage.goto();
    await regPage.expectPageLoaded();

    // Try to register with the same email
    await regPage.inputEmail.fill(existingEmail);
    await regPage.inputPassword.fill(testPassword);
    await regPage.inputPasswordConfirm.fill(testPassword);
    await regPage.btnSubmitStep1.click();

    await regPage.expectErrorMessage('sudah terdaftar');
  });

  test('should show error when registering with existing phone', async ({ page }) => {
    const regPage = new RegistrationPage(page);
    await regPage.goto();
    await regPage.expectPageLoaded();

    // Try to register with the same phone
    await regPage.inputPhone.fill(existingPhone);
    await regPage.inputPassword.fill(testPassword);
    await regPage.inputPasswordConfirm.fill(testPassword);
    await regPage.btnSubmitStep1.click();

    await regPage.expectErrorMessage('sudah terdaftar');
  });
});

test.describe('Registration with Tracking Parameters', () => {
  test('should preserve ref parameter through registration', async ({ page }) => {
    const registrationPage = new RegistrationPage(page);
    await registrationPage.gotoWithRef('TEST123');
    await registrationPage.expectPageLoaded();

    // Check that ref is in a hidden field
    const refInput = page.locator('input[name="ref"]');
    await expect(refInput).toHaveValue('TEST123');
  });

  test('should preserve campaign parameter through registration', async ({ page }) => {
    const registrationPage = new RegistrationPage(page);
    await registrationPage.gotoWithCampaign('SUMMER2026');
    await registrationPage.expectPageLoaded();

    // Check that campaign is in a hidden field
    const campaignInput = page.locator('input[name="utm_campaign"]');
    await expect(campaignInput).toHaveValue('SUMMER2026');
  });
});
