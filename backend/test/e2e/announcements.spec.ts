import { test, expect } from '@playwright/test';
import { AnnouncementsPage } from './pages';

test.describe('Admin Announcements Management', () => {
  let announcementsPage: AnnouncementsPage;

  test.beforeEach(async ({ page }) => {
    announcementsPage = new AnnouncementsPage(page);
    // Login as admin before each test
    await announcementsPage.login('admin');
    await announcementsPage.goto();
    await announcementsPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display announcements page with section', async () => {
      await expect(announcementsPage.announcementsSection).toBeVisible();
    });

    test('should display add announcement button', async () => {
      await expect(announcementsPage.addAnnouncementButton).toBeVisible();
    });
  });

  test.describe('Announcement CRUD', () => {
    // Run CRUD tests serially to avoid race conditions
    test.describe.configure({ mode: 'serial' });

    test('should open add announcement modal', async () => {
      await announcementsPage.openAddAnnouncementModal();
      await expect(announcementsPage.addAnnouncementModal).toBeVisible();
      await expect(announcementsPage.inputTitle).toBeVisible();
      await expect(announcementsPage.inputContent).toBeVisible();
      await expect(announcementsPage.selectTargetStatus).toBeVisible();
      await expect(announcementsPage.selectTargetProdi).toBeVisible();
    });

    test('should add new announcement via HTMX', async ({ page }) => {
      // Generate unique title
      const timestamp = Date.now().toString().slice(-6);
      const newTitle = `Test Announcement ${timestamp}`;
      const newContent = `This is test content for announcement ${timestamp}`;

      // Get current announcement count
      const announcementIdsBefore = await announcementsPage.getAllAnnouncementIds();
      const countBefore = announcementIdsBefore.length;

      // Add new announcement
      await announcementsPage.addAnnouncement(newTitle, newContent);

      // Verify new announcement appears
      const announcementIdsAfter = await announcementsPage.getAllAnnouncementIds();
      expect(announcementIdsAfter.length).toBe(countBefore + 1);

      // Find the new announcement
      const newAnnouncementId = announcementIdsAfter.find(id => !announcementIdsBefore.includes(id));
      expect(newAnnouncementId).toBeTruthy();

      if (newAnnouncementId) {
        await announcementsPage.expectAnnouncementTitleValue(newAnnouncementId, newTitle);
        // New announcements should be draft by default
        await announcementsPage.expectAnnouncementStatusText(newAnnouncementId, 'Draft');

        // Reload and verify persistence
        await page.reload();
        await announcementsPage.expectPageLoaded();
        await announcementsPage.expectAnnouncementDisplayed(newAnnouncementId);
      }
    });

    test('should publish and unpublish announcement via HTMX', async () => {
      // First create a new announcement to test with
      const timestamp = Date.now().toString().slice(-6);
      const newTitle = `Publish Test ${timestamp}`;
      const newContent = `Content for publish test ${timestamp}`;

      const announcementIdsBefore = await announcementsPage.getAllAnnouncementIds();
      await announcementsPage.addAnnouncement(newTitle, newContent);

      const announcementIdsAfter = await announcementsPage.getAllAnnouncementIds();
      const newAnnouncementId = announcementIdsAfter.find(id => !announcementIdsBefore.includes(id));
      expect(newAnnouncementId).toBeTruthy();

      if (newAnnouncementId) {
        // Verify initial draft status
        await announcementsPage.expectAnnouncementStatusText(newAnnouncementId, 'Draft');

        // Publish announcement
        await announcementsPage.publishAnnouncement(newAnnouncementId);
        await announcementsPage.expectAnnouncementStatusText(newAnnouncementId, 'Terbit');

        // Unpublish announcement
        await announcementsPage.unpublishAnnouncement(newAnnouncementId);
        await announcementsPage.expectAnnouncementStatusText(newAnnouncementId, 'Draft');
      }
    });

    test('should edit announcement via HTMX', async ({ page }) => {
      // First create a new announcement to test with
      const timestamp = Date.now().toString().slice(-6);
      const originalTitle = `Edit Test ${timestamp}`;
      const originalContent = `Original content ${timestamp}`;

      const announcementIdsBefore = await announcementsPage.getAllAnnouncementIds();
      await announcementsPage.addAnnouncement(originalTitle, originalContent);

      const announcementIdsAfter = await announcementsPage.getAllAnnouncementIds();
      const newAnnouncementId = announcementIdsAfter.find(id => !announcementIdsBefore.includes(id));
      expect(newAnnouncementId).toBeTruthy();

      if (newAnnouncementId) {
        // Edit the announcement
        const updatedTitle = `${originalTitle} Updated`;
        const updatedContent = `Updated content ${timestamp}`;
        await announcementsPage.editAnnouncement(newAnnouncementId, updatedTitle, updatedContent);

        // Verify title changed
        await announcementsPage.expectAnnouncementTitleValue(newAnnouncementId, updatedTitle);

        // Reload and verify persistence
        await page.reload();
        await announcementsPage.expectPageLoaded();
        await announcementsPage.expectAnnouncementTitleValue(newAnnouncementId, updatedTitle);
      }
    });

    test('should delete announcement via HTMX', async () => {
      // First create a new announcement to test with
      const timestamp = Date.now().toString().slice(-6);
      const newTitle = `Delete Test ${timestamp}`;
      const newContent = `Content to delete ${timestamp}`;

      const announcementIdsBefore = await announcementsPage.getAllAnnouncementIds();
      await announcementsPage.addAnnouncement(newTitle, newContent);

      const announcementIdsAfter = await announcementsPage.getAllAnnouncementIds();
      const newAnnouncementId = announcementIdsAfter.find(id => !announcementIdsBefore.includes(id));
      expect(newAnnouncementId).toBeTruthy();

      if (newAnnouncementId) {
        // Delete the announcement
        await announcementsPage.deleteAnnouncement(newAnnouncementId);

        // Verify announcement is removed
        await expect(announcementsPage.getAnnouncementRow(newAnnouncementId)).not.toBeVisible();
      }
    });
  });

  test.describe('Announcement Targeting', () => {
    test('should add announcement with target status', async () => {
      const timestamp = Date.now().toString().slice(-6);
      const newTitle = `Targeted Status ${timestamp}`;
      const newContent = `Content for targeted announcement ${timestamp}`;

      const announcementIdsBefore = await announcementsPage.getAllAnnouncementIds();
      await announcementsPage.addAnnouncement(newTitle, newContent, 'registered');

      const announcementIdsAfter = await announcementsPage.getAllAnnouncementIds();
      expect(announcementIdsAfter.length).toBe(announcementIdsBefore.length + 1);
    });
  });
});

test.describe('Admin Navigation to Announcements', () => {
  test('should navigate to announcements from sidebar', async ({ page }) => {
    // Login as admin
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);

    // Click announcements link in sidebar
    await page.getByTestId('nav-announcements').click();
    await page.waitForURL(/\/admin\/announcements/);

    // Verify page loaded
    await expect(page.getByTestId('settings-announcements-page')).toBeVisible();
  });
});
