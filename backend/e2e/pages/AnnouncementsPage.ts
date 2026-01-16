import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class AnnouncementsPage extends BasePage {
  readonly path = '/admin/announcements';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    await this.page.waitForURL(/\/admin\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-announcements-page');
  }

  get announcementsSection(): Locator {
    return this.page.getByTestId('announcements-section');
  }

  get announcementsList(): Locator {
    return this.page.getByTestId('announcements-list');
  }

  get addAnnouncementButton(): Locator {
    return this.page.getByTestId('add-announcement-button');
  }

  // Add Announcement Modal elements
  get addAnnouncementModal(): Locator {
    return this.page.getByTestId('add-announcement-modal');
  }

  get inputTitle(): Locator {
    return this.page.getByTestId('input-title');
  }

  get inputContent(): Locator {
    return this.page.getByTestId('input-content');
  }

  get selectTargetStatus(): Locator {
    return this.page.getByTestId('select-target-status');
  }

  get selectTargetProdi(): Locator {
    return this.page.getByTestId('select-target-prodi');
  }

  get submitAddAnnouncementButton(): Locator {
    return this.page.getByTestId('submit-add-announcement');
  }

  // Edit Announcement Modal elements
  get editAnnouncementModal(): Locator {
    return this.page.getByTestId('edit-announcement-modal');
  }

  get submitEditAnnouncementButton(): Locator {
    return this.page.getByTestId('submit-edit-announcement');
  }

  // Announcement row elements by ID
  getAnnouncementRow(announcementId: string): Locator {
    return this.page.getByTestId(`announcement-row-${announcementId}`);
  }

  getAnnouncementTitle(announcementId: string): Locator {
    return this.getAnnouncementRow(announcementId).getByTestId('announcement-title');
  }

  getAnnouncementStatus(announcementId: string): Locator {
    return this.getAnnouncementRow(announcementId).getByTestId('announcement-status');
  }

  getAnnouncementReadCount(announcementId: string): Locator {
    return this.getAnnouncementRow(announcementId).getByTestId('announcement-read-count');
  }

  getEditButton(announcementId: string): Locator {
    return this.page.getByTestId(`edit-announcement-${announcementId}`);
  }

  getPublishButton(announcementId: string): Locator {
    return this.page.getByTestId(`publish-announcement-${announcementId}`);
  }

  getUnpublishButton(announcementId: string): Locator {
    return this.page.getByTestId(`unpublish-announcement-${announcementId}`);
  }

  getDeleteButton(announcementId: string): Locator {
    return this.page.getByTestId(`delete-announcement-${announcementId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.announcementsSection).toBeVisible();
  }

  // Get all announcement IDs
  async getAllAnnouncementIds(): Promise<string[]> {
    const rows = this.announcementsList.locator('tr[data-testid^="announcement-row-"]');
    const count = await rows.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await rows.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('announcement-row-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectAnnouncementDisplayed(announcementId: string): Promise<void> {
    await expect(this.getAnnouncementRow(announcementId)).toBeVisible();
    await expect(this.getAnnouncementTitle(announcementId)).toBeVisible();
  }

  async expectAnnouncementTitleValue(announcementId: string, expectedTitle: string): Promise<void> {
    const title = await this.getAnnouncementTitle(announcementId).textContent();
    expect(title?.trim()).toBe(expectedTitle);
  }

  async expectAnnouncementStatusText(announcementId: string, expectedStatus: 'Terbit' | 'Draft'): Promise<void> {
    const statusText = await this.getAnnouncementStatus(announcementId).textContent();
    expect(statusText?.trim()).toBe(expectedStatus);
  }

  // Announcement Actions
  async openAddAnnouncementModal(): Promise<void> {
    await this.addAnnouncementButton.click();
    await expect(this.addAnnouncementModal).toBeVisible();
  }

  async addAnnouncement(
    title: string,
    content: string,
    targetStatus?: string,
    targetProdiId?: string
  ): Promise<void> {
    await this.openAddAnnouncementModal();
    await this.inputTitle.fill(title);
    await this.inputContent.fill(content);

    if (targetStatus) {
      await this.selectTargetStatus.selectOption(targetStatus);
    }

    if (targetProdiId) {
      await this.selectTargetProdi.selectOption(targetProdiId);
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/announcements') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/publish') &&
      !response.url().includes('/unpublish') &&
      response.status() === 200
    );
    await this.submitAddAnnouncementButton.click();
    await responsePromise;
  }

  async openEditAnnouncementModal(announcementId: string): Promise<void> {
    // Trigger HTMX to load edit form
    await this.getEditButton(announcementId).click();
    await expect(this.editAnnouncementModal).toBeVisible();
    // Wait for content to be loaded via HTMX
    await expect(this.inputTitle).toBeVisible();
  }

  async editAnnouncement(
    announcementId: string,
    title: string,
    content: string,
    targetStatus?: string,
    targetProdiId?: string
  ): Promise<void> {
    await this.openEditAnnouncementModal(announcementId);
    await this.inputTitle.fill(title);
    await this.inputContent.fill(content);

    if (targetStatus !== undefined) {
      await this.selectTargetStatus.selectOption(targetStatus);
    }

    if (targetProdiId !== undefined) {
      await this.selectTargetProdi.selectOption(targetProdiId);
    }

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/announcements/${announcementId}`) &&
      response.request().method() === 'PUT' &&
      response.status() === 200
    );
    await this.submitEditAnnouncementButton.click();
    await responsePromise;
  }

  async publishAnnouncement(announcementId: string): Promise<void> {
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/announcements/${announcementId}/publish`) &&
      response.status() === 200
    );
    await this.getPublishButton(announcementId).click();
    await responsePromise;
  }

  async unpublishAnnouncement(announcementId: string): Promise<void> {
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/announcements/${announcementId}/unpublish`) &&
      response.status() === 200
    );
    await this.getUnpublishButton(announcementId).click();
    await responsePromise;
  }

  async deleteAnnouncement(announcementId: string): Promise<void> {
    // Accept dialog confirmation
    this.page.once('dialog', dialog => dialog.accept());

    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/announcements/${announcementId}`) &&
      response.request().method() === 'DELETE' &&
      response.status() === 200
    );
    await this.getDeleteButton(announcementId).click();
    await responsePromise;
  }
}

export class PortalAnnouncementsPage extends BasePage {
  readonly path = '/portal/announcements';

  // Login as candidate
  async loginAsCandidate(): Promise<void> {
    // Use test login endpoint for candidate
    await this.page.goto('/test/login/candidate');
    await this.page.waitForURL(/\/portal\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('announcements-page');
  }

  get announcementsList(): Locator {
    return this.page.getByTestId('announcements-list');
  }

  get unreadCount(): Locator {
    return this.page.getByTestId('unread-count');
  }

  // Announcement card elements by ID
  getAnnouncementCard(announcementId: string): Locator {
    return this.page.getByTestId(`announcement-card-${announcementId}`);
  }

  getAnnouncementTitle(announcementId: string): Locator {
    return this.getAnnouncementCard(announcementId).getByTestId('announcement-title');
  }

  getAnnouncementContent(announcementId: string): Locator {
    return this.getAnnouncementCard(announcementId).getByTestId('announcement-content');
  }

  getUnreadBadge(announcementId: string): Locator {
    return this.getAnnouncementCard(announcementId).getByTestId('unread-badge');
  }

  getMarkReadButton(announcementId: string): Locator {
    return this.getAnnouncementCard(announcementId).getByTestId('mark-read-button');
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
  }

  // Get all announcement IDs
  async getAllAnnouncementIds(): Promise<string[]> {
    const cards = this.announcementsList.locator('div[data-testid^="announcement-card-"]');
    const count = await cards.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await cards.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('announcement-card-', ''));
      }
    }
    return ids;
  }

  // Mark announcement as read
  async markAsRead(announcementId: string): Promise<void> {
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/portal/announcements/${announcementId}/read`) &&
      response.status() === 200
    );
    await this.getMarkReadButton(announcementId).click();
    await responsePromise;
  }

  // Verification helpers
  async expectAnnouncementDisplayed(announcementId: string): Promise<void> {
    await expect(this.getAnnouncementCard(announcementId)).toBeVisible();
  }

  async expectAnnouncementUnread(announcementId: string): Promise<void> {
    await expect(this.getUnreadBadge(announcementId)).toBeVisible();
    await expect(this.getMarkReadButton(announcementId)).toBeVisible();
  }

  async expectAnnouncementRead(announcementId: string): Promise<void> {
    await expect(this.getUnreadBadge(announcementId)).not.toBeVisible();
    await expect(this.getMarkReadButton(announcementId)).not.toBeVisible();
  }
}
