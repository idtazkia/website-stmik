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

  // Add modal elements
  get addProgramModal(): Locator {
    return this.page.getByTestId('add-program-modal');
  }

  get inputProgramName(): Locator {
    return this.page.getByTestId('input-program-name');
  }

  get inputProgramCode(): Locator {
    return this.page.getByTestId('input-program-code');
  }

  get inputProgramDegree(): Locator {
    return this.page.getByTestId('input-program-degree');
  }

  get submitAddProgramButton(): Locator {
    return this.page.getByTestId('submit-add-program');
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

  getProgramStatusToggle(programId: string): Locator {
    return this.page.getByTestId(`program-status-toggle-${programId}`);
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

  // Edit modal elements by program ID
  getEditProgramModal(programId: string): Locator {
    return this.page.getByTestId(`edit-program-modal-${programId}`);
  }

  getEditProgramName(programId: string): Locator {
    return this.page.getByTestId(`edit-program-name-${programId}`);
  }

  getEditProgramCode(programId: string): Locator {
    return this.page.getByTestId(`edit-program-code-${programId}`);
  }

  getEditProgramDegree(programId: string): Locator {
    return this.page.getByTestId(`edit-program-degree-${programId}`);
  }

  getSubmitEditProgramButton(programId: string): Locator {
    return this.page.getByTestId(`submit-edit-program-${programId}`);
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

  // Actions
  async openAddProgramModal(): Promise<void> {
    await this.addProgramButton.click();
    await expect(this.addProgramModal).toBeVisible();
  }

  async addProgram(name: string, code: string, degree: 'S1' | 'D3'): Promise<void> {
    await this.openAddProgramModal();
    await this.inputProgramName.fill(name);
    await this.inputProgramCode.fill(code);
    await this.inputProgramDegree.selectOption(degree);

    // Start waiting for response before triggering the action
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes('/admin/settings/programs') &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.submitAddProgramButton.click();
    await responsePromise;
  }

  async openEditProgramModal(programId: string): Promise<void> {
    await this.getProgramEditButton(programId).click();
    await expect(this.getEditProgramModal(programId)).toBeVisible();
  }

  async editProgram(programId: string, name: string, code: string, degree: 'S1' | 'D3'): Promise<void> {
    await this.openEditProgramModal(programId);
    await this.getEditProgramName(programId).fill(name);
    await this.getEditProgramCode(programId).fill(code);
    await this.getEditProgramDegree(programId).selectOption(degree);

    // Start waiting for response before triggering the action
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/programs/${programId}`) &&
      response.request().method() === 'POST' &&
      !response.url().includes('/toggle-active') &&
      response.status() === 200
    );
    await this.getSubmitEditProgramButton(programId).click();
    await responsePromise;
  }

  async toggleProgramStatus(programId: string): Promise<void> {
    const button = this.getProgramStatusToggle(programId);
    // Start waiting for response before triggering the action
    const responsePromise = this.page.waitForResponse(response =>
      response.url().includes(`/admin/settings/programs/${programId}/toggle-active`) &&
      response.status() === 200
    );
    await button.click();
    await responsePromise;
  }

  // Status verification
  async expectProgramStatus(programId: string, expectedStatus: 'active' | 'inactive'): Promise<void> {
    const button = this.getProgramStatusToggle(programId);
    const text = await button.textContent();
    if (expectedStatus === 'active') {
      expect(text?.trim()).toBe('Aktif');
    } else {
      expect(text?.trim()).toBe('Nonaktif');
    }
  }

  async expectProgramName(programId: string, expectedName: string): Promise<void> {
    const name = await this.getProgramName(programId).textContent();
    expect(name?.trim()).toBe(expectedName);
  }

  async expectProgramCode(programId: string, expectedCode: string): Promise<void> {
    const code = await this.getProgramCode(programId).textContent();
    expect(code?.trim()).toBe(expectedCode);
  }

  async expectProgramLevel(programId: string, expectedLevel: string): Promise<void> {
    const level = await this.getProgramLevel(programId).textContent();
    expect(level?.trim()).toBe(expectedLevel);
  }
}
