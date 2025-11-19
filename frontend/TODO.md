# Frontend - Implementation TODO

## Phase 1: Project Setup & Configuration

### Initialize Project
- [x] Move logo to `public/images/`
- [ ] Initialize Astro project
- [ ] Configure TypeScript (strict mode)
- [ ] Install and configure Tailwind CSS
- [ ] Set up Prettier and ESLint
- [ ] Configure `astro.config.mjs`
- [ ] Configure `tailwind.config.mjs`
- [ ] Create `.dev.vars.example` for environment variables

### Directory Structure
- [ ] Create `src/content/` directory structure
  - [ ] `src/content/programs/`
  - [ ] `src/content/about/`
  - [ ] `src/content/admissions/`
  - [ ] `src/content/config.ts` (content collections schema)
- [ ] Create `src/pages/` directory
- [ ] Create `src/layouts/` directory
- [ ] Create `src/components/` directory
- [ ] Create `src/scripts/` directory
- [ ] Create `src/styles/` directory
- [ ] Create `functions/` directory (Cloudflare Workers)

### Dependencies
- [ ] Install Astro core packages
- [ ] Install Tailwind CSS and plugins
- [ ] Install Astro integrations:
  - [ ] `@astrojs/tailwind`
  - [ ] `@astrojs/mdx` (for enhanced markdown)
  - [ ] `@astrojs/sitemap` (for SEO)
- [ ] Install TypeScript types:
  - [ ] `@types/node`
- [ ] Install utility libraries:
  - [ ] `clsx` (conditional classnames)
  - [ ] `date-fns` (date formatting)

---

## Phase 2: Layouts & Base Components

### Layouts
- [ ] Create `BaseLayout.astro`
  - [ ] HTML structure
  - [ ] Meta tags (SEO)
  - [ ] Google Analytics (optional)
  - [ ] Font loading
- [ ] Create `MarketingLayout.astro`
  - [ ] Header component
  - [ ] Footer component
  - [ ] Navigation
- [ ] Create `DashboardLayout.astro`
  - [ ] Authenticated layout
  - [ ] Sidebar navigation
  - [ ] User menu

### Base Components
- [ ] Create `Header.astro`
  - [ ] Logo
  - [ ] Navigation menu
  - [ ] Mobile menu toggle
  - [ ] Login/Register buttons (if not logged in)
  - [ ] User dropdown (if logged in)
- [ ] Create `Footer.astro`
  - [ ] Campus information
  - [ ] Quick links
  - [ ] Social media links
  - [ ] Copyright
- [ ] Create `Navigation.astro`
  - [ ] Desktop menu
  - [ ] Mobile hamburger menu
  - [ ] Active link highlighting

### Utility Components
- [ ] Create `Button.astro` (reusable button)
- [ ] Create `Card.astro` (reusable card)
- [ ] Create `Container.astro` (content container)
- [ ] Create `Section.astro` (page section wrapper)

---

## Phase 3: Marketing Pages (Static Content)

### Content Files (Markdown)
- [ ] Write `src/content/programs/computer-science.md`
- [ ] Write `src/content/programs/business.md`
- [ ] Write `src/content/programs/engineering.md`
- [ ] Write `src/content/about/history.md`
- [ ] Write `src/content/about/vision-mission.md`
- [ ] Write `src/content/about/facilities.md`
- [ ] Write `src/content/admissions/requirements.md`
- [ ] Write `src/content/admissions/process.md`
- [ ] Write `src/content/admissions/calendar.md`

### Homepage
- [ ] Create `src/pages/index.astro`
  - [ ] Hero section with CTA
  - [ ] Featured programs
  - [ ] Why choose us section
  - [ ] Latest news/updates
  - [ ] Quick stats (students, programs, etc.)
  - [ ] Call-to-action (Apply Now)

### Programs
- [ ] Create `src/pages/programs/index.astro`
  - [ ] List all programs
  - [ ] Filter by type/duration
  - [ ] Program cards with images
- [ ] Create `src/pages/programs/[slug].astro`
  - [ ] Program details
  - [ ] Curriculum overview
  - [ ] Career prospects
  - [ ] Admission requirements
  - [ ] Apply CTA button

### Static Pages
- [ ] Create `src/pages/about.astro`
  - [ ] Campus history
  - [ ] Vision & mission
  - [ ] Leadership team
  - [ ] Facilities
- [ ] Create `src/pages/admissions.astro`
  - [ ] Admission process
  - [ ] Requirements
  - [ ] Important dates
  - [ ] FAQ
