import { Locator, expect, Page } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsUsersPage extends BasePage {
  readonly path = '/admin/settings/users';

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-users-page');
  }

  get usersTable(): Locator {
    return this.page.getByTestId('users-table');
  }

  get usersList(): Locator {
    return this.page.getByTestId('users-list');
  }

  // Stats elements
  get statTotal(): Locator {
    return this.page.getByTestId('stat-total-value');
  }

  get statAdmin(): Locator {
    return this.page.getByTestId('stat-admin-value');
  }

  get statSupervisor(): Locator {
    return this.page.getByTestId('stat-supervisor-value');
  }

  get statConsultant(): Locator {
    return this.page.getByTestId('stat-consultant-value');
  }

  // User row elements by ID
  getUserRow(userId: string): Locator {
    return this.page.getByTestId(`user-row-${userId}`);
  }

  getUserName(userId: string): Locator {
    return this.page.getByTestId(`user-name-${userId}`);
  }

  getUserEmail(userId: string): Locator {
    return this.page.getByTestId(`user-email-${userId}`);
  }

  getUserRoleSelect(userId: string): Locator {
    return this.page.getByTestId(`user-role-select-${userId}`);
  }

  getUserSupervisorSelect(userId: string): Locator {
    return this.page.getByTestId(`user-supervisor-select-${userId}`);
  }

  getUserStatusToggle(userId: string): Locator {
    return this.page.getByTestId(`user-status-toggle-${userId}`);
  }

  getUserLastLogin(userId: string): Locator {
    return this.page.getByTestId(`user-last-login-${userId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.usersTable).toBeVisible();
  }

  async expectStatsVisible(): Promise<void> {
    await expect(this.statTotal).toBeVisible();
    await expect(this.statAdmin).toBeVisible();
    await expect(this.statSupervisor).toBeVisible();
    await expect(this.statConsultant).toBeVisible();
  }

  // Actions
  async changeUserRole(userId: string, newRole: 'admin' | 'supervisor' | 'consultant'): Promise<void> {
    const select = this.getUserRoleSelect(userId);
    // Start waiting for response before triggering the action
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/users/${userId}/role`) &&
      response.status() === 200
    );
    await select.selectOption(newRole);
    await responsePromise;
  }

  async changeUserSupervisor(userId: string, supervisorId: string): Promise<void> {
    const select = this.getUserSupervisorSelect(userId);
    // Start waiting for response before triggering the action
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/users/${userId}/supervisor`) &&
      response.status() === 200
    );
    await select.selectOption(supervisorId);
    await responsePromise;
  }

  async toggleUserStatus(userId: string): Promise<void> {
    const button = this.getUserStatusToggle(userId);
    // Start waiting for response before triggering the action
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/users/${userId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  // Verification helpers
  async expectUserRole(userId: string, expectedRole: string): Promise<void> {
    const select = this.getUserRoleSelect(userId);
    await expect(select).toHaveValue(expectedRole);
  }

  async expectUserStatus(userId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getUserStatusToggle(userId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  async expectUserSupervisor(userId: string, expectedSupervisorId: string): Promise<void> {
    const select = this.getUserSupervisorSelect(userId);
    await expect(select).toHaveValue(expectedSupervisorId);
  }

  // Get all user IDs from the table
  async getAllUserIds(): Promise<string[]> {
    const rows = this.usersList.locator('tr[data-testid^="user-row-"]');
    const count = await rows.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await rows.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('user-row-', ''));
      }
    }
    return ids;
  }
}
