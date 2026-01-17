import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsCategoriesPage extends BasePage {
  readonly path = '/admin/settings/categories';


  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-categories-page');
  }

  get categoriesSection(): Locator {
    return this.page.getByTestId('categories-section');
  }

  get categoriesList(): Locator {
    return this.page.getByTestId('categories-list');
  }

  get addCategoryButton(): Locator {
    return this.page.getByTestId('add-category-button');
  }

  get obstaclesSection(): Locator {
    return this.page.getByTestId('obstacles-section');
  }

  get obstaclesList(): Locator {
    return this.page.getByTestId('obstacles-list');
  }

  get addObstacleButton(): Locator {
    return this.page.getByTestId('add-obstacle-button');
  }

  // Add Category Modal elements
  get addCategoryModal(): Locator {
    return this.page.getByTestId('add-category-modal');
  }

  get inputCategoryName(): Locator {
    return this.page.getByTestId('input-category-name');
  }

  get inputCategorySentiment(): Locator {
    return this.page.getByTestId('input-category-sentiment');
  }

  get submitAddCategoryButton(): Locator {
    return this.page.getByTestId('submit-add-category');
  }

  // Add Obstacle Modal elements
  get addObstacleModal(): Locator {
    return this.page.getByTestId('add-obstacle-modal');
  }

  get inputObstacleName(): Locator {
    return this.page.getByTestId('input-obstacle-name');
  }

  get submitAddObstacleButton(): Locator {
    return this.page.getByTestId('submit-add-obstacle');
  }

  // Category item elements by ID
  getCategoryItem(categoryId: string): Locator {
    return this.page.getByTestId(`category-item-${categoryId}`);
  }

  getCategoryIcon(categoryId: string): Locator {
    return this.page.getByTestId(`category-icon-${categoryId}`);
  }

  getCategoryName(categoryId: string): Locator {
    return this.page.getByTestId(`category-name-${categoryId}`);
  }

  getCategorySentiment(categoryId: string): Locator {
    return this.page.getByTestId(`category-sentiment-${categoryId}`);
  }

  getCategoryCount(categoryId: string): Locator {
    return this.page.getByTestId(`category-count-${categoryId}`);
  }

  getCategoryStatusToggle(categoryId: string): Locator {
    return this.page.getByTestId(`category-status-toggle-${categoryId}`);
  }

  getCategoryEditButton(categoryId: string): Locator {
    return this.page.getByTestId(`category-edit-${categoryId}`);
  }

  // Edit Category Modal elements by ID
  getEditCategoryModal(categoryId: string): Locator {
    return this.page.getByTestId(`edit-category-modal-${categoryId}`);
  }

  getEditCategoryName(categoryId: string): Locator {
    return this.page.getByTestId(`edit-category-name-${categoryId}`);
  }

  getEditCategorySentiment(categoryId: string): Locator {
    return this.page.getByTestId(`edit-category-sentiment-${categoryId}`);
  }

  getSubmitEditCategoryButton(categoryId: string): Locator {
    return this.page.getByTestId(`submit-edit-category-${categoryId}`);
  }

  // Obstacle item elements by ID
  getObstacleItem(obstacleId: string): Locator {
    return this.page.getByTestId(`obstacle-item-${obstacleId}`);
  }

  getObstacleName(obstacleId: string): Locator {
    return this.page.getByTestId(`obstacle-name-${obstacleId}`);
  }

  getObstacleCount(obstacleId: string): Locator {
    return this.page.getByTestId(`obstacle-count-${obstacleId}`);
  }

  getObstacleStatusToggle(obstacleId: string): Locator {
    return this.page.getByTestId(`obstacle-status-toggle-${obstacleId}`);
  }

  getObstacleEditButton(obstacleId: string): Locator {
    return this.page.getByTestId(`obstacle-edit-${obstacleId}`);
  }

  // Edit Obstacle Modal elements by ID
  getEditObstacleModal(obstacleId: string): Locator {
    return this.page.getByTestId(`edit-obstacle-modal-${obstacleId}`);
  }

  getEditObstacleName(obstacleId: string): Locator {
    return this.page.getByTestId(`edit-obstacle-name-${obstacleId}`);
  }

  getSubmitEditObstacleButton(obstacleId: string): Locator {
    return this.page.getByTestId(`submit-edit-obstacle-${obstacleId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.categoriesSection).toBeVisible();
    await expect(this.obstaclesSection).toBeVisible();
  }

  // Get all category IDs
  async getAllCategoryIds(): Promise<string[]> {
    const items = this.categoriesList.locator('div[data-testid^="category-item-"]');
    const count = await items.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await items.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('category-item-', ''));
      }
    }
    return ids;
  }

  // Get all obstacle IDs
  async getAllObstacleIds(): Promise<string[]> {
    const items = this.obstaclesList.locator('div[data-testid^="obstacle-item-"]');
    const count = await items.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await items.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('obstacle-item-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectCategoryDisplayed(categoryId: string): Promise<void> {
    await expect(this.getCategoryItem(categoryId)).toBeVisible();
    await expect(this.getCategoryName(categoryId)).toBeVisible();
  }

  async expectObstacleDisplayed(obstacleId: string): Promise<void> {
    await expect(this.getObstacleItem(obstacleId)).toBeVisible();
    await expect(this.getObstacleName(obstacleId)).toBeVisible();
  }

  // Category Actions
  async openAddCategoryModal(): Promise<void> {
    await this.addCategoryButton.click();
    await expect(this.addCategoryModal).toBeVisible();
  }

  async addCategory(name: string, sentiment: 'positive' | 'neutral' | 'negative'): Promise<void> {
    await this.openAddCategoryModal();
    await this.inputCategoryName.fill(name);
    await this.inputCategorySentiment.selectOption(sentiment);

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/categories') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.submitAddCategoryButton.click();
    await responsePromise;
  }

  async openEditCategoryModal(categoryId: string): Promise<void> {
    await this.getCategoryEditButton(categoryId).click();
    await expect(this.getEditCategoryModal(categoryId)).toBeVisible();
  }

  async editCategory(categoryId: string, name: string, sentiment: 'positive' | 'neutral' | 'negative'): Promise<void> {
    await this.openEditCategoryModal(categoryId);
    await this.getEditCategoryName(categoryId).fill(name);
    await this.getEditCategorySentiment(categoryId).selectOption(sentiment);

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/categories/${categoryId}`) &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.getSubmitEditCategoryButton(categoryId).click();
    await responsePromise;
  }

  async toggleCategoryStatus(categoryId: string): Promise<void> {
    const button = this.getCategoryStatusToggle(categoryId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/categories/${categoryId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  async expectCategoryStatus(categoryId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getCategoryStatusToggle(categoryId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  async expectCategoryNameValue(categoryId: string, expectedName: string): Promise<void> {
    const name = await this.getCategoryName(categoryId).textContent();
    expect(name?.trim()).toBe(expectedName);
  }

  // Obstacle Actions
  async openAddObstacleModal(): Promise<void> {
    await this.addObstacleButton.click();
    await expect(this.addObstacleModal).toBeVisible();
  }

  async addObstacle(name: string): Promise<void> {
    await this.openAddObstacleModal();
    await this.inputObstacleName.fill(name);

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/obstacles') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.submitAddObstacleButton.click();
    await responsePromise;
  }

  async openEditObstacleModal(obstacleId: string): Promise<void> {
    await this.getObstacleEditButton(obstacleId).click();
    await expect(this.getEditObstacleModal(obstacleId)).toBeVisible();
  }

  async editObstacle(obstacleId: string, name: string): Promise<void> {
    await this.openEditObstacleModal(obstacleId);
    await this.getEditObstacleName(obstacleId).fill(name);

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/obstacles/${obstacleId}`) &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.getSubmitEditObstacleButton(obstacleId).click();
    await responsePromise;
  }

  async toggleObstacleStatus(obstacleId: string): Promise<void> {
    const button = this.getObstacleStatusToggle(obstacleId);
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/obstacles/${obstacleId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  async expectObstacleStatus(obstacleId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getObstacleStatusToggle(obstacleId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  async expectObstacleNameValue(obstacleId: string, expectedName: string): Promise<void> {
    const name = await this.getObstacleName(obstacleId).textContent();
    expect(name?.trim()).toBe(expectedName);
  }
}
