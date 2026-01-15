import { Page, Locator, expect } from '@playwright/test';

export abstract class BasePage {
  constructor(protected readonly page: Page) {}

  abstract readonly path: string;

  async goto(): Promise<void> {
    await this.page.goto(this.path);
  }

  get versionInfo(): Locator {
    return this.page.getByTestId('version-info');
  }

  async expectVersionVisible(): Promise<void> {
    await expect(this.versionInfo).toBeVisible();
    const text = await this.versionInfo.textContent();
    expect(text).toBeTruthy();
    expect(text!.length).toBeGreaterThan(0);
  }
}
