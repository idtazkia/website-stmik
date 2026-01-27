import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsFeesPage extends BasePage {
  readonly path = '/admin/settings/fees';


  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-fees-page');
  }

  get feesSection(): Locator {
    return this.page.getByTestId('fees-section');
  }

  get feesTable(): Locator {
    return this.page.getByTestId('fees-table');
  }

  get feesList(): Locator {
    return this.page.getByTestId('fees-list');
  }

  get academicYearSelect(): Locator {
    return this.page.getByTestId('academic-year-select');
  }

  get addFeeButton(): Locator {
    return this.page.getByTestId('add-fee-button');
  }

  // Add Fee Modal elements
  get addFeeModal(): Locator {
    return this.page.getByTestId('add-fee-modal');
  }

  get inputFeeType(): Locator {
    return this.page.getByTestId('input-fee-type');
  }

  get inputFeeProdi(): Locator {
    return this.page.getByTestId('input-fee-prodi');
  }

  get inputFeeAmount(): Locator {
    return this.page.getByTestId('input-fee-amount');
  }

  get submitAddFeeButton(): Locator {
    return this.page.getByTestId('submit-add-fee');
  }

  // Fee row elements by ID
  getFeeRow(feeId: string): Locator {
    return this.page.getByTestId(`fee-row-${feeId}`);
  }

  getFeeTypeName(feeId: string): Locator {
    return this.page.getByTestId(`fee-type-name-${feeId}`);
  }

  getFeeProdi(feeId: string): Locator {
    return this.page.getByTestId(`fee-prodi-${feeId}`);
  }

  getFeeAmount(feeId: string): Locator {
    return this.page.getByTestId(`fee-amount-${feeId}`);
  }

  getFeeStatusToggle(feeId: string): Locator {
    return this.page.getByTestId(`fee-status-toggle-${feeId}`);
  }

  getFeeEditButton(feeId: string): Locator {
    return this.page.getByTestId(`fee-edit-${feeId}`);
  }

  // Edit Fee Modal elements by ID
  getEditFeeModal(feeId: string): Locator {
    return this.page.getByTestId(`edit-fee-modal-${feeId}`);
  }

  getEditFeeType(feeId: string): Locator {
    return this.page.getByTestId(`edit-fee-type-${feeId}`);
  }

  getEditFeeProdi(feeId: string): Locator {
    return this.page.getByTestId(`edit-fee-prodi-${feeId}`);
  }

  getEditFeeAmount(feeId: string): Locator {
    return this.page.getByTestId(`edit-fee-amount-${feeId}`);
  }

  getSubmitEditFeeButton(feeId: string): Locator {
    return this.page.getByTestId(`submit-edit-fee-${feeId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.feesSection).toBeVisible();
    await expect(this.feesTable).toBeVisible();
  }

  // Get all fee IDs
  async getAllFeeIds(): Promise<string[]> {
    const rows = this.feesList.locator('tr[data-testid^="fee-row-"]');
    const count = await rows.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await rows.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('fee-row-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectFeeDisplayed(feeId: string): Promise<void> {
    await expect(this.getFeeRow(feeId)).toBeVisible();
    await expect(this.getFeeTypeName(feeId)).toBeVisible();
    await expect(this.getFeeAmount(feeId)).toBeVisible();
  }

  // Fee Actions
  async openAddFeeModal(): Promise<void> {
    await this.addFeeButton.click();
    await expect(this.addFeeModal).toBeVisible();
  }

  async addFee(feeTypeIndex: number, prodiIndex: number | null, amount: number): Promise<void> {
    await this.openAddFeeModal();
    // Select fee type by index (first option is placeholder)
    const feeTypeOptions = await this.inputFeeType.locator('option').all();
    if (feeTypeIndex > 0 && feeTypeIndex < feeTypeOptions.length) {
      const value = await feeTypeOptions[feeTypeIndex].getAttribute('value');
      if (value) {
        await this.inputFeeType.selectOption(value);
      }
    }

    // Select prodi by index if specified (0 = all prodi)
    if (prodiIndex !== null && prodiIndex > 0) {
      const prodiOptions = await this.inputFeeProdi.locator('option').all();
      if (prodiIndex < prodiOptions.length) {
        const value = await prodiOptions[prodiIndex].getAttribute('value');
        if (value) {
          await this.inputFeeProdi.selectOption(value);
        }
      }
    }

    await this.inputFeeAmount.fill(amount.toString());

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/fees') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.submitAddFeeButton.click();
    await responsePromise;
  }

  async addFeeWithTypeId(feeTypeId: string, prodiId: string | null, amount: number): Promise<void> {
    await this.openAddFeeModal();
    await this.inputFeeType.selectOption(feeTypeId);

    if (prodiId) {
      await this.inputFeeProdi.selectOption(prodiId);
    }

    await this.inputFeeAmount.fill(amount.toString());

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/fees') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.submitAddFeeButton.click();
    await responsePromise;
  }

  async openEditFeeModal(feeId: string): Promise<void> {
    await this.getFeeEditButton(feeId).click();
    await expect(this.getEditFeeModal(feeId)).toBeVisible();
  }

  async editFeeAmount(feeId: string, amount: number): Promise<void> {
    await this.openEditFeeModal(feeId);
    await this.getEditFeeAmount(feeId).fill(amount.toString());

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/fees/${feeId}`) &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.getSubmitEditFeeButton(feeId).click();
    await responsePromise;
  }

  async toggleFeeStatus(feeId: string): Promise<void> {
    const button = this.getFeeStatusToggle(feeId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/fees/${feeId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  async expectFeeStatus(feeId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getFeeStatusToggle(feeId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  async expectFeeAmountContains(feeId: string, expectedAmount: string): Promise<void> {
    const amount = await this.getFeeAmount(feeId).textContent();
    expect(amount).toContain(expectedAmount);
  }
}
