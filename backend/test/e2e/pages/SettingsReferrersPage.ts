import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsReferrersPage extends BasePage {
  readonly path = '/admin/settings/referrers';


  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-referrers-page');
  }

  get referrersSection(): Locator {
    return this.page.getByTestId('referrers-section');
  }

  get referrersList(): Locator {
    return this.page.getByTestId('referrers-list');
  }

  get referrerStats(): Locator {
    return this.page.getByTestId('referrer-stats');
  }

  get addReferrerButton(): Locator {
    return this.page.getByTestId('add-referrer-button');
  }

  // Stats elements
  get statTotal(): Locator {
    return this.page.getByTestId('stat-total');
  }

  get statAlumni(): Locator {
    return this.page.getByTestId('stat-alumni');
  }

  get statTeacher(): Locator {
    return this.page.getByTestId('stat-teacher');
  }

  get statStudent(): Locator {
    return this.page.getByTestId('stat-student');
  }

  get statPartner(): Locator {
    return this.page.getByTestId('stat-partner');
  }

  get statStaff(): Locator {
    return this.page.getByTestId('stat-staff');
  }

  // Add Referrer Modal elements
  get addReferrerModal(): Locator {
    return this.page.getByTestId('add-referrer-modal');
  }

  get inputReferrerName(): Locator {
    return this.page.getByTestId('input-referrer-name');
  }

  get inputReferrerType(): Locator {
    return this.page.getByTestId('input-referrer-type');
  }

  get inputReferrerInstitution(): Locator {
    return this.page.getByTestId('input-referrer-institution');
  }

  get inputReferrerPhone(): Locator {
    return this.page.getByTestId('input-referrer-phone');
  }

  get inputReferrerEmail(): Locator {
    return this.page.getByTestId('input-referrer-email');
  }

  get inputReferrerCode(): Locator {
    return this.page.getByTestId('input-referrer-code');
  }

  get inputReferrerCommission(): Locator {
    return this.page.getByTestId('input-referrer-commission');
  }

  get inputReferrerPayout(): Locator {
    return this.page.getByTestId('input-referrer-payout');
  }

  get inputReferrerBank(): Locator {
    return this.page.getByTestId('input-referrer-bank');
  }

  get inputReferrerAccount(): Locator {
    return this.page.getByTestId('input-referrer-account');
  }

  get inputReferrerHolder(): Locator {
    return this.page.getByTestId('input-referrer-holder');
  }

  get submitAddReferrerButton(): Locator {
    return this.page.getByTestId('submit-add-referrer');
  }

  // Referrer row elements by ID
  getReferrerRow(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-row-${referrerId}`);
  }

  getReferrerName(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-name-${referrerId}`);
  }

  getReferrerType(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-type-${referrerId}`);
  }

  getReferrerInstitution(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-institution-${referrerId}`);
  }

  getReferrerCode(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-code-${referrerId}`);
  }

  getReferrerCommission(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-commission-${referrerId}`);
  }

  getReferrerStatusToggle(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-status-${referrerId}`);
  }

  getReferrerEditButton(referrerId: string): Locator {
    return this.page.getByTestId(`referrer-edit-${referrerId}`);
  }

  // Edit Referrer Modal elements by ID
  getEditReferrerModal(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-modal-${referrerId}`);
  }

  getEditReferrerName(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-name-${referrerId}`);
  }

  getEditReferrerType(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-type-${referrerId}`);
  }

  getEditReferrerInstitution(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-institution-${referrerId}`);
  }

  getEditReferrerPhone(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-phone-${referrerId}`);
  }

  getEditReferrerEmail(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-email-${referrerId}`);
  }

  getEditReferrerCode(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-code-${referrerId}`);
  }

  getEditReferrerCommission(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-commission-${referrerId}`);
  }

  getEditReferrerPayout(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-payout-${referrerId}`);
  }

  getEditReferrerBank(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-bank-${referrerId}`);
  }

  getEditReferrerAccount(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-account-${referrerId}`);
  }

  getEditReferrerHolder(referrerId: string): Locator {
    return this.page.getByTestId(`edit-referrer-holder-${referrerId}`);
  }

  getSubmitEditReferrerButton(referrerId: string): Locator {
    return this.page.getByTestId(`submit-edit-referrer-${referrerId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.referrersSection).toBeVisible();
  }

  // Get all referrer IDs
  async getAllReferrerIds(): Promise<string[]> {
    const rows = this.referrersList.locator('tr[data-testid^="referrer-row-"]');
    const count = await rows.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await rows.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('referrer-row-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectReferrerDisplayed(referrerId: string): Promise<void> {
    await expect(this.getReferrerRow(referrerId)).toBeVisible();
    await expect(this.getReferrerName(referrerId)).toBeVisible();
  }

  // Referrer Actions
  async openAddReferrerModal(): Promise<void> {
    await this.addReferrerButton.click();
    await expect(this.addReferrerModal).toBeVisible();
  }

  async addReferrer(
    name: string,
    type: string,
    institution?: string,
    phone?: string,
    email?: string,
    code?: string,
    commission?: number,
    payoutPreference?: string,
    bankName?: string,
    bankAccount?: string,
    accountHolder?: string
  ): Promise<void> {
    await this.openAddReferrerModal();
    await this.inputReferrerName.fill(name);
    await this.inputReferrerType.selectOption(type);

    if (institution) {
      await this.inputReferrerInstitution.fill(institution);
    }

    if (phone) {
      await this.inputReferrerPhone.fill(phone);
    }

    if (email) {
      await this.inputReferrerEmail.fill(email);
    }

    if (code) {
      await this.inputReferrerCode.fill(code);
    }

    if (commission !== undefined) {
      await this.inputReferrerCommission.fill(commission.toString());
    }

    if (payoutPreference) {
      await this.inputReferrerPayout.selectOption(payoutPreference);
    }

    if (bankName) {
      await this.inputReferrerBank.fill(bankName);
    }

    if (bankAccount) {
      await this.inputReferrerAccount.fill(bankAccount);
    }

    if (accountHolder) {
      await this.inputReferrerHolder.fill(accountHolder);
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/referrers') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      !response.url().match(/\/referrers\/[a-f0-9-]+$/) &&
      response.status() === 200
    );
    await this.submitAddReferrerButton.click();
    await responsePromise;
  }

  async openEditReferrerModal(referrerId: string): Promise<void> {
    await this.getReferrerEditButton(referrerId).click();
    await expect(this.getEditReferrerModal(referrerId)).toBeVisible();
  }

  async editReferrerName(referrerId: string, name: string): Promise<void> {
    await this.openEditReferrerModal(referrerId);
    await this.getEditReferrerName(referrerId).fill(name);

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/referrers/${referrerId}`) &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.getSubmitEditReferrerButton(referrerId).click();
    await responsePromise;
  }

  async editReferrerInstitution(referrerId: string, institution: string): Promise<void> {
    await this.openEditReferrerModal(referrerId);
    await this.getEditReferrerInstitution(referrerId).fill(institution);

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/referrers/${referrerId}`) &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.getSubmitEditReferrerButton(referrerId).click();
    await responsePromise;
  }

  async toggleReferrerStatus(referrerId: string): Promise<void> {
    const button = this.getReferrerStatusToggle(referrerId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/referrers/${referrerId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  async expectReferrerStatus(referrerId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getReferrerStatusToggle(referrerId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  async expectReferrerNameContains(referrerId: string, expectedName: string): Promise<void> {
    const name = await this.getReferrerName(referrerId).textContent();
    expect(name).toContain(expectedName);
  }

  async expectReferrerInstitutionContains(referrerId: string, expectedInstitution: string): Promise<void> {
    const institution = await this.getReferrerInstitution(referrerId).textContent();
    expect(institution).toContain(expectedInstitution);
  }
}
