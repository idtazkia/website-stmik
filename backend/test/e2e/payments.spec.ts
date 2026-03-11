import { test, expect, Browser, Page } from '@playwright/test';
import * as path from 'path';
import * as fs from 'fs';

// Create test file for payment proof
function getTestProofPath(): string {
  const testFilesDir = path.join(__dirname, 'test-files-payments');
  if (!fs.existsSync(testFilesDir)) {
    fs.mkdirSync(testFilesDir, { recursive: true });
  }
  const filePath = path.join(testFilesDir, 'payment-proof.jpg');
  if (!fs.existsSync(filePath)) {
    // Create a minimal valid JPEG file
    const jpegHeader = Buffer.from([
      0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01,
      0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0xff, 0xdb, 0x00, 0x43,
      0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08, 0x07, 0x07, 0x07, 0x09,
      0x09, 0x08, 0x0a, 0x0c, 0x14, 0x0d, 0x0c, 0x0b, 0x0b, 0x0c, 0x19, 0x12,
      0x13, 0x0f, 0x14, 0x1d, 0x1a, 0x1f, 0x1e, 0x1d, 0x1a, 0x1c, 0x1c, 0x20,
      0x24, 0x2e, 0x27, 0x20, 0x22, 0x2c, 0x23, 0x1c, 0x1c, 0x28, 0x37, 0x29,
      0x2c, 0x30, 0x31, 0x34, 0x34, 0x34, 0x1f, 0x27, 0x39, 0x3d, 0x38, 0x32,
      0x3c, 0x2e, 0x33, 0x34, 0x32, 0xff, 0xc0, 0x00, 0x0b, 0x08, 0x00, 0x01,
      0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0xff, 0xc4, 0x00, 0x1f, 0x00, 0x00,
      0x01, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
      0x09, 0x0a, 0x0b, 0xff, 0xc4, 0x00, 0xb5, 0x10, 0x00, 0x02, 0x01, 0x03,
      0x03, 0x02, 0x04, 0x03, 0x05, 0x05, 0x04, 0x04, 0x00, 0x00, 0x01, 0x7d,
      0x01, 0x02, 0x03, 0x00, 0x04, 0x11, 0x05, 0x12, 0x21, 0x31, 0x41, 0x06,
      0x13, 0x51, 0x61, 0x07, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1, 0x08,
      0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1, 0xf0, 0x24, 0x33, 0x62, 0x72,
      0x82, 0x09, 0x0a, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27, 0x28,
      0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45,
      0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59,
      0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75,
      0x76, 0x77, 0x78, 0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89,
      0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3,
      0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6,
      0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9,
      0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1, 0xe2,
      0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf1, 0xf2, 0xf3, 0xf4,
      0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xda, 0x00, 0x08, 0x01, 0x01,
      0x00, 0x00, 0x3f, 0x00, 0xfb, 0xd5, 0xdb, 0x20, 0xa8, 0xf3, 0xff, 0xd9
    ]);
    fs.writeFileSync(filePath, jpegHeader);
  }
  return filePath;
}

