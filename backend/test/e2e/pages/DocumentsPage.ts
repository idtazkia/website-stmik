import { Page, Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class DocumentsPage extends BasePage {
  readonly path = '/portal/documents';

  // Main containers
  readonly documentsPage: Locator;
  readonly documentProgress: Locator;
  readonly documentsList: Locator;

  // Progress elements
  readonly progressBar: Locator;
  readonly progressText: Locator;

  constructor(page: Page) {
    super(page);
    this.documentsPage = page.getByTestId('documents-page');
    this.documentProgress = page.getByTestId('document-progress');
    this.documentsList = page.getByTestId('documents-list');

    this.progressBar = this.documentProgress.locator('.bg-green-600');
    this.progressText = this.documentProgress.locator('span.text-sm.text-gray-500');
  }

  async expectPageLoaded() {
    await expect(this.documentsPage).toBeVisible();
    await expect(this.documentProgress).toBeVisible();
    await expect(this.documentsList).toBeVisible();
  }

  async getProgressText(): Promise<string> {
    return await this.progressText.textContent() || '';
  }

  async getDocumentCount(): Promise<number> {
    return await this.documentsList.locator('> div').count();
  }

  getDocumentCard(documentType: string): Locator {
    return this.documentsList.locator(`form input[value="${documentType}"]`).locator('..').locator('..').locator('..').locator('..');
  }

  getUploadForm(documentType: string): Locator {
    return this.documentsList.locator(`form:has(input[value="${documentType}"])`);
  }

  getFileInput(documentType: string): Locator {
    return this.getUploadForm(documentType).locator('input[type="file"]');
  }

  getUploadButton(documentType: string): Locator {
    return this.getUploadForm(documentType).locator('button[type="submit"]');
  }

  async uploadDocument(documentType: string, filePath: string) {
    const fileInput = this.getFileInput(documentType);
    await fileInput.setInputFiles(filePath);
    const uploadButton = this.getUploadButton(documentType);
    await uploadButton.click();
    await this.page.waitForURL('/portal/documents');
  }

  async expectDocumentStatus(documentType: string, status: string) {
    const statusBadge = this.documentsList.locator(`form input[value="${documentType}"]`).locator('..').locator('..').locator('..').locator('..').locator('span.rounded-full');
    if (status === 'approved') {
      await expect(statusBadge).toContainText('Disetujui');
    } else if (status === 'pending') {
      await expect(statusBadge).toContainText('Menunggu Review');
    } else if (status === 'rejected') {
      await expect(statusBadge).toContainText('Ditolak');
    } else if (status === 'not_uploaded') {
      await expect(statusBadge).toContainText('Belum Upload');
    }
  }

  async expectDocumentHasUploadForm(documentType: string) {
    const form = this.getUploadForm(documentType);
    await expect(form).toBeVisible();
  }

  async expectDocumentNoUploadForm(documentType: string) {
    const form = this.getUploadForm(documentType);
    await expect(form).not.toBeVisible();
  }
}
