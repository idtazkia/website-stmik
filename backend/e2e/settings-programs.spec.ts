import { test, expect } from '@playwright/test';
import { SettingsProgramsPage } from './pages';

test.describe('Settings - Programs Management', () => {
  let programsPage: SettingsProgramsPage;

  test.beforeEach(async ({ page }) => {
    programsPage = new SettingsProgramsPage(page);
    // Login as admin before each test
    await programsPage.login('admin');
    await programsPage.goto();
    await programsPage.expectPageLoaded();
  });

  test.describe('Page Load', () => {
    test('should display programs page with grid', async () => {
      await expect(programsPage.programsGrid).toBeVisible();
    });

    test('should display add program button', async () => {
      await expect(programsPage.addProgramButton).toBeVisible();
    });

    test('should display program list from database', async () => {
      // Verify at least one program is displayed
      const programIds = await programsPage.getAllProgramIds();
      expect(programIds.length).toBeGreaterThan(0);
    });
  });

  test.describe('Program Display', () => {
    test('should display program details correctly', async () => {
      const programIds = await programsPage.getAllProgramIds();
      expect(programIds.length).toBeGreaterThan(0);

      // Check first program displays all required fields
      const firstProgramId = programIds[0];
      await programsPage.expectProgramDisplayed(firstProgramId);

      // Verify additional fields are visible
      await expect(programsPage.getProgramStatus(firstProgramId)).toBeVisible();
      await expect(programsPage.getProgramSpp(firstProgramId)).toBeVisible();
      await expect(programsPage.getProgramStudents(firstProgramId)).toBeVisible();
    });

    test('should display edit and curriculum buttons for each program', async () => {
      const programIds = await programsPage.getAllProgramIds();
      expect(programIds.length).toBeGreaterThan(0);

      for (const programId of programIds) {
        await expect(programsPage.getProgramEditButton(programId)).toBeVisible();
        await expect(programsPage.getProgramCurriculumButton(programId)).toBeVisible();
      }
    });

    test('should display program code, name, and level', async () => {
      const programIds = await programsPage.getAllProgramIds();
      expect(programIds.length).toBeGreaterThan(0);

      for (const programId of programIds) {
        const code = await programsPage.getProgramCode(programId).textContent();
        const name = await programsPage.getProgramName(programId).textContent();
        const level = await programsPage.getProgramLevel(programId).textContent();

        expect(code).toBeTruthy();
        expect(name).toBeTruthy();
        expect(level).toBeTruthy();
        expect(level).toMatch(/^(S1|D3)$/);
      }
    });
  });

  test.describe('Program Status', () => {
    test('should display status badge for each program', async () => {
      const programIds = await programsPage.getAllProgramIds();
      expect(programIds.length).toBeGreaterThan(0);

      for (const programId of programIds) {
        const statusBadge = programsPage.getProgramStatus(programId);
        const statusText = await statusBadge.textContent();
        expect(statusText).toMatch(/^(Aktif|Nonaktif)$/);
      }
    });
  });
});