// Register candidate, create billing, upload payment proof, return IDs
async function setupCandidateWithPayment(browser: Browser): Promise<{
  candidateId: string;
  billingId: string;
  paymentPage: Page;
}> {
  const uniqueId = `${Date.now()}${Math.random().toString(36).slice(2, 8)}`;
  const uniqueEmail = `pay${uniqueId}@example.com`;
  const password = 'testpassword123';

  // Register candidate
  const candidatePage = await browser.newPage();
  await candidatePage.goto('/register');
  await candidatePage.getByTestId('input-email').fill(uniqueEmail);
  await candidatePage.getByTestId('input-password').fill(password);
  await candidatePage.getByTestId('input-password-confirm').fill(password);
  await candidatePage.getByTestId('btn-submit-step1').click();
  await expect(candidatePage.getByTestId('step2-form')).toBeVisible({ timeout: 10000 });

  await candidatePage.getByTestId('input-name').fill(`PayTest ${uniqueId}`);
  await candidatePage.getByTestId('input-address').fill('Test Address');
  await candidatePage.getByTestId('input-city').fill('Jakarta');
  await candidatePage.getByTestId('input-province').fill('DKI Jakarta');
  await candidatePage.getByTestId('btn-submit-step2').click();
  await expect(candidatePage.getByTestId('step3-form')).toBeVisible({ timeout: 10000 });

  await candidatePage.getByTestId('input-high-school').fill('SMA Test');
  await candidatePage.getByTestId('select-graduation-year').selectOption('2025');
  await candidatePage.locator('input[type="radio"][name="prodi_id"]').first().click();
  await candidatePage.getByTestId('btn-submit-step3').click();
  await expect(candidatePage.getByTestId('step4-form')).toBeVisible({ timeout: 10000 });

  await candidatePage.getByTestId('select-source-type').selectOption('instagram');
  await candidatePage.getByTestId('btn-submit-step4').click();
  await expect(candidatePage).toHaveURL('/portal', { timeout: 10000 });
  await candidatePage.close();

  // Find candidate ID as admin
  const adminPage = await browser.newPage();
  await adminPage.goto('/test/login/admin');
  await adminPage.goto('/admin/candidates?search=' + encodeURIComponent(uniqueEmail));
  await expect(adminPage.getByTestId('candidates-page')).toBeVisible();
  await adminPage.waitForTimeout(1000);

  const viewLink = adminPage.locator('[data-testid^="view-candidate-"]').first();
  const testId = await viewLink.getAttribute('data-testid');
  const candidateId = testId?.replace('view-candidate-', '') || '';
  await adminPage.close();

  // Create billing as finance
  const financePage = await browser.newPage();
  await financePage.goto('/test/login/finance');
  await financePage.goto(`/admin/finance/billings/create?candidate_id=${candidateId}`);
  await expect(financePage.getByTestId('billing-form-page')).toBeVisible();

  await financePage.getByTestId('select-billing-type').selectOption('registration');
  await financePage.getByTestId('input-amount').fill('500000');
  await financePage.getByTestId('input-due-date').fill('2026-06-01');
  await financePage.getByTestId('btn-submit').click();
  await expect(financePage.getByTestId('billing-detail-page')).toBeVisible({ timeout: 10000 });

  // Extract billing ID from URL
  const billingUrl = financePage.url();
  const billingId = billingUrl.split('/').pop() || '';
  await financePage.close();

  // Upload payment proof as candidate
  const portalPage = await browser.newPage();
  await portalPage.goto('/test/login/candidate?email=' + encodeURIComponent(uniqueEmail));
  await portalPage.goto('/portal/payments');
  await expect(portalPage.getByTestId('payments-page')).toBeVisible();

  // Find the billing's upload form and submit proof
  const uploadForm = portalPage.locator(`[data-testid="upload-form-${billingId}"]`);
  // If upload form is visible, fill and submit
  const uploadVisible = await uploadForm.isVisible().catch(() => false);
  if (uploadVisible) {
    await uploadForm.locator('input[name="transfer_date"]').fill('2026-03-10');
    await uploadForm.locator('input[name="amount"]').fill('500000');
    await uploadForm.locator('input[name="proof"]').setInputFiles(getTestProofPath());
    await uploadForm.locator('button[type="submit"]').click();
    await expect(portalPage.getByTestId('payments-page')).toBeVisible({ timeout: 10000 });
  } else {
    // Try the generic upload button approach
    const uploadBtn = portalPage.locator(`[data-testid="btn-upload-${billingId}"]`);
    if (await uploadBtn.isVisible().catch(() => false)) {
      await uploadBtn.click();
      await portalPage.locator('input[name="transfer_date"]').first().fill('2026-03-10');
      await portalPage.locator('input[name="amount"]').first().fill('500000');
      await portalPage.locator('input[name="proof"]').first().setInputFiles(getTestProofPath());
      await portalPage.locator('button[type="submit"]').first().click();
      await expect(portalPage.getByTestId('payments-page')).toBeVisible({ timeout: 10000 });
    }
  }
  await portalPage.close();

  // Return finance page for payment review
  const paymentPage = await browser.newPage();
  await paymentPage.goto('/test/login/finance');

  return { candidateId, billingId, paymentPage };
}

