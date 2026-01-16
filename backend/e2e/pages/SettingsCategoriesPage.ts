import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsCategoriesPage extends BasePage {
  readonly path = '/admin/settings/categories';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    // Wait for redirect to admin dashboard (with or without trailing slash)
    await this.page.waitForURL(/\/admin\/?$/);
  }

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

  getCategoryEditButton(categoryId: string): Locator {
    return this.page.getByTestId(`category-edit-${categoryId}`);
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

  getObstacleEditButton(obstacleId: string): Locator {
    return this.page.getByTestId(`obstacle-edit-${obstacleId}`);
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
}
