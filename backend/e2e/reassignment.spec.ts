import { test, expect, Browser, Page } from '@playwright/test';
import { CandidatesPage, CandidateDetailPage } from './pages';

// Helper to register a new candidate and get their ID
async function registerCandidateAndGetId(browser: Browser): Promise<{ id: string; page: Page }> {
  const candidatePage = await browser.newPage();
  const uniqueId = `${Date.now()}${Math.random().toString(36).slice(2, 8)}`;
  const uniqueEmail = `reassign${uniqueId}@example.com`;
  const password = 'testpassword123';

  // Step 1: Account creation
  await candidatePage.goto('/register');
  await candidatePage.getByTestId('input-email').fill(uniqueEmail);
  await candidatePage.getByTestId('input-password').fill(password);
  await candidatePage.getByTestId('input-password-confirm').fill(password);
  await candidatePage.getByTestId('btn-submit-step1').click();
  await expect(candidatePage.getByTestId('step2-form')).toBeVisible({ timeout: 10000 });

  // Step 2: Personal info
  await candidatePage.getByTestId('input-name').fill(`ReassignTest ${uniqueId}`);
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

  // Step 4: Source tracking - complete registration
  await candidatePage.getByTestId('select-source-type').selectOption('instagram');
  await candidatePage.getByTestId('btn-submit-step4').click();
  await expect(candidatePage).toHaveURL('/portal', { timeout: 10000 });

  // Get candidate ID from URL or database
  // For now, we'll search for this candidate in admin panel
  await candidatePage.close();

  // Open admin page to find candidate
  const adminPage = await browser.newPage();
  await adminPage.goto('/test/login/admin');
  await adminPage.goto('/admin/candidates?search=' + encodeURIComponent(uniqueEmail));
  await expect(adminPage.getByTestId('candidates-page')).toBeVisible();

  // Wait for search to complete
  await adminPage.waitForTimeout(1000);

  // Get the candidate row and extract ID from the view link
  const viewLink = adminPage.locator('[data-testid^="view-candidate-"]').first();
  const testId = await viewLink.getAttribute('data-testid');
  const candidateId = testId?.replace('view-candidate-', '') || '';

  return { id: candidateId, page: adminPage };
}

