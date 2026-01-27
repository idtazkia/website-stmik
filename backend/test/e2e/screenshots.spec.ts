import { test, expect } from '@playwright/test';

// Screenshot tests for user manual documentation
// Run with: npx playwright test e2e/screenshots.spec.ts
// Screenshots saved to: docs/screenshots/

const SCREENSHOT_DIR = 'docs/screenshots';

test.describe('User Manual Screenshots - Admin Settings', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    // Login as admin
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
  });

  // Authentication
  test('01 - Login Page', async ({ page }) => {
    await page.context().clearCookies();
    await page.goto('/admin/login');
    await expect(page.locator('text=Masuk dengan Google')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/01-login.png`, fullPage: true });
  });

  // Dashboard (mockup - no specific test ID)
  test('02 - Admin Dashboard', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.locator('text=Total Kandidat')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/02-dashboard.png`, fullPage: true });
  });

  // Settings - Users
  test('03 - Settings Users', async ({ page }) => {
    await page.goto('/admin/settings/users');
    await expect(page.getByTestId('settings-users-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/03-settings-users.png`, fullPage: true });
  });

  // Settings - Programs
  test('04 - Settings Programs', async ({ page }) => {
    await page.goto('/admin/settings/programs');
    await expect(page.getByTestId('settings-programs-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/04-settings-programs.png`, fullPage: true });
  });

  test('04a - Settings Programs - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/programs');
    await expect(page.getByTestId('settings-programs-page')).toBeVisible();
    await page.getByTestId('add-program-button').click();
    await expect(page.getByTestId('add-program-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/04a-settings-programs-add.png`, fullPage: true });
  });

  // Settings - Fees
  test('05 - Settings Fees', async ({ page }) => {
    await page.goto('/admin/settings/fees');
    await expect(page.getByTestId('settings-fees-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/05-settings-fees.png`, fullPage: true });
  });

  test('05a - Settings Fees - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/fees');
    await expect(page.getByTestId('settings-fees-page')).toBeVisible();
    await page.getByTestId('add-fee-button').click();
    await expect(page.getByTestId('add-fee-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/05a-settings-fees-add.png`, fullPage: true });
  });

  // Settings - Categories & Obstacles
  test('06 - Settings Categories', async ({ page }) => {
    await page.goto('/admin/settings/categories');
    await expect(page.getByTestId('settings-categories-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/06-settings-categories.png`, fullPage: true });
  });

  test('06a - Settings Categories - Add Category Modal', async ({ page }) => {
    await page.goto('/admin/settings/categories');
    await expect(page.getByTestId('settings-categories-page')).toBeVisible();
    await page.getByTestId('add-category-button').click();
    await expect(page.getByTestId('add-category-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/06a-settings-categories-add.png`, fullPage: true });
  });

  test('06b - Settings Categories - Add Obstacle Modal', async ({ page }) => {
    await page.goto('/admin/settings/categories');
    await expect(page.getByTestId('settings-categories-page')).toBeVisible();
    await page.getByTestId('add-obstacle-button').click();
    await expect(page.getByTestId('add-obstacle-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/06b-settings-obstacles-add.png`, fullPage: true });
  });

  // Settings - Campaigns
  test('07 - Settings Campaigns', async ({ page }) => {
    await page.goto('/admin/settings/campaigns');
    await expect(page.getByTestId('settings-campaigns-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/07-settings-campaigns.png`, fullPage: true });
  });

  test('07a - Settings Campaigns - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/campaigns');
    await expect(page.getByTestId('settings-campaigns-page')).toBeVisible();
    await page.getByTestId('add-campaign-button').click();
    await expect(page.getByTestId('add-campaign-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/07a-settings-campaigns-add.png`, fullPage: true });
  });

  // Settings - Rewards
  test('08 - Settings Rewards', async ({ page }) => {
    await page.goto('/admin/settings/rewards');
    await expect(page.getByTestId('settings-rewards-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/08-settings-rewards.png`, fullPage: true });
  });

  test('08a - Settings Rewards - Add Reward Modal', async ({ page }) => {
    await page.goto('/admin/settings/rewards');
    await expect(page.getByTestId('settings-rewards-page')).toBeVisible();
    await page.getByTestId('add-reward-button').click();
    await expect(page.getByTestId('add-reward-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/08a-settings-rewards-add.png`, fullPage: true });
  });

  test('08b - Settings Rewards - Add MGM Modal', async ({ page }) => {
    await page.goto('/admin/settings/rewards');
    await expect(page.getByTestId('settings-rewards-page')).toBeVisible();
    await page.getByTestId('add-mgm-reward-button').click();
    await expect(page.getByTestId('add-mgm-reward-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/08b-settings-mgm-add.png`, fullPage: true });
  });

  // Settings - Referrers
  test('09 - Settings Referrers', async ({ page }) => {
    await page.goto('/admin/settings/referrers');
    await expect(page.getByTestId('settings-referrers-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/09-settings-referrers.png`, fullPage: true });
  });

  test('09a - Settings Referrers - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/referrers');
    await expect(page.getByTestId('settings-referrers-page')).toBeVisible();
    await page.getByTestId('add-referrer-button').click();
    await expect(page.getByTestId('add-referrer-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/09a-settings-referrers-add.png`, fullPage: true });
  });

  // Settings - Assignment Algorithm
  test('10 - Settings Assignment', async ({ page }) => {
    await page.goto('/admin/settings/assignment');
    await expect(page.getByTestId('settings-assignment-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/10-settings-assignment.png`, fullPage: true });
  });

  // Settings - Document Types
  test('11 - Settings Document Types', async ({ page }) => {
    await page.goto('/admin/settings/document-types');
    await expect(page.getByTestId('settings-documents-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/11-settings-document-types.png`, fullPage: true });
  });

  test('11a - Settings Document Types - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/document-types');
    await expect(page.getByTestId('settings-documents-page')).toBeVisible();
    await page.getByTestId('add-document-type-button').click();
    await expect(page.getByTestId('add-document-type-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/11a-settings-document-types-add.png`, fullPage: true });
  });

  // Settings - Lost Reasons
  test('12 - Settings Lost Reasons', async ({ page }) => {
    await page.goto('/admin/settings/lost-reasons');
    await expect(page.getByTestId('settings-lost-reasons-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/12-settings-lost-reasons.png`, fullPage: true });
  });

  test('12a - Settings Lost Reasons - Add Modal', async ({ page }) => {
    await page.goto('/admin/settings/lost-reasons');
    await expect(page.getByTestId('settings-lost-reasons-page')).toBeVisible();
    await page.getByTestId('add-lost-reason-button').click();
    await expect(page.getByTestId('add-lost-reason-modal')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/12a-settings-lost-reasons-add.png`, fullPage: true });
  });
});

// Generate unique identifiers for test data
function generateUniqueEmail(): string {
  const timestamp = Date.now();
  return `screenshot${timestamp}@example.com`;
}

test.describe('User Manual Screenshots - Candidate Registration', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
  });

  // Registration Step 1 - Account Creation (Empty)
  test('20 - Registration Step 1 - Account Empty', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('registration-form')).toBeVisible();
    await expect(page.getByTestId('step1-form')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/20-registration-step1-empty.png`, fullPage: true });
  });

  // Registration Step 1 - Account Creation (Filled)
  test('20a - Registration Step 1 - Account Filled', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('step1-form')).toBeVisible();

    await page.getByTestId('input-email').fill('contoh@email.com');
    await page.getByTestId('input-phone').fill('081234567890');
    await page.getByTestId('input-password').fill('password123');
    await page.getByTestId('input-password-confirm').fill('password123');
    await page.screenshot({ path: `${SCREENSHOT_DIR}/20a-registration-step1-filled.png`, fullPage: true });
  });

  // Registration Step 2 - Personal Info (Empty)
  test('21 - Registration Step 2 - Personal Info Empty', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('step1-form')).toBeVisible();

    const email = generateUniqueEmail();
    await page.getByTestId('input-email').fill(email);
    await page.getByTestId('input-password').fill('testpassword123');
    await page.getByTestId('input-password-confirm').fill('testpassword123');
    await page.getByTestId('btn-submit-step1').click();

    await expect(page.getByTestId('step2-form')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/21-registration-step2-empty.png`, fullPage: true });
  });

  // Registration Step 2 - Personal Info (Filled)
  test('21a - Registration Step 2 - Personal Info Filled', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('step1-form')).toBeVisible();

    const email = generateUniqueEmail();
    await page.getByTestId('input-email').fill(email);
    await page.getByTestId('input-password').fill('testpassword123');
    await page.getByTestId('input-password-confirm').fill('testpassword123');
    await page.getByTestId('btn-submit-step1').click();

    await expect(page.getByTestId('step2-form')).toBeVisible();
    await page.getByTestId('input-name').fill('Budi Santoso');
    await page.getByTestId('input-address').fill('Jl. Merdeka No. 123, RT 01/RW 02, Kelurahan Sukamaju');
    await page.getByTestId('input-city').fill('Jakarta Selatan');
    await page.getByTestId('input-province').fill('DKI Jakarta');
    await page.screenshot({ path: `${SCREENSHOT_DIR}/21a-registration-step2-filled.png`, fullPage: true });
  });

  // Registration Step 3 - Education (Empty)
  test('22 - Registration Step 3 - Education Empty', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('step1-form')).toBeVisible();

    const email = generateUniqueEmail();
    await page.getByTestId('input-email').fill(email);
    await page.getByTestId('input-password').fill('testpassword123');
    await page.getByTestId('input-password-confirm').fill('testpassword123');
    await page.getByTestId('btn-submit-step1').click();

    await expect(page.getByTestId('step2-form')).toBeVisible();
    await page.getByTestId('input-name').fill('Test Candidate');
    await page.getByTestId('input-address').fill('Jl. Test No. 123');
    await page.getByTestId('input-city').fill('Jakarta');
    await page.getByTestId('input-province').fill('DKI Jakarta');
    await page.getByTestId('btn-submit-step2').click();

    await expect(page.getByTestId('step3-form')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/22-registration-step3-empty.png`, fullPage: true });
  });

  // Registration Step 3 - Education (Filled)
  test('22a - Registration Step 3 - Education Filled', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('step1-form')).toBeVisible();

    const email = generateUniqueEmail();
    await page.getByTestId('input-email').fill(email);
    await page.getByTestId('input-password').fill('testpassword123');
    await page.getByTestId('input-password-confirm').fill('testpassword123');
    await page.getByTestId('btn-submit-step1').click();

    await expect(page.getByTestId('step2-form')).toBeVisible();
    await page.getByTestId('input-name').fill('Test Candidate');
    await page.getByTestId('input-address').fill('Jl. Test No. 123');
    await page.getByTestId('input-city').fill('Jakarta');
    await page.getByTestId('input-province').fill('DKI Jakarta');
    await page.getByTestId('btn-submit-step2').click();

    await expect(page.getByTestId('step3-form')).toBeVisible();

    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    const radioCount = await prodiRadios.count();

    if (radioCount === 0) {
      test.skip();
      return;
    }

    await page.getByTestId('input-high-school').fill('SMA Negeri 1 Jakarta');
    await page.getByTestId('select-graduation-year').selectOption('2025');
    await prodiRadios.first().click();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/22a-registration-step3-filled.png`, fullPage: true });
  });

  // Registration Step 4 - Source Tracking (Empty)
  test('23 - Registration Step 4 - Source Tracking Empty', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('step1-form')).toBeVisible();

    const email = generateUniqueEmail();
    await page.getByTestId('input-email').fill(email);
    await page.getByTestId('input-password').fill('testpassword123');
    await page.getByTestId('input-password-confirm').fill('testpassword123');
    await page.getByTestId('btn-submit-step1').click();

    await expect(page.getByTestId('step2-form')).toBeVisible();
    await page.getByTestId('input-name').fill('Test Candidate');
    await page.getByTestId('input-address').fill('Jl. Test No. 123');
    await page.getByTestId('input-city').fill('Jakarta');
    await page.getByTestId('input-province').fill('DKI Jakarta');
    await page.getByTestId('btn-submit-step2').click();

    await expect(page.getByTestId('step3-form')).toBeVisible();

    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    const radioCount = await prodiRadios.count();

    if (radioCount === 0) {
      test.skip();
      return;
    }

    await page.getByTestId('input-high-school').fill('SMA Negeri 1 Jakarta');
    await page.getByTestId('select-graduation-year').selectOption('2025');
    await prodiRadios.first().click();
    await page.getByTestId('btn-submit-step3').click();

    await expect(page.getByTestId('step4-form')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/23-registration-step4-empty.png`, fullPage: true });
  });

  // Registration Step 4 - Source Tracking (Filled)
  test('23a - Registration Step 4 - Source Tracking Filled', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByTestId('step1-form')).toBeVisible();

    const email = generateUniqueEmail();
    await page.getByTestId('input-email').fill(email);
    await page.getByTestId('input-password').fill('testpassword123');
    await page.getByTestId('input-password-confirm').fill('testpassword123');
    await page.getByTestId('btn-submit-step1').click();

    await expect(page.getByTestId('step2-form')).toBeVisible();
    await page.getByTestId('input-name').fill('Test Candidate');
    await page.getByTestId('input-address').fill('Jl. Test No. 123');
    await page.getByTestId('input-city').fill('Jakarta');
    await page.getByTestId('input-province').fill('DKI Jakarta');
    await page.getByTestId('btn-submit-step2').click();

    await expect(page.getByTestId('step3-form')).toBeVisible();

    const prodiRadios = page.locator('[data-testid^="radio-prodi-"]');
    const radioCount = await prodiRadios.count();

    if (radioCount === 0) {
      test.skip();
      return;
    }

    await page.getByTestId('input-high-school').fill('SMA Negeri 1 Jakarta');
    await page.getByTestId('select-graduation-year').selectOption('2025');
    await prodiRadios.first().click();
    await page.getByTestId('btn-submit-step3').click();

    await expect(page.getByTestId('step4-form')).toBeVisible();
    await page.getByTestId('select-source-type').selectOption('friend_family');
    await page.getByTestId('input-source-detail').fill('Kakak tingkat di kampus');
    await page.screenshot({ path: `${SCREENSHOT_DIR}/23a-registration-step4-filled.png`, fullPage: true });
  });

  // Candidate Login Page
  test('24 - Candidate Login', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByTestId('portal-login-form')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/24-candidate-login.png`, fullPage: true });
  });
});

test.describe('User Manual Screenshots - Candidate Portal', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    // Login as test candidate
    await page.goto('/test/login/candidate');
    await page.waitForURL(/\/portal\/?$/);
  });

  // Portal Dashboard
  test('25 - Portal Dashboard', async ({ page }) => {
    await expect(page.getByTestId('portal-dashboard')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/25-portal-dashboard.png`, fullPage: true });
  });

  // Portal Documents
  test('26 - Portal Documents', async ({ page }) => {
    await page.goto('/portal/documents');
    await expect(page.getByTestId('documents-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/26-portal-documents.png`, fullPage: true });
  });

  // Portal Payments
  test('27 - Portal Payments', async ({ page }) => {
    await page.goto('/portal/payments');
    await expect(page.getByTestId('payments-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/27-portal-payments.png`, fullPage: true });
  });

  // Portal Announcements
  test('28 - Portal Announcements', async ({ page }) => {
    await page.goto('/portal/announcements');
    await expect(page.getByTestId('announcements-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/28-portal-announcements.png`, fullPage: true });
  });
});

test.describe('User Manual Screenshots - Admin Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
  });

  // Admin Dashboard - Overview
  test('30 - Admin Dashboard Overview', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.getByTestId('admin-dashboard')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/30-admin-dashboard-overview.png`, fullPage: true });
  });

  // Admin Dashboard - Stats Cards
  test('30a - Admin Dashboard Stats Cards', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.getByTestId('stats-cards')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/30a-admin-dashboard-stats.png`,
      clip: await page.getByTestId('stats-cards').boundingBox() || undefined
    });
  });

  // Admin Dashboard - Overdue Section
  test('30b - Admin Dashboard Overdue Section', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.getByTestId('overdue-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/30b-admin-dashboard-overdue.png`,
      clip: await page.getByTestId('overdue-section').boundingBox() || undefined
    });
  });

  // Admin Dashboard - Today Tasks Section
  test('30c - Admin Dashboard Today Tasks', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.getByTestId('today-tasks-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/30c-admin-dashboard-today-tasks.png`,
      clip: await page.getByTestId('today-tasks-section').boundingBox() || undefined
    });
  });

  // Admin Dashboard - Funnel Section
  test('30d - Admin Dashboard Funnel', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.getByTestId('funnel-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/30d-admin-dashboard-funnel.png`,
      clip: await page.getByTestId('funnel-section').boundingBox() || undefined
    });
  });

  // Admin Dashboard - Recent Candidates
  test('30e - Admin Dashboard Recent Candidates', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.getByTestId('recent-candidates-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/30e-admin-dashboard-recent-candidates.png`,
      clip: await page.getByTestId('recent-candidates-section').boundingBox() || undefined
    });
  });
});

test.describe('User Manual Screenshots - Candidate Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
  });

  // Candidates List - Full View
  test('31 - Candidates List', async ({ page }) => {
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('candidates-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/31-candidates-list.png`, fullPage: true });
  });

  // Candidates List - Stats Section
  test('31a - Candidates Stats', async ({ page }) => {
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('candidate-stats')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/31a-candidates-stats.png`,
      clip: await page.getByTestId('candidate-stats').boundingBox() || undefined
    });
  });

  // Candidates List - Filters Section
  test('31b - Candidates Filters', async ({ page }) => {
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('filters-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/31b-candidates-filters.png`,
      clip: await page.getByTestId('filters-section').boundingBox() || undefined
    });
  });

  // Candidates List - Filter by Status
  test('31c - Candidates Filter by Status', async ({ page }) => {
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('filters-section')).toBeVisible();
    await page.getByTestId('filter-status').selectOption('prospecting');
    await page.waitForTimeout(500); // Wait for HTMX to update
    await page.screenshot({ path: `${SCREENSHOT_DIR}/31c-candidates-filter-status.png`, fullPage: true });
  });

  // Candidates List - Search
  test('31d - Candidates Search', async ({ page }) => {
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('filters-section')).toBeVisible();
    await page.getByTestId('filter-search').fill('test');
    await page.waitForTimeout(600); // Wait for debounced HTMX
    await page.screenshot({ path: `${SCREENSHOT_DIR}/31d-candidates-search.png`, fullPage: true });
  });

  // Candidates List - Table
  test('31e - Candidates Table', async ({ page }) => {
    await page.goto('/admin/candidates');
    await expect(page.getByTestId('candidates-table')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/31e-candidates-table.png`,
      clip: await page.getByTestId('candidates-table').boundingBox() || undefined
    });
  });

  // Candidate Detail - Full View
  test('32 - Candidate Detail', async ({ page }) => {
    // Find first candidate and navigate
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      await expect(page.getByTestId('candidate-detail-page')).toBeVisible();
      await page.screenshot({ path: `${SCREENSHOT_DIR}/32-candidate-detail.png`, fullPage: true });
    } else {
      test.skip();
    }
  });

  // Candidate Detail - Header
  test('32a - Candidate Detail Header', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      await expect(page.getByTestId('candidate-header')).toBeVisible();
      await page.screenshot({
        path: `${SCREENSHOT_DIR}/32a-candidate-detail-header.png`,
        clip: await page.getByTestId('candidate-header').boundingBox() || undefined
      });
    } else {
      test.skip();
    }
  });

  // Candidate Detail - Personal Info Section
  test('32b - Candidate Personal Info', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      await expect(page.getByTestId('section-personal-info')).toBeVisible();
      await page.screenshot({
        path: `${SCREENSHOT_DIR}/32b-candidate-personal-info.png`,
        clip: await page.getByTestId('section-personal-info').boundingBox() || undefined
      });
    } else {
      test.skip();
    }
  });

  // Candidate Detail - Education Section
  test('32c - Candidate Education', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      await expect(page.getByTestId('section-education')).toBeVisible();
      await page.screenshot({
        path: `${SCREENSHOT_DIR}/32c-candidate-education.png`,
        clip: await page.getByTestId('section-education').boundingBox() || undefined
      });
    } else {
      test.skip();
    }
  });

  // Candidate Detail - Source & Assignment Section
  test('32d - Candidate Source Assignment', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      await expect(page.getByTestId('section-source-assignment')).toBeVisible();
      await page.screenshot({
        path: `${SCREENSHOT_DIR}/32d-candidate-source-assignment.png`,
        clip: await page.getByTestId('section-source-assignment').boundingBox() || undefined
      });
    } else {
      test.skip();
    }
  });

  // Candidate Detail - Timeline Section
  test('32e - Candidate Timeline', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      await expect(page.getByTestId('section-timeline')).toBeVisible();
      await page.screenshot({
        path: `${SCREENSHOT_DIR}/32e-candidate-timeline.png`,
        clip: await page.getByTestId('section-timeline').boundingBox() || undefined
      });
    } else {
      test.skip();
    }
  });

  // Candidate Detail - Interaction Modal
  test('32f - Candidate Interaction Modal', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      await expect(page.getByTestId('btn-log-interaction')).toBeVisible();
      await page.getByTestId('btn-log-interaction').click();
      await expect(page.getByTestId('modal-interaction')).toBeVisible();
      await page.screenshot({ path: `${SCREENSHOT_DIR}/32f-candidate-interaction-modal.png`, fullPage: true });
    } else {
      test.skip();
    }
  });
});

test.describe('User Manual Screenshots - AC (Academic Consultant) Features', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await page.goto('/test/login/consultant');
    await page.waitForURL(/\/admin\/?$/);
  });

  // AC Dashboard - Full View
  test('40 - AC Dashboard', async ({ page }) => {
    await page.goto('/admin/my-dashboard');
    await expect(page.getByTestId('consultant-dashboard')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/40-ac-dashboard.png`, fullPage: true });
  });

  // AC Dashboard - Welcome Section
  test('40a - AC Dashboard Welcome Section', async ({ page }) => {
    await page.goto('/admin/my-dashboard');
    await expect(page.getByTestId('welcome-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/40a-ac-dashboard-welcome.png`,
      clip: await page.getByTestId('welcome-section').boundingBox() || undefined
    });
  });

  // AC Dashboard - Personal Stats
  test('40b - AC Dashboard Personal Stats', async ({ page }) => {
    await page.goto('/admin/my-dashboard');
    await expect(page.getByTestId('personal-stats')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/40b-ac-dashboard-stats.png`,
      clip: await page.getByTestId('personal-stats').boundingBox() || undefined
    });
  });

  // AC Dashboard - Overdue Section
  test('40c - AC Dashboard Overdue', async ({ page }) => {
    await page.goto('/admin/my-dashboard');
    await expect(page.getByTestId('ac-overdue-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/40c-ac-dashboard-overdue.png`,
      clip: await page.getByTestId('ac-overdue-section').boundingBox() || undefined
    });
  });

  // AC Dashboard - Today Tasks Section
  test('40d - AC Dashboard Today Tasks', async ({ page }) => {
    await page.goto('/admin/my-dashboard');
    await expect(page.getByTestId('ac-today-tasks-section')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/40d-ac-dashboard-today-tasks.png`,
      clip: await page.getByTestId('ac-today-tasks-section').boundingBox() || undefined
    });
  });

  // AC Dashboard - Supervisor Suggestions
  test('40e - AC Dashboard Supervisor Suggestions', async ({ page }) => {
    await page.goto('/admin/my-dashboard');
    const section = page.getByTestId('supervisor-suggestions-section');
    await expect(section).toBeVisible();
    await section.scrollIntoViewIfNeeded();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/40e-ac-dashboard-suggestions.png`,
      fullPage: true
    });
  });

  // AC Dashboard - Monthly Performance
  test('40f - AC Dashboard Monthly Performance', async ({ page }) => {
    await page.goto('/admin/my-dashboard');
    const section = page.getByTestId('monthly-performance-section');
    await expect(section).toBeVisible();
    await section.scrollIntoViewIfNeeded();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/40f-ac-dashboard-monthly-performance.png`,
      fullPage: true
    });
  });

  // Interaction Form - Full View (Empty)
  test('41 - Interaction Form', async ({ page }) => {
    // Navigate to a candidate's interaction form
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      // Click "Log Interaksi" link or look for interaction form page
      const interactionLink = page.locator('a:has-text("Log")').first();
      if (await interactionLink.count() > 0) {
        await interactionLink.click();
        if (await page.getByTestId('interaction-form-page').count() > 0) {
          await expect(page.getByTestId('interaction-form-page')).toBeVisible();
          await page.screenshot({ path: `${SCREENSHOT_DIR}/41-interaction-form.png`, fullPage: true });
        } else {
          test.skip();
        }
      } else {
        test.skip();
      }
    } else {
      test.skip();
    }
  });

  // Interaction Form - Channel Options
  test('41a - Interaction Form Channel Options', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      const interactionLink = page.locator('a:has-text("Log")').first();
      if (await interactionLink.count() > 0) {
        await interactionLink.click();
        if (await page.getByTestId('channel-section').count() > 0) {
          await expect(page.getByTestId('channel-section')).toBeVisible();
          await page.screenshot({
            path: `${SCREENSHOT_DIR}/41a-interaction-form-channels.png`,
            clip: await page.getByTestId('channel-section').boundingBox() || undefined
          });
        } else {
          test.skip();
        }
      } else {
        test.skip();
      }
    } else {
      test.skip();
    }
  });

  // Interaction Form - Category Options
  test('41b - Interaction Form Category Options', async ({ page }) => {
    await page.goto('/admin/candidates');
    const candidateLink = page.locator('[data-testid="candidate-name"]').first();
    if (await candidateLink.count() > 0) {
      await candidateLink.click();
      const interactionLink = page.locator('a:has-text("Log")').first();
      if (await interactionLink.count() > 0) {
        await interactionLink.click();
        if (await page.getByTestId('category-section').count() > 0) {
          await expect(page.getByTestId('category-section')).toBeVisible();
          await page.screenshot({
            path: `${SCREENSHOT_DIR}/41b-interaction-form-categories.png`,
            clip: await page.getByTestId('category-section').boundingBox() || undefined
          });
        } else {
          test.skip();
        }
      } else {
        test.skip();
      }
    } else {
      test.skip();
    }
  });
});

test.describe('User Manual Screenshots - Consultant Performance Reports', () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await page.goto('/test/login/admin');
    await page.waitForURL(/\/admin\/?$/);
  });

  // Consultant Report - Full View
  test('50 - Consultant Report', async ({ page }) => {
    await page.goto('/admin/reports/consultants');
    await expect(page.getByTestId('consultant-report-page')).toBeVisible();
    await page.screenshot({ path: `${SCREENSHOT_DIR}/50-consultant-report.png`, fullPage: true });
  });

  // Consultant Report - Filter Section
  test('50a - Consultant Report Filters', async ({ page }) => {
    await page.goto('/admin/reports/consultants');
    await expect(page.getByTestId('report-filter')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/50a-consultant-report-filters.png`,
      clip: await page.getByTestId('report-filter').boundingBox() || undefined
    });
  });

  // Consultant Report - Summary Cards
  test('50b - Consultant Report Summary', async ({ page }) => {
    await page.goto('/admin/reports/consultants');
    await expect(page.getByTestId('report-summary')).toBeVisible();
    await page.screenshot({
      path: `${SCREENSHOT_DIR}/50b-consultant-report-summary.png`,
      clip: await page.getByTestId('report-summary').boundingBox() || undefined
    });
  });
});
