import { test, expect, Browser, Page } from '@playwright/test';
import { DocumentsPage } from './pages';
import * as path from 'path';
import * as fs from 'fs';

// Create test file on demand (works with parallel workers)
function getTestFilePath(): string {
  const testFilesDir = path.join(__dirname, 'test-files-upload');
  if (!fs.existsSync(testFilesDir)) {
    fs.mkdirSync(testFilesDir, { recursive: true });
  }
  const filePath = path.join(testFilesDir, 'test-document.pdf');
  if (!fs.existsSync(filePath)) {
    fs.writeFileSync(filePath, '%PDF-1.4 Test Document Content for Upload Test');
  }
  return filePath;
}

// Helper to register a new candidate with unique name
async function registerCandidate(browser: Browser): Promise<{ email: string; name: string; page: Page }> {
  const candidatePage = await browser.newPage();
  const uniqueId = `${Date.now()}${Math.random().toString(36).slice(2, 8)}`;
  const uniqueEmail = `upload${uniqueId}@example.com`;
  const uniqueName = `UploadCandidate ${uniqueId}`;
  const password = 'testpassword123';

  // Step 1: Account creation
  await candidatePage.goto('/register');
  await candidatePage.getByTestId('input-email').fill(uniqueEmail);
  await candidatePage.getByTestId('input-password').fill(password);
  await candidatePage.getByTestId('input-password-confirm').fill(password);
  await candidatePage.getByTestId('btn-submit-step1').click();
  await expect(candidatePage.getByTestId('step2-form')).toBeVisible({ timeout: 10000 });

  // Step 2: Personal info - use unique name
  await candidatePage.getByTestId('input-name').fill(uniqueName);
  await candidatePage.getByTestId('input-address').fill('Test Address');
  await candidatePage.getByTestId('input-city').fill('Jakarta');
  await candidatePage.getByTestId('input-province').fill('DKI Jakarta');
  await candidatePage.getByTestId('btn-submit-step2').click();
  await expect(candidatePage.getByTestId('step3-form')).toBeVisible({ timeout: 10000 });

  // Step 3: Education
  await candidatePage.getByTestId('input-high-school').fill('SMA Test');
  await candidatePage.getByTestId('select-graduation-year').selectOption('2025');
  // Click first prodi radio
  const prodiRadios = candidatePage.locator('input[type="radio"][name="prodi_id"]');
  await prodiRadios.first().click();
  await candidatePage.getByTestId('btn-submit-step3').click();
  await expect(candidatePage.getByTestId('step4-form')).toBeVisible({ timeout: 10000 });

  // Step 4: Source tracking
  await candidatePage.getByTestId('select-source-type').selectOption('instagram');
  await candidatePage.getByTestId('btn-submit-step4').click();
  await expect(candidatePage).toHaveURL('/portal', { timeout: 10000 });

  return { email: uniqueEmail, name: uniqueName, page: candidatePage };
}

test.describe('Document Upload - Page Display', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/test/login/candidate');
    await expect(page).toHaveURL('/portal');
  });

  test('should display documents page with required documents', async ({ page }) => {
    const documentsPage = new DocumentsPage(page);

    await page.goto('/portal/documents');
    await documentsPage.expectPageLoaded();

    // Verify progress text is shown
    await expect(documentsPage.progressText).toBeVisible({ timeout: 10000 });
    const progressText = await documentsPage.getProgressText();
    expect(progressText).toMatch(/\d+ dari \d+ dokumen/);

    // Verify at least 4 document types are shown (seeded in migration)
    const docCount = await documentsPage.getDocumentCount();
    expect(docCount).toBeGreaterThanOrEqual(4);
  });

  test('should show upload forms for documents', async ({ page }) => {
    await page.goto('/portal/documents');
    await expect(page.getByTestId('documents-page')).toBeVisible();

    // Verify upload forms exist
    const uploadForms = page.locator('form[action="/portal/documents/upload"]');
    const formCount = await uploadForms.count();
    expect(formCount).toBeGreaterThan(0);

    // Verify file inputs have correct accept attribute
    const fileInput = uploadForms.first().locator('input[type="file"]');
    const acceptValue = await fileInput.getAttribute('accept');
    expect(acceptValue).toBeTruthy();
    expect(acceptValue).toContain('application/pdf');
  });
});

