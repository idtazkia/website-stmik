# Backend TODO - STMIK Tazkia Admission System

Feature-wise implementation plan. Each feature is a vertical slice delivering end-to-end functionality.

---

## Feature 1: Health Check API ✅

Basic server setup with health endpoint.

- [x] Project structure (`cmd/`, `handler/`, `model/`, `config/`)
- [x] `config/config.go` - Environment loading
- [x] `model/db.go` - PostgreSQL connection pool
- [x] `cmd/server/main.go` - HTTP server with graceful shutdown
- [x] `GET /health` endpoint
- [x] Docker Compose for local PostgreSQL
- [x] Integration test with testcontainers

---

## Feature 2: Portal User Registration

Registrants can create account with email/password.

- [ ] Migration: `004_create_users_extended.up.sql` (add provider, provider_id, is_active)
- [ ] `model/user.go` - User struct, Create, FindByEmail, FindByID
- [ ] `auth/password.go` - HashPassword, VerifyPassword (bcrypt)
- [ ] `auth/jwt.go` - GenerateToken, ValidateToken
- [ ] `handler/middleware.go` - RequireAuth middleware
- [ ] `templates/layout.templ` - Base HTML layout
- [ ] `templates/portal/register.templ` - Registration form
- [ ] `handler/portal_auth.go` - GET/POST /portal/register
- [ ] Test: Registration flow

---

## Feature 3: Portal User Login

Registrants can login with email/password.

- [ ] `templates/portal/login.templ` - Login form
- [ ] `handler/portal_auth.go` - GET/POST /portal/login, POST /portal/logout
- [ ] Cookie-based session (HttpOnly JWT)
- [ ] Test: Login/logout flow

---

## Feature 4: Portal Google OAuth

Registrants can login/register with Google.

- [ ] `auth/google.go` - OAuth flow (GetAuthURL, ExchangeCode, GetUserInfo)
- [ ] `handler/portal_auth.go` - GET /portal/auth/google, GET /portal/auth/google/callback
- [ ] Auto-create user if not exists
- [ ] Test: Google OAuth flow

---

## Feature 5: Lead Capture API

Landing page can submit prospect data.

- [ ] Migration: `005_create_intakes.up.sql`
- [ ] Migration: `006_create_prospects.up.sql` (with UTM fields)
- [ ] `model/intake.go` - Intake struct, FindActive, FindByID
- [ ] `model/prospect.go` - Prospect struct, Create, FindByEmail
- [ ] `handler/api.go` - POST /api/prospects
- [ ] Validate required fields (name, email, whatsapp, intake_id)
- [ ] Capture UTM parameters (source, medium, campaign, term, content)
- [ ] Capture landing_page, device_type
- [ ] Return prospect ID
- [ ] Test: Lead capture with various UTM combinations

---

## Feature 6: Admin Google Login

Staff can login with Google (domain-restricted).

- [ ] `templates/layout_admin.templ` - Admin layout
- [ ] `templates/admin/login.templ` - Login page
- [ ] `handler/admin_auth.go` - GET /admin/login, GET /admin/auth/google, callback
- [ ] Domain check (STAFF_EMAIL_DOMAIN)
- [ ] Auto-create staff user if valid domain
- [ ] `handler/middleware.go` - RequireRole("staff") middleware
- [ ] Test: Staff login with valid/invalid domain

---

## Feature 7: Admin Dashboard

Staff can see funnel overview.

- [ ] `model/stats.go` - GetFunnelStats(intakeID)
- [ ] `templates/admin/dashboard.templ` - Stats cards, funnel visualization
- [ ] `handler/admin.go` - GET /admin
- [ ] Show: new, contacted, applicant, approved, enrolled counts
- [ ] Test: Dashboard with sample data

---

## Feature 8: Prospect List

Staff can view and filter prospects.

