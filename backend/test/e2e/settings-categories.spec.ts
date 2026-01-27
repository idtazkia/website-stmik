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

  test.describe('Category CRUD', () => {
    // Run CRUD tests serially to avoid race conditions
    test.describe.configure({ mode: 'serial' });

    test('should open add category modal', async () => {
      await categoriesPage.openAddCategoryModal();
      await expect(categoriesPage.addCategoryModal).toBeVisible();
      await expect(categoriesPage.inputCategoryName).toBeVisible();
      await expect(categoriesPage.inputCategorySentiment).toBeVisible();
    });

    test('should add new category via HTMX', async ({ page }) => {
      // Generate unique name
      const timestamp = Date.now().toString().slice(-4);
      const newName = `Test Category ${timestamp}`;

      // Get current category count
      const categoryIdsBefore = await categoriesPage.getAllCategoryIds();
      const countBefore = categoryIdsBefore.length;

      // Add new category
      await categoriesPage.addCategory(newName, 'positive');

      // Verify new category appears
      const categoryIdsAfter = await categoriesPage.getAllCategoryIds();
      expect(categoryIdsAfter.length).toBe(countBefore + 1);

      // Find the new category
      const newCategoryId = categoryIdsAfter.find(id => !categoryIdsBefore.includes(id));
      expect(newCategoryId).toBeTruthy();

      if (newCategoryId) {
        await categoriesPage.expectCategoryNameValue(newCategoryId, newName);

        // Reload and verify persistence
        await page.reload();
        await categoriesPage.expectPageLoaded();
        await categoriesPage.expectCategoryDisplayed(newCategoryId);
      }
    });

    test('should toggle category status via HTMX', async () => {
      const categoryIds = await categoriesPage.getAllCategoryIds();
      expect(categoryIds.length).toBeGreaterThan(0);

      const categoryId = categoryIds[0];

      // Get current status
      const statusBefore = await categoriesPage.getCategoryStatusToggle(categoryId).textContent();
      const isActiveBefore = statusBefore?.trim() === 'Aktif';

      // Toggle status
      await categoriesPage.toggleCategoryStatus(categoryId);

      // Verify status changed
      if (isActiveBefore) {
        await categoriesPage.expectCategoryStatus(categoryId, 'inactive');
      } else {
        await categoriesPage.expectCategoryStatus(categoryId, 'active');
      }

      // Toggle back to restore original state
      await categoriesPage.toggleCategoryStatus(categoryId);
    });

    test('should edit category name via HTMX', async ({ page }) => {
      const categoryIds = await categoriesPage.getAllCategoryIds();
      expect(categoryIds.length).toBeGreaterThan(0);

      const categoryId = categoryIds[0];

      // Get current name
      const currentName = await categoriesPage.getCategoryName(categoryId).textContent();

      // Edit with new name
      const newName = `${currentName?.trim()} Updated`;
      await categoriesPage.editCategory(categoryId, newName, 'positive');

      // Verify name changed
      await categoriesPage.expectCategoryNameValue(categoryId, newName);

      // Reload and verify persistence
      await page.reload();
      await categoriesPage.expectPageLoaded();
      await categoriesPage.expectCategoryNameValue(categoryId, newName);

      // Restore original name
      await categoriesPage.editCategory(categoryId, currentName?.trim() || '', 'positive');
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

  test.describe('Obstacle CRUD', () => {
    // Run CRUD tests serially to avoid race conditions
    test.describe.configure({ mode: 'serial' });

    test('should open add obstacle modal', async () => {
      await categoriesPage.openAddObstacleModal();
      await expect(categoriesPage.addObstacleModal).toBeVisible();
      await expect(categoriesPage.inputObstacleName).toBeVisible();
    });

    test('should add new obstacle via HTMX', async ({ page }) => {
      // Generate unique name
      const timestamp = Date.now().toString().slice(-4);
      const newName = `Test Obstacle ${timestamp}`;

      // Get current obstacle count
      const obstacleIdsBefore = await categoriesPage.getAllObstacleIds();
      const countBefore = obstacleIdsBefore.length;

      // Add new obstacle
      await categoriesPage.addObstacle(newName);

      // Verify new obstacle appears
      const obstacleIdsAfter = await categoriesPage.getAllObstacleIds();
      expect(obstacleIdsAfter.length).toBe(countBefore + 1);

      // Find the new obstacle
      const newObstacleId = obstacleIdsAfter.find(id => !obstacleIdsBefore.includes(id));
      expect(newObstacleId).toBeTruthy();

      if (newObstacleId) {
        await categoriesPage.expectObstacleNameValue(newObstacleId, newName);

        // Reload and verify persistence
        await page.reload();
        await categoriesPage.expectPageLoaded();
        await categoriesPage.expectObstacleDisplayed(newObstacleId);
      }
    });

    test('should toggle obstacle status via HTMX', async () => {
      const obstacleIds = await categoriesPage.getAllObstacleIds();
      expect(obstacleIds.length).toBeGreaterThan(0);

      const obstacleId = obstacleIds[0];

      // Get current status
      const statusBefore = await categoriesPage.getObstacleStatusToggle(obstacleId).textContent();
      const isActiveBefore = statusBefore?.trim() === 'Aktif';

      // Toggle status
      await categoriesPage.toggleObstacleStatus(obstacleId);

      // Verify status changed
      if (isActiveBefore) {
        await categoriesPage.expectObstacleStatus(obstacleId, 'inactive');
      } else {
        await categoriesPage.expectObstacleStatus(obstacleId, 'active');
      }

      // Toggle back to restore original state
      await categoriesPage.toggleObstacleStatus(obstacleId);
    });

    test('should edit obstacle name via HTMX', async ({ page }) => {
      const obstacleIds = await categoriesPage.getAllObstacleIds();
      expect(obstacleIds.length).toBeGreaterThan(0);

      const obstacleId = obstacleIds[0];

      // Get current name
      const currentName = await categoriesPage.getObstacleName(obstacleId).textContent();

      // Edit with new name
      const newName = `${currentName?.trim()} Updated`;
      await categoriesPage.editObstacle(obstacleId, newName);

      // Verify name changed
      await categoriesPage.expectObstacleNameValue(obstacleId, newName);

      // Reload and verify persistence
      await page.reload();
      await categoriesPage.expectPageLoaded();
      await categoriesPage.expectObstacleNameValue(obstacleId, newName);

      // Restore original name
      await categoriesPage.editObstacle(obstacleId, currentName?.trim() || '');
    });
  });
});
