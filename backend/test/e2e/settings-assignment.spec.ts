import { test, expect } from '@playwright/test';
import { SettingsAssignmentPage } from './pages';

test.describe('Settings - Assignment Algorithm', () => {
  let assignmentPage: SettingsAssignmentPage;

  test.beforeEach(async ({ page }) => {
    assignmentPage = new SettingsAssignmentPage(page);
    await assignmentPage.login('admin');
    await assignmentPage.goto();
    await assignmentPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display assignment settings page', async () => {
      await expect(assignmentPage.pageContainer).toBeVisible();
      await expect(assignmentPage.algorithmsSection).toBeVisible();
    });

    test('should display algorithms list', async () => {
      await expect(assignmentPage.algorithmsList).toBeVisible();
    });

    test('should display seeded algorithms', async () => {
      const algorithmIds = await assignmentPage.getAllAlgorithmIds();
      expect(algorithmIds.length).toBeGreaterThan(0);
    });

    test('should have exactly one active algorithm', async () => {
      const activeId = await assignmentPage.getActiveAlgorithmId();
      expect(activeId).not.toBeNull();
    });
  });

  test.describe('Algorithm Display', () => {
    test('should display algorithm details', async () => {
      const algorithmIds = await assignmentPage.getAllAlgorithmIds();
      if (algorithmIds.length === 0) {
        test.skip();
        return;
      }

      for (const id of algorithmIds) {
        await assignmentPage.expectAlgorithmDisplayed(id);
        await expect(assignmentPage.getAlgorithmCode(id)).toBeVisible();
        await expect(assignmentPage.getAlgorithmDescription(id)).toBeVisible();
      }
    });

    test('should display activate button for inactive algorithms', async () => {
      const algorithmIds = await assignmentPage.getAllAlgorithmIds();
      const activeId = await assignmentPage.getActiveAlgorithmId();

      for (const id of algorithmIds) {
        if (id === activeId) {
          await assignmentPage.expectAlgorithmActive(id);
        } else {
          await assignmentPage.expectAlgorithmInactive(id);
        }
      }
    });
  });

  test.describe('Algorithm Switching', () => {
    test.describe.configure({ mode: 'serial' });

    test('should switch active algorithm via HTMX', async ({ page }) => {
      const algorithmIds = await assignmentPage.getAllAlgorithmIds();
      if (algorithmIds.length < 2) {
        test.skip();
        return;
      }

      // Get current active algorithm
      const activeIdBefore = await assignmentPage.getActiveAlgorithmId();
      expect(activeIdBefore).not.toBeNull();

      // Find an inactive algorithm
      const inactiveId = algorithmIds.find(id => id !== activeIdBefore);
      expect(inactiveId).toBeDefined();

      // Activate the inactive algorithm
      await assignmentPage.activateAlgorithm(inactiveId!);

      // Verify the new algorithm is active
      await assignmentPage.expectAlgorithmActive(inactiveId!);

      // Verify the old algorithm is now inactive
      await assignmentPage.expectAlgorithmInactive(activeIdBefore!);

      // Reload and verify persistence
      await page.reload();
      await assignmentPage.expectPageLoaded();

      const activeIdAfter = await assignmentPage.getActiveAlgorithmId();
      expect(activeIdAfter).toBe(inactiveId);
    });

    test('should restore original algorithm', async ({ page }) => {
      // This test switches back to round_robin to restore state
      const algorithmIds = await assignmentPage.getAllAlgorithmIds();
      if (algorithmIds.length === 0) {
        test.skip();
        return;
      }

      // Find round_robin algorithm by checking code
      let roundRobinId: string | null = null;
      for (const id of algorithmIds) {
        const code = await assignmentPage.getAlgorithmCode(id).textContent();
        if (code === 'round_robin') {
          roundRobinId = id;
          break;
        }
      }

      if (!roundRobinId) {
        test.skip();
        return;
      }

      const activeId = await assignmentPage.getActiveAlgorithmId();
      if (activeId === roundRobinId) {
        // Already active, no need to switch
        return;
      }

      // Activate round_robin
      await assignmentPage.activateAlgorithm(roundRobinId);
      await assignmentPage.expectAlgorithmActive(roundRobinId);

      // Verify persistence
      await page.reload();
      await assignmentPage.expectPageLoaded();
      await assignmentPage.expectAlgorithmActive(roundRobinId);
    });
  });
});
