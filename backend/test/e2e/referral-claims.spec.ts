import { test, expect, Browser, Page } from '@playwright/test';

// Helper to create a candidate with referral claim
async function createCandidateWithReferralClaim(browser: Browser): Promise<{ candidateId: string; page: Page }> {
  const candidatePage = await browser.newPage();
  const uniqueId = `${Date.now()}${Math.random().toString(36).slice(2, 8)}`;
  const uniqueEmail = `referral${uniqueId}@example.com`;
  const password = 'testpassword123';

  // Step 1: Account creation
  await candidatePage.goto('/register');
  await candidatePage.getByTestId('input-email').fill(uniqueEmail);
  await candidatePage.getByTestId('input-password').fill(password);
  await candidatePage.getByTestId('input-password-confirm').fill(password);
  await candidatePage.getByTestId('btn-submit-step1').click();
  await expect(candidatePage.getByTestId('step2-form')).toBeVisible({ timeout: 10000 });

  // Step 2: Personal info
  await candidatePage.getByTestId('input-name').fill(`ReferralTest ${uniqueId}`);
  await candidatePage.getByTestId('input-address').fill('Test Address');
  await candidatePage.getByTestId('input-city').fill('Jakarta');
  await candidatePage.getByTestId('input-province').fill('DKI Jakarta');
  await candidatePage.getByTestId('btn-submit-step2').click();
  await expect(candidatePage.getByTestId('step3-form')).toBeVisible({ timeout: 10000 });

  // Step 3: Education
  await candidatePage.getByTestId('input-high-school').fill('SMA Test');
  await candidatePage.getByTestId('select-graduation-year').selectOption('2025');
  const prodiRadios = candidatePage.locator('input[type="radio"][name="prodi_id"]');
  await prodiRadios.first().click();
  await candidatePage.getByTestId('btn-submit-step3').click();
  await expect(candidatePage.getByTestId('step4-form')).toBeVisible({ timeout: 10000 });

  // Step 4: Source tracking with referral claim (teacher_alumni triggers referral claim)
  await candidatePage.getByTestId('select-source-type').selectOption('teacher_alumni');
  // Wait for detail container to be visible
  await expect(candidatePage.locator('#source_detail_container')).toBeVisible();
  await candidatePage.getByTestId('input-source-detail').fill('Pak Ahmad dari SMA Negeri 1');
  await candidatePage.getByTestId('btn-submit-step4').click();
  await expect(candidatePage).toHaveURL('/portal', { timeout: 10000 });

  await candidatePage.close();

  // Login as admin to get candidate ID
  const adminPage = await browser.newPage();
  await adminPage.goto('/test/login/admin');
  await adminPage.goto('/admin/candidates?search=' + encodeURIComponent(uniqueEmail));
  await expect(adminPage.getByTestId('candidates-page')).toBeVisible();
  await adminPage.waitForTimeout(1000);

  // Get candidate ID
  const viewLink = adminPage.locator('[data-testid^="view-candidate-"]').first();
  const testId = await viewLink.getAttribute('data-testid');
  const candidateId = testId?.replace('view-candidate-', '') || '';

  return { candidateId, page: adminPage };
}

test.describe('Referral Claims Management', () => {
  test.describe('Referral Claims List', () => {
    test('admin can access referral claims page', async ({ browser }) => {
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/referral-claims');

      await expect(adminPage.getByTestId('referral-claims-page')).toBeVisible();

      await adminPage.close();
    });

    test('referral claims page shows unverified claims', async ({ browser }) => {
      // Create a candidate with referral claim
      const { candidateId, page: adminPage } = await createCandidateWithReferralClaim(browser);

      // Go to referral claims page
      await adminPage.goto('/admin/referral-claims');
      await expect(adminPage.getByTestId('referral-claims-page')).toBeVisible();

      // Should see the claim in the list
      const claimRow = adminPage.getByTestId(`claim-row-${candidateId}`);
      await expect(claimRow).toBeVisible();
      await expect(claimRow).toContainText('Pak Ahmad dari SMA Negeri 1');

      await adminPage.close();
    });
  });

  test.describe('Link Referral Claim', () => {
    test('admin can open link referrer modal', async ({ browser }) => {
      const { candidateId, page: adminPage } = await createCandidateWithReferralClaim(browser);

      await adminPage.goto('/admin/referral-claims');
      await expect(adminPage.getByTestId('referral-claims-page')).toBeVisible();

      // Click link button
      await adminPage.getByTestId(`btn-link-${candidateId}`).click();

      // Modal should appear
      await expect(adminPage.locator('#link-modal')).toBeVisible();
      await expect(adminPage.locator('#modal-claim-text')).toContainText('Pak Ahmad');

      await adminPage.close();
    });

    test('admin can link claim to existing referrer', async ({ browser }) => {
      // First create a referrer
      const adminPage = await browser.newPage();
      await adminPage.goto('/test/login/admin');
      await adminPage.goto('/admin/settings/referrers');

      const uniqueId = `${Date.now()}`;
      const referrerName = `TestReferrer ${uniqueId}`;

      // Open add referrer modal
      await adminPage.getByTestId('add-referrer-button').click();
      await expect(adminPage.getByTestId('add-referrer-modal')).toBeVisible();

      // Create referrer
      await adminPage.getByTestId('input-referrer-name').fill(referrerName);
      await adminPage.getByTestId('input-referrer-type').selectOption('teacher');
      await adminPage.getByTestId('input-referrer-institution').fill('SMA Negeri 1');
      await adminPage.getByTestId('submit-add-referrer').click();
      await adminPage.waitForTimeout(1000);

      await adminPage.close();

      // Create candidate with referral claim
      const { candidateId, page: candidatePage } = await createCandidateWithReferralClaim(browser);

      // Go to referral claims and link
      await candidatePage.goto('/admin/referral-claims');
      await candidatePage.getByTestId(`btn-link-${candidateId}`).click();
      await expect(candidatePage.locator('#link-modal')).toBeVisible();

      // Select the referrer - find option containing the name
      const referrerSelect = candidatePage.locator('#referrer-select');
      const options = await referrerSelect.locator('option').all();
      let referrerValue = '';
      for (const option of options) {
        const text = await option.textContent();
        if (text?.includes(referrerName)) {
          referrerValue = await option.getAttribute('value') || '';
          break;
        }
      }
      expect(referrerValue).not.toBe('');
      await referrerSelect.selectOption(referrerValue);

      // Submit
      await candidatePage.getByTestId('btn-confirm-link').click();

      // Page should refresh and claim should be removed from list
      await candidatePage.waitForTimeout(2000);
      await expect(candidatePage.getByTestId(`claim-row-${candidateId}`)).not.toBeVisible();

      await candidatePage.close();
    });
  });

  test.describe('Invalid Referral Claim', () => {
    test('admin can mark claim as invalid', async ({ browser }) => {
      const { candidateId, page: adminPage } = await createCandidateWithReferralClaim(browser);

      await adminPage.goto('/admin/referral-claims');
      await expect(adminPage.getByTestId('referral-claims-page')).toBeVisible();

      // Accept the confirmation dialog
      adminPage.on('dialog', dialog => dialog.accept());

      // Click invalid button
      await adminPage.getByTestId(`btn-invalid-${candidateId}`).click();

      // Page should refresh and claim should be removed
      await adminPage.waitForTimeout(2000);
      await expect(adminPage.getByTestId(`claim-row-${candidateId}`)).not.toBeVisible();

      await adminPage.close();
    });
  });
});