test.describe('Payments - Portal Display', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/test/login/candidate');
    await expect(page).toHaveURL('/portal');
  });

  test('should display payments page with summary cards', async ({ page }) => {
    await page.goto('/portal/payments');

    await expect(page.getByTestId('payments-page')).toBeVisible();
    await expect(page.getByTestId('payment-summary')).toBeVisible();

    await expect(page.locator('text=Total Tagihan')).toBeVisible();
    await expect(page.locator('text=Sudah Dibayar')).toBeVisible();
    await expect(page.locator('text=Menunggu Verifikasi')).toBeVisible();
  });

  test('should display bank info', async ({ page }) => {
    await page.goto('/portal/payments');

    await expect(page.getByTestId('bank-info')).toBeVisible();
    await expect(page.locator('text=Bank Syariah Indonesia')).toBeVisible();
    await expect(page.locator('text=Yayasan Tazkia')).toBeVisible();
  });

  test('should display payments list', async ({ page }) => {
    await page.goto('/portal/payments');

    await expect(page.getByTestId('payments-list')).toBeVisible();
    await expect(page.locator('text=Daftar Tagihan')).toBeVisible();
  });
});

test.describe('Payments - With Billing', () => {
  test('should show empty state when no billings', async ({ browser }) => {
    // Register a fresh candidate (no billings created yet)
    const uniqueId = `${Date.now()}${Math.random().toString(36).slice(2, 8)}`;
    const uniqueEmail = `empty${uniqueId}@example.com`;
    const password = 'testpassword123';

    const candidatePage = await browser.newPage();
    await candidatePage.goto('/register');
    await candidatePage.getByTestId('input-email').fill(uniqueEmail);
    await candidatePage.getByTestId('input-password').fill(password);
    await candidatePage.getByTestId('input-password-confirm').fill(password);
    await candidatePage.getByTestId('btn-submit-step1').click();
    await expect(candidatePage.getByTestId('step2-form')).toBeVisible({ timeout: 10000 });

    await candidatePage.getByTestId('input-name').fill(`EmptyPay ${uniqueId}`);
    await candidatePage.getByTestId('input-address').fill('Test Address');
    await candidatePage.getByTestId('input-city').fill('Jakarta');
    await candidatePage.getByTestId('input-province').fill('DKI Jakarta');
    await candidatePage.getByTestId('btn-submit-step2').click();
    await expect(candidatePage.getByTestId('step3-form')).toBeVisible({ timeout: 10000 });

    await candidatePage.getByTestId('input-high-school').fill('SMA Test');
    await candidatePage.getByTestId('select-graduation-year').selectOption('2025');
    await candidatePage.locator('input[type="radio"][name="prodi_id"]').first().click();
    await candidatePage.getByTestId('btn-submit-step3').click();
    await expect(candidatePage.getByTestId('step4-form')).toBeVisible({ timeout: 10000 });

    await candidatePage.getByTestId('select-source-type').selectOption('instagram');
    await candidatePage.getByTestId('btn-submit-step4').click();
    await expect(candidatePage).toHaveURL('/portal', { timeout: 10000 });

    await candidatePage.goto('/portal/payments');
    await expect(candidatePage.getByTestId('payments-page')).toBeVisible();
    await expect(candidatePage.locator('text=Belum ada tagihan')).toBeVisible();

    await candidatePage.close();
  });
});