- [ ] `model/prospect.go` - List with filters (status, intake, assigned_to)
- [ ] `templates/admin/prospects_list.templ` - Table with filters
- [ ] `templates/components/table.templ` - Reusable table component
- [ ] `templates/components/pagination.templ` - Pagination component
- [ ] `handler/admin.go` - GET /admin/prospects
- [ ] HTMX: Filter without full page reload
- [ ] Test: List with various filters

---

## Feature 9: Prospect Detail

Staff can view prospect details and timeline.

- [ ] Migration: `007_create_activity_log.up.sql`
- [ ] `model/activity.go` - Activity struct, Create, ListByEntity
- [ ] `templates/admin/prospect_detail.templ` - Detail view with timeline
- [ ] `templates/components/timeline.templ` - Timeline component
- [ ] `handler/admin.go` - GET /admin/prospects/{id}
- [ ] Show: contact info, status, assigned staff, activity history
- [ ] Test: Detail view with activities

---

## Feature 10: Prospect Assignment

Staff can be assigned to prospects (round-robin or manual).

- [ ] `model/user.go` - ListActiveStaff, UpdateLastAssigned
- [ ] `model/prospect.go` - Assign, GetNextStaffRoundRobin
- [ ] `handler/admin.go` - POST /admin/prospects/{id}/assign
- [ ] HTMX: Assign dropdown, update without reload
- [ ] Log activity on assignment
- [ ] Test: Round-robin distribution, manual assignment

---

## Feature 11: Prospect Status Update

Staff can update prospect status.

- [ ] `model/prospect.go` - UpdateStatus
- [ ] `handler/admin.go` - POST /admin/prospects/{id}/status
- [ ] HTMX: Status dropdown, update without reload
- [ ] Log activity on status change
- [ ] Test: Status transitions

---

## Feature 12: Portal Dashboard

Registrants can see their application status.

- [ ] `templates/portal/dashboard.templ` - Status overview
- [ ] `handler/portal.go` - GET /portal
- [ ] Show: current status, next steps, documents checklist
- [ ] Test: Dashboard for various statuses

---

## Feature 13: Application Form

Registrants can create/update application.

- [ ] Migration: `008_create_programs.up.sql`
- [ ] Migration: `009_create_tracks.up.sql`
- [ ] Migration: `010_create_applications.up.sql`
- [ ] `model/program.go` - Program struct, ListActive
- [ ] `model/track.go` - Track struct, ListActive
- [ ] `model/application.go` - Application struct, Create, Update, FindByUserID
- [ ] `templates/portal/application.templ` - Program/track selection form
- [ ] `handler/portal.go` - GET/POST /portal/application
- [ ] Link application to prospect (by email)
- [ ] Test: Create and update application

---

## Feature 14: Document Upload

Registrants can upload required documents.

- [ ] Migration: `011_create_documents.up.sql`
- [ ] `model/document.go` - Document struct, Create, Delete, ListByApplication
- [ ] `templates/portal/documents.templ` - Upload form with preview
- [ ] `templates/components/file_upload.templ` - Drag-drop upload component
- [ ] `handler/portal.go` - GET /portal/documents, POST /portal/documents, DELETE /portal/documents/{id}
- [ ] Validate: file type (PDF, JPG, PNG), size (max 5MB)
- [ ] Generate unique filename, save to UPLOAD_DIR
- [ ] HTMX: Upload without page reload
- [ ] Test: Upload valid/invalid files

---

## Feature 15: Document Review

Staff can review documents with checklist.

- [ ] Migration: `012_create_document_checklists.up.sql`
- [ ] Migration: `013_create_document_reviews.up.sql`
- [ ] `model/checklist.go` - Checklist struct, ListByDocType
- [ ] `model/document.go` - Review, UpdateStatus
- [ ] `templates/admin/document_review.templ` - Review modal with checklist
- [ ] `templates/components/checklist.templ` - Checklist component
- [ ] `handler/admin.go` - GET/POST /admin/applications/{id}/documents/{docId}/review
- [ ] HTMX: Modal, submit review without reload
- [ ] Auto-set document status based on checklist results
- [ ] Log activity on review
- [ ] Test: Review with pass/fail items

