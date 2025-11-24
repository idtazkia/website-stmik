/**
 * LinkedIn Profile Scraper using Playwright
 *
 * Usage:
 *   1. First time - Login and save session:
 *      npx tsx tests/scrape-linkedin.ts --login
 *
 *   2. Scrape a profile (after login):
 *      npx tsx tests/scrape-linkedin.ts "https://www.linkedin.com/in/username"
 *
 * The script saves browser session to .linkedin-session/ directory
 * so you only need to login once.
 */

import { chromium, BrowserContext } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';

const SESSION_DIR = path.join(process.cwd(), '.linkedin-session');
const STORAGE_STATE_PATH = path.join(SESSION_DIR, 'storage-state.json');

async function ensureSessionDir() {
  if (!fs.existsSync(SESSION_DIR)) {
    fs.mkdirSync(SESSION_DIR, { recursive: true });
  }
}

async function loginToLinkedIn() {
  console.log('\n=== LinkedIn Login ===\n');
  console.log('A browser window will open. Please log in to LinkedIn manually.');
  console.log('After logging in, the session will be saved automatically.\n');

  await ensureSessionDir();

  const browser = await chromium.launch({
    headless: false,
  });

  const context = await browser.newContext({
    viewport: { width: 1280, height: 800 },
  });

  const page = await context.newPage();

  await page.goto('https://www.linkedin.com/login');

  console.log('Waiting for you to log in...');
  console.log('(The script will continue automatically once you reach the feed page)\n');

  // Wait for successful login - user reaches the feed
  await page.waitForURL('**/feed/**', { timeout: 300000 }); // 5 minutes timeout

  console.log('Login successful! Saving session...');

  // Save the session state
  await context.storageState({ path: STORAGE_STATE_PATH });

  console.log(`Session saved to: ${STORAGE_STATE_PATH}`);
  console.log('\nYou can now run the scraper without logging in again.');

  await browser.close();
}

