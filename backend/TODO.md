# Backend TODO - STMIK Tazkia Admission CRM

Sales-driven admission system with CRM capabilities. Features ordered by deployment sequence.

---

## Completed

### Feature 1: Health Check API ✅

- [x] Project structure (`cmd/`, `handler/`, `model/`, `config/`)
- [x] `config/config.go` - Environment loading
- [x] `model/db.go` - PostgreSQL connection pool
- [x] `cmd/server/main.go` - HTTP server with graceful shutdown
- [x] `GET /health` endpoint with version info
- [x] Docker Compose for local PostgreSQL
- [x] Integration test with testcontainers

### Feature 2: UI Infrastructure ✅

- [x] Templ templates with Portal/Admin layouts
- [x] Self-hosted HTMX 2.0.8 and Alpine.js CSP 3.15.3
- [x] Tailwind CSS v4 with brand colors
- [x] Git version indicator in UI
- [x] CSRF protection (Go 1.25 stdlib CrossOriginProtection)
- [x] Playwright E2E tests with page object pattern
- [x] Makefile with build-time version injection

---

## Data Model

### Core Entities

```
USER
├─ id, created_at, updated_at
├─ email, name
├─ google_id (for OAuth)
├─ role (admin, supervisor, consultant)
├─ supervisor_id (nullable, for consultant hierarchy)
├─ is_active
└─ last_login_at

PRODI (Program Studi)
├─ id, name, code
├─ degree (S1, D3)
└─ is_active

FEE_TYPE
├─ id, name (registration, tuition, dormitory)
├─ is_recurring (tuition=yes, registration=no)
└─ installment_options (JSON: [1] or [1,2,10])

FEE_STRUCTURE
├─ id, fee_type_id, prodi_id (nullable)
├─ academic_year, amount
└─ is_active

INTERACTION_CATEGORY
├─ id, name
├─ sentiment (positive, neutral, negative)
└─ is_active

OBSTACLE
├─ id, name
├─ suggested_response (template)
└─ is_active

ASSIGNMENT_ALGORITHM
├─ id, name, description
└─ is_active (only one active)

DOCUMENT_TYPE
├─ id, name (ktp, ijazah, transcript, photo)
├─ is_required, can_defer
└─ max_file_size_mb

LOST_REASON
├─ id, name
└─ is_active

CAMPAIGN
├─ id, name
├─ type (promo, event, ads)
├─ source_channel (instagram, google, expo, etc)
├─ start_date, end_date
├─ registration_fee_override (nullable)
└─ is_active

REFERRER
├─ id, name
├─ type (alumni, teacher, student, partner)
├─ institution (helps match claims)
├─ phone, email, code (all optional)
├─ bank_name, bank_account, account_holder
├─ commission_per_enrollment
├─ payout_preference (monthly, per_enrollment)
└─ is_active

CANDIDATE
├─ id, created_at, updated_at
├─ name, email, phone, whatsapp
├─ address, city, province
├─ high_school, graduation_year
├─ prodi_id
├─ source_type, source_detail (referral claim text)
├─ campaign_id, referrer_id, referrer_verified_at
├─ status (registered → prospecting → committed → enrolled / lost)
├─ assigned_consultant_id
├─ registration_fee_paid_at
├─ lost_at, lost_reason_id
└─ enrolled_at, nim

BILLING
├─ id, candidate_id, fee_type_id
├─ academic_year, semester
├─ total_amount, installment_count
├─ discount_amount, discount_reason
└─ status (pending, partial, paid, cancelled)

PAYMENT
├─ id, billing_id
├─ installment_number, amount, due_date
├─ paid_at, payment_method, proof_url
└─ verified_by, verified_at

INTERACTION
├─ id, candidate_id, consultant_id
├─ channel, category_id, obstacle_id
├─ remarks, next_followup_date
├─ supervisor_suggestion, suggestion_read_at
└─ created_at

DOCUMENT
├─ id, candidate_id, document_type_id
├─ file_url, file_name, file_size
├─ status (pending, approved, rejected)
├─ reviewed_by, reviewed_at, rejection_reason
└─ created_at

COMMISSION_LEDGER
├─ id, referrer_id, candidate_id
├─ amount, status (pending, approved, paid)
├─ approved_at, paid_at
└─ payout_batch_id

NOTIFICATION_LOG
├─ id, candidate_id
├─ channel (whatsapp, email)
├─ template, message, status
└─ sent_at, error_message
```

