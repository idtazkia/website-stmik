# User Manual Creation Guidelines

Guidelines for writing and maintaining user manual markdown files and their corresponding section definitions in `cmd/usermanual/main.go`.

## Architecture

```
docs/user-manual/*.md          → Markdown source files
cmd/usermanual/main.go         → Section definitions + HTML generation
test/e2e/screenshot_test.ts    → Screenshot capture via Playwright
Functional tests               → Auto-capture screenshots via takeManualScreenshot()
```

The generator reads markdown files, extracts specific H2 sections based on section definitions, and produces a single-page HTML manual.

## Section Extraction Rules

`extractSectionContent()` extracts content from a markdown file based on the section title. It applies these rules in order:

1. **H2 match**: If the section title matches an `## H2` heading in the file, return that H2's content (up to the next H2 or end of file).

2. **H1 match / fallback (aggregate)**: If no H2 matches, return content from the start of the file (after H1), **excluding** any H2 sections that have their own separate section definitions (sibling sections referencing the same file).

### CRITICAL: Avoid Duplicate Section Rendering

**Problem**: When multiple `Section` entries reference the same markdown file, a section whose title matches the file's H1 heading will aggregate all H2 content — EXCEPT H2 sections that are separately defined as their own `Section` entries. If a new H2 is added to the markdown but not registered as a separate section, it will be included in the H1-matching section's content.

**Rules to prevent duplicates**:

1. **Each section title must match exactly one H2 heading** in the markdown file. The title matching uses `titlesMatch()` (case-insensitive, contains, keyword overlap).

2. **When multiple sections reference the same file**, every H2 that should be rendered separately MUST have its own `Section` entry. The "primary" section (matching H1 or acting as catch-all) will automatically exclude H2 sections that have their own entries.

3. **Never use the H1 title as a section title** if you also want to split H2s into separate sections AND want precise control over which H2s are included in the primary section. The H1-matching section acts as a catch-all for un-mapped H2s.

4. **After adding a new H2 to a markdown file**, check `getSectionGroups()` to verify it will be rendered correctly — either as part of an existing aggregate section or as its own new `Section` entry.

### Example: Correct Multi-Section File

```go
// 04-interaksi.md has H1 "Pencatatan Interaksi"
// and H2s: "Mencatat Interaksi", "Kategori Interaksi", ..., "Mark as Lost", ...

Section{ID: "catat-interaksi", Title: "Mencatat Interaksi", File: "04-interaksi.md"},
Section{ID: "mark-lost", Title: "Mark as Lost", File: "04-interaksi.md"},
```

Result:
- Section "Mencatat Interaksi" renders only that specific H2
- Section "Mark as Lost" renders only that specific H2

### Example: Single-Section File

```go
// 12-keamanan.md — entire file rendered as one section
Section{ID: "keamanan", Title: "Keamanan & Enkripsi", File: "12-keamanan.md"}
```

Result: Entire file content rendered (no H2 exclusions needed since no siblings).

## Title Matching Rules

`titlesMatch()` uses flexible matching:

1. **Exact** (case-insensitive): "Mencatat Interaksi" == "mencatat interaksi"
2. **Contains**: "Upload Dokumen" contains/is-contained-by the H2 title
3. **Keyword overlap**: All significant words (length >= 4) from the shorter title must appear in the longer title

Be careful with short or generic titles that might accidentally match multiple H2 headings.

## Adding Screenshots

1. Define screenshot routes in `test/e2e/screenshot-capture.spec.ts`
2. Run screenshot capture: `npx playwright test --config=playwright.testrunner.config.ts screenshot-capture`
3. Screenshots saved to `docs/user-manual/screenshots/`
4. Add screenshot IDs to the `Section.Screenshots` list in `cmd/usermanual/main.go`
5. Reference in markdown: `![Alt text](screenshots/section/screenshot-id.png)`

## Markdown File Conventions

- H1 (`#`) is the file title — stripped during extraction
- H2 (`##`) are the section boundaries used for extraction
- H3+ are subsections within an H2 — included with their parent H2
- `---` horizontal rules are visual separators, not section boundaries
- Image references use relative paths: `screenshots/category/name.png`

## Checklist: Adding a New Section

1. Write content in the appropriate `docs/user-manual/*.md` file
2. Add `Section` entry in `getSectionGroups()` with title matching the H2 heading
3. If the markdown file already has other sections defined, verify no title conflicts
4. Add screenshot definitions if needed (Playwright test)
5. Run `make user-manual` locally to verify output
6. Check the generated HTML for duplicate content before committing

## Running the Generator

```bash
# Generate user manual
make user-manual

# Output: docs/user-manual/output/index.html
```

## File Structure

```
docs/user-manual/
├── 01-login-dashboard.md      → Login & Dashboard
├── 02-registrasi.md           → Registrasi Calon Mahasiswa
├── 03-manajemen-kandidat.md   → Manajemen Kandidat
├── 04-interaksi.md            → Pencatatan Interaksi
├── 05-dokumen.md              → Manajemen Dokumen
├── 06-keuangan.md             → Keuangan & Pembayaran
├── 07-marketing-referral.md   → Marketing & Referral
├── 08-portal.md               → Portal Calon Mahasiswa
├── 09-laporan.md              → Laporan & Analisis
├── 10-pengumuman.md           → Pengumuman
├── 11-pengaturan.md           → Pengaturan Sistem
├── 12-keamanan.md             → Keamanan & Enkripsi
├── screenshots/               → Auto-captured screenshots
│   ├── admin/                 → Admin panel screenshots
│   └── portal/                → Portal screenshots
├── output/                    → Generated HTML (gitignored)
│   └── index.html
└── user-manual-creation-guidelines.md → This file
```
