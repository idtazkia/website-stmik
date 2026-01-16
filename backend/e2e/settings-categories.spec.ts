import { test, expect } from '@playwright/test';
import { SettingsCategoriesPage } from './pages';

test.describe('Settings - Categories Management', () => {
  let categoriesPage: SettingsCategoriesPage;

  test.beforeEach(async ({ page }) => {
    categoriesPage = new SettingsCategoriesPage(page);
    // Login as admin before each test
    await categoriesPage.login('admin');
    await categoriesPage.goto();
    await categoriesPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display categories page with both sections', async () => {
      await expect(categoriesPage.categoriesSection).toBeVisible();
      await expect(categoriesPage.obstaclesSection).toBeVisible();
    });

    test('should display add buttons for categories and obstacles', async () => {
      await expect(categoriesPage.addCategoryButton).toBeVisible();
      await expect(categoriesPage.addObstacleButton).toBeVisible();
    });
  });

  test.describe('Categories Display', () => {
    test('should display category list from database', async () => {
      // Verify at least one category is displayed
      const categoryIds = await categoriesPage.getAllCategoryIds();
      expect(categoryIds.length).toBeGreaterThan(0);
    });

    test('should display category details correctly', async () => {
      const categoryIds = await categoriesPage.getAllCategoryIds();
      expect(categoryIds.length).toBeGreaterThan(0);

      // Check first category displays required fields
      const firstCategoryId = categoryIds[0];
      await categoriesPage.expectCategoryDisplayed(firstCategoryId);

      // Verify icon, sentiment, and count are visible
      await expect(categoriesPage.getCategoryIcon(firstCategoryId)).toBeVisible();
      await expect(categoriesPage.getCategorySentiment(firstCategoryId)).toBeVisible();
      await expect(categoriesPage.getCategoryCount(firstCategoryId)).toBeVisible();
    });

    test('should display sentiment label correctly', async () => {
      const categoryIds = await categoriesPage.getAllCategoryIds();
      expect(categoryIds.length).toBeGreaterThan(0);

      for (const categoryId of categoryIds) {
        const sentimentText = await categoriesPage.getCategorySentiment(categoryId).textContent();
        expect(sentimentText?.trim()).toMatch(/^(Positif|Netral|Negatif)$/);
      }
    });

    test('should display edit button for each category', async () => {
      const categoryIds = await categoriesPage.getAllCategoryIds();
      expect(categoryIds.length).toBeGreaterThan(0);

      for (const categoryId of categoryIds) {
        await expect(categoriesPage.getCategoryEditButton(categoryId)).toBeVisible();
      }
    });
  });

  test.describe('Obstacles Display', () => {
    test('should display obstacle list from database', async () => {
      // Verify at least one obstacle is displayed
      const obstacleIds = await categoriesPage.getAllObstacleIds();
      expect(obstacleIds.length).toBeGreaterThan(0);
    });

    test('should display obstacle details correctly', async () => {
      const obstacleIds = await categoriesPage.getAllObstacleIds();
      expect(obstacleIds.length).toBeGreaterThan(0);

      // Check first obstacle displays required fields
      const firstObstacleId = obstacleIds[0];
      await categoriesPage.expectObstacleDisplayed(firstObstacleId);

      // Verify count is visible
      await expect(categoriesPage.getObstacleCount(firstObstacleId)).toBeVisible();
    });

    test('should display edit button for each obstacle', async () => {
      const obstacleIds = await categoriesPage.getAllObstacleIds();
      expect(obstacleIds.length).toBeGreaterThan(0);

      for (const obstacleId of obstacleIds) {
        await expect(categoriesPage.getObstacleEditButton(obstacleId)).toBeVisible();
      }
    });

    test('should display obstacle name and count', async () => {
      const obstacleIds = await categoriesPage.getAllObstacleIds();
      expect(obstacleIds.length).toBeGreaterThan(0);

      for (const obstacleId of obstacleIds) {
        const name = await categoriesPage.getObstacleName(obstacleId).textContent();
        const countText = await categoriesPage.getObstacleCount(obstacleId).textContent();

        expect(name).toBeTruthy();
        expect(countText).toMatch(/\d+x dilaporkan/);
      }
    });
  });
});