### Migration Plan (by dependency)

```
Level 0 - Lookup tables (no dependencies):
  004_create_users.sql
  005_create_prodis.sql
  006_create_fee_types.sql
  007_create_interaction_categories.sql
  008_create_obstacles.sql
  009_create_assignment_algorithms.sql
  010_create_document_types.sql
  011_create_lost_reasons.sql
  012_create_campaigns.sql
  013_create_referrers.sql

Level 1 - Depends on lookup tables:
  014_create_fee_structures.sql       → fee_types, prodis
  015_create_candidates.sql           → prodis, campaigns, referrers, users, lost_reasons

Level 2 - Depends on candidates:
  016_create_billings.sql             → candidates, fee_types
  017_create_payments.sql             → billings, users
  018_create_interactions.sql         → candidates, users, categories, obstacles
  019_create_documents.sql            → candidates, document_types
  020_create_commission_ledger.sql    → candidates, referrers
  021_create_notification_logs.sql    → candidates

Seed data:
  022_seed_data.sql                   → all lookup tables
```

---

# Phase 1: Admin Foundation

Must complete before opening registration.

---

## Feature 3: Staff Login (Google OAuth)

All staff (admin, supervisor, consultant) login with domain-restricted Google.

**Migrations:** 004

- [ ] `model/user.go` - Create, FindByEmail, FindByGoogleID
- [ ] `auth/google.go` - OAuth flow (GetAuthURL, ExchangeCode, GetUserInfo)
- [ ] `handler/admin_auth.go` - GET /admin/login, /admin/auth/google, callback
- [ ] Domain check (STAFF_EMAIL_DOMAIN env var)
- [ ] Auto-create user with role=consultant if valid domain
- [ ] `handler/middleware.go` - RequireAuth, RequireRole
- [ ] Cookie-based session (HttpOnly JWT)
- [ ] Test: Login with valid/invalid domain

---

## Feature 4: Staff Management

Admin manages staff accounts (admin, supervisor, consultant).

**Migrations:** 004

- [ ] `model/user.go` - List, UpdateRole, ToggleActive, SetSupervisor
- [ ] `templates/admin/settings/users.templ` - User list with role dropdown
- [ ] `handler/admin.go` - GET /admin/settings/users, POST update role/active
- [ ] Roles: admin (full access), supervisor (team + suggestions), consultant (own candidates)
- [ ] Assign supervisor to consultants (hierarchy)
- [ ] Toggle active status (for assignment pool)
- [ ] Test: Role changes, supervisor assignment

---

## Feature 5: Settings - Prodi Management

Admin configures available programs.

**Migrations:** 005

- [ ] `model/prodi.go` - CRUD, ListActive
- [ ] `templates/admin/settings/prodis.templ` - Prodi list with inline edit
- [ ] `handler/admin.go` - CRUD /admin/settings/prodis
- [ ] Fields: name, code, degree (S1/D3), is_active
- [ ] HTMX: Inline edit without reload
- [ ] Test: CRUD operations

---

## Feature 6: Settings - Fee Structure

Admin configures fees per prodi and academic year.

**Migrations:** 006, 014

- [ ] `model/fee_type.go` - List (seeded: registration, tuition, dormitory)
- [ ] `model/fee_structure.go` - CRUD, FindByTypeAndProdi
- [ ] `templates/admin/settings/fees.templ` - Fee matrix (prodi x fee_type)
- [ ] `handler/admin.go` - CRUD /admin/settings/fees
- [ ] Set: registration fee (global), tuition per prodi, dormitory (global)
- [ ] Installment options per fee type
- [ ] Test: CRUD, fee lookup

---

## Feature 7: Settings - Categories & Obstacles

