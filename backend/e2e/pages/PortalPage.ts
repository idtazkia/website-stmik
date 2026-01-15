import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class PortalPage extends BasePage {
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
