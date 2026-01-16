import { Page, Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

// Test portal page (for CSRF tests)
export class TestPortalPage extends BasePage {
  readonly path = '/test/portal';

  get header(): Locator {
    return this.page.getByTestId('portal-header');
  }

  get pageHeading(): Locator {
    return this.page.getByTestId('page-heading');
  }

  get testForm(): Locator {
    return this.page.getByTestId('test-form');
  }

  get testInput(): Locator {
    return this.page.getByTestId('test-input');
  }

  get submitButton(): Locator {
    return this.page.getByTestId('test-submit');
  }

  get navDashboard(): Locator {
    return this.page.getByTestId('nav-dashboard');
  }

  get navApplication(): Locator {
    return this.page.getByTestId('nav-application');
  }

  get navDocuments(): Locator {
    return this.page.getByTestId('nav-documents');
  }

  get logoutButton(): Locator {
    return this.page.getByTestId('btn-logout');
  }

  async expectPageLoaded(): Promise<void> {
    await expect(this.header).toBeVisible();
    await expect(this.pageHeading).toHaveText('Test Portal Page');
  }

  async expectNavigationVisible(): Promise<void> {
    await expect(this.navDashboard).toBeVisible();
    await expect(this.navApplication).toBeVisible();
    await expect(this.navDocuments).toBeVisible();
    await expect(this.logoutButton).toBeVisible();
  }

  async fillAndSubmitForm(value: string): Promise<Response | null> {
    await this.testInput.fill(value);
    const responsePromise = this.page.waitForResponse('/test/submit');
    await this.submitButton.click();
    return responsePromise;
  }
}

// Actual candidate portal page
export class PortalPage extends BasePage {
  readonly path = '/portal';

  // Welcome banner
  readonly welcomeBanner: Locator;
  readonly candidateName: Locator;
  readonly programName: Locator;
  readonly statusBadge: Locator;

  // Checklist
  readonly checklist: Locator;
  readonly checklistItems: Locator;

  // Announcements
  readonly announcements: Locator;

  // Consultant info
  readonly consultantSection: Locator;
  readonly consultantName: Locator;
  readonly emailButton: Locator;

  // Navigation
  readonly navDocuments: Locator;
  readonly navPayments: Locator;
  readonly navAnnouncements: Locator;
  readonly logoutButton: Locator;

  constructor(page: Page) {
    super(page);
    this.welcomeBanner = page.locator('.bg-gradient-to-r.from-primary-600');
    this.candidateName = this.welcomeBanner.locator('h1');
    this.programName = this.welcomeBanner.locator('p').first();
    this.statusBadge = this.welcomeBanner.locator('span.inline-block');

    this.checklist = page.locator('text=Checklist Pendaftaran').locator('..');
    this.checklistItems = page.locator('[class*="flex items-center justify-between p-3 rounded-lg"]');

    this.announcements = page.locator('text=Pengumuman Terbaru').locator('..');

    this.consultantSection = page.locator('text=Konsultan Anda').locator('..').locator('..');
    this.consultantName = this.consultantSection.locator('p.font-medium');
    this.emailButton = page.locator('text=Hubungi via Email');

    this.navDocuments = page.locator('a[href="/portal/documents"]').first();
    this.navPayments = page.locator('a[href="/portal/payments"]').first();
    this.navAnnouncements = page.locator('a[href="/portal/announcements"]').last();
    this.logoutButton = page.getByTestId('btn-logout');
  }

  async expectPageLoaded() {
    await expect(this.welcomeBanner).toBeVisible();
    await expect(this.checklist).toBeVisible();
  }

  async expectWelcomeMessage(name: string) {
    await expect(this.candidateName).toContainText(name);
  }

  async expectProgram(programName: string) {
    await expect(this.programName).toContainText(programName);
  }

  async expectStatus(status: string) {
    await expect(this.statusBadge).toContainText(status);
  }

  async expectChecklistVisible() {
    await expect(this.checklist).toBeVisible();
  }

  async getChecklistItemCount(): Promise<number> {
    return await this.checklistItems.count();
  }

  async expectConsultantVisible(name: string) {
    await expect(this.consultantName).toContainText(name);
    await expect(this.emailButton).toBeVisible();
  }

  async expectNoConsultantAssigned() {
    await expect(this.page.locator('text=Konsultan akan segera ditugaskan')).toBeVisible();
  }

  async expectAnnouncementsVisible() {
    await expect(this.announcements).toBeVisible();
  }

  async clickDocuments() {
    await this.navDocuments.click();
  }

  async clickPayments() {
    await this.navPayments.click();
  }

  async logout() {
    await this.logoutButton.click();
  }
}