Supervisor manages interaction categories and obstacles.

**Migrations:** 007, 008, 022 (seed)

- [ ] `model/interaction_category.go` - CRUD, ListActive
- [ ] `model/obstacle.go` - CRUD, ListActive
- [ ] `templates/admin/settings/categories.templ` - Category CRUD
- [ ] `templates/admin/settings/obstacles.templ` - Obstacle CRUD with suggested response
- [ ] `handler/admin.go` - CRUD for categories, obstacles
- [ ] Seed default categories: interested, considering, hesitant, cold, unreachable
- [ ] Seed default obstacles: price, location, parents, timing, competitor
- [ ] Test: CRUD operations

---

# Phase 2: Configuration

Setup before opening registration.

---

## Feature 8: Campaign Management

Admin manages campaigns with promo pricing.

**Migrations:** 012

- [ ] `model/campaign.go` - CRUD, FindActive, FindByCode
- [ ] `templates/admin/campaigns.templ` - Campaign list
- [ ] `templates/admin/campaign_form.templ` - Create/edit form
- [ ] `handler/admin.go` - CRUD /admin/campaigns
- [ ] Fields: name, type, source_channel, dates, registration_fee_override
- [ ] Fee override: fixed amount or percentage discount
- [ ] Generate UTM-compatible tracking code
- [ ] Test: CRUD, fee override calculation

---

## Feature 9: Referrer Management

Admin manages referrers for commission tracking.

**Migrations:** 013

- [ ] `model/referrer.go` - CRUD, GenerateCode, FindByCode, SearchByName
- [ ] `templates/admin/referrers.templ` - Referrer list
- [ ] `templates/admin/referrer_form.templ` - Create/edit form
- [ ] `handler/admin.go` - CRUD /admin/referrers
- [ ] Fields: name, type, institution, contact, bank details, commission rate
- [ ] Generate optional referral code (for partners who want trackable links)
- [ ] Test: CRUD, code generation, search

---

## Feature 10: Settings - Assignment Algorithm

Configure consultant assignment algorithm.

**Migrations:** 009, 022 (seed)

- [ ] `model/assignment_algorithm.go` - List, SetActive
- [ ] `templates/admin/settings/assignment.templ` - Algorithm selection
- [ ] `handler/admin.go` - GET/POST /admin/settings/assignment
- [ ] Algorithms: round_robin, load_balanced, performance_weighted, activity_based
- [ ] Only one algorithm active at a time
- [ ] Test: Switch algorithms

---

## Feature 11: Settings - Document Types

Configure required documents.

**Migrations:** 010, 022 (seed)

- [ ] `model/document_type.go` - CRUD, ListActive
- [ ] `templates/admin/settings/documents.templ` - Document type list
- [ ] `handler/admin.go` - CRUD /admin/settings/document-types
- [ ] Fields: name, is_required, can_defer, max_file_size_mb
- [ ] Seed: KTP (required), Photo (required), Ijazah (required, can_defer), Transcript (required, can_defer)
- [ ] Test: CRUD operations

---

# Phase 3: Public Registration

Open the app for candidates.

---

## Feature 12: Candidate Registration

Public registration form with source tracking.

**Migrations:** 015

- [ ] `model/candidate.go` - Create, FindByEmail, FindByPhone
- [ ] `templates/public/register.templ` - Registration form
- [ ] `handler/public.go` - GET/POST /register
- [ ] Form fields:
  - Personal: name, email, phone, whatsapp, address, city, province
  - Education: high_school, graduation_year, prodi
  - Source: source_type (dropdown), source_detail (free text if referral)
- [ ] Source types: instagram, google, tiktok, youtube, expo, school_visit, friend_family, teacher_alumni, walkin, other
- [ ] If source is friend_family/teacher_alumni, show "Siapa yang mereferensikan?" field
- [ ] If URL has `?ref=CODE`, auto-link to referrer
- [ ] If URL has `?utm_campaign=X`, auto-link to campaign
- [ ] Auto-assign to consultant (using active algorithm)
- [ ] Show registration fee (from fee_structure or campaign override)
- [ ] Test: Registration with various sources

