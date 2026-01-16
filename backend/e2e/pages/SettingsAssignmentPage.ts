import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsAssignmentPage extends BasePage {
  readonly path = '/admin/settings/assignment';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    await this.page.waitForURL(/\/admin\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-assignment-page');
  }

  get algorithmsSection(): Locator {
    return this.page.getByTestId('algorithms-section');
  }

  get algorithmsList(): Locator {
    return this.page.getByTestId('algorithms-list');
  }

  // Algorithm card elements by ID
  getAlgorithmCard(algorithmId: string): Locator {
    return this.page.getByTestId(`algorithm-card-${algorithmId}`);
  }

  getAlgorithmName(algorithmId: string): Locator {
    return this.page.getByTestId(`algorithm-name-${algorithmId}`);
  }

  getAlgorithmCode(algorithmId: string): Locator {
    return this.page.getByTestId(`algorithm-code-${algorithmId}`);
  }

  getAlgorithmDescription(algorithmId: string): Locator {
    return this.page.getByTestId(`algorithm-description-${algorithmId}`);
  }

  getAlgorithmStatus(algorithmId: string): Locator {
    return this.page.getByTestId(`algorithm-status-${algorithmId}`);
  }

  getAlgorithmActivateButton(algorithmId: string): Locator {
    return this.page.getByTestId(`algorithm-activate-${algorithmId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.algorithmsSection).toBeVisible();
  }

  // Get all algorithm IDs
  async getAllAlgorithmIds(): Promise<string[]> {
    const cards = this.algorithmsList.locator('div[data-testid^="algorithm-card-"]');
    const count = await cards.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await cards.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('algorithm-card-', ''));
      }
    }
    return ids;
  }

  // Get active algorithm ID
  async getActiveAlgorithmId(): Promise<string | null> {
    const ids = await this.getAllAlgorithmIds();
    for (const id of ids) {
      const statusLocator = this.getAlgorithmStatus(id);
      if (await statusLocator.isVisible()) {
        return id;
      }
    }
    return null;
  }

  // Verification helpers
  async expectAlgorithmDisplayed(algorithmId: string): Promise<void> {
    await expect(this.getAlgorithmCard(algorithmId)).toBeVisible();
    await expect(this.getAlgorithmName(algorithmId)).toBeVisible();
  }

  async expectAlgorithmActive(algorithmId: string): Promise<void> {
    await expect(this.getAlgorithmStatus(algorithmId)).toBeVisible();
    await expect(this.getAlgorithmStatus(algorithmId)).toContainText('Aktif');
  }

  async expectAlgorithmInactive(algorithmId: string): Promise<void> {
    await expect(this.getAlgorithmActivateButton(algorithmId)).toBeVisible();
  }

  // Actions
  async activateAlgorithm(algorithmId: string): Promise<void> {
    const button = this.getAlgorithmActivateButton(algorithmId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/assignment/${algorithmId}/activate`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }
}
