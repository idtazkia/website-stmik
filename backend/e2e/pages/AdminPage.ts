import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class AdminPage extends BasePage {
  readonly path = '/test/admin';

  get sidebar(): Locator {
    return this.page.getByTestId('admin-sidebar');
  }

  get navDashboard(): Locator {
    return this.page.getByTestId('nav-dashboard');
  }

  get navCandidates(): Locator {
    return this.page.getByTestId('nav-candidates');
  }

  get navDocuments(): Locator {
    return this.page.getByTestId('nav-documents');
  }

  get navCampaigns(): Locator {
    return this.page.getByTestId('nav-campaigns');
  }

  get navReferrers(): Locator {
    return this.page.getByTestId('nav-referrers');
  }

  get navUsers(): Locator {
    return this.page.getByTestId('nav-users');
  }

  get navPrograms(): Locator {
    return this.page.getByTestId('nav-programs');
  }

  async expectPageLoaded(): Promise<void> {
    await expect(this.sidebar).toBeVisible();
  }

  async expectNavigationVisible(): Promise<void> {
    await expect(this.navDashboard).toBeVisible();
    await expect(this.navCandidates).toBeVisible();
    await expect(this.navDocuments).toBeVisible();
    await expect(this.navCampaigns).toBeVisible();
    await expect(this.navReferrers).toBeVisible();
    await expect(this.navUsers).toBeVisible();
    await expect(this.navPrograms).toBeVisible();
  }
}