---

## Feature 16: Application Approval

Staff can approve applications (when all docs approved).

- [ ] `model/application.go` - Approve, CheckAllDocsApproved
- [ ] `handler/admin.go` - POST /admin/applications/{id}/approve
- [ ] Validate: all documents must be approved
- [ ] Generate VA number (placeholder for payment integration)
- [ ] Update status: pending_review → approved
- [ ] Log activity
- [ ] Test: Approve with all docs approved, reject when docs pending

---

## Feature 17: Application Cancellation

Registrants or staff can cancel application.

- [ ] Migration: `014_create_cancel_reasons.up.sql`
- [ ] `model/cancel_reason.go` - CancelReason struct, ListActive
- [ ] `model/application.go` - Cancel
- [ ] `model/prospect.go` - Cancel
- [ ] `templates/portal/cancel.templ` - Cancel confirmation with reason
- [ ] `handler/portal.go` - GET/POST /portal/cancel
- [ ] `handler/admin.go` - POST /admin/prospects/{id}/cancel, POST /admin/applications/{id}/cancel
- [ ] Log activity with reason
- [ ] Test: Cancel from portal and admin

---

## Feature 18: WhatsApp Notifications

System sends WhatsApp messages at key events.

- [ ] Migration: `015_create_communication_log.up.sql`
- [ ] `model/communication.go` - CommunicationLog struct, Create
- [ ] `integration/whatsapp.go` - Client, SendTemplate
- [ ] Templates: welcome, followup, document_reminder, revision_request, approved, enrolled
- [ ] `handler/admin.go` - POST /admin/prospects/{id}/whatsapp (manual send)
- [ ] Hook into: prospect creation, document revision, approval
- [ ] Log all communications
- [ ] Test: Send with mock WhatsApp API

---

## Feature 19: Kafka Payment Integration

Payment events update application status to enrolled.

- [ ] `integration/kafka.go` - Consumer for payment.completed topic
- [ ] Match VA number to application
- [ ] Update status: approved → enrolled
- [ ] Log payment details in activity
- [ ] Send enrolled notification (WhatsApp)
- [ ] Test: Consume payment event, verify status update

---

## Feature 20: Referral Tracking

Track referrals from students, alumni, partners.

- [ ] Migration: `016_create_referrers.up.sql`
- [ ] `model/referrer.go` - Referrer struct, Create, FindByCode, ListActive
- [ ] `model/prospect.go` - Link referrer on creation
- [ ] `handler/api.go` - GET /api/referrers/{code} (validate code for landing page)
- [ ] `templates/admin/referrers.templ` - Referrer management
- [ ] `handler/admin.go` - CRUD /admin/referrers
- [ ] Test: Create prospect with referral code

---

## Feature 21: Campaign Tracking

Track ad campaign performance via UTM parameters.

- [ ] Migration: `017_create_campaigns.up.sql`
- [ ] `model/campaign.go` - Campaign struct, Create, FindByUTM, ListActive
- [ ] `templates/admin/campaigns.templ` - Campaign management
- [ ] `handler/admin.go` - CRUD /admin/campaigns
- [ ] Link prospects to campaigns by utm_campaign
- [ ] Test: Campaign CRUD

---

## Feature 22: Reports - Funnel

Staff can view funnel conversion report.

- [ ] `model/stats.go` - GetFunnelByIntake, GetFunnelByDateRange
- [ ] `templates/admin/reports/funnel.templ` - Funnel chart
- [ ] `handler/admin.go` - GET /admin/reports/funnel
- [ ] Compare current vs previous intake
- [ ] HTMX: Filter by intake/date without reload
- [ ] Test: Funnel with sample data

---