- [ ] Create `src/pages/contact.astro`
  - [ ] Contact form
  - [ ] Campus location (map)
  - [ ] Phone/email/address
  - [ ] Social media

### Components for Marketing Pages
- [ ] Create `Hero.astro` (hero section)
- [ ] Create `ProgramCard.astro` (program display)
- [ ] Create `FeatureCard.astro` (feature highlights)
- [ ] Create `StatsSection.astro` (statistics display)
- [ ] Create `Testimonial.astro` (student testimonials)
- [ ] Create `CTA.astro` (call-to-action section)

---

## Phase 4: Authentication Pages

### Login Page
- [ ] Create `src/pages/login.astro`
  - [ ] Email/password form
  - [ ] "Sign in with Google" button
  - [ ] "Forgot password" link
  - [ ] "Don't have an account? Register" link
  - [ ] Form validation (client-side)
  - [ ] Error message display
  - [ ] Loading state

### Register Page
- [ ] Create `src/pages/register.astro`
  - [ ] Registration form (name, email, password)
  - [ ] "Sign up with Google" button
  - [ ] Password strength indicator
  - [ ] Terms & conditions checkbox
  - [ ] "Already have an account? Login" link
  - [ ] Form validation
  - [ ] Success/error messages

### Client-Side Auth Utilities
- [ ] Create `src/scripts/auth.ts`
  - [ ] `login(email, password)` function
  - [ ] `register(name, email, password)` function
  - [ ] `loginWithGoogle()` function
  - [ ] `logout()` function
  - [ ] `getCurrentUser()` function
  - [ ] `isAuthenticated()` function
- [ ] Create `src/scripts/api.ts`
  - [ ] API client wrapper
  - [ ] Error handling
  - [ ] Automatic token inclusion (via cookies)

---

## Phase 5: Application Portal (Authenticated Pages)

### Dashboard
- [ ] Create `src/pages/dashboard.astro`
  - [ ] Welcome message with user name
  - [ ] Application status summary
  - [ ] Quick actions (Apply, View Applications)
  - [ ] Profile completion progress
  - [ ] Notifications/updates

### Application Form
- [ ] Create `src/pages/apply.astro`
  - [ ] Multi-step form wizard
  - [ ] Step 1: Program selection
  - [ ] Step 2: Personal information
  - [ ] Step 3: Educational background
  - [ ] Step 4: Document upload
  - [ ] Step 5: Review and submit
  - [ ] Auto-save draft functionality
  - [ ] Form validation (client-side)
  - [ ] File upload with preview
  - [ ] Progress indicator
  - [ ] Submit confirmation modal

### Application Status
- [ ] Create `src/pages/applications.astro`
  - [ ] List user's applications
  - [ ] Status badges (pending, approved, rejected)
  - [ ] View application details
  - [ ] Download submitted documents
  - [ ] Edit draft applications

### Components for Application Portal
- [ ] Create `ApplicationForm.astro`
  - [ ] Form steps
  - [ ] Field validation
  - [ ] File upload component
- [ ] Create `StatusBadge.astro`
  - [ ] Color-coded status (pending/approved/rejected)
- [ ] Create `FileUpload.astro`
  - [ ] Drag & drop upload
  - [ ] File preview
  - [ ] Delete uploaded file
  - [ ] File size/type validation
- [ ] Create `FormStep.astro`
  - [ ] Step indicator
  - [ ] Previous/Next buttons
  - [ ] Validation messages

### Form Validation
- [ ] Create `src/scripts/form-validation.ts`
  - [ ] Email validation
  - [ ] Password strength check
  - [ ] Required field validation
  - [ ] File type/size validation
  - [ ] Display error messages

---

## Phase 6: Admin Interface

### Application Review
- [ ] Create `src/pages/admin/applications.astro`
  - [ ] List all applications
  - [ ] Filter by status/program/date
  - [ ] Search by name/email
  - [ ] Pagination
  - [ ] Quick actions (Approve/Reject)
- [ ] Create `src/pages/admin/applications/[id].astro`
  - [ ] Full application details
  - [ ] Applicant information
  - [ ] Uploaded documents viewer
  - [ ] Decision form (Approve/Reject with notes)
  - [ ] Activity log

### User Management
- [ ] Create `src/pages/admin/users.astro`
  - [ ] List all users
  - [ ] Filter by role (registrant/staff)
  - [ ] Search users
  - [ ] Edit user roles
  - [ ] Deactivate users

### Admin Components
- [ ] Create `DataTable.astro`
  - [ ] Sortable columns
  - [ ] Pagination
  - [ ] Row selection