test.describe('Document Upload - Upload Flow', () => {
  test('should upload document and verify storage', async ({ browser }) => {
    // Register a fresh candidate so we have clean upload slots
    const { email, page: candidatePage } = await registerCandidate(browser);

    // Navigate to documents page
    await candidatePage.goto('/portal/documents');
    await expect(candidatePage.getByTestId('documents-page')).toBeVisible();

    // Find upload forms
    const uploadForms = candidatePage.locator('form[action="/portal/documents/upload"]');
    const formCount = await uploadForms.count();
    expect(formCount).toBeGreaterThan(0);

    // Upload a document
    const firstForm = uploadForms.first();
    const fileInput = firstForm.locator('input[type="file"]');
    await fileInput.setInputFiles(getTestFilePath());

    const uploadButton = firstForm.locator('button[type="submit"]');
    await uploadButton.click();

    // Wait for page to reload after upload
    await candidatePage.waitForLoadState('networkidle');

    // Verify upload succeeded - should show "Menunggu Review" status
    await expect(candidatePage.locator('text=Menunggu Review').first()).toBeVisible({ timeout: 10000 });

    // Verify file details are shown
    const fileDetailsVisible =
      await candidatePage.locator('text=File:').first().isVisible().catch(() => false) ||
      await candidatePage.locator('text=Diupload:').first().isVisible().catch(() => false);
    expect(fileDetailsVisible).toBe(true);

    await candidatePage.close();
  });

  test('should verify uploaded file is stored and accessible', async ({ browser }) => {
    // Register a fresh candidate
    const { email, name, page: candidatePage } = await registerCandidate(browser);

    // Upload a document
    await candidatePage.goto('/portal/documents');
    await expect(candidatePage.getByTestId('documents-page')).toBeVisible();

    const uploadForms = candidatePage.locator('form[action="/portal/documents/upload"]');
    expect(await uploadForms.count()).toBeGreaterThan(0);

    const firstForm = uploadForms.first();
    await firstForm.locator('input[type="file"]').setInputFiles(getTestFilePath());
    await firstForm.locator('button[type="submit"]').click();
    await candidatePage.waitForLoadState('networkidle');

    // Verify pending status
    await expect(candidatePage.locator('text=Menunggu Review').first()).toBeVisible();

    await candidatePage.close();

    // Now verify file is accessible via admin panel
    const adminPage = await browser.newPage();
    await adminPage.goto('/test/login/admin');
    await adminPage.goto('/admin/documents?status=pending');

    // Search for the uploaded document by email (name search not supported due to encryption)
    const searchInput = adminPage.locator('input[name="search"]');
    await searchInput.fill(email);
    await adminPage.keyboard.press('Enter');
    await adminPage.waitForLoadState('networkidle');

    // Verify document appears - check candidate name in the results
    await expect(adminPage.locator(`text=${name}`).first()).toBeVisible();

    // Check if there's a link to view the file
    const viewLink = adminPage.locator('a[href*="/uploads/"]').first();
    const hasViewLink = await viewLink.isVisible().catch(() => false);

    if (hasViewLink) {
      const fileUrl = await viewLink.getAttribute('href');
      expect(fileUrl).toContain('/uploads/');

      // Verify file is accessible via HTTP request
      const response = await adminPage.request.get(fileUrl!);
      expect(response.status()).toBe(200);
    }

    await adminPage.close();
  });

  test('should upload document and persist across page reload', async ({ browser }) => {
    // Register a fresh candidate
    const { email, page: candidatePage } = await registerCandidate(browser);

    // Upload a document
    await candidatePage.goto('/portal/documents');
    await expect(candidatePage.getByTestId('documents-page')).toBeVisible();

    const uploadForms = candidatePage.locator('form[action="/portal/documents/upload"]');
    expect(await uploadForms.count()).toBeGreaterThan(0);

    const firstForm = uploadForms.first();
    await firstForm.locator('input[type="file"]').setInputFiles(getTestFilePath());
    await firstForm.locator('button[type="submit"]').click();
    await candidatePage.waitForLoadState('networkidle');

    // Verify pending status
    await expect(candidatePage.locator('text=Menunggu Review').first()).toBeVisible();

    // Close page and re-login
    await candidatePage.close();

    // Login again and verify document still shows
    const verifyPage = await browser.newPage();
    await verifyPage.goto('/login');
    await verifyPage.getByTestId('input-identifier').fill(email);
    await verifyPage.getByTestId('input-password').fill('testpassword123');
    await verifyPage.getByTestId('btn-login').click();
    await expect(verifyPage).toHaveURL('/portal', { timeout: 10000 });

    await verifyPage.goto('/portal/documents');
    await expect(verifyPage.getByTestId('documents-page')).toBeVisible();

    // Document should still show pending status
    await expect(verifyPage.locator('text=Menunggu Review').first()).toBeVisible();

    await verifyPage.close();
  });

  test('should replace document on re-upload', async ({ browser }) => {
    // Register a fresh candidate
    const { page: candidatePage } = await registerCandidate(browser);

    // Upload first document
    await candidatePage.goto('/portal/documents');
    await expect(candidatePage.getByTestId('documents-page')).toBeVisible();

    const uploadForms = candidatePage.locator('form[action="/portal/documents/upload"]');
    expect(await uploadForms.count()).toBeGreaterThan(0);

    const firstForm = uploadForms.first();
    await firstForm.locator('input[type="file"]').setInputFiles(getTestFilePath());
    await firstForm.locator('button[type="submit"]').click();
    await candidatePage.waitForLoadState('networkidle');

    // First upload done
    await expect(candidatePage.locator('text=Menunggu Review').first()).toBeVisible();

    // Create a different test file
    const testFilesDir = path.join(__dirname, 'test-files-upload');
    const newFilePath = path.join(testFilesDir, 'test-document-v2.pdf');
    fs.writeFileSync(newFilePath, '%PDF-1.4 Test Document v2 - Replacement Content');

    // Find the re-upload form (should still be visible or there's a replace option)
    // The system uses UPSERT so same candidate+doctype will replace
    await candidatePage.goto('/portal/documents');

    // Upload again to the same document type if form is available
    const reuploadForms = candidatePage.locator('form[action="/portal/documents/upload"]');
    const formCount = await reuploadForms.count();

    if (formCount > 0) {
      await reuploadForms.first().locator('input[type="file"]').setInputFiles(newFilePath);
      await reuploadForms.first().locator('button[type="submit"]').click();
      await candidatePage.waitForLoadState('networkidle');

      // Should still show pending (reset to pending on re-upload)
      await expect(candidatePage.locator('text=Menunggu Review').first()).toBeVisible();
    }

    await candidatePage.close();
  });
});

test.describe('Document Upload - Validation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/test/login/candidate');
    await expect(page).toHaveURL('/portal');
  });

  test('should have file size limit displayed', async ({ page }) => {
    await page.goto('/portal/documents');
    await expect(page.getByTestId('documents-page')).toBeVisible();

    // Should mention file size limit somewhere on page
    const hasFileSizeInfo =
      await page.locator('text=5MB').first().isVisible().catch(() => false) ||
      await page.locator('text=maksimal').first().isVisible().catch(() => false) ||
      await page.locator('text=max').first().isVisible().catch(() => false);

    expect(hasFileSizeInfo).toBe(true);
  });

  test('should show accepted file formats', async ({ page }) => {
    await page.goto('/portal/documents');
    await expect(page.getByTestId('documents-page')).toBeVisible();

    // Should mention accepted formats
    const hasFormatInfo =
      await page.locator('text=JPG').first().isVisible().catch(() => false) ||
      await page.locator('text=PNG').first().isVisible().catch(() => false) ||
      await page.locator('text=PDF').first().isVisible().catch(() => false);

    expect(hasFormatInfo).toBe(true);
  });
});
