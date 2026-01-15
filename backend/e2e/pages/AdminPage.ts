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

  get navProspects(): Locator {
    return this.page.getByTestId('nav-prospects');
  }

  get navApplications(): Locator {
    return this.page.getByTestId('nav-applications');
  }

  get navReferrers(): Locator {
    return this.page.getByTestId('nav-referrers');
  }

  get navCampaigns(): Locator {
    return this.page.getByTestId('nav-campaigns');
  }

  get navReports(): Locator {
    return this.page.getByTestId('nav-reports');
  }

  get navSettings(): Locator {
    return this.page.getByTestId('nav-settings');
  }

  async expectPageLoaded(): Promise<void> {
    await expect(this.sidebar).toBeVisible();
  }

  async expectNavigationVisible(): Promise<void> {
    await expect(this.navDashboard).toBeVisible();
    await expect(this.navProspects).toBeVisible();
    await expect(this.navApplications).toBeVisible();
    await expect(this.navReferrers).toBeVisible();
    await expect(this.navCampaigns).toBeVisible();
    await expect(this.navReports).toBeVisible();
    await expect(this.navSettings).toBeVisible();
  }
}
