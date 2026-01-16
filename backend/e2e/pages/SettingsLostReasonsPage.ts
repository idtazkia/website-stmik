import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsLostReasonsPage extends BasePage {
  readonly path = '/admin/settings/lost-reasons';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    await this.page.waitForURL(/\/admin\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-lost-reasons-page');
  }

  get lostReasonsSection(): Locator {
    return this.page.getByTestId('lost-reasons-section');
  }

  get lostReasonsList(): Locator {
    return this.page.getByTestId('lost-reasons-list');
  }

  get addLostReasonButton(): Locator {
    return this.page.getByTestId('add-lost-reason-button');
  }

  get addLostReasonModal(): Locator {
    return this.page.getByTestId('add-lost-reason-modal');
  }

  // Form inputs for add modal
  get inputName(): Locator {
    return this.page.getByTestId('input-lost-reason-name');
  }

  get inputDescription(): Locator {
    return this.page.getByTestId('input-lost-reason-description');
  }

  get inputDisplayOrder(): Locator {
    return this.page.getByTestId('input-lost-reason-order');
  }

  get submitAddButton(): Locator {
    return this.page.getByTestId('submit-add-lost-reason');
  }

  // Lost reason item elements by ID
  getLostReasonItem(id: string): Locator {
    return this.page.getByTestId(`lost-reason-item-${id}`);
  }

  getLostReasonName(id: string): Locator {
    return this.page.getByTestId(`lost-reason-name-${id}`);
  }

  getLostReasonDescription(id: string): Locator {
    return this.page.getByTestId(`lost-reason-description-${id}`);
  }

  getLostReasonStatusToggle(id: string): Locator {
    return this.page.getByTestId(`lost-reason-status-toggle-${id}`);
  }

  getLostReasonEditButton(id: string): Locator {
    return this.page.getByTestId(`lost-reason-edit-${id}`);
  }

  getEditModal(id: string): Locator {
    return this.page.getByTestId(`edit-lost-reason-modal-${id}`);
  }

  // Edit form inputs by ID
  getEditInputName(id: string): Locator {
    return this.page.getByTestId(`edit-lost-reason-name-${id}`);
  }

  getEditInputDescription(id: string): Locator {
    return this.page.getByTestId(`edit-lost-reason-description-${id}`);
  }

  getEditInputOrder(id: string): Locator {
    return this.page.getByTestId(`edit-lost-reason-order-${id}`);
  }

  getEditSubmitButton(id: string): Locator {
    return this.page.getByTestId(`submit-edit-lost-reason-${id}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.lostReasonsSection).toBeVisible();
  }

  // Get all lost reason IDs
  async getAllLostReasonIds(): Promise<string[]> {
    const items = this.lostReasonsList.locator('div[data-testid^="lost-reason-item-"]');
    const count = await items.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await items.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('lost-reason-item-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectLostReasonDisplayed(id: string): Promise<void> {
    await expect(this.getLostReasonItem(id)).toBeVisible();
    await expect(this.getLostReasonName(id)).toBeVisible();
  }

  async expectLostReasonActive(id: string): Promise<void> {
    await expect(this.getLostReasonStatusToggle(id)).toContainText('Aktif');
  }

  async expectLostReasonInactive(id: string): Promise<void> {
    await expect(this.getLostReasonStatusToggle(id)).toContainText('Nonaktif');
  }

  // Actions
  async openAddModal(): Promise<void> {
    await this.addLostReasonButton.click();
    await expect(this.addLostReasonModal).toBeVisible();
  }

  async addLostReason(data: {
    name: string;
    description?: string;
    displayOrder?: number;
  }): Promise<void> {
    await this.openAddModal();
    await this.inputName.fill(data.name);
    if (data.description) {
      await this.inputDescription.fill(data.description);
    }
    if (data.displayOrder !== undefined) {
      await this.inputDisplayOrder.fill(String(data.displayOrder));
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/lost-reasons') &&
      response.request().method() === 'POST' &&
      response.status() === 200
    );
    await this.submitAddButton.click();
    await responsePromise;
  }

  async openEditModal(id: string): Promise<void> {
    await this.getLostReasonEditButton(id).click();
    await expect(this.getEditModal(id)).toBeVisible();
  }

  async editLostReason(id: string, data: {
    name?: string;
    description?: string;
    displayOrder?: number;
  }): Promise<void> {
    await this.openEditModal(id);
    if (data.name) {
      await this.getEditInputName(id).fill(data.name);
    }
    if (data.description !== undefined) {
      await this.getEditInputDescription(id).fill(data.description);
    }
    if (data.displayOrder !== undefined) {
      await this.getEditInputOrder(id).fill(String(data.displayOrder));
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/lost-reasons/${id}`) &&
      !response.url().includes('toggle-active') &&
      response.request().method() === 'POST' &&
      response.status() === 200
    );
    await this.getEditSubmitButton(id).click();
    await responsePromise;
  }

  async toggleLostReasonStatus(id: string): Promise<void> {
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/lost-reasons/${id}/toggle-active`) &&
      response.status() === 200
    );
    await this.getLostReasonStatusToggle(id).click();
    await responsePromise;
  }
}