## Feature 23: Reports - Source & Campaign

Staff can view conversion by source and campaign.

- [ ] `model/stats.go` - GetConversionBySource, GetConversionByCampaign
- [ ] `templates/admin/reports/sources.templ` - Source breakdown
- [ ] `templates/admin/reports/campaigns.templ` - Campaign performance
- [ ] `handler/admin.go` - GET /admin/reports/sources, GET /admin/reports/campaigns
- [ ] Show: leads, conversions, conversion rate per source/campaign
- [ ] Test: Reports with various sources/campaigns

---

## Feature 24: Reports - Referrer Leaderboard

Staff can view referrer performance.

- [ ] `model/stats.go` - GetReferrerLeaderboard
- [ ] `templates/admin/reports/referrers.templ` - Referrer leaderboard
- [ ] `handler/admin.go` - GET /admin/reports/referrers
- [ ] Show: referrer name, leads, conversions, conversion rate
- [ ] Test: Leaderboard with sample data

---

## Feature 25: CSV Export

Staff can export data to CSV.

- [ ] `handler/admin.go` - GET /admin/reports/export
- [ ] Export options: prospects, applications, funnel report
- [ ] Filter by intake, date range
- [ ] Test: Export with filters

---

## Feature 26: Settings - Intakes

Admin can manage intake periods.

- [ ] `templates/admin/settings/intakes.templ` - Intake CRUD
- [ ] `handler/admin.go` - GET /admin/settings/intakes, POST (create), PUT (update)
- [ ] HTMX: Inline edit without reload
- [ ] Test: Intake CRUD

---

## Feature 27: Settings - Staff Management

Admin can manage staff active status.

- [ ] `templates/admin/settings/staff.templ` - Staff list with toggle
- [ ] `handler/admin.go` - GET /admin/settings/staff, POST /admin/settings/staff/{id}/toggle
- [ ] Toggle active/inactive for round-robin assignment
- [ ] Test: Toggle staff status

---

## Feature 28: Settings - Checklists

Admin can manage document checklists.

- [ ] `templates/admin/settings/checklists.templ` - Checklist management
- [ ] `handler/admin.go` - CRUD /admin/settings/checklists
- [ ] Group by document type (KTP, Ijazah)
- [ ] Test: Checklist CRUD

---

## Feature 29: Seed Data

Initial data for lookup tables.

- [ ] Migration: `018_seed_data.up.sql`
- [ ] Programs: SI, TI
- [ ] Tracks: Regular, KIP-K, LPDP, Internal Scholarship
- [ ] Cancel reasons: no_response, chose_other, financial, changed_mind, etc.
- [ ] Document checklists: KTP items, Ijazah items
- [ ] Test: Verify seed data exists

---

## Feature 30: Static Assets & Styling

Tailwind CSS, HTMX, Alpine.js setup.

- [ ] `static/css/input.css` - Tailwind input
- [ ] `static/css/output.css` - Compiled Tailwind
- [ ] `static/js/htmx.min.js` - HTMX
- [ ] `static/js/alpine.min.js` - Alpine.js (CSP build)
- [ ] Tailwind config with brand colors (primary: #194189, secondary: #EE7B1D)
- [ ] Build script for Tailwind

---

## Success Criteria

- [ ] Prospect can register from landing page with UTM tracking
- [ ] Prospect can login (email/password or Google)
- [ ] Prospect can select program/track and upload documents
- [ ] Staff can login via Google (domain-restricted)
- [ ] Staff can view/filter/assign prospects
- [ ] Staff can review documents with checklist
- [ ] Staff can approve/cancel applications
- [ ] Round-robin assignment works
- [ ] WhatsApp notifications sent at key events
- [ ] Kafka payment events update status to enrolled
- [ ] Referral tracking works end-to-end
- [ ] Campaign/source reports available
- [ ] All HTMX interactions work (no full page reloads)
- [ ] Mobile responsive