---

## Feature 13: Registration Fee Payment

Candidate pays registration fee.

**Migrations:** 016, 017

- [ ] `model/billing.go` - Create, FindByCandidate, UpdateStatus
- [ ] `model/payment.go` - Create, UploadProof
- [ ] `templates/public/payment.templ` - Payment instructions, proof upload
- [ ] `handler/public.go` - GET /payment/{token}, POST /payment/{token}/proof
- [ ] Generate billing on registration
- [ ] Amount from fee_structure, apply campaign discount if any
- [ ] If 100% discount, auto-mark as paid
- [ ] Payment proof upload (image)
- [ ] Email/token-based access (no login required)
- [ ] Test: Payment flow, promo discount, proof upload

---

## Feature 14: Document Upload

Candidate uploads required documents.

**Migrations:** 019

- [ ] `model/document.go` - Upload, ListByCandidate
- [ ] `templates/public/documents.templ` - Upload form with checklist
- [ ] `handler/public.go` - GET/POST /documents/{token}
- [ ] Show required vs optional, uploaded vs pending
- [ ] Mark deferrable documents (ijazah, transcript)
- [ ] File validation: type (PDF, JPG, PNG), size limit
- [ ] Email/token-based access (no login required)
- [ ] Test: Upload, validation, defer

---

# Phase 4: CRM Operations

Day-to-day sales operations.

---

## Feature 15: Candidate List & Filters

Admin/consultant views candidates.

- [ ] `model/candidate.go` - List with filters, pagination
- [ ] `templates/admin/candidates_list.templ` - Table with filters
- [ ] `handler/admin.go` - GET /admin/candidates
- [ ] Filters: status, assigned consultant, prodi, campaign, source_type, date range
- [ ] Sort: newest, oldest, next followup due, last interaction
- [ ] Highlight overdue followups (red if > 3 days)
- [ ] Consultant sees only their candidates, supervisor sees team, admin sees all
- [ ] HTMX: Filter without reload
- [ ] Test: Various filter combinations, role-based visibility

---

## Feature 16: Candidate Detail & Timeline

View candidate info and history.

- [ ] `templates/admin/candidate_detail.templ` - Info + timeline
- [ ] `handler/admin.go` - GET /admin/candidates/{id}
- [ ] Show: personal info, prodi, source, campaign/referrer, status, assigned consultant
- [ ] Timeline: interactions, payments, documents, status changes
- [ ] Quick actions: log interaction, reassign, change status
- [ ] Test: Detail view, timeline ordering

---

## Feature 17: Interaction Logging

Consultants log each contact.

**Migrations:** 018

- [ ] `model/interaction.go` - Create, ListByCandidate, ListByConsultant
- [ ] `templates/admin/interaction_form.templ` - Log interaction modal
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/interactions
- [ ] Fields: channel, category, obstacle (optional), remarks, next_followup_date
- [ ] Channels: call, whatsapp, email, campus_visit, home_visit
- [ ] Auto-update candidate last_interaction_at
- [ ] Test: Create interaction, list

---

## Feature 18: Supervisor Suggestions

Supervisor reviews and provides guidance.

- [ ] `templates/admin/candidate_detail.templ` - Suggestion field in timeline
- [ ] `handler/admin.go` - POST /admin/interactions/{id}/suggestion
- [ ] Consultant sees suggestion, marks as read
- [ ] Notification badge for unread suggestions
- [ ] Test: Add suggestion, mark as read

---

## Feature 19: Consultant Assignment

Manual reassignment of candidates.

- [ ] `model/candidate.go` - Assign, GetAssignmentStats
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/assign
- [ ] Supervisor/admin can reassign candidates
- [ ] Show consultant workload in dropdown
- [ ] Log assignment change in timeline
- [ ] Test: Reassignment, workload display

---

## Feature 20: Referral Claim Verification

Link referral claims to referrers.

