import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsCampaignsPage extends BasePage {
  readonly path = '/admin/settings/campaigns';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    await this.page.waitForURL(/\/admin\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-campaigns-page');
  }

  get campaignsSection(): Locator {
    return this.page.getByTestId('campaigns-section');
  }

  get campaignsList(): Locator {
    return this.page.getByTestId('campaigns-list');
  }

  get addCampaignButton(): Locator {
    return this.page.getByTestId('add-campaign-button');
  }

  // Add Campaign Modal elements
  get addCampaignModal(): Locator {
    return this.page.getByTestId('add-campaign-modal');
  }

  get inputCampaignName(): Locator {
    return this.page.getByTestId('input-campaign-name');
  }

  get inputCampaignType(): Locator {
    return this.page.getByTestId('input-campaign-type');
  }

  get inputCampaignChannel(): Locator {
    return this.page.getByTestId('input-campaign-channel');
  }

  get inputCampaignStart(): Locator {
    return this.page.getByTestId('input-campaign-start');
  }

  get inputCampaignEnd(): Locator {
    return this.page.getByTestId('input-campaign-end');
  }

  get inputCampaignFee(): Locator {
    return this.page.getByTestId('input-campaign-fee');
  }

  get inputCampaignDescription(): Locator {
    return this.page.getByTestId('input-campaign-description');
  }

  get submitAddCampaignButton(): Locator {
    return this.page.getByTestId('submit-add-campaign');
  }

  // Campaign row elements by ID
  getCampaignRow(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-row-${campaignId}`);
  }

  getCampaignName(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-name-${campaignId}`);
  }

  getCampaignType(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-type-${campaignId}`);
  }

  getCampaignChannel(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-channel-${campaignId}`);
  }

  getCampaignPeriod(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-period-${campaignId}`);
  }

  getCampaignFee(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-fee-${campaignId}`);
  }

  getCampaignStatusToggle(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-status-${campaignId}`);
  }

  getCampaignEditButton(campaignId: string): Locator {
    return this.page.getByTestId(`campaign-edit-${campaignId}`);
  }

  // Edit Campaign Modal elements by ID
  getEditCampaignModal(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-modal-${campaignId}`);
  }

  getEditCampaignName(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-name-${campaignId}`);
  }

  getEditCampaignType(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-type-${campaignId}`);
  }

  getEditCampaignChannel(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-channel-${campaignId}`);
  }

  getEditCampaignStart(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-start-${campaignId}`);
  }

  getEditCampaignEnd(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-end-${campaignId}`);
  }

  getEditCampaignFee(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-fee-${campaignId}`);
  }

  getEditCampaignDescription(campaignId: string): Locator {
    return this.page.getByTestId(`edit-campaign-description-${campaignId}`);
  }

  getSubmitEditCampaignButton(campaignId: string): Locator {
    return this.page.getByTestId(`submit-edit-campaign-${campaignId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.campaignsSection).toBeVisible();
  }

  // Get all campaign IDs
  async getAllCampaignIds(): Promise<string[]> {
    const rows = this.campaignsList.locator('tr[data-testid^="campaign-row-"]');
    const count = await rows.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await rows.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('campaign-row-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectCampaignDisplayed(campaignId: string): Promise<void> {
    await expect(this.getCampaignRow(campaignId)).toBeVisible();
    await expect(this.getCampaignName(campaignId)).toBeVisible();
  }

  // Campaign Actions
  async openAddCampaignModal(): Promise<void> {
    await this.addCampaignButton.click();
    await expect(this.addCampaignModal).toBeVisible();
  }

  async addCampaign(
    name: string,
    type: string,
    channel?: string,
    startDate?: string,
    endDate?: string,
    feeOverride?: number,
    description?: string
  ): Promise<void> {
    await this.openAddCampaignModal();
    await this.inputCampaignName.fill(name);
    await this.inputCampaignType.selectOption(type);

    if (channel) {
      await this.inputCampaignChannel.selectOption(channel);
    }

    if (startDate) {
      await this.inputCampaignStart.fill(startDate);
    }

    if (endDate) {
      await this.inputCampaignEnd.fill(endDate);
    }

    if (feeOverride !== undefined) {
      await this.inputCampaignFee.fill(feeOverride.toString());
    }

    if (description) {
      await this.inputCampaignDescription.fill(description);
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/campaigns') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      !response.url().match(/\/campaigns\/[a-f0-9-]+$/) &&
      response.status() === 200
    );
    await this.submitAddCampaignButton.click();
    await responsePromise;
  }

  async openEditCampaignModal(campaignId: string): Promise<void> {
    await this.getCampaignEditButton(campaignId).click();
    await expect(this.getEditCampaignModal(campaignId)).toBeVisible();
  }

  async editCampaignName(campaignId: string, name: string): Promise<void> {
    await this.openEditCampaignModal(campaignId);
    await this.getEditCampaignName(campaignId).fill(name);

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/campaigns/${campaignId}`) &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.getSubmitEditCampaignButton(campaignId).click();
    await responsePromise;
  }

  async toggleCampaignStatus(campaignId: string): Promise<void> {
    const button = this.getCampaignStatusToggle(campaignId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/campaigns/${campaignId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  async expectCampaignStatus(campaignId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getCampaignStatusToggle(campaignId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  async expectCampaignNameContains(campaignId: string, expectedName: string): Promise<void> {
    const name = await this.getCampaignName(campaignId).textContent();
    expect(name).toContain(expectedName);
  }
}