- [ ] Create `FilterBar.astro`
  - [ ] Status filter
  - [ ] Date range filter
  - [ ] Program filter
- [ ] Create `ActionButtons.astro`
  - [ ] Approve/Reject buttons
  - [ ] Confirm modal

---

## Phase 7: BFF Layer (Cloudflare Workers)

### Authentication Handlers
- [ ] Create `functions/auth/login.ts`
  - [ ] Validate input
  - [ ] Call backend API
  - [ ] Set HttpOnly cookie
  - [ ] Return user data (no token)
- [ ] Create `functions/auth/register.ts`
  - [ ] Validate input
  - [ ] Call backend API
  - [ ] Auto-login after registration
- [ ] Create `functions/auth/google/login.ts`
  - [ ] Generate Google OAuth URL
  - [ ] Redirect to Google
- [ ] Create `functions/auth/google/callback.ts`
  - [ ] Exchange code for tokens
  - [ ] Verify Google ID token
  - [ ] Call backend to create/update user
  - [ ] Set HttpOnly cookie
  - [ ] Redirect to dashboard
- [ ] Create `functions/auth/logout.ts`
  - [ ] Clear cookie
  - [ ] Return success

### Application Handlers
- [ ] Create `functions/applications/submit.ts`
  - [ ] Extract token from cookie
  - [ ] Forward to backend API
  - [ ] Return response
- [ ] Create `functions/applications/list.ts`
  - [ ] Get user's applications
  - [ ] Forward to backend
- [ ] Create `functions/applications/status.ts`
  - [ ] Get application status
  - [ ] Forward to backend

### User Handlers
- [ ] Create `functions/users/me.ts`
  - [ ] Get current user info
  - [ ] Extract from cookie
  - [ ] Call backend
- [ ] Create `functions/users/update.ts`
  - [ ] Update user profile
  - [ ] Forward to backend

### File Upload
- [ ] Create `functions/files/upload.ts`
  - [ ] Handle multipart/form-data
  - [ ] Validate file type/size
  - [ ] Upload to storage (R2 or VPS)
  - [ ] Return file URL

### BFF Middleware
- [ ] Create `functions/_middleware/auth.ts`
  - [ ] Check if token exists in cookie
  - [ ] Return 401 if not authenticated
- [ ] Create `functions/_middleware/cors.ts`
  - [ ] Set CORS headers
  - [ ] Handle OPTIONS requests
- [ ] Create `functions/_middleware/error-handler.ts`
  - [ ] Catch errors
  - [ ] Return formatted error response
  - [ ] Log errors

### BFF Utilities
- [ ] Create cookie parser
- [ ] Create token extractor
- [ ] Create API client (to backend)

---

## Phase 8: Styling & UI Polish

### Tailwind Configuration
- [ ] Configure custom colors (brand colors)
- [ ] Configure custom fonts
- [ ] Configure breakpoints
- [ ] Configure spacing scale
- [ ] Add custom utilities

### Global Styles
- [ ] Create `src/styles/global.css`
  - [ ] Tailwind directives
  - [ ] Custom CSS variables
  - [ ] Typography styles
  - [ ] Form styles
  - [ ] Animation keyframes

### Responsive Design
- [ ] Test all pages on mobile
- [ ] Test all pages on tablet
- [ ] Test all pages on desktop
- [ ] Fix layout issues
- [ ] Optimize touch targets

### Accessibility
- [ ] Add proper heading hierarchy
- [ ] Add ARIA labels
- [ ] Ensure keyboard navigation works
- [ ] Test with screen reader
- [ ] Add focus indicators
- [ ] Ensure color contrast (WCAG AA)

### Performance
- [ ] Optimize images
  - [ ] Use Astro Image component
  - [ ] Provide multiple sizes
  - [ ] Lazy loading
- [ ] Optimize fonts
  - [ ] Use font-display: swap
  - [ ] Subset fonts if possible
- [ ] Code splitting
  - [ ] Dynamic imports for heavy components
- [ ] Minimize JavaScript
  - [ ] Use Astro islands for interactivity

---

## Phase 9: SEO & Meta Tags

### SEO Setup
- [ ] Configure `astro.config.mjs` for SEO
  - [ ] Add sitemap integration
  - [ ] Add robots.txt
- [ ] Create `public/robots.txt`
- [ ] Create `public/sitemap.xml` (auto-generated)

### Meta Tags (Per Page)
- [ ] Homepage meta tags
  - [ ] Title
  - [ ] Description
  - [ ] Open Graph tags
  - [ ] Twitter Card tags
  - [ ] Canonical URL