- [ ] `templates/admin/referral_claims.templ` - List unverified claims
- [ ] `handler/admin.go` - GET /admin/referral-claims, POST link
- [ ] Show candidates with source_detail (referral claim) but no referrer_id
- [ ] Search existing referrers by name/institution
- [ ] Actions: link to existing referrer, create new referrer then link, mark as invalid
- [ ] Test: Claim verification flow

---

# Phase 5: Commitment & Enrollment

Convert candidates to students.

---

## Feature 21: Commitment & Tuition Billing

Generate billing when candidate commits.

- [ ] `model/candidate.go` - Commit (change status)
- [ ] `model/billing.go` - CreateTuitionBilling, CreateDormitoryBilling
- [ ] `templates/admin/commitment_form.templ` - Commitment modal
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/commit
- [ ] Select: tuition installments (1x), dormitory (1x, 2x, or 10x)
- [ ] Generate billing records with due dates
- [ ] Change status: prospecting → committed
- [ ] Test: Commit with various installment options

---

## Feature 22: Payment Tracking

Track and verify installment payments.

- [ ] `model/payment.go` - ListByBilling, RecordPayment, VerifyPayment
- [ ] `templates/admin/payments.templ` - Payment list
- [ ] `templates/admin/payment_verify.templ` - Verification modal
- [ ] `handler/admin.go` - GET /admin/candidates/{id}/payments, POST verify
- [ ] Show: installment number, due date, amount, status (pending/paid/overdue)
- [ ] Candidate uploads proof via public link
- [ ] Admin verifies payment, marks as paid
- [ ] Overdue highlighting
- [ ] Test: Payment lifecycle, overdue detection

---

## Feature 23: Document Review

Admin reviews uploaded documents.

- [ ] `model/document.go` - UpdateStatus
- [ ] `templates/admin/document_review.templ` - Review modal
- [ ] `handler/admin.go` - GET /admin/candidates/{id}/documents, POST approve/reject
- [ ] View document, approve or reject with reason
- [ ] Notify candidate of rejection (for re-upload)
- [ ] Test: Review flow, rejection

---

## Feature 24: Enrollment

Mark candidate as enrolled.

- [ ] `model/candidate.go` - Enroll, GenerateNIM
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/enroll
- [ ] Validation:
  - Registration fee: paid
  - Tuition: at least 1st installment paid
  - Documents: KTP + photo approved (ijazah/transcript can be pending)
- [ ] Generate NIM: YYYY + PRODI_CODE + SEQUENCE (e.g., 2026SI001)
- [ ] Change status: committed → enrolled
- [ ] Trigger commission creation if referred
- [ ] Test: Enrollment validation, NIM generation

---

## Feature 25: Lost Candidate

Mark candidate as lost.

**Migrations:** 011, 022 (seed)

- [ ] `model/lost_reason.go` - ListActive
- [ ] `model/candidate.go` - MarkLost
- [ ] `templates/admin/lost_form.templ` - Lost modal with reason
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/lost
- [ ] Select reason: no_response, chose_competitor, financial, not_qualified, timing, other
- [ ] Change status: any → lost
- [ ] Record lost_at timestamp
- [ ] Test: Mark as lost with reason

---

# Phase 6: Commissions

Track referrer commissions.

---

## Feature 26: Commission Tracking

Auto-create and track commissions.

**Migrations:** 020

- [ ] `model/commission.go` - Create, ListByReferrer, ListPending
- [ ] Auto-create commission when referred candidate enrolls
- [ ] Amount from referrer.commission_per_enrollment
- [ ] Status: pending → approved → paid
- [ ] Test: Auto-creation on enrollment

---

## Feature 27: Commission Payout

Approve and pay commissions.

- [ ] `templates/admin/commissions.templ` - Commission list
- [ ] `handler/admin.go` - GET /admin/commissions, POST approve, POST mark-paid
- [ ] Filter by: referrer, status, date range
- [ ] Batch approve, batch mark as paid
- [ ] Export for bank transfer
- [ ] Test: Approval flow, batch operations

---

# Phase 7: Reporting & Analytics

Insights for decision making.

---

## Feature 28: Dashboard - Consultant

Consultant's daily view.

