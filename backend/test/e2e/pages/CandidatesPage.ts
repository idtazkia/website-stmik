import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class CandidatesPage extends BasePage {
  readonly path = '/admin/candidates';

  // Page sections
  get candidatesPage(): Locator {
    return this.page.getByTestId('candidates-page');
  }

  get statsSection(): Locator {
    return this.page.getByTestId('candidate-stats');
  }

  get filtersSection(): Locator {
    return this.page.getByTestId('filters-section');
  }

  get candidatesTable(): Locator {
    return this.page.getByTestId('candidates-table');
  }

  get candidatesList(): Locator {
    return this.page.getByTestId('candidates-list');
  }

  get pagination(): Locator {
    return this.page.getByTestId('pagination');
  }

  // Stats
  get statTotal(): Locator {
    return this.page.getByTestId('stat-total');
  }

  get statRegistered(): Locator {
    return this.page.getByTestId('stat-registered');
  }

  get statProspecting(): Locator {
    return this.page.getByTestId('stat-prospecting');
  }

  get statCommitted(): Locator {
    return this.page.getByTestId('stat-committed');
  }

  get statEnrolled(): Locator {
    return this.page.getByTestId('stat-enrolled');
  }

  get statLost(): Locator {
    return this.page.getByTestId('stat-lost');
  }

  // Filters
  get filterStatus(): Locator {
    return this.page.getByTestId('filter-status');
  }

  get filterConsultant(): Locator {
    return this.page.getByTestId('filter-consultant');
  }

  get filterProdi(): Locator {
    return this.page.getByTestId('filter-prodi');
  }

  get filterCampaign(): Locator {
    return this.page.getByTestId('filter-campaign');
  }

  get filterSource(): Locator {
    return this.page.getByTestId('filter-source');
  }

  get filterSearch(): Locator {
    return this.page.getByTestId('filter-search');
  }

  // Pagination
  get prevPageButton(): Locator {
    return this.page.getByTestId('prev-page');
  }

  get nextPageButton(): Locator {
    return this.page.getByTestId('next-page');
  }

  // Helper methods
  async expectPageLoaded(): Promise<void> {
    await expect(this.candidatesPage).toBeVisible();
    await expect(this.statsSection).toBeVisible();
    await expect(this.filtersSection).toBeVisible();
    await expect(this.candidatesTable).toBeVisible();
  }

  getCandidateRow(id: string): Locator {
    return this.page.getByTestId(`candidate-row-${id}`);
  }

  getViewCandidateLink(id: string): Locator {
    return this.page.getByTestId(`view-candidate-${id}`);
  }

  async getAllCandidateRows(): Promise<Locator[]> {
    return this.candidatesList.locator('tr[data-testid^="candidate-row-"]').all();
  }

  async getCandidateRowCount(): Promise<number> {
    const rows = await this.getAllCandidateRows();
    return rows.length;
  }

  async selectStatus(status: string): Promise<void> {
    await this.filterStatus.selectOption(status);
    // Wait for HTMX request to complete
    await this.page.waitForResponse(resp => resp.url().includes('/admin/candidates') && resp.status() === 200);
  }

  async selectConsultant(consultantId: string): Promise<void> {
    await this.filterConsultant.selectOption(consultantId);
    await this.page.waitForResponse(resp => resp.url().includes('/admin/candidates') && resp.status() === 200);
  }

  async selectProdi(prodiId: string): Promise<void> {
    await this.filterProdi.selectOption(prodiId);
    await this.page.waitForResponse(resp => resp.url().includes('/admin/candidates') && resp.status() === 200);
  }

  async selectCampaign(campaignId: string): Promise<void> {
    await this.filterCampaign.selectOption(campaignId);
    await this.page.waitForResponse(resp => resp.url().includes('/admin/candidates') && resp.status() === 200);
  }

  async selectSourceType(sourceType: string): Promise<void> {
    await this.filterSource.selectOption(sourceType);
    await this.page.waitForResponse(resp => resp.url().includes('/admin/candidates') && resp.status() === 200);
  }

  async searchCandidates(query: string): Promise<void> {
    // Set up response listener before filling the input
    const responsePromise = this.page.waitForResponse(
      resp => resp.url().includes('/admin/candidates') && resp.status() === 200,
      { timeout: 10000 }
    );
    await this.filterSearch.fill(query);
    // Wait for debounced HTMX request to complete
    await responsePromise;
  }

  async clearFilters(): Promise<void> {
    await this.filterStatus.selectOption('');
    await this.filterConsultant.selectOption('');
    await this.filterProdi.selectOption('');
    await this.filterCampaign.selectOption('');
    await this.filterSource.selectOption('');
    await this.filterSearch.fill('');
  }

  async goToNextPage(): Promise<void> {
    await this.nextPageButton.click();
    await this.page.waitForResponse(resp => resp.url().includes('/admin/candidates') && resp.status() === 200);
  }

  async goToPrevPage(): Promise<void> {
    await this.prevPageButton.click();
    await this.page.waitForResponse(resp => resp.url().includes('/admin/candidates') && resp.status() === 200);
  }

  async expectStatValue(stat: 'total' | 'registered' | 'prospecting' | 'committed' | 'enrolled' | 'lost', value: string): Promise<void> {
    const statLocator = {
      total: this.statTotal,
      registered: this.statRegistered,
      prospecting: this.statProspecting,
      committed: this.statCommitted,
      enrolled: this.statEnrolled,
      lost: this.statLost,
    }[stat];
    await expect(statLocator).toHaveText(value);
  }

  async expectCandidateInList(id: string): Promise<void> {
    await expect(this.getCandidateRow(id)).toBeVisible();
  }

  async expectCandidateNotInList(id: string): Promise<void> {
    await expect(this.getCandidateRow(id)).not.toBeVisible();
  }

  async expectEmptyList(): Promise<void> {
    await expect(this.page.getByTestId('empty-candidates-message')).toBeVisible();
  }

  async viewCandidateDetail(id: string): Promise<void> {
    await this.getViewCandidateLink(id).click();
    await this.page.waitForURL(/\/admin\/candidates\/[a-f0-9-]+/);
  }
}
