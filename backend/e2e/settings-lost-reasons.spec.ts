import { test, expect } from '@playwright/test';
import { SettingsLostReasonsPage } from './pages/SettingsLostReasonsPage';

test.describe('Settings - Lost Reasons', () => {
  let page: SettingsLostReasonsPage;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new SettingsLostReasonsPage(browserPage);
    await page.login('admin');
    await page.goto(page.path);
    await page.expectPageLoaded();
  });

  test.describe('Page Display', () => {
    test('should display lost reasons page with seeded data', async () => {
      // Check page structure
      await expect(page.pageContainer).toBeVisible();
      await expect(page.lostReasonsSection).toBeVisible();
      await expect(page.lostReasonsList).toBeVisible();
      await expect(page.addLostReasonButton).toBeVisible();

      // Should have seeded lost reasons
      const ids = await page.getAllLostReasonIds();
      expect(ids.length).toBeGreaterThanOrEqual(8); // 8 seeded reasons
    });

    test('should display seeded lost reasons with correct structure', async () => {
      const ids = await page.getAllLostReasonIds();

      // Find "Tidak ada respon" reason
      let targetId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes('Tidak ada respon')) {
          targetId = id;
          break;
        }
      }

      expect(targetId).not.toBeNull();
      if (targetId) {
        await expect(page.getLostReasonName(targetId)).toBeVisible();
        await expect(page.getLostReasonStatusToggle(targetId)).toBeVisible();
        await expect(page.getLostReasonEditButton(targetId)).toBeVisible();
      }
    });

    test('should display description for reasons that have one', async () => {
      const ids = await page.getAllLostReasonIds();

      // Find reason with description
      let reasonWithDesc: string | null = null;
      for (const id of ids) {
        try {
          const descElement = page.getLostReasonDescription(id);
          if (await descElement.isVisible()) {
            reasonWithDesc = id;
            break;
          }
        } catch {
          // Element not found, continue
        }
      }

      // At least one seeded reason should have a description
      expect(reasonWithDesc).not.toBeNull();
    });
  });

  test.describe('Add Lost Reason', () => {
    test('should open add modal when clicking add button', async () => {
      await page.openAddModal();
      await expect(page.addLostReasonModal).toBeVisible();
      await expect(page.inputName).toBeVisible();
      await expect(page.inputDescription).toBeVisible();
      await expect(page.inputDisplayOrder).toBeVisible();
    });

    test('should add new lost reason with all fields', async () => {
      const uniqueName = `Test Reason ${Date.now()}`;
      const newReason = {
        name: uniqueName,
        description: 'Alasan untuk testing',
        displayOrder: 50
      };

      await page.addLostReason(newReason);

      // Modal should close
      await expect(page.addLostReasonModal).not.toBeVisible();

      // New lost reason should appear in list
      const ids = await page.getAllLostReasonIds();
      let newId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          newId = id;
          break;
        }
      }

      expect(newId).not.toBeNull();
      if (newId) {
        await expect(page.getLostReasonName(newId)).toContainText(uniqueName);
        await expect(page.getLostReasonDescription(newId)).toContainText('Alasan untuk testing');
      }
    });

    test('should add lost reason without description', async () => {
      const uniqueName = `No Desc Reason ${Date.now()}`;
      const newReason = {
        name: uniqueName,
        displayOrder: 60
      };

      await page.addLostReason(newReason);

      const ids = await page.getAllLostReasonIds();
      let newId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          newId = id;
          break;
        }
      }

      expect(newId).not.toBeNull();
      if (newId) {
        await expect(page.getLostReasonName(newId)).toContainText(uniqueName);
      }
    });
  });

  test.describe('Edit Lost Reason', () => {
    test('should open edit modal with current values', async () => {
      const ids = await page.getAllLostReasonIds();
      expect(ids.length).toBeGreaterThan(0);

      const id = ids[0];
      const originalName = await page.getLostReasonName(id).textContent();

      await page.openEditModal(id);
      await expect(page.getEditModal(id)).toBeVisible();

      // Check that input has current value
      const inputValue = await page.getEditInputName(id).inputValue();
      expect(inputValue).toBe(originalName?.trim());
    });

    test('should update lost reason name', async () => {
      // Create a new reason to edit (avoids conflicts with parallel tests)
      const uniqueName = `Edit Test ${Date.now()}`;
      await page.addLostReason({ name: uniqueName });

      // Reload to get proper edit modal
      await page.goto(page.path);
      await page.expectPageLoaded();

      const ids = await page.getAllLostReasonIds();
      let testId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          testId = id;
          break;
        }
      }

      expect(testId).not.toBeNull();
      if (testId) {
        await page.editLostReason(testId, {
          name: `${uniqueName} Updated`
        });

        await expect(page.getEditModal(testId).first()).not.toBeVisible();
        await expect(page.getLostReasonName(testId)).toContainText('Updated');
      }
    });

    test('should update lost reason description', async () => {
      // Create a new reason to edit
      const uniqueName = `Desc Edit Test ${Date.now()}`;
      await page.addLostReason({ name: uniqueName, description: 'Original' });

      await page.goto(page.path);
      await page.expectPageLoaded();

      const ids = await page.getAllLostReasonIds();
      let testId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          testId = id;
          break;
        }
      }

      expect(testId).not.toBeNull();
      if (testId) {
        await page.editLostReason(testId, {
          description: 'Updated description'
        });

        await expect(page.getLostReasonDescription(testId)).toContainText('Updated description');
      }
    });
  });

  test.describe('Toggle Status', () => {
    test('should toggle lost reason from active to inactive', async () => {
      // Create a new reason to toggle
      const uniqueName = `Toggle Test ${Date.now()}`;
      await page.addLostReason({ name: uniqueName });

      const ids = await page.getAllLostReasonIds();
      let testId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          testId = id;
          break;
        }
      }

      expect(testId).not.toBeNull();
      if (testId) {
        // Should be active initially
        await page.expectLostReasonActive(testId);

        // Toggle to inactive
        await page.toggleLostReasonStatus(testId);
        await page.expectLostReasonInactive(testId);

        // Toggle back to active
        await page.toggleLostReasonStatus(testId);
        await page.expectLostReasonActive(testId);
      }
    });

    test('should persist status after page reload', async () => {
      // Create a new reason
      const uniqueName = `Persist Toggle ${Date.now()}`;
      await page.addLostReason({ name: uniqueName });

      const ids = await page.getAllLostReasonIds();
      let testId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          testId = id;
          break;
        }
      }

      expect(testId).not.toBeNull();
      if (testId) {
        // Toggle to inactive
        await page.toggleLostReasonStatus(testId);
        await page.expectLostReasonInactive(testId);

        // Reload page
        await page.goto(page.path);
        await page.expectPageLoaded();

        // Find the same reason again
        const newIds = await page.getAllLostReasonIds();
        let persistedId: string | null = null;
        for (const id of newIds) {
          const nameText = await page.getLostReasonName(id).textContent();
          if (nameText?.includes(uniqueName)) {
            persistedId = id;
            break;
          }
        }

        expect(persistedId).not.toBeNull();
        if (persistedId) {
          await page.expectLostReasonInactive(persistedId);
        }
      }
    });
  });

  test.describe('Database Persistence', () => {
    test('should persist new lost reason after page reload', async () => {
      const uniqueName = `Persist Test ${Date.now()}`;
      await page.addLostReason({
        name: uniqueName,
        description: 'Testing persistence'
      });

      // Reload page
      await page.goto(page.path);
      await page.expectPageLoaded();

      // Find the reason
      const ids = await page.getAllLostReasonIds();
      let foundId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          foundId = id;
          break;
        }
      }

      expect(foundId).not.toBeNull();
      if (foundId) {
        await expect(page.getLostReasonName(foundId)).toContainText(uniqueName);
        await expect(page.getLostReasonDescription(foundId)).toContainText('Testing persistence');
      }
    });

    test('should persist edited lost reason after page reload', async () => {
      // Create then edit
      const uniqueName = `Edit Persist ${Date.now()}`;
      await page.addLostReason({ name: uniqueName });

      await page.goto(page.path);
      await page.expectPageLoaded();

      const ids = await page.getAllLostReasonIds();
      let testId: string | null = null;
      for (const id of ids) {
        const nameText = await page.getLostReasonName(id).textContent();
        if (nameText?.includes(uniqueName)) {
          testId = id;
          break;
        }
      }

      expect(testId).not.toBeNull();
      if (testId) {
        await page.editLostReason(testId, {
          name: `${uniqueName} Final`,
          description: 'Edited description'
        });

        // Reload page
        await page.goto(page.path);
        await page.expectPageLoaded();

        // Find the reason again
        const newIds = await page.getAllLostReasonIds();
        let editedId: string | null = null;
        for (const id of newIds) {
          const nameText = await page.getLostReasonName(id).textContent();
          if (nameText?.includes(`${uniqueName} Final`)) {
            editedId = id;
            break;
          }
        }

        expect(editedId).not.toBeNull();
        if (editedId) {
          await expect(page.getLostReasonName(editedId)).toContainText('Final');
          await expect(page.getLostReasonDescription(editedId)).toContainText('Edited description');
        }
      }
    });
  });
});
