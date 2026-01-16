import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsDocumentTypesPage extends BasePage {
  readonly path = '/admin/settings/document-types';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    await this.page.waitForURL(/\/admin\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-documents-page');
  }

  get documentTypesSection(): Locator {
    return this.page.getByTestId('document-types-section');
  }

  get documentTypesTable(): Locator {
    return this.page.getByTestId('document-types-table');
  }

  get documentTypesList(): Locator {
    return this.page.getByTestId('document-types-list');
  }

  get addDocumentTypeButton(): Locator {
    return this.page.getByTestId('add-document-type-button');
  }

  get addDocumentTypeModal(): Locator {
    return this.page.getByTestId('add-document-type-modal');
  }

  // Form inputs for add modal
  get inputName(): Locator {
    return this.page.getByTestId('input-doctype-name');
  }

  get inputCode(): Locator {
    return this.page.getByTestId('input-doctype-code');
  }

  get inputDescription(): Locator {
    return this.page.getByTestId('input-doctype-description');
  }

  get inputMaxSize(): Locator {
    return this.page.getByTestId('input-doctype-maxsize');
  }

  get inputDisplayOrder(): Locator {
    return this.page.getByTestId('input-doctype-order');
  }

  get inputIsRequired(): Locator {
    return this.page.getByTestId('input-doctype-required');
  }

  get inputCanDefer(): Locator {
    return this.page.getByTestId('input-doctype-defer');
  }

  get submitAddButton(): Locator {
    return this.page.getByTestId('submit-add-doctype');
  }

  // Document type row elements by ID
  getDocumentTypeRow(id: string): Locator {
    return this.page.getByTestId(`doctype-row-${id}`);
  }

  getDocumentTypeName(id: string): Locator {
    return this.page.getByTestId(`doctype-name-${id}`);
  }

  getDocumentTypeCode(id: string): Locator {
    return this.page.getByTestId(`doctype-code-${id}`);
  }

  getDocumentTypeRequired(id: string): Locator {
    return this.page.getByTestId(`doctype-required-${id}`);
  }

  getDocumentTypeDefer(id: string): Locator {
    return this.page.getByTestId(`doctype-defer-${id}`);
  }

  getDocumentTypeMaxSize(id: string): Locator {
    return this.page.getByTestId(`doctype-maxsize-${id}`);
  }

  getDocumentTypeStatusToggle(id: string): Locator {
    return this.page.getByTestId(`doctype-status-toggle-${id}`);
  }

  getDocumentTypeEditButton(id: string): Locator {
    return this.page.getByTestId(`doctype-edit-${id}`);
  }

  getEditModal(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-modal-${id}`);
  }

  // Edit form inputs by ID
  getEditInputName(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-name-${id}`);
  }

  getEditInputCode(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-code-${id}`);
  }

  getEditInputDescription(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-description-${id}`);
  }

  getEditInputMaxSize(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-maxsize-${id}`);
  }

  getEditInputDisplayOrder(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-order-${id}`);
  }

  getEditInputRequired(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-required-${id}`);
  }

  getEditInputDefer(id: string): Locator {
    return this.page.getByTestId(`edit-doctype-defer-${id}`);
  }

  getEditSubmitButton(id: string): Locator {
    return this.page.getByTestId(`submit-edit-doctype-${id}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.documentTypesSection).toBeVisible();
  }

  // Get all document type IDs
  async getAllDocumentTypeIds(): Promise<string[]> {
    const rows = this.documentTypesList.locator('tr[data-testid^="doctype-row-"]');
    const count = await rows.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await rows.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('doctype-row-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectDocumentTypeDisplayed(id: string): Promise<void> {
    await expect(this.getDocumentTypeRow(id)).toBeVisible();
    await expect(this.getDocumentTypeName(id)).toBeVisible();
  }

  async expectDocumentTypeActive(id: string): Promise<void> {
    await expect(this.getDocumentTypeStatusToggle(id)).toContainText('Aktif');
  }

  async expectDocumentTypeInactive(id: string): Promise<void> {
    await expect(this.getDocumentTypeStatusToggle(id)).toContainText('Nonaktif');
  }

  // Actions
  async openAddModal(): Promise<void> {
    await this.addDocumentTypeButton.click();
    await expect(this.addDocumentTypeModal).toBeVisible();
  }

  async addDocumentType(data: {
    name: string;
    code: string;
    description?: string;
    maxFileSizeMB?: number;
    displayOrder?: number;
    isRequired?: boolean;
    canDefer?: boolean;
  }): Promise<void> {
    await this.openAddModal();
    await this.inputName.fill(data.name);
    await this.inputCode.fill(data.code);
    if (data.description) {
      await this.inputDescription.fill(data.description);
    }
    if (data.maxFileSizeMB !== undefined) {
      await this.inputMaxSize.fill(String(data.maxFileSizeMB));
    }
    if (data.displayOrder !== undefined) {
      await this.inputDisplayOrder.fill(String(data.displayOrder));
    }
    // Handle checkboxes
    if (data.isRequired === false) {
      await this.inputIsRequired.uncheck();
    } else if (data.isRequired === true) {
      await this.inputIsRequired.check();
    }
    if (data.canDefer === true) {
      await this.inputCanDefer.check();
    } else if (data.canDefer === false) {
      await this.inputCanDefer.uncheck();
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/document-types') &&
      response.request().method() === 'POST' &&
      response.status() === 200
    );
    await this.submitAddButton.click();
    await responsePromise;
  }

  async openEditModal(id: string): Promise<void> {
    await this.getDocumentTypeEditButton(id).click();
    await expect(this.getEditModal(id)).toBeVisible();
  }

  async editDocumentType(id: string, data: {
    name?: string;
    code?: string;
    description?: string;
    maxFileSizeMB?: number;
    displayOrder?: number;
    isRequired?: boolean;
    canDefer?: boolean;
  }): Promise<void> {
    await this.openEditModal(id);
    if (data.name) {
      await this.getEditInputName(id).fill(data.name);
    }
    if (data.code) {
      await this.getEditInputCode(id).fill(data.code);
    }
    if (data.description !== undefined) {
      await this.getEditInputDescription(id).fill(data.description);
    }
    if (data.maxFileSizeMB !== undefined) {
      await this.getEditInputMaxSize(id).fill(String(data.maxFileSizeMB));
    }
    if (data.displayOrder !== undefined) {
      await this.getEditInputDisplayOrder(id).fill(String(data.displayOrder));
    }
    if (data.isRequired === false) {
      await this.getEditInputRequired(id).uncheck();
    } else if (data.isRequired === true) {
      await this.getEditInputRequired(id).check();
    }
    if (data.canDefer === true) {
      await this.getEditInputDefer(id).check();
    } else if (data.canDefer === false) {
      await this.getEditInputDefer(id).uncheck();
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/document-types/${id}`) &&
      !response.url().includes('toggle-active') &&
      response.request().method() === 'POST' &&
      response.status() === 200
    );
    await this.getEditSubmitButton(id).click();
    await responsePromise;
  }

  async toggleDocumentTypeStatus(id: string): Promise<void> {
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/document-types/${id}/toggle-active`) &&
      response.status() === 200
    );
    await this.getDocumentTypeStatusToggle(id).click();
    await responsePromise;
  }
}
