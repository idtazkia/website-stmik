import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsRewardsPage extends BasePage {
  readonly path = '/admin/settings/rewards';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    await this.page.waitForURL(/\/admin\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-rewards-page');
  }

  get rewardsSection(): Locator {
    return this.page.getByTestId('rewards-section');
  }

  get rewardsList(): Locator {
    return this.page.getByTestId('rewards-list');
  }

  get mgmRewardsSection(): Locator {
    return this.page.getByTestId('mgm-rewards-section');
  }

  get mgmRewardsList(): Locator {
    return this.page.getByTestId('mgm-rewards-list');
  }

  get addRewardButton(): Locator {
    return this.page.getByTestId('add-reward-button');
  }

  get addMGMRewardButton(): Locator {
    return this.page.getByTestId('add-mgm-reward-button');
  }

  // Add Reward Modal elements
  get addRewardModal(): Locator {
    return this.page.getByTestId('add-reward-modal');
  }

  get inputRewardReferrerType(): Locator {
    return this.page.getByTestId('input-reward-referrer-type');
  }

  get inputRewardType(): Locator {
    return this.page.getByTestId('input-reward-type');
  }

  get inputRewardTrigger(): Locator {
    return this.page.getByTestId('input-reward-trigger');
  }

  get inputRewardAmount(): Locator {
    return this.page.getByTestId('input-reward-amount');
  }

  get inputRewardPercentage(): Locator {
    return this.page.getByTestId('input-reward-percentage');
  }

  get inputRewardDescription(): Locator {
    return this.page.getByTestId('input-reward-description');
  }

  get submitAddRewardButton(): Locator {
    return this.page.getByTestId('submit-add-reward');
  }

  // Add MGM Reward Modal elements
  get addMGMRewardModal(): Locator {
    return this.page.getByTestId('add-mgm-reward-modal');
  }

  get inputMGMYear(): Locator {
    return this.page.getByTestId('input-mgm-year');
  }

  get inputMGMRewardType(): Locator {
    return this.page.getByTestId('input-mgm-reward-type');
  }

  get inputMGMTrigger(): Locator {
    return this.page.getByTestId('input-mgm-trigger');
  }

  get inputMGMReferrerAmount(): Locator {
    return this.page.getByTestId('input-mgm-referrer-amount');
  }

  get inputMGMRefereeAmount(): Locator {
    return this.page.getByTestId('input-mgm-referee-amount');
  }

  get inputMGMDescription(): Locator {
    return this.page.getByTestId('input-mgm-description');
  }

  get submitAddMGMRewardButton(): Locator {
    return this.page.getByTestId('submit-add-mgm-reward');
  }

  // Reward card elements by ID
  getRewardCard(rewardId: string): Locator {
    return this.page.getByTestId(`reward-card-${rewardId}`);
  }

  getRewardReferrerType(rewardId: string): Locator {
    return this.page.getByTestId(`reward-referrer-type-${rewardId}`);
  }

  getRewardType(rewardId: string): Locator {
    return this.page.getByTestId(`reward-type-${rewardId}`);
  }

  getRewardAmount(rewardId: string): Locator {
    return this.page.getByTestId(`reward-amount-${rewardId}`);
  }

  getRewardStatusToggle(rewardId: string): Locator {
    return this.page.getByTestId(`reward-status-${rewardId}`);
  }

  getRewardEditButton(rewardId: string): Locator {
    return this.page.getByTestId(`reward-edit-${rewardId}`);
  }

  // MGM Reward card elements by ID
  getMGMRewardCard(mgmRewardId: string): Locator {
    return this.page.getByTestId(`mgm-reward-card-${mgmRewardId}`);
  }

  getMGMYear(mgmRewardId: string): Locator {
    return this.page.getByTestId(`mgm-year-${mgmRewardId}`);
  }

  getMGMRewardType(mgmRewardId: string): Locator {
    return this.page.getByTestId(`mgm-reward-type-${mgmRewardId}`);
  }

  getMGMReferrerAmount(mgmRewardId: string): Locator {
    return this.page.getByTestId(`mgm-referrer-amount-${mgmRewardId}`);
  }

  getMGMStatusToggle(mgmRewardId: string): Locator {
    return this.page.getByTestId(`mgm-status-${mgmRewardId}`);
  }

  getMGMEditButton(mgmRewardId: string): Locator {
    return this.page.getByTestId(`mgm-edit-${mgmRewardId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.rewardsSection).toBeVisible();
    await expect(this.mgmRewardsSection).toBeVisible();
  }

  // Get all reward IDs
  async getAllRewardIds(): Promise<string[]> {
    const cards = this.rewardsList.locator('div[data-testid^="reward-card-"]');
    const count = await cards.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await cards.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('reward-card-', ''));
      }
    }
    return ids;
  }

  // Get all MGM reward IDs
  async getAllMGMRewardIds(): Promise<string[]> {
    const cards = this.mgmRewardsList.locator('div[data-testid^="mgm-reward-card-"]');
    const count = await cards.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await cards.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('mgm-reward-card-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectRewardDisplayed(rewardId: string): Promise<void> {
    await expect(this.getRewardCard(rewardId)).toBeVisible();
    await expect(this.getRewardReferrerType(rewardId)).toBeVisible();
  }

  async expectMGMRewardDisplayed(mgmRewardId: string): Promise<void> {
    await expect(this.getMGMRewardCard(mgmRewardId)).toBeVisible();
    await expect(this.getMGMYear(mgmRewardId)).toBeVisible();
  }

  // Reward Actions
  async openAddRewardModal(): Promise<void> {
    await this.addRewardButton.click();
    await expect(this.addRewardModal).toBeVisible();
  }

  async addReward(
    referrerType: string,
    rewardType: string,
    triggerEvent: string,
    amount: number,
    isPercentage: boolean = false,
    description?: string
  ): Promise<void> {
    await this.openAddRewardModal();
    await this.inputRewardReferrerType.selectOption(referrerType);
    await this.inputRewardType.selectOption(rewardType);
    await this.inputRewardTrigger.selectOption(triggerEvent);
    await this.inputRewardAmount.fill(amount.toString());

    if (isPercentage) {
      await this.inputRewardPercentage.check();
    }

    if (description) {
      await this.inputRewardDescription.fill(description);
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().endsWith('/admin/settings/rewards') &&
      response.request().method() === 'POST'
    );
    await this.submitAddRewardButton.click();
    await responsePromise;
  }

  async toggleRewardStatus(rewardId: string): Promise<void> {
    const button = this.getRewardStatusToggle(rewardId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/rewards/${rewardId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  async expectRewardStatus(rewardId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getRewardStatusToggle(rewardId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  // MGM Reward Actions
  async openAddMGMRewardModal(): Promise<void> {
    await this.addMGMRewardButton.click();
    await expect(this.addMGMRewardModal).toBeVisible();
  }

  async addMGMReward(
    academicYear: string,
    rewardType: string,
    triggerEvent: string,
    referrerAmount: number,
    refereeAmount?: number,
    description?: string
  ): Promise<void> {
    await this.openAddMGMRewardModal();
    await this.inputMGMYear.fill(academicYear);
    await this.inputMGMRewardType.selectOption(rewardType);
    await this.inputMGMTrigger.selectOption(triggerEvent);
    await this.inputMGMReferrerAmount.fill(referrerAmount.toString());

    if (refereeAmount !== undefined) {
      await this.inputMGMRefereeAmount.fill(refereeAmount.toString());
    }

    if (description) {
      await this.inputMGMDescription.fill(description);
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().endsWith('/admin/settings/mgm-rewards') &&
      response.request().method() === 'POST'
    );
    await this.submitAddMGMRewardButton.click();
    await responsePromise;
  }

  async toggleMGMRewardStatus(mgmRewardId: string): Promise<void> {
    const button = this.getMGMStatusToggle(mgmRewardId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/mgm-rewards/${mgmRewardId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  async expectMGMRewardStatus(mgmRewardId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getMGMStatusToggle(mgmRewardId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }
}