- [ ] Programs pages meta tags
- [ ] Static pages meta tags
- [ ] Application portal meta tags (noindex for authenticated pages)

### Structured Data
- [ ] Add Organization schema
- [ ] Add EducationalOrganization schema
- [ ] Add Course schema (for programs)

---

## Phase 10: Testing

### Manual Testing
- [ ] Test all static pages
- [ ] Test login flow (email/password)
- [ ] Test login flow (Google OIDC)
- [ ] Test registration flow
- [ ] Test application submission
- [ ] Test file upload
- [ ] Test admin interface
- [ ] Test responsive design
- [ ] Test cross-browser (Chrome, Firefox, Safari)

### Type Checking
- [ ] Run `npm run typecheck`
- [ ] Fix all TypeScript errors

### Build Testing
- [ ] Run `npm run build`
- [ ] Fix build errors
- [ ] Test preview build (`npm run preview`)

### Lighthouse Audit
- [ ] Run Lighthouse on homepage
- [ ] Run Lighthouse on programs page
- [ ] Run Lighthouse on application portal
- [ ] Fix performance issues
- [ ] Fix accessibility issues
- [ ] Achieve 90+ scores

---

## Phase 11: Deployment Preparation

### Environment Variables
- [ ] Create `.dev.vars.example`
- [ ] Document all required environment variables
- [ ] Set up production environment variables in Cloudflare dashboard

### Wrangler Configuration
- [ ] Create `wrangler.toml`
  - [ ] Configure Workers routes
  - [ ] Set compatibility date
  - [ ] Configure environment variables
  - [ ] Set build output directory

### GitHub Actions
- [ ] Verify workflow file exists (`.github/workflows/deploy-frontend.yml`)
- [ ] Test deployment workflow
- [ ] Set up GitHub secrets
  - [ ] `CLOUDFLARE_API_TOKEN`
  - [ ] `CLOUDFLARE_ACCOUNT_ID`

### Pre-Deployment Checklist
- [ ] All pages load correctly
- [ ] All forms work
- [ ] All links work
- [ ] No console errors
- [ ] No build warnings
- [ ] Environment variables configured
- [ ] Favicon present
- [ ] robots.txt configured
- [ ] Sitemap generated

---

## Phase 12: Launch

### Initial Deployment
- [ ] Deploy to Cloudflare Pages (staging)
- [ ] Test on staging URL
- [ ] Deploy to production domain
- [ ] Verify DNS settings
- [ ] Verify SSL certificate

### Post-Launch
- [ ] Monitor Cloudflare Analytics
- [ ] Check for errors in Workers logs
- [ ] Test all critical paths
- [ ] Set up uptime monitoring
- [ ] Document any issues

---

## Future Enhancements

### Performance
- [ ] Implement service worker for offline support
- [ ] Add caching strategy
- [ ] Optimize bundle size

### Features
- [ ] Email notifications
- [ ] Real-time application status updates
- [ ] Document preview in browser
- [ ] Application analytics dashboard
- [ ] Bulk application operations (admin)
- [ ] Export applications to CSV

### UX Improvements
- [ ] Add loading skeletons
- [ ] Add success/error toast notifications
- [ ] Add confirmation dialogs
- [ ] Add keyboard shortcuts
- [ ] Add dark mode toggle

---

## Notes

### Estimated Timeline
- **Phase 1-2:** 3-5 days (Setup, layouts, base components)
- **Phase 3:** 3-4 days (Marketing pages and content)
- **Phase 4:** 2-3 days (Authentication pages)
- **Phase 5:** 5-7 days (Application portal)
- **Phase 6:** 3-4 days (Admin interface)
- **Phase 7:** 4-5 days (BFF layer)
- **Phase 8-9:** 3-4 days (Styling, SEO)
- **Phase 10-11:** 2-3 days (Testing, deployment prep)
- **Phase 12:** 1 day (Launch)

**Total: 26-38 days (~5-8 weeks)**

### Priority Levels
- **P0 (Critical):** Phases 1-7, 10-12
- **P1 (High):** Phases 8-9
- **P2 (Medium):** Future enhancements

### Success Criteria
- [ ] All marketing pages are live and look professional
- [ ] Users can register and login (both methods work)
- [ ] Users can submit applications with file uploads
- [ ] Admins can review and approve/reject applications
- [ ] Site loads in <2 seconds
- [ ] Mobile-responsive
- [ ] Zero console errors in production
- [ ] Lighthouse score 90+ across all categories
