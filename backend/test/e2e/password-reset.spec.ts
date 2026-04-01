import { test, expect } from '@playwright/test';

test.describe('Password Reset', () => {
  test.describe('Forgot Password Page', () => {
    test('should display forgot password form', async ({ page }) => {
      await page.goto('/forgot-password');
      await expect(page.getByTestId('forgot-password-form')).toBeVisible();
      await expect(page.getByTestId('input-email')).toBeVisible();
      await expect(page.getByTestId('btn-send-reset')).toBeVisible();
    });

    test('should have link from login page', async ({ page }) => {
      await page.goto('/login');
      const forgotLink = page.getByTestId('link-forgot-password');
      await expect(forgotLink).toBeVisible();
      await forgotLink.click();
      await expect(page).toHaveURL('/forgot-password');
    });

    test('should redirect to reset password page after submitting email', async ({ page }) => {
      await page.goto('/forgot-password');
      await page.getByTestId('input-email').fill('nonexistent@example.com');
      await page.getByTestId('btn-send-reset').click();
      // Should redirect to reset-password page regardless of email existence (anti-enumeration)
      await expect(page).toHaveURL(/\/reset-password/);
    });

    test('should require email field', async ({ page }) => {
      await page.goto('/forgot-password');
      await page.getByTestId('btn-send-reset').click();
      // HTML5 validation should prevent submission — stays on same page
      await expect(page).toHaveURL('/forgot-password');
    });
  });

  test.describe('Reset Password Page', () => {
    test('should display reset password form with OTP and password fields', async ({ page }) => {
      await page.goto('/reset-password?email=test@example.com');
      await expect(page.getByTestId('reset-password-form')).toBeVisible();
      await expect(page.getByTestId('input-otp')).toBeVisible();
      await expect(page.getByTestId('input-password')).toBeVisible();
      await expect(page.getByTestId('input-password-confirm')).toBeVisible();
      await expect(page.getByTestId('btn-reset-password')).toBeVisible();
    });

    test('should show error for mismatched passwords', async ({ page }) => {
      await page.goto('/reset-password?email=test@example.com');
      await page.getByTestId('input-otp').fill('123456');
      await page.getByTestId('input-password').fill('newpassword123');
      await page.getByTestId('input-password-confirm').fill('differentpassword');
      await page.getByTestId('btn-reset-password').click();

      await expect(page.getByTestId('error-message')).toBeVisible();
      await expect(page.getByTestId('error-message')).toContainText('tidak cocok');
    });

    test('should show error for invalid OTP', async ({ page }) => {
      await page.goto('/reset-password?email=test@example.com');
      await page.getByTestId('input-otp').fill('000000');
      await page.getByTestId('input-password').fill('newpassword123');
      await page.getByTestId('input-password-confirm').fill('newpassword123');
      await page.getByTestId('btn-reset-password').click();

      await expect(page.getByTestId('error-message')).toBeVisible();
      await expect(page.getByTestId('error-message')).toContainText('tidak valid');
    });

    test('should have back to login link', async ({ page }) => {
      await page.goto('/reset-password?email=test@example.com');
      const loginLink = page.locator('a[href="/login"]');
      await expect(loginLink).toBeVisible();
    });
  });
});
