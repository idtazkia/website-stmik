# Adding Lecturer Profiles

This document describes how to add new lecturer profiles to the STMIK Tazkia website using LinkedIn as the data source.

## Prerequisites

- Node.js and npm installed
- Playwright installed (`npx playwright install chromium`)
- LinkedIn account credentials
- Access to the lecturer's LinkedIn profile URL

## Method 1: Manual Scraping with Playwright Script

### 1. Login to LinkedIn (First Time Only)

The scraper needs an authenticated session. Run the login command:

```bash
npx tsx tests/scrape-linkedin.ts --login
```

This will:
- Open a browser window with LinkedIn login page
- Wait for you to log in manually
- Save the session to `.linkedin-session/` directory

The session is saved locally, so you only need to do this once (or when the session expires).

### 2. Run the Scraper

```bash
npx tsx tests/scrape-linkedin.ts "https://www.linkedin.com/in/username"
```

This will:
- Load the saved LinkedIn session
- Navigate to the profile
- Scroll to load all content
- Extract and print profile data (raw text + JSON)

### 3. Copy the Output

The script outputs:
1. Raw text content from the profile page
2. Structured JSON data with name, headline, experience, education, skills, certifications

### 4. Create the Markdown File

Using the scraped data, create a new file in `frontend/src/content/lecturers/`:

```bash
touch frontend/src/content/lecturers/lecturer-name.md
```

Then populate it with the profile data structure (see below).

**Note:** The `.linkedin-session/` directory is git-ignored to protect your credentials.

## Method 2: Using Claude Code

### 1. Prepare LinkedIn URL

Get the LinkedIn profile URL of the lecturer. Examples:
- `https://www.linkedin.com/in/username/`
- `https://id.linkedin.com/in/username`

### 2. Use Claude Code to Extract Profile

Run Claude Code and provide the LinkedIn URL:

```bash
claude
```

Then in the conversation:

```
Add lecturer profile from: https://www.linkedin.com/in/username
```

Claude Code will:
1. Fetch the LinkedIn profile using WebFetch
2. Extract relevant information (name, education, experience, skills)
3. Create a markdown file in `frontend/src/content/lecturers/`

### 3. Profile Data Structure

Lecturer profiles are stored as Markdown files with YAML frontmatter:

```markdown
---
name: Full Name
title: Dosen Tetap  # or "Dosen Luar Biasa"
position: Current Position at Company
expertise:
  - Skill 1
  - Skill 2
  - Skill 3
education:
  - degree: S2 Field of Study
    institution: University Name
    year: "2020-2022"  # optional
  - degree: S1 Field of Study
    institution: University Name
linkedin: https://www.linkedin.com/in/username/
github: github-username  # optional
website: https://example.com  # optional
email: email@example.com  # optional
youtube: https://youtube.com/@channel  # optional
photo: /images/lecturers/name.jpg  # optional
order: 1  # for sorting on listing page
---

## Tentang

Brief biography paragraph.

## Pengalaman Profesional

### Job Title - Company Name (Year - Year)
Description of role and achievements.

### Previous Job Title - Previous Company (Year - Year)
Description.

## Bidang Keahlian

### Category 1
- **Subcategory**: Details
- **Subcategory**: Details

### Category 2
- **Subcategory**: Details

## Pengajaran

Sebagai dosen tetap di STMIK Tazkia, [Name] mengajar mata kuliah:
- Course 1
- Course 2
- Course 3
```

### 4. File Naming Convention

Files should be named using kebab-case of the lecturer's name:
- `endy-muhardin.md`
- `hendri-karisma.md`
- `ricky-setiadi.md`

### 5. Redacting Company Names (Optional)

If the lecturer prefers to keep company names private, replace specific company names with generic descriptions:

| Original | Redacted |
|----------|----------|
| Blibli.com | E-commerce Company |
| tiket.com | E-commerce Travel Company |
| Gojek | Ride-hailing Company |
| Bank XYZ | Banking Institution |

When requesting profile creation, specify:
```
Add lecturer profile from: https://www.linkedin.com/in/username
Please redact company names in work history.
```

### 6. Verify Build

After creating the profile, verify the build succeeds:

```bash
cd frontend
npm run build
```

The lecturer should appear at:
- Indonesian: `/lecturers/lecturer-slug/`
- English: `/en/lecturers/lecturer-slug/`

### 7. Commit Changes

```bash
git add frontend/src/content/lecturers/new-lecturer.md
git commit -m "Add [Lecturer Name] lecturer profile"
```

## Current Lecturers

| Name | Expertise | File |
|------|-----------|------|
| Endy Muhardin | Software Engineering, DevOps | `endy-muhardin.md` |
| Hendri Karisma | AI/ML, Data Engineering | `hendri-karisma.md` |
| Q Fadlan | Data Science, AI | `q-fadlan.md` |
| Agus Sulaiman | DevOps, Cloud Engineering | `agus-sulaiman.md` |
| Ricky Setiadi | Information Security, Cybersecurity | `ricky-setiadi.md` |

## Schema Reference

The lecturer content collection schema is defined in `frontend/src/content/config.ts`:

```typescript
const lecturers = defineCollection({
  type: 'content',
  schema: z.object({
    name: z.string(),
    title: z.string(),
    position: z.string().optional(),
    expertise: z.array(z.string()),
    education: z.array(z.object({
      degree: z.string(),
      institution: z.string(),
      year: z.string().optional(),
    })),
    email: z.string().optional(),
    phone: z.string().optional(),
    website: z.string().optional(),
    github: z.string().optional(),
    linkedin: z.string().optional(),
    youtube: z.string().optional(),
    photo: z.string().optional(),
    order: z.number().default(999),
  }),
});
```

## Troubleshooting

### LinkedIn Profile Not Accessible
- Ensure the profile is public
- Try using the `id.linkedin.com` domain for Indonesian profiles

### Build Fails After Adding Profile
- Check YAML frontmatter syntax (proper indentation, quotes around special characters)
- Verify all required fields are present (`name`, `title`, `expertise`, `education`)
- Run `npm run typecheck` for detailed errors

### Profile Not Showing in Listing
- Check the `order` field - lower numbers appear first
- Verify the file extension is `.md`
- Ensure the file is in `frontend/src/content/lecturers/` directory
