import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class SettingsProgramsPage extends BasePage {
  readonly path = '/admin/settings/programs';

  async login(role: 'admin' | 'supervisor' | 'consultant' = 'admin'): Promise<void> {
    await this.page.goto(`/test/login/${role}`);
    // Wait for redirect to admin dashboard (with or without trailing slash)
    await this.page.waitForURL(/\/admin\/?$/);
  }

  // Page elements
  get pageContainer(): Locator {
    return this.page.getByTestId('settings-programs-page');
  }

  get programsGrid(): Locator {
    return this.page.getByTestId('programs-grid');
  }

  get addProgramButton(): Locator {
    return this.page.getByTestId('add-program-button');
  }

  // Program card elements by ID
  getProgramCard(programId: string): Locator {
    return this.page.getByTestId(`program-card-${programId}`);
  }

  getProgramCode(programId: string): Locator {
    return this.page.getByTestId(`program-code-${programId}`);
  }

  getProgramName(programId: string): Locator {
    return this.page.getByTestId(`program-name-${programId}`);
  }

  getProgramLevel(programId: string): Locator {
    return this.page.getByTestId(`program-level-${programId}`);
  }

  getProgramStatus(programId: string): Locator {
    return this.page.getByTestId(`program-status-${programId}`);
  }

  getProgramSpp(programId: string): Locator {
    return this.page.getByTestId(`program-spp-${programId}`);
  }

  getProgramStudents(programId: string): Locator {
    return this.page.getByTestId(`program-students-${programId}`);
  }

  getProgramEditButton(programId: string): Locator {
    return this.page.getByTestId(`program-edit-${programId}`);
  }

  getProgramCurriculumButton(programId: string): Locator {
    return this.page.getByTestId(`program-curriculum-${programId}`);
  }

  // Page assertions
  async expectPageLoaded(): Promise<void> {
    await expect(this.pageContainer).toBeVisible();
    await expect(this.programsGrid).toBeVisible();
  }

  // Get all program IDs from the grid
  async getAllProgramIds(): Promise<string[]> {
    const cards = this.programsGrid.locator('div[data-testid^="program-card-"]');
    const count = await cards.count();
    const ids: string[] = [];
    for (let i = 0; i < count; i++) {
      const testId = await cards.nth(i).getAttribute('data-testid');
      if (testId) {
        ids.push(testId.replace('program-card-', ''));
      }
    }
    return ids;
  }

  // Verification helpers
  async expectProgramDisplayed(programId: string): Promise<void> {
    await expect(this.getProgramCard(programId)).toBeVisible();
    await expect(this.getProgramCode(programId)).toBeVisible();
    await expect(this.getProgramName(programId)).toBeVisible();
    await expect(this.getProgramLevel(programId)).toBeVisible();
  }
}
