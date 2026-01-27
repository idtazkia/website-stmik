import { test, expect } from '@playwright/test';
import { SettingsUsersPage } from './pages';

test.describe('Settings - User Management', () => {
  let usersPage: SettingsUsersPage;

  test.beforeEach(async ({ page }) => {
    usersPage = new SettingsUsersPage(page);
    // Login as admin before each test
    await usersPage.login('admin');
    await usersPage.goto();
    await usersPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display users page with stats and table', async () => {
      await usersPage.expectStatsVisible();
      await expect(usersPage.usersTable).toBeVisible();
    });

    test('should display user list from database', async () => {
      // Verify at least one user is displayed
      const userIds = await usersPage.getAllUserIds();
      expect(userIds.length).toBeGreaterThan(0);
    });

    test('should display correct stats totals', async () => {
      // Get counts from stats
      const totalText = await usersPage.statTotal.textContent();
      const adminText = await usersPage.statAdmin.textContent();
      const supervisorText = await usersPage.statSupervisor.textContent();
      const consultantText = await usersPage.statConsultant.textContent();
      const financeText = await usersPage.statFinance.textContent();
      const academicText = await usersPage.statAcademic.textContent();

      // Parse as numbers
      const total = parseInt(totalText || '0');
      const admin = parseInt(adminText || '0');
      const supervisor = parseInt(supervisorText || '0');
      const consultant = parseInt(consultantText || '0');
      const finance = parseInt(financeText || '0');
      const academic = parseInt(academicText || '0');

      // Total should equal sum of roles
      expect(total).toBe(admin + supervisor + consultant + finance + academic);
    });
  });

  test.describe('Role Change', () => {
    test('should change user role and verify UI update', async ({ page }) => {
      // Get first user ID that is not the logged-in user
      const userIds = await usersPage.getAllUserIds();
      expect(userIds.length).toBeGreaterThan(0);

      // Find a user to test with (skip test users)
      let testUserId: string | undefined;
      for (const id of userIds) {
        if (!id.startsWith('test-')) {
          testUserId = id;
          break;
        }
      }

      if (!testUserId) {
        test.skip();
        return;
      }

      // Get current role
      const roleSelect = usersPage.getUserRoleSelect(testUserId);
      const originalRole = await roleSelect.inputValue();

      // Change role to a different one
      const newRole = originalRole === 'consultant' ? 'supervisor' : 'consultant';
      await usersPage.changeUserRole(testUserId, newRole as 'admin' | 'supervisor' | 'consultant');

      // Verify UI shows new role
      await usersPage.expectUserRole(testUserId, newRole);

      // Reload page and verify persistence
      await page.reload();
      await usersPage.expectPageLoaded();
      await usersPage.expectUserRole(testUserId, newRole);

      // Restore original role
      await usersPage.changeUserRole(testUserId, originalRole as 'admin' | 'supervisor' | 'consultant');
      await usersPage.expectUserRole(testUserId, originalRole);
    });
  });

  test.describe('Status Toggle', () => {
    test('should toggle user active status and verify UI update', async ({ page }) => {
      // Get first user ID
      const userIds = await usersPage.getAllUserIds();
      expect(userIds.length).toBeGreaterThan(0);

      // Find a non-test user
      let testUserId: string | undefined;
      for (const id of userIds) {
        if (!id.startsWith('test-')) {
          testUserId = id;
          break;
        }
      }

      if (!testUserId) {
        test.skip();
        return;
      }

      // Get current status
      const statusButton = usersPage.getUserStatusToggle(testUserId);
      const originalStatus = (await statusButton.textContent())?.trim();
      const isCurrentlyActive = originalStatus === 'Aktif';

      // Toggle status
      await usersPage.toggleUserStatus(testUserId);

      // Verify UI shows new status
      await usersPage.expectUserStatus(testUserId, isCurrentlyActive ? 'inactive' : 'active');

      // Reload page and verify persistence
      await page.reload();
      await usersPage.expectPageLoaded();
      await usersPage.expectUserStatus(testUserId, isCurrentlyActive ? 'inactive' : 'active');

      // Restore original status
      await usersPage.toggleUserStatus(testUserId);
      await usersPage.expectUserStatus(testUserId, isCurrentlyActive ? 'active' : 'inactive');
    });
  });

  test.describe('Supervisor Assignment', () => {
    test('should assign supervisor to user and verify UI update', async ({ page }) => {
      // Get all user IDs
      const userIds = await usersPage.getAllUserIds();
      expect(userIds.length).toBeGreaterThan(0);

      // Find a consultant user to assign supervisor
      let consultantId: string | undefined;
      for (const id of userIds) {
        if (!id.startsWith('test-')) {
          const roleSelect = usersPage.getUserRoleSelect(id);
          const role = await roleSelect.inputValue();
          if (role === 'consultant') {
            consultantId = id;
            break;
          }
        }
      }

      if (!consultantId) {
        test.skip();
        return;
      }

      // Get supervisor select and find available supervisors
      const supervisorSelect = usersPage.getUserSupervisorSelect(consultantId);
      const options = await supervisorSelect.locator('option').all();

      // Skip first option (no supervisor)
      if (options.length <= 1) {
        test.skip();
        return;
      }

      // Get original supervisor
      const originalSupervisorId = await supervisorSelect.inputValue();

      // Select a different supervisor (or none if currently has one)
      let newSupervisorId: string;
      if (originalSupervisorId === '') {
        // Currently no supervisor, assign first available
        newSupervisorId = await options[1].getAttribute('value') || '';
      } else {
        // Currently has supervisor, clear it
        newSupervisorId = '';
      }

      // Change supervisor
      await usersPage.changeUserSupervisor(consultantId, newSupervisorId);

      // Verify UI shows new supervisor
      await usersPage.expectUserSupervisor(consultantId, newSupervisorId);

      // Reload page and verify persistence
      await page.reload();
      await usersPage.expectPageLoaded();
      await usersPage.expectUserSupervisor(consultantId, newSupervisorId);

      // Restore original supervisor
      await usersPage.changeUserSupervisor(consultantId, originalSupervisorId);
      await usersPage.expectUserSupervisor(consultantId, originalSupervisorId);
    });
  });
});
