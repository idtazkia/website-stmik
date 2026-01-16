import { Page, Locator, expect } from '@playwright/test';
import { BasePage } from './BasePage';

export class RegistrationPage extends BasePage {
  // Registration form
  readonly registrationForm: Locator;
  readonly errorMessage: Locator;

  // Step 1: Account
  readonly inputEmail: Locator;
  readonly inputPhone: Locator;
  readonly inputPassword: Locator;
  readonly inputPasswordConfirm: Locator;
  readonly btnSubmitStep1: Locator;

  // Step 2: Personal Info
  readonly inputName: Locator;
  readonly inputAddress: Locator;
  readonly inputCity: Locator;
  readonly inputProvince: Locator;
  readonly btnSubmitStep2: Locator;

  // Step 3: Education
  readonly inputHighSchool: Locator;
  readonly selectGraduationYear: Locator;
  readonly btnSubmitStep3: Locator;

  // Step 4: Source Tracking
  readonly selectSourceType: Locator;
  readonly inputSourceDetail: Locator;
  readonly btnSubmitStep4: Locator;

  constructor(page: Page) {
    super(page);
    this.registrationForm = page.getByTestId('registration-form');
    this.errorMessage = page.getByTestId('error-message');

    // Step 1
    this.inputEmail = page.getByTestId('input-email');
    this.inputPhone = page.getByTestId('input-phone');
    this.inputPassword = page.getByTestId('input-password');
    this.inputPasswordConfirm = page.getByTestId('input-password-confirm');
    this.btnSubmitStep1 = page.getByTestId('btn-submit-step1');

    // Step 2
    this.inputName = page.getByTestId('input-name');
    this.inputAddress = page.getByTestId('input-address');
    this.inputCity = page.getByTestId('input-city');
    this.inputProvince = page.getByTestId('input-province');
    this.btnSubmitStep2 = page.getByTestId('btn-submit-step2');

    // Step 3
    this.inputHighSchool = page.getByTestId('input-high-school');
    this.selectGraduationYear = page.getByTestId('select-graduation-year');
    this.btnSubmitStep3 = page.getByTestId('btn-submit-step3');

    // Step 4
    this.selectSourceType = page.getByTestId('select-source-type');
    this.inputSourceDetail = page.getByTestId('input-source-detail');
    this.btnSubmitStep4 = page.getByTestId('btn-submit-step4');
  }

  async goto() {
    await this.page.goto('/register');
  }

  async gotoWithRef(refCode: string) {
    await this.page.goto(`/register?ref=${refCode}`);
  }

  async gotoWithCampaign(campaignCode: string) {
    await this.page.goto(`/register?utm_campaign=${campaignCode}`);
  }

  async expectPageLoaded() {
    await expect(this.registrationForm).toBeVisible();
  }

  async expectStep1Visible() {
    await expect(this.page.getByTestId('step1-form')).toBeVisible();
  }

  async expectStep2Visible() {
    await expect(this.page.getByTestId('step2-form')).toBeVisible();
  }

  async expectStep3Visible() {
    await expect(this.page.getByTestId('step3-form')).toBeVisible();
  }

  async expectStep4Visible() {
    await expect(this.page.getByTestId('step4-form')).toBeVisible();
  }

  async expectErrorMessage(message: string) {
    await expect(this.errorMessage).toBeVisible();
    await expect(this.errorMessage).toContainText(message);
  }

  async fillStep1WithEmail(email: string, password: string) {
    await this.inputEmail.fill(email);
    await this.inputPassword.fill(password);
    await this.inputPasswordConfirm.fill(password);
    await this.btnSubmitStep1.click();
  }

  async fillStep1WithPhone(phone: string, password: string) {
    await this.inputPhone.fill(phone);
    await this.inputPassword.fill(password);
    await this.inputPasswordConfirm.fill(password);
    await this.btnSubmitStep1.click();
  }

  async fillStep1WithBoth(email: string, phone: string, password: string) {
    await this.inputEmail.fill(email);
    await this.inputPhone.fill(phone);
    await this.inputPassword.fill(password);
    await this.inputPasswordConfirm.fill(password);
    await this.btnSubmitStep1.click();
  }

  async fillStep2(name: string, address: string, city: string, province: string) {
    await this.inputName.fill(name);
    await this.inputAddress.fill(address);
    await this.inputCity.fill(city);
    await this.inputProvince.fill(province);
    await this.btnSubmitStep2.click();
  }

  async fillStep3(highSchool: string, graduationYear: string, prodiCode: string) {
    await this.inputHighSchool.fill(highSchool);
    await this.selectGraduationYear.selectOption(graduationYear);
    await this.page.getByTestId(`radio-prodi-${prodiCode}`).click();
    await this.btnSubmitStep3.click();
  }

  async fillStep4(sourceType: string, sourceDetail?: string) {
    await this.selectSourceType.selectOption(sourceType);
    if (sourceDetail) {
      await this.inputSourceDetail.fill(sourceDetail);
    }
    await this.btnSubmitStep4.click();
  }

  async expectRedirectToPortal() {
    await expect(this.page).toHaveURL('/portal');
  }
}

export class LoginPage extends BasePage {
  readonly loginForm: Locator;
  readonly inputIdentifier: Locator;
  readonly inputPassword: Locator;
  readonly btnLogin: Locator;
  readonly errorMessage: Locator;

  constructor(page: Page) {
    super(page);
    this.loginForm = page.getByTestId('portal-login-form');
    this.inputIdentifier = page.getByTestId('input-identifier');
    this.inputPassword = page.getByTestId('input-password');
    this.btnLogin = page.getByTestId('btn-login');
    this.errorMessage = page.getByTestId('error-message');
  }

  async goto() {
    await this.page.goto('/login');
  }

  async expectPageLoaded() {
    await expect(this.loginForm).toBeVisible();
  }

  async login(identifier: string, password: string) {
    await this.inputIdentifier.fill(identifier);
    await this.inputPassword.fill(password);
    await this.btnLogin.click();
  }

  async expectErrorMessage(message: string) {
    await expect(this.errorMessage).toBeVisible();
    await expect(this.errorMessage).toContainText(message);
  }

  async expectRedirectToPortal() {
    await expect(this.page).toHaveURL('/portal');
  }
}
