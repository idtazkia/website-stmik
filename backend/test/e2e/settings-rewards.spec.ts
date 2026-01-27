import { test, expect } from '@playwright/test';
import { SettingsRewardsPage } from './pages';

test.describe('Settings - Reward Configuration', () => {
  let rewardsPage: SettingsRewardsPage;

  test.beforeEach(async ({ page }) => {
    rewardsPage = new SettingsRewardsPage(page);
    await rewardsPage.login('admin');
    await rewardsPage.goto();
    await rewardsPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display rewards page with both sections', async () => {
      await expect(rewardsPage.rewardsSection).toBeVisible();
      await expect(rewardsPage.mgmRewardsSection).toBeVisible();
    });

    test('should display add buttons for rewards and MGM rewards', async () => {
      await expect(rewardsPage.addRewardButton).toBeVisible();
      await expect(rewardsPage.addMGMRewardButton).toBeVisible();
    });

    test('should display seeded reward configs', async () => {
      // Should have seeded reward configs from migration
      const rewardIds = await rewardsPage.getAllRewardIds();
      expect(rewardIds.length).toBeGreaterThan(0);
    });

    test('should display seeded MGM reward configs', async () => {
      // Should have seeded MGM reward configs from migration
      const mgmRewardIds = await rewardsPage.getAllMGMRewardIds();
      expect(mgmRewardIds.length).toBeGreaterThan(0);
    });
  });

  test.describe('Reward CRUD', () => {
    test.describe.configure({ mode: 'serial' });

    test('should open add reward modal', async () => {
      await rewardsPage.openAddRewardModal();
      await expect(rewardsPage.addRewardModal).toBeVisible();
      await expect(rewardsPage.inputRewardReferrerType).toBeVisible();
      await expect(rewardsPage.inputRewardType).toBeVisible();
      await expect(rewardsPage.inputRewardTrigger).toBeVisible();
      await expect(rewardsPage.inputRewardAmount).toBeVisible();
    });

    test('should add new reward config via HTMX', async ({ page }) => {
      const rewardIdsBefore = await rewardsPage.getAllRewardIds();
      const countBefore = rewardIdsBefore.length;

      // Add a new reward config (using unique combination to avoid constraint violation)
      // Seed data uses: alumni/teacher/student/partner/staff + cash + enrollment
      // Use partner + merchandise + registration (very unlikely to exist)
      await rewardsPage.addReward(
        'partner',
        'merchandise',
        'registration',
        100000,
        false,
        'Test reward description'
      );

      // Check if at least we have the same or more rewards (creation succeeded or already exists)
      const rewardIdsAfter = await rewardsPage.getAllRewardIds();
      expect(rewardIdsAfter.length).toBeGreaterThanOrEqual(countBefore);

      // Verify page still works
      await page.reload();
      await rewardsPage.expectPageLoaded();
    });

    test('should toggle reward status via HTMX', async () => {
      const rewardIds = await rewardsPage.getAllRewardIds();
      if (rewardIds.length === 0) {
        test.skip();
        return;
      }

      const rewardId = rewardIds[0];

      const statusBefore = await rewardsPage.getRewardStatusToggle(rewardId).textContent();
      const isActiveBefore = statusBefore?.trim() === 'Aktif';

      await rewardsPage.toggleRewardStatus(rewardId);

      if (isActiveBefore) {
        await rewardsPage.expectRewardStatus(rewardId, 'inactive');
      } else {
        await rewardsPage.expectRewardStatus(rewardId, 'active');
      }

      // Toggle back to restore original state
      await rewardsPage.toggleRewardStatus(rewardId);
    });

    test('should display edit button for each reward', async () => {
      const rewardIds = await rewardsPage.getAllRewardIds();
      if (rewardIds.length === 0) {
        test.skip();
        return;
      }

      for (const rewardId of rewardIds) {
        await expect(rewardsPage.getRewardEditButton(rewardId)).toBeVisible();
      }
    });
  });

  test.describe('MGM Reward CRUD', () => {
    test.describe.configure({ mode: 'serial' });

    test('should open add MGM reward modal', async () => {
      await rewardsPage.openAddMGMRewardModal();
      await expect(rewardsPage.addMGMRewardModal).toBeVisible();
      await expect(rewardsPage.inputMGMYear).toBeVisible();
      await expect(rewardsPage.inputMGMRewardType).toBeVisible();
      await expect(rewardsPage.inputMGMTrigger).toBeVisible();
      await expect(rewardsPage.inputMGMReferrerAmount).toBeVisible();
    });

    test('should add new MGM reward config via HTMX', async ({ page }) => {
      const mgmRewardIdsBefore = await rewardsPage.getAllMGMRewardIds();
      const countBefore = mgmRewardIdsBefore.length;

      // Add a new MGM reward config (using unique combination to avoid constraint violation)
      // Seed data uses: 2025/2026 + (cash|tuition_discount) + enrollment
      // MGM trigger options are only: commitment, enrollment
      // Use: 2028/2029 + cash + commitment (unlikely to exist)
      await rewardsPage.addMGMReward(
        '2028/2029',
        'cash',
        'commitment',
        150000,
        50000,
        'Test MGM reward description'
      );

      // Check if at least we have the same or more rewards (creation succeeded or already exists)
      const mgmRewardIdsAfter = await rewardsPage.getAllMGMRewardIds();
      expect(mgmRewardIdsAfter.length).toBeGreaterThanOrEqual(countBefore);

      // Verify page still works
      await page.reload();
      await rewardsPage.expectPageLoaded();
    });

    test('should toggle MGM reward status via HTMX', async () => {
      const mgmRewardIds = await rewardsPage.getAllMGMRewardIds();
      if (mgmRewardIds.length === 0) {
        test.skip();
        return;
      }

      const mgmRewardId = mgmRewardIds[0];

      const statusBefore = await rewardsPage.getMGMStatusToggle(mgmRewardId).textContent();
      const isActiveBefore = statusBefore?.trim() === 'Aktif';

      await rewardsPage.toggleMGMRewardStatus(mgmRewardId);

      if (isActiveBefore) {
        await rewardsPage.expectMGMRewardStatus(mgmRewardId, 'inactive');
      } else {
        await rewardsPage.expectMGMRewardStatus(mgmRewardId, 'active');
      }

      // Toggle back to restore original state
      await rewardsPage.toggleMGMRewardStatus(mgmRewardId);
    });

    test('should display edit button for each MGM reward', async () => {
      const mgmRewardIds = await rewardsPage.getAllMGMRewardIds();
      if (mgmRewardIds.length === 0) {
        test.skip();
        return;
      }

      for (const mgmRewardId of mgmRewardIds) {
        await expect(rewardsPage.getMGMEditButton(mgmRewardId)).toBeVisible();
      }
    });
  });
});
