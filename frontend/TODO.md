# Frontend TODO - Marketing Site Phase

## Current Status
**Phase 3: Marketing Site - 30% Complete**

### ✅ Completed
- Astro 5.x + Tailwind CSS 4.x setup
- Custom i18n system (Indonesian/English)
- BaseLayout and MarketingLayout
- Reusable components (Card, Button, Container, Section, Header, Footer, Navigation)
- Homepage (ID/EN)
- About page (ID/EN)
- Lecturer profiles (ID/EN) with content collections
- Responsive design
- SEO optimization (meta tags, sitemap, Open Graph)
- Cloudflare Pages deployment

### 🚧 Next Steps (Phase 3)
- [ ] Programs listing page
- [ ] Programs detail pages
- [ ] Contact page with form
- [ ] Admissions information page
- [ ] News/blog system

---

## Phase 3: Marketing Pages (Static Content)

### Content Collections
- [x] Lecturers collection
- [ ] Programs collection (content exists, pages needed)
- [ ] Admissions collection (content exists, pages needed)
- [ ] About collection (content exists, pages needed)
- [ ] News/blog collection

### Programs Pages
- [ ] Create `src/pages/programs/index.astro`
  - [ ] List all programs from content collection
  - [ ] Filter/search functionality
  - [ ] Program cards with images
- [ ] Create `src/pages/programs/[slug].astro`
  - [ ] Program details from markdown
  - [ ] Curriculum overview
  - [ ] Career prospects
  - [ ] Admission requirements
  - [ ] CTA: "Apply Now" button
- [ ] Create English versions: `src/pages/en/programs/...`

### Contact Page
- [ ] Create `src/pages/contact.astro`
  - [ ] Contact information (phone, email, address)
  - [ ] Embedded map (Google Maps)
  - [ ] Social media links
  - [ ] Contact form (static - form submission deferred to Phase 2)
  - [ ] Office hours
- [ ] Create `src/pages/en/contact.astro`

### Admissions Page
- [ ] Create `src/pages/admissions.astro`
  - [ ] Admission requirements
  - [ ] Application process steps
  - [ ] Important dates/calendar
  - [ ] Required documents list
  - [ ] FAQ section
  - [ ] CTA: "Apply Now" (links to application portal when ready)
- [ ] Create `src/pages/en/admissions.astro`

### News/Blog System
- [ ] Create content collection for news/blog posts
- [ ] Create `src/pages/news/index.astro`
  - [ ] List all news posts
  - [ ] Pagination
  - [ ] Categories/tags
  - [ ] Search functionality
- [ ] Create `src/pages/news/[slug].astro`
  - [ ] Blog post detail page
  - [ ] Metadata (author, date, category)
  - [ ] Related posts
- [ ] Create English versions: `src/pages/en/news/...`

---

## Phase 3b: UI/UX Enhancements (Optional)

### Performance Optimizations
- [ ] Optimize images with Astro Image component
- [ ] Lazy loading for images
- [ ] Preload critical fonts
- [ ] Code splitting for heavy components

### Accessibility Improvements
- [ ] ARIA labels for interactive elements
- [ ] Keyboard navigation testing
- [ ] Screen reader testing
- [ ] Color contrast verification (WCAG AA)
- [ ] Focus indicators

### SEO Enhancements
- [x] Add structured data (JSON-LD)
  - [x] Organization schema
  - [x] EducationalOrganization schema
  - [x] Course schema for programs
- [x] Optimize meta descriptions per page
- [x] Add breadcrumb navigation
- [x] Technical SEO meta tags
  - [x] Canonical URL (`<link rel="canonical">`)
  - [x] Open Graph image (`og:image`, `og:url`)
  - [x] Twitter card image (`twitter:image`)
  - [x] Per-page dynamic hreflang (id, en, x-default)
  - [x] Default OG image (1200x630px)
- [ ] Improve internal linking

---

## Phase 4-5: Authentication & Application Portal (Deferred)

**Status:** Not started - waiting for backend (Phase 2)

These features require backend API:
- Login/register pages
- Dashboard
- Application form with file upload
- Application status tracking
- Admin interface

See root `TODO.md` Phase 4-5 for details.

---

## Phase 7: Testing & Polish (Frontend)

### Manual Testing
- [ ] Test all pages on mobile devices
- [ ] Test on different browsers (Chrome, Firefox, Safari, Edge)
- [ ] Test all navigation links
- [ ] Test language switcher
- [ ] Test responsive breakpoints
- [ ] Verify no console errors

### Automated Testing
- [ ] Playwright tests for critical paths
- [ ] Visual regression tests (optional)

### Performance Testing
- [ ] Lighthouse audit (target 90+ scores)
  - [ ] Homepage
  - [ ] Programs page
  - [ ] About page
  - [ ] Lecturers page
- [ ] Core Web Vitals optimization
- [ ] Page load time <2s

---

## Notes

### Priority
**P0 (Critical for Phase 3 completion):**
- Programs pages
- Contact page
- Admissions page

**P1 (High priority):**
- News/blog system
- SEO enhancements

**P2 (Medium priority):**
- Performance optimizations
- Accessibility improvements

### Estimated Timeline
- Programs pages: 2-3 days
- Contact page: 1 day
- Admissions page: 1-2 days
- News/blog: 2-3 days
- Testing & polish: 2-3 days

**Total: 8-12 days to complete Phase 3**