- [ ] `model/stats.go` - GetConsultantStats
- [ ] `templates/admin/dashboard_consultant.templ`
- [ ] `handler/admin.go` - GET /admin (role-based)
- [ ] Show: my candidates by status, overdue followups, today's tasks
- [ ] Quick access to pending followups
- [ ] Test: Dashboard data accuracy

---

## Feature 29: Dashboard - Supervisor

Supervisor's team view.

- [ ] `model/stats.go` - GetTeamStats, GetFunnelStats
- [ ] `templates/admin/dashboard_supervisor.templ`
- [ ] Show: team funnel, consultant leaderboard, stuck candidates (> 7 days no interaction)
- [ ] Common obstacles this period
- [ ] Test: Dashboard data accuracy

---

## Feature 30: Reports - Funnel

Conversion funnel analysis.

- [ ] `model/stats.go` - GetFunnelByDateRange, GetFunnelByProdi
- [ ] `templates/admin/reports/funnel.templ`
- [ ] `handler/admin.go` - GET /admin/reports/funnel
- [ ] Filter by: date range, prodi, campaign
- [ ] Show: registered → prospecting → committed → enrolled, with conversion rates
- [ ] Test: Report accuracy

---

## Feature 31: Reports - Consultant Performance

Individual performance metrics.

- [ ] `model/stats.go` - GetConsultantPerformance
- [ ] `templates/admin/reports/consultants.templ`
- [ ] `handler/admin.go` - GET /admin/reports/consultants
- [ ] Metrics: candidates handled, success rate, avg days to commit, interaction frequency
- [ ] Ranking by success rate
- [ ] Test: Metrics calculation

---

## Feature 32: Reports - Campaign ROI

Campaign effectiveness.

- [ ] `model/stats.go` - GetCampaignStats
- [ ] `templates/admin/reports/campaigns.templ`
- [ ] `handler/admin.go` - GET /admin/reports/campaigns
- [ ] Show: leads, commits, enrollments, conversion rate per campaign
- [ ] Cost per enrollment (if cost data available)
- [ ] Test: Report accuracy

---

## Feature 33: Reports - Referrer Leaderboard

Referrer performance.

- [ ] `model/stats.go` - GetReferrerStats
- [ ] `templates/admin/reports/referrers.templ`
- [ ] `handler/admin.go` - GET /admin/reports/referrers
- [ ] Show: referrals, enrollments, conversion rate, commission earned/paid
- [ ] Test: Report accuracy

---

## Feature 34: CSV Export

Export data for external analysis.

- [ ] `handler/admin.go` - GET /admin/export/candidates, /admin/export/interactions
- [ ] Filter by: date range, status, consultant, campaign
- [ ] Include all relevant fields
- [ ] Test: Export with filters

---

# Phase 8: Notifications

Communication automation.

---

## Feature 35: WhatsApp Notifications

Send notifications at key events.

**Migrations:** 021

- [ ] `integration/whatsapp.go` - Send via API (Fonnte/similar)
- [ ] `model/notification.go` - Create, ListByCandidate
- [ ] Templates: registration_confirmed, payment_reminder, document_reminder, enrolled
- [ ] Manual send from candidate detail
- [ ] Log all sent messages
- [ ] Test: Send with mock API

---

## Success Criteria

- [ ] Admin can login, configure prodis/fees/campaigns before opening registration
- [ ] Candidate can register with source tracking, pay fee, upload documents
- [ ] Registration fee waived during promo campaigns
- [ ] Candidates auto-assigned to consultants
- [ ] Consultants log interactions with category/obstacle/remarks
- [ ] Supervisors provide suggestions on interactions
- [ ] Commitment generates tuition billing with installments
- [ ] Payment tracking with proof upload and verification
- [ ] Documents can be deferred (ijazah/transcript)
- [ ] Enrollment validates requirements, generates NIM
- [ ] Referrer commissions auto-created and tracked
- [ ] Campaign ROI trackable via reports
- [ ] All admin interactions use HTMX (no full page reloads)
- [ ] Mobile responsive