async function scrapeLinkedInProfile(url: string) {
  console.log(`\n=== Scraping LinkedIn Profile ===`);
  console.log(`URL: ${url}\n`);

  // Check if session exists
  if (!fs.existsSync(STORAGE_STATE_PATH)) {
    console.error('Error: No LinkedIn session found.');
    console.error('Please run with --login first to authenticate:\n');
    console.error('  npx tsx tests/scrape-linkedin.ts --login\n');
    process.exit(1);
  }

  const browser = await chromium.launch({
    headless: false, // Set to true for headless scraping after testing
  });

  const context = await browser.newContext({
    storageState: STORAGE_STATE_PATH,
    viewport: { width: 1280, height: 800 },
  });

  const page = await context.newPage();

  try {
    // Navigate to LinkedIn profile
    await page.goto(url, { waitUntil: 'networkidle', timeout: 30000 });

    // Check if we're redirected to login (session expired)
    if (page.url().includes('/login')) {
      console.error('\nError: Session expired. Please login again:');
      console.error('  npx tsx tests/scrape-linkedin.ts --login\n');
      await browser.close();
      process.exit(1);
    }

    // Wait for profile to load
    await page.waitForSelector('h1', { timeout: 10000 });

    // Scroll down to load more content
    await page.evaluate(async () => {
      for (let i = 0; i < 5; i++) {
        window.scrollBy(0, 800);
        await new Promise(resolve => setTimeout(resolve, 500));
      }
      window.scrollTo(0, 0);
    });

    // Wait a bit for dynamic content
    await page.waitForTimeout(2000);

    // Extract profile data
    console.log('=== Extracting profile data ===\n');

    const profileData = await page.evaluate(() => {
      // Get main profile info
      const name = document.querySelector('h1')?.textContent?.trim() || '';
      const headline = document.querySelector('.text-body-medium')?.textContent?.trim() || '';

      // Get About section
      const aboutSection = document.querySelector('#about')?.closest('section');
      const about = aboutSection?.querySelector('.inline-show-more-text')?.textContent?.trim() || '';

      // Get Experience
      const experienceSection = document.querySelector('#experience')?.closest('section');
      const experiences: any[] = [];

      if (experienceSection) {
        const expItems = experienceSection.querySelectorAll(':scope > div > ul > li');
        expItems.forEach(item => {
          const title = item.querySelector('.t-bold span')?.textContent?.trim() || '';
          const company = item.querySelector('.t-normal span')?.textContent?.trim() || '';
          const duration = item.querySelector('.t-normal.t-black--light span')?.textContent?.trim() || '';
          const description = item.querySelector('.inline-show-more-text')?.textContent?.trim() || '';

          if (title || company) {
            experiences.push({ title, company, duration, description });
          }
        });
      }

      // Get Education
      const educationSection = document.querySelector('#education')?.closest('section');
      const education: any[] = [];

      if (educationSection) {
        const eduItems = educationSection.querySelectorAll(':scope > div > ul > li');
        eduItems.forEach(item => {
          const institution = item.querySelector('.t-bold span')?.textContent?.trim() || '';
          const degree = item.querySelector('.t-normal span')?.textContent?.trim() || '';
          const years = item.querySelector('.t-normal.t-black--light span')?.textContent?.trim() || '';

          if (institution) {
            education.push({ institution, degree, years });
          }
        });
      }

      // Get Skills
      const skillsSection = document.querySelector('#skills')?.closest('section');
      const skills: string[] = [];

      if (skillsSection) {
        const skillItems = skillsSection.querySelectorAll('.t-bold span[aria-hidden="true"]');
        skillItems.forEach(item => {
          const skill = item.textContent?.trim();
          if (skill) skills.push(skill);
        });
      }

      // Get Certifications
      const certsSection = document.querySelector('#licenses_and_certifications')?.closest('section');
      const certifications: any[] = [];

      if (certsSection) {
        const certItems = certsSection.querySelectorAll(':scope > div > ul > li');
        certItems.forEach(item => {
          const certName = item.querySelector('.t-bold span')?.textContent?.trim() || '';
          const issuer = item.querySelector('.t-normal span')?.textContent?.trim() || '';

          if (certName) {
            certifications.push({ name: certName, issuer });
          }
        });
      }

      return {
        name,
        headline,
        about,
        experiences,
        education,
        skills,
        certifications,
      };
    });

    // Also get raw text content for backup
    const rawContent = await page.textContent('main') || '';

    console.log(rawContent);

    console.log('\n=== Structured Data (JSON) ===\n');
    console.log(JSON.stringify(profileData, null, 2));

    // Save session state again (in case cookies were refreshed)
    await context.storageState({ path: STORAGE_STATE_PATH });

    return { profileData, rawContent };

  } catch (error) {
    console.error('Error scraping profile:', error);
    throw error;
  } finally {
    await browser.close();
  }
}

// Main execution
const args = process.argv.slice(2);

if (args.includes('--login')) {
  loginToLinkedIn()
    .then(() => {
      console.log('\n=== Login completed ===');
      process.exit(0);
    })
    .catch((error) => {
      console.error('Login failed:', error);
      process.exit(1);
    });
} else if (args.length > 0 && args[0].includes('linkedin.com')) {
  scrapeLinkedInProfile(args[0])
    .then(() => {
      console.log('\n=== Scraping completed ===');
      process.exit(0);
    })
    .catch((error) => {
      console.error('Scraping failed:', error);
      process.exit(1);
    });
} else {
  console.log(`
LinkedIn Profile Scraper

Usage:
  1. Login first (one-time):
     npx tsx tests/scrape-linkedin.ts --login

  2. Scrape a profile:
     npx tsx tests/scrape-linkedin.ts "https://www.linkedin.com/in/username"

Options:
  --login    Open browser to login and save session
  <url>      LinkedIn profile URL to scrape
`);
  process.exit(0);
}
