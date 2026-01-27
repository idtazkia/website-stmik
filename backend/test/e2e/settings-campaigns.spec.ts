import { test, expect } from '@playwright/test';
import { SettingsCampaignsPage } from './pages';

test.describe('Settings - Campaign Management', () => {
  let campaignsPage: SettingsCampaignsPage;

  test.beforeEach(async ({ page }) => {
    campaignsPage = new SettingsCampaignsPage(page);
    await campaignsPage.login('admin');
    await campaignsPage.goto();
    await campaignsPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display campaigns page', async () => {
      await expect(campaignsPage.pageContainer).toBeVisible();
      await expect(campaignsPage.campaignsSection).toBeVisible();
    });

    test('should display add campaign button', async () => {
      await expect(campaignsPage.addCampaignButton).toBeVisible();
    });
  });

  test.describe('Campaign CRUD', () => {
    test.describe.configure({ mode: 'serial' });

    test('should open add campaign modal', async () => {
      await campaignsPage.openAddCampaignModal();
      await expect(campaignsPage.addCampaignModal).toBeVisible();
      await expect(campaignsPage.inputCampaignName).toBeVisible();
      await expect(campaignsPage.inputCampaignType).toBeVisible();
      await expect(campaignsPage.inputCampaignChannel).toBeVisible();
    });

    test('should add new campaign via HTMX', async ({ page }) => {
      const campaignIdsBefore = await campaignsPage.getAllCampaignIds();
      const countBefore = campaignIdsBefore.length;

      const uniqueName = `Test Campaign ${Date.now()}`;
      await campaignsPage.addCampaign(
        uniqueName,
        'promo',
        'instagram',
        '2026-01-01',
        '2026-03-31',
        0,
        'Test campaign description'
      );

      const campaignIdsAfter = await campaignsPage.getAllCampaignIds();
      expect(campaignIdsAfter.length).toBe(countBefore + 1);

      const newCampaignId = campaignIdsAfter.find(id => !campaignIdsBefore.includes(id));
      expect(newCampaignId).toBeTruthy();

      if (newCampaignId) {
        await campaignsPage.expectCampaignDisplayed(newCampaignId);
        await campaignsPage.expectCampaignNameContains(newCampaignId, uniqueName);

        await page.reload();
        await campaignsPage.expectPageLoaded();
        await campaignsPage.expectCampaignDisplayed(newCampaignId);
      }
    });

    test('should toggle campaign status via HTMX', async () => {
      const campaignIds = await campaignsPage.getAllCampaignIds();
      if (campaignIds.length === 0) {
        test.skip();
        return;
      }

      const campaignId = campaignIds[0];

      const statusBefore = await campaignsPage.getCampaignStatusToggle(campaignId).textContent();
      const isActiveBefore = statusBefore?.trim() === 'Aktif';

      await campaignsPage.toggleCampaignStatus(campaignId);

      if (isActiveBefore) {
        await campaignsPage.expectCampaignStatus(campaignId, 'inactive');
      } else {
        await campaignsPage.expectCampaignStatus(campaignId, 'active');
      }

      await campaignsPage.toggleCampaignStatus(campaignId);
    });

    test('should edit campaign name via HTMX', async ({ page }) => {
      const campaignIds = await campaignsPage.getAllCampaignIds();
      if (campaignIds.length === 0) {
        test.skip();
        return;
      }

      const campaignId = campaignIds[0];
      const newName = `Updated Campaign ${Date.now()}`;

      await campaignsPage.editCampaignName(campaignId, newName);

      await campaignsPage.expectCampaignNameContains(campaignId, newName);

      await page.reload();
      await campaignsPage.expectPageLoaded();
      await campaignsPage.expectCampaignNameContains(campaignId, newName);
    });

    test('should display edit button for each campaign', async () => {
      const campaignIds = await campaignsPage.getAllCampaignIds();
      if (campaignIds.length === 0) {
        test.skip();
        return;
      }

      for (const campaignId of campaignIds) {
        await expect(campaignsPage.getCampaignEditButton(campaignId)).toBeVisible();
      }
    });
  });
});