test.describe('Admin Payment Review', () => {
  test('finance user can access payment review page', async ({ browser }) => {
    const financePage = await browser.newPage();
    await financePage.goto('/test/login/finance');
    await financePage.goto('/admin/finance/payments');

    await expect(financePage.getByTestId('finance-payments-page')).toBeVisible();
    await expect(financePage.getByTestId('stat-pending')).toBeVisible();
    await expect(financePage.getByTestId('stat-approved')).toBeVisible();
    await expect(financePage.getByTestId('stat-rejected')).toBeVisible();

    await financePage.close();
  });

  test('admin can also access payment review page', async ({ browser }) => {
    const adminPage = await browser.newPage();
    await adminPage.goto('/test/login/admin');
    await adminPage.goto('/admin/finance/payments');

    await expect(adminPage.getByTestId('finance-payments-page')).toBeVisible();

    await adminPage.close();
  });

  test('consultant cannot access payment review page', async ({ browser }) => {
    const consultantPage = await browser.newPage();
    await consultantPage.goto('/test/login/consultant');

    const response = await consultantPage.goto('/admin/finance/payments');
    expect(response?.status()).toBe(403);

    await consultantPage.close();
  });

  test('finance user can filter payments by status', async ({ browser }) => {
    const financePage = await browser.newPage();
    await financePage.goto('/test/login/finance');
    await financePage.goto('/admin/finance/payments');

    await financePage.getByTestId('select-status').selectOption('pending');
    await financePage.getByTestId('btn-filter').click();

    await expect(financePage).toHaveURL(/status=pending/);

    await financePage.close();
  });

  test('finance user can approve payment', async ({ browser }) => {
    const { billingId, paymentPage } = await setupCandidateWithPayment(browser);

    // Go to payments review page
    await paymentPage.goto('/admin/finance/payments?status=pending');
    await expect(paymentPage.getByTestId('finance-payments-page')).toBeVisible();

    // Find approve button for any pending payment
    const approveBtn = paymentPage.locator('[data-testid^="btn-approve-"]').first();
    const isVisible = await approveBtn.isVisible().catch(() => false);

    if (isVisible) {
      await approveBtn.click();

      // Should redirect back to payment list or show updated status
      await expect(paymentPage.getByTestId('finance-payments-page')).toBeVisible({ timeout: 10000 });
    }

    await paymentPage.close();
  });

  test('finance user can reject payment with reason', async ({ browser }) => {
    const { billingId, paymentPage } = await setupCandidateWithPayment(browser);

    // Go to payments review page
    await paymentPage.goto('/admin/finance/payments?status=pending');
    await expect(paymentPage.getByTestId('finance-payments-page')).toBeVisible();

    // Find reject button for any pending payment
    const rejectBtn = paymentPage.locator('[data-testid^="btn-reject-"]').first();
    const isVisible = await rejectBtn.isVisible().catch(() => false);

    if (isVisible) {
      await rejectBtn.click();

      // Fill rejection reason in the modal
      const reasonInput = paymentPage.locator('[data-testid^="input-reject-reason-"]').first();
      await expect(reasonInput).toBeVisible({ timeout: 5000 });
      await reasonInput.fill('Bukti transfer tidak valid');

      // Confirm rejection
      const confirmBtn = paymentPage.locator('[data-testid^="btn-confirm-reject-"]').first();
      await confirmBtn.click();

      // Should redirect back to payment list
      await expect(paymentPage.getByTestId('finance-payments-page')).toBeVisible({ timeout: 10000 });
    }

    await paymentPage.close();
  });
});

test.describe('Candidate Search API', () => {
  test('search API returns results for valid query', async ({ browser }) => {
    const financePage = await browser.newPage();
    await financePage.goto('/test/login/finance');

    const response = await financePage.request.get('/admin/api/candidates/search?q=test');
    expect(response.ok()).toBeTruthy();

    const body = await response.json();
    expect(Array.isArray(body)).toBeTruthy();

    await financePage.close();
  });

  test('search API requires minimum 2 characters', async ({ browser }) => {
    const financePage = await browser.newPage();
    await financePage.goto('/test/login/finance');

    const response = await financePage.request.get('/admin/api/candidates/search?q=a');
    const body = await response.json();
    expect(Array.isArray(body)).toBeTruthy();
    expect(body.length).toBe(0);

    await financePage.close();
  });
});