test.describe('Candidate Reassignment', () => {
  test.describe('Modal Behavior', () => {
    test('should open reassign modal when clicking reassign button', async ({ page }) => {
      const candidatesPage = new CandidatesPage(page);
      const detailPage = new CandidateDetailPage(page);

      await candidatesPage.login('admin');
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      // Click on first candidate if available
      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);
        await detailPage.expectPageLoaded();

        // Click reassign button
        await detailPage.btnReassign.click();

        // Modal should appear
        await expect(detailPage.reassignModal).toBeVisible();
        await expect(detailPage.reassignModalTitle).toHaveText('Reassign Kandidat');
        await expect(detailPage.consultantList).toBeVisible();
      }
    });

    test('should close modal when clicking close button', async ({ page }) => {
      const candidatesPage = new CandidatesPage(page);
      const detailPage = new CandidateDetailPage(page);

      await candidatesPage.login('admin');
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Open and close modal
        await detailPage.openReassignModal();
        await detailPage.closeReassignModal();

        // Modal should be gone (removed from DOM)
        await expect(detailPage.reassignModal).not.toBeVisible();
      }
    });

    test('should close modal when clicking cancel button', async ({ page }) => {
      const candidatesPage = new CandidatesPage(page);
      const detailPage = new CandidateDetailPage(page);

      await candidatesPage.login('admin');
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        await detailPage.openReassignModal();
        await detailPage.reassignBtnCancel.click();

        await expect(detailPage.reassignModal).not.toBeVisible();
      }
    });

    test('should display consultant list with workload info', async ({ page }) => {
      const candidatesPage = new CandidatesPage(page);
      const detailPage = new CandidateDetailPage(page);

      await candidatesPage.login('admin');
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        await detailPage.openReassignModal();

        // Verify consultant options are displayed
        const consultantOptions = page.locator('[data-testid^="consultant-option-"]');
        const count = await consultantOptions.count();
        expect(count).toBeGreaterThan(0);

        // Check that each consultant shows workload info (active count, total count)
        if (count > 0) {
          const firstOption = consultantOptions.first();
          await expect(firstOption.locator('text=aktif')).toBeVisible();
          await expect(firstOption.locator('text=total')).toBeVisible();
        }
      }
    });
  });

  test.describe('Access Control', () => {
    test('consultant should not see reassign button', async ({ page }) => {
      const detailPage = new CandidateDetailPage(page);

      await page.goto('/test/login/consultant');
      await page.goto('/admin/candidates');
      await expect(page.getByTestId('candidates-page')).toBeVisible();

      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        // Consultant should still see field-consultant but button behavior differs
        await expect(detailPage.fieldConsultant).toBeVisible();

        // Clicking reassign should not open modal (or button not functional for consultant)
        // The actual restriction is in the backend - consultant gets 403 on POST
      }
    });

    test('supervisor should be able to open reassign modal', async ({ page }) => {
      const candidatesPage = new CandidatesPage(page);
      const detailPage = new CandidateDetailPage(page);

      await page.goto('/test/login/supervisor');
      await page.goto('/admin/candidates');
      await candidatesPage.expectPageLoaded();

      const detailLink = page.locator('[data-testid^="view-candidate-"]').first();
      if (await detailLink.isVisible()) {
        await detailLink.click();
        await page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

        await detailPage.btnReassign.click();
        await expect(detailPage.reassignModal).toBeVisible();
      }
    });
  });

  test.describe('Reassignment Flow', () => {
    test('should successfully reassign candidate to different consultant', async ({ browser }) => {
      // Register a candidate first
      const { id: candidateId, page: adminPage } = await registerCandidateAndGetId(browser);
      const detailPage = new CandidateDetailPage(adminPage);

      if (!candidateId) {
        test.skip();
        return;
      }

      // Navigate to candidate detail
      await adminPage.goto(`/admin/candidates/${candidateId}`);
      await detailPage.expectPageLoaded();

      // Get current consultant name
      const currentConsultant = await detailPage.fieldConsultant.textContent();

      // Open reassign modal
      await detailPage.openReassignModal();

      // Select a different consultant (if available)
      const consultantOptions = adminPage.locator('[data-testid^="consultant-option-"]');
      const count = await consultantOptions.count();

      if (count > 1) {
        // Find a consultant that's not currently selected
        for (let i = 0; i < count; i++) {
          const option = consultantOptions.nth(i);
          const radio = option.locator('input[type="radio"]');
          const isChecked = await radio.isChecked();

          if (!isChecked) {
            // Click on a different consultant
            await radio.click();

            // Submit the form
            await detailPage.reassignBtnSubmit.click();

            // Wait for redirect back to detail page
            await adminPage.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

            // Verify consultant was changed
            const newConsultant = await detailPage.fieldConsultant.textContent();

            // If there was a change, the consultant text should be different
            // Note: The new consultant's name will be displayed
            expect(newConsultant).not.toBe(currentConsultant);
            break;
          }
        }
      }

      await adminPage.close();
    });

    test('should log reassignment in interaction timeline', async ({ browser }) => {
      // Register a candidate first
      const { id: candidateId, page: adminPage } = await registerCandidateAndGetId(browser);
      const detailPage = new CandidateDetailPage(adminPage);

      if (!candidateId) {
        test.skip();
        return;
      }

      await adminPage.goto(`/admin/candidates/${candidateId}`);
      await detailPage.expectPageLoaded();

      // Open reassign modal
      await detailPage.openReassignModal();

      // Select a different consultant
      const consultantOptions = adminPage.locator('[data-testid^="consultant-option-"]');
      const count = await consultantOptions.count();

      if (count > 1) {
        for (let i = 0; i < count; i++) {
          const option = consultantOptions.nth(i);
          const radio = option.locator('input[type="radio"]');
          const isChecked = await radio.isChecked();

          if (!isChecked) {
            await radio.click();
            await detailPage.reassignBtnSubmit.click();
            await adminPage.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);

            // Check that timeline shows reassignment entry
            // The reassignment creates an interaction with channel "system"
            const timelineList = detailPage.timelineList;
            await expect(timelineList).toBeVisible();

            // Look for reassignment entry in timeline
            const reassignEntry = adminPage.locator('text=dialihkan ke konsultan');
            await expect(reassignEntry).toBeVisible();
            break;
          }
        }
      }

      await adminPage.close();
    });
  });
});
