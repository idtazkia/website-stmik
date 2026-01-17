import { Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class CandidateDetailPage extends BasePage {
  readonly path = '/admin/candidates';

  // Page container
  get detailPage(): Locator {
    return this.page.getByTestId('candidate-detail-page');
  }

  // Header elements
  get candidateHeader(): Locator {
    return this.page.getByTestId('candidate-header');
  }

  get backLink(): Locator {
    return this.page.getByTestId('back-link');
  }

  get candidateName(): Locator {
    return this.page.getByTestId('candidate-name');
  }

  // Action buttons
  get actionButtons(): Locator {
    return this.page.getByTestId('action-buttons');
  }

  get btnLogInteraction(): Locator {
    return this.page.getByTestId('btn-log-interaction');
  }

  get btnCommit(): Locator {
    return this.page.getByTestId('btn-commit');
  }

  get btnEnroll(): Locator {
    return this.page.getByTestId('btn-enroll');
  }

  get btnMarkLost(): Locator {
    return this.page.getByTestId('btn-mark-lost');
  }

  // Personal Info Section
  get sectionPersonalInfo(): Locator {
    return this.page.getByTestId('section-personal-info');
  }

  get sectionTitlePersonalInfo(): Locator {
    return this.page.getByTestId('section-title-personal-info');
  }

  get fieldEmail(): Locator {
    return this.page.getByTestId('field-email');
  }

  get fieldPhone(): Locator {
    return this.page.getByTestId('field-phone');
  }

  get fieldWhatsapp(): Locator {
    return this.page.getByTestId('field-whatsapp');
  }

  get fieldAddress(): Locator {
    return this.page.getByTestId('field-address');
  }

  // Education Section
  get sectionEducation(): Locator {
    return this.page.getByTestId('section-education');
  }

  get sectionTitleEducation(): Locator {
    return this.page.getByTestId('section-title-education');
  }

  get fieldHighSchool(): Locator {
    return this.page.getByTestId('field-high-school');
  }

  get fieldGraduationYear(): Locator {
    return this.page.getByTestId('field-graduation-year');
  }

  get fieldProdi(): Locator {
    return this.page.getByTestId('field-prodi');
  }

  // Source & Assignment Section
  get sectionSourceAssignment(): Locator {
    return this.page.getByTestId('section-source-assignment');
  }

  get sectionTitleSourceAssignment(): Locator {
    return this.page.getByTestId('section-title-source-assignment');
  }

  get fieldSourceInfo(): Locator {
    return this.page.getByTestId('field-source-info');
  }

  get fieldCampaign(): Locator {
    return this.page.getByTestId('field-campaign');
  }

  get fieldReferrer(): Locator {
    return this.page.getByTestId('field-referrer');
  }

  get fieldConsultant(): Locator {
    return this.page.getByTestId('field-consultant');
  }

  get btnReassign(): Locator {
    return this.page.getByTestId('btn-reassign');
  }

  // Payment Status Section
  get sectionPaymentStatus(): Locator {
    return this.page.getByTestId('section-payment-status');
  }

  get sectionTitlePaymentStatus(): Locator {
    return this.page.getByTestId('section-title-payment-status');
  }

  get fieldRegistrationFee(): Locator {
    return this.page.getByTestId('field-registration-fee');
  }

  get feeStatusPaid(): Locator {
    return this.page.getByTestId('fee-status-paid');
  }

  get feeStatusUnpaid(): Locator {
    return this.page.getByTestId('fee-status-unpaid');
  }

  // Documents Section
  get sectionDocuments(): Locator {
    return this.page.getByTestId('section-documents');
  }

  get sectionTitleDocuments(): Locator {
    return this.page.getByTestId('section-title-documents');
  }

  get docKtp(): Locator {
    return this.page.getByTestId('doc-ktp');
  }

  get docFoto(): Locator {
    return this.page.getByTestId('doc-foto');
  }

  get docIjazah(): Locator {
    return this.page.getByTestId('doc-ijazah');
  }

  get docTranskrip(): Locator {
    return this.page.getByTestId('doc-transkrip');
  }

  // Timeline Section
  get sectionTimeline(): Locator {
    return this.page.getByTestId('section-timeline');
  }

  get sectionTitleTimeline(): Locator {
    return this.page.getByTestId('section-title-timeline');
  }

  get timelineContent(): Locator {
    return this.page.getByTestId('timeline-content');
  }

  get timelineEmpty(): Locator {
    return this.page.getByTestId('timeline-empty');
  }

  get timelineList(): Locator {
    return this.page.getByTestId('timeline-list');
  }

  getInteraction(id: string): Locator {
    return this.page.getByTestId(`interaction-${id}`);
  }

  // Interaction Modal
  get modalInteraction(): Locator {
    return this.page.getByTestId('modal-interaction');
  }

  get modalTitle(): Locator {
    return this.page.getByTestId('modal-title');
  }

  get modalClose(): Locator {
    return this.page.getByTestId('modal-close');
  }

  get formInteraction(): Locator {
    return this.page.getByTestId('form-interaction');
  }

  get selectChannel(): Locator {
    return this.page.getByTestId('select-channel');
  }

  get selectCategory(): Locator {
    return this.page.getByTestId('select-category');
  }

  get selectObstacle(): Locator {
    return this.page.getByTestId('select-obstacle');
  }

  get inputRemarks(): Locator {
    return this.page.getByTestId('input-remarks');
  }

  get inputNextFollowup(): Locator {
    return this.page.getByTestId('input-next-followup');
  }

  get btnCancel(): Locator {
    return this.page.getByTestId('btn-cancel');
  }

  get btnSubmit(): Locator {
    return this.page.getByTestId('btn-submit');
  }

  // Helper methods
  async expectPageLoaded(): Promise<void> {
    await expect(this.detailPage).toBeVisible();
    await expect(this.candidateHeader).toBeVisible();
    await expect(this.sectionPersonalInfo).toBeVisible();
  }

  async goToCandidateDetail(candidateId: string): Promise<void> {
    await this.page.goto(`/admin/candidates/${candidateId}`);
  }

  async openInteractionModal(): Promise<void> {
    await this.btnLogInteraction.click();
    await expect(this.modalInteraction).toBeVisible();
  }

  async closeInteractionModal(): Promise<void> {
    await this.modalClose.click();
    await expect(this.modalInteraction).toBeHidden();
  }

  async goBackToCandidatesList(): Promise<void> {
    await this.backLink.click();
    await this.page.waitForURL(/\/admin\/candidates\/?$/);
  }
}
