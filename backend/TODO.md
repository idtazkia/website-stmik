# Backend TODO - STMIK Tazkia Admission CRM

Sales-driven admission system with CRM capabilities. Each feature is a vertical slice.

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
CANDIDATE
├─ id, created_at, updated_at
├─ name, email, phone, whatsapp
├─ address, city, province
├─ high_school, graduation_year
├─ prodi_id (program choice)
├─ source_type (instagram, google, expo, referral, walkin)
├─ source_detail (free text - "Pak Asep guru BK SMAN 1")
├─ campaign_id (if from tracked campaign)
├─ referrer_id (nullable - linked by admin after verification)
├─ referrer_verified_at
├─ status (registered → prospecting → committed → enrolled / lost)
├─ assigned_consultant_id
├─ registration_fee_paid_at
├─ lost_at, lost_reason_id
└─ enrolled_at, nim

ACADEMIC_CONSULTANT (extends USER)
├─ supervisor_id (nullable, for hierarchy)
├─ is_active (for assignment pool)
├─ Calculated: active_candidate_count, success_rate, activity_score

INTERACTION
├─ id, created_at
├─ candidate_id, consultant_id
├─ channel (call, whatsapp, email, campus_visit, home_visit)
├─ category_id (interested, hesitant, cold, etc)
├─ obstacle_id (nullable)
├─ remarks (free text)
├─ next_followup_date
├─ supervisor_suggestion (filled by supervisor)
└─ suggestion_read_at

INTERACTION_CATEGORY (managed by supervisor)
├─ id, name
├─ sentiment (positive, neutral, negative)
└─ is_active

OBSTACLE (managed by supervisor)
├─ id, name
├─ suggested_response (template)
└─ is_active

CAMPAIGN
├─ id, name
├─ type (promo, event, ads)
├─ source_channel (instagram, google, expo, referral, walkin)
├─ start_date, end_date
├─ registration_fee_override (nullable, for promo)
└─ is_active

REFERRER
├─ id, created_at
├─ name
├─ type (alumni, teacher, student, partner)
├─ institution (SMAN 1 Bogor, etc - helps admin match claims)
├─ phone, email (optional)
├─ code (optional - for partners who want trackable links)
├─ bank_name, bank_account, account_holder
├─ commission_per_enrollment (amount)
├─ payout_preference (monthly, per_enrollment)
├─ notes (admin notes)
└─ is_active

SOURCE_TYPE (enum for registration form)
├─ instagram, google, tiktok, youtube
├─ expo, school_visit
├─ friend_family, teacher_alumni (triggers referral claim field)
├─ walkin, other

COMMISSION_LEDGER
├─ id, created_at
├─ referrer_id, candidate_id
├─ amount
├─ status (pending, approved, paid)
├─ approved_at, paid_at
└─ payout_batch_id

PRODI (Program Studi)
├─ id, name, code
├─ degree (S1, D3)
├─ tuition_per_semester
└─ is_active

FEE_TYPE
├─ id, name (registration, tuition, dormitory)
├─ is_recurring (tuition=yes, registration=no)
└─ installment_options (JSON: [1] or [1,2,10])

FEE_STRUCTURE
├─ id, fee_type_id
├─ prodi_id (nullable - some fees are global)
├─ academic_year
├─ amount
└─ is_active

BILLING
├─ id, created_at
├─ candidate_id, fee_type_id
├─ academic_year, semester (for tuition)
├─ total_amount, installment_count
├─ discount_amount, discount_reason
└─ status (pending, partial, paid, cancelled)

PAYMENT
├─ id, created_at
├─ billing_id
├─ installment_number (1, 2, 3...)
├─ amount, due_date
├─ paid_at, payment_method
├─ proof_url
└─ verified_by, verified_at

DOCUMENT_TYPE
├─ id, name (ktp, ijazah, transcript, photo, etc)
├─ is_required
├─ can_defer (ijazah/transcript can be uploaded later)
└─ max_file_size_mb

DOCUMENT
├─ id, created_at
├─ candidate_id, document_type_id
├─ file_url, file_name, file_size
├─ status (pending, approved, rejected)
├─ reviewed_by, reviewed_at
└─ rejection_reason

ASSIGNMENT_ALGORITHM
├─ id, name
├─ description
└─ is_active (only one active at a time)

Algorithms:
- ROUND_ROBIN: Simple rotation
- LOAD_BALANCED: Lowest active candidates first
- PERFORMANCE_WEIGHTED: Higher success rate prioritized
- ACTIVITY_BASED: Higher follow-up activity prioritized
```

### Migration Plan (by dependency)

```
Level 0 - Lookup tables (no dependencies):
  004_create_prodis.sql
  005_create_fee_types.sql
  006_create_campaigns.sql
  007_create_referrers.sql
  008_create_users.sql
  009_create_assignment_algorithms.sql
  010_create_interaction_categories.sql
  011_create_obstacles.sql
  012_create_document_types.sql
  013_create_lost_reasons.sql

Level 1 - Depends on lookup tables:
  014_create_fee_structures.sql       → fee_types, prodis
  015_create_candidates.sql           → prodis, campaigns, referrers, users, lost_reasons

Level 2 - Depends on candidates:
  016_create_billings.sql             → candidates, fee_types
  017_create_interactions.sql         → candidates, users, categories, obstacles
  018_create_documents.sql            → candidates, document_types
  019_create_commission_ledger.sql    → candidates, referrers

Level 3 - Depends on level 2:
  020_create_payments.sql             → billings, users
  021_create_notification_logs.sql    → candidates

Seed data:
  022_seed_data.sql                   → prodis, fee_types, algorithms, categories,
                                         obstacles, document_types, lost_reasons
```

---

## Feature 3: Candidate Registration

Public registration form with source tracking and referral claims.

**Migrations:** 004, 006, 007, 015

- [ ] `model/prodi.go` - Prodi CRUD, ListActive
- [ ] `model/campaign.go` - Campaign CRUD, FindActive
- [ ] `model/referrer.go` - FindByCode (for ?ref=CODE links)
- [ ] `model/candidate.go` - Create, FindByEmail, FindByPhone
- [ ] `templates/public/register.templ` - Registration form
- [ ] `handler/public.go` - GET/POST /register
- [ ] Form fields:
  - Personal: name, email, phone, whatsapp, address
  - Education: high_school, graduation_year, prodi
  - Source: source_type (dropdown), source_detail (free text if referral)
- [ ] If URL has `?ref=CODE`, auto-link to referrer
- [ ] If source_type is teacher_alumni/friend_family, show "Nama yang mereferensikan" field
- [ ] Show registration fee (default or campaign override)
- [ ] Test: Registration with various source types

---

## Feature 4: Registration Fee Payment

Candidate pays registration fee (can be waived during promo).

**Migrations:** 005, 014, 016, 020

- [ ] `model/fee_structure.go` - FindByType, GetRegistrationFee
- [ ] `model/billing.go` - Create, FindByCandidate, UpdateStatus
- [ ] `model/payment.go` - Create, MarkPaid, UploadProof
- [ ] `templates/public/payment.templ` - Payment instructions, proof upload
- [ ] `handler/public.go` - GET /payment/{token}, POST /payment/{token}/proof
- [ ] Generate billing on registration (amount from fee_structure or campaign override)
- [ ] If campaign has 100% discount, auto-mark as paid
- [ ] Test: Payment flow, promo discount

---

## Feature 5: Admin Login (Google OAuth)

Staff login with domain-restricted Google.

**Migrations:** 008

- [ ] `model/user.go` - Create, FindByEmail, FindByGoogleID
- [ ] `auth/google.go` - OAuth flow
- [ ] `handler/admin_auth.go` - GET /admin/login, /admin/auth/google, callback
- [ ] Domain check (STAFF_EMAIL_DOMAIN env var)
- [ ] Auto-create user with role=consultant if valid domain
- [ ] `handler/middleware.go` - RequireAuth, RequireRole
- [ ] Test: Login with valid/invalid domain

---

## Feature 6: Consultant Assignment

Auto-assign new candidates to consultants.

**Migrations:** 009, 022 (seed)

- [ ] `model/assignment.go` - GetNextConsultant(algorithm), Assign
- [ ] `model/user.go` - ListActiveConsultants, GetConsultantStats
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/assign (manual override)
- [ ] Auto-assign on registration (using active algorithm)
- [ ] Supervisor can reassign
- [ ] Test: Each algorithm, manual override

---

## Feature 7: Interaction Logging

Consultants log each contact with candidate.

**Migrations:** 010, 011, 017, 022 (seed)

- [ ] `model/interaction.go` - Create, ListByCandidate, ListByConsultant
- [ ] `model/interaction_category.go` - CRUD (supervisor only)
- [ ] `model/obstacle.go` - CRUD (supervisor only)
- [ ] `templates/admin/interaction_form.templ` - Log interaction modal
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/interactions
- [ ] Fields: channel, category, obstacle, remarks, next_followup_date
- [ ] Test: Create interaction, list by candidate

---

## Feature 8: Supervisor Suggestions

Supervisor reviews interactions and provides guidance.

- [ ] `templates/admin/candidate_detail.templ` - Interaction timeline with suggestion field
- [ ] `handler/admin.go` - POST /admin/interactions/{id}/suggestion
- [ ] Consultant sees suggestion, marks as read
- [ ] Test: Add suggestion, mark as read

---

## Feature 9: Candidate List & Filters

Admin views candidates with filters.

- [ ] `model/candidate.go` - List with filters (status, consultant, prodi, date range)
- [ ] `templates/admin/candidates_list.templ` - Table with filters
- [ ] `handler/admin.go` - GET /admin/candidates
- [ ] Filters: status, assigned consultant, prodi, campaign, date range
- [ ] Sort: newest, oldest, next followup due
- [ ] Highlight overdue followups
- [ ] HTMX: Filter without reload
- [ ] Test: Various filter combinations

---

## Feature 10: Candidate Detail & Timeline

View candidate info and full interaction history.

- [ ] `templates/admin/candidate_detail.templ` - Info + timeline
- [ ] `handler/admin.go` - GET /admin/candidates/{id}
- [ ] Show: personal info, prodi, campaign/referrer source, status
- [ ] Timeline: all interactions, payments, status changes
- [ ] Quick actions: log interaction, change status, reassign
- [ ] Test: Detail view with interactions

---

## Feature 11: Commitment & Tuition Billing

When candidate commits, generate tuition billing.

- [ ] `model/candidate.go` - Commit (change status, create billing)
- [ ] `model/billing.go` - CreateTuitionBilling (with installment plan)
- [ ] `templates/admin/commitment_form.templ` - Select installment count
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/commit
- [ ] Generate billing: tuition (per prodi amount), dormitory (optional)
- [ ] Support installments: 1x, 2x, or 10x for dormitory
- [ ] Test: Commit with various installment options

---

## Feature 12: Payment Tracking

Track installment payments.

- [ ] `model/payment.go` - ListByBilling, RecordPayment, VerifyPayment
- [ ] `templates/admin/payments.templ` - Payment list with status
- [ ] `templates/admin/payment_detail.templ` - Record/verify payment
- [ ] `handler/admin.go` - GET /admin/billings/{id}/payments, POST record, POST verify
- [ ] Show: due date, amount, status (pending/paid/overdue)
- [ ] Upload payment proof
- [ ] Admin verifies payment
- [ ] Test: Record payment, verify, overdue detection

---

## Feature 13: Document Upload (Deferred)

Candidate uploads documents (some can be deferred).

**Migrations:** 012, 018, 022 (seed)

- [ ] `model/document_type.go` - ListActive, FindByID
- [ ] `model/document.go` - Upload, ListByCandidate, UpdateStatus
- [ ] `templates/public/documents.templ` - Upload form
- [ ] `handler/public.go` - GET/POST /documents/{token}
- [ ] Mark required vs optional, deferrable vs not
- [ ] Test: Upload, defer ijazah

---

## Feature 14: Document Review

Admin reviews uploaded documents.

- [ ] `templates/admin/documents_review.templ` - Review interface
- [ ] `handler/admin.go` - GET /admin/candidates/{id}/documents, POST approve/reject
- [ ] Approve/reject with reason
- [ ] Test: Review flow

---

## Feature 15: Enrollment

Mark candidate as enrolled when requirements met.

- [ ] `model/candidate.go` - Enroll (validate payments, generate NIM)
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/enroll
- [ ] Requirements: registration fee paid, tuition (at least 1st installment) paid
- [ ] Documents: KTP and photo required, ijazah/transcript can be pending
- [ ] Generate NIM (format: YYYY-PRODI-SEQUENCE)
- [ ] Test: Enroll with various states

---

## Feature 16: Lost Candidate

Mark candidate as lost with reason.

**Migrations:** 013, 022 (seed)

- [ ] `model/lost_reason.go` - ListActive, FindByID
- [ ] `model/candidate.go` - MarkLost
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/lost
- [ ] Test: Mark as lost

---

## Feature 17: Referrer Management

Manage referrers and verify referral claims.

**Migrations:** 007, 019

- [ ] `model/referrer.go` - CRUD, GenerateCode, FindByCode, SearchByName
- [ ] `templates/admin/referrers.templ` - Referrer list
- [ ] `templates/admin/referrer_form.templ` - Create/edit referrer
- [ ] `handler/admin.go` - CRUD /admin/referrers
- [ ] Generate unique referral code (optional - for partners who want trackable links)
- [ ] Test: CRUD, code generation

### Referral Claim Verification

- [ ] `templates/admin/referral_claims.templ` - List unverified claims
- [ ] `handler/admin.go` - GET /admin/referral-claims
- [ ] Show candidates with source_detail but no referrer_id
- [ ] Search existing referrers by name/institution
- [ ] Link to existing referrer OR create new then link
- [ ] Test: Claim verification flow

---

## Feature 18: Commission Tracking

Track and pay referrer commissions.

- [ ] `model/commission.go` - CreateOnEnrollment, Approve, MarkPaid, ListByReferrer
- [ ] `templates/admin/commissions.templ` - Commission list
- [ ] `handler/admin.go` - GET /admin/commissions, POST approve, POST pay
- [ ] Auto-create commission entry when referred candidate enrolls
- [ ] Payout modes: monthly batch or per-enrollment
- [ ] Test: Commission lifecycle

---

## Feature 19: Campaign Management

Manage campaigns with promo pricing.

- [ ] `templates/admin/campaigns.templ` - Campaign list
- [ ] `templates/admin/campaign_form.templ` - Create/edit
- [ ] `handler/admin.go` - CRUD /admin/campaigns
- [ ] Set registration fee override (discount or waive)
- [ ] Track source channel (instagram, google, expo, etc)
- [ ] Test: CRUD, fee override

---

## Feature 20: Settings - Categories & Obstacles

Supervisor manages interaction categories and obstacles.

- [ ] `templates/admin/settings/categories.templ` - Category CRUD
- [ ] `templates/admin/settings/obstacles.templ` - Obstacle CRUD with suggested response
- [ ] `handler/admin.go` - CRUD for categories, obstacles
- [ ] Test: CRUD

---

## Feature 21: Settings - Assignment Algorithm

Configure which assignment algorithm is active.

- [ ] `templates/admin/settings/assignment.templ` - Algorithm selection
- [ ] `handler/admin.go` - GET/POST /admin/settings/assignment
- [ ] Only one algorithm active at a time
- [ ] Test: Switch algorithms

---

## Feature 22: Settings - Fee Structure

Manage fee amounts per prodi and academic year.

- [ ] `templates/admin/settings/fees.templ` - Fee structure management
- [ ] `handler/admin.go` - CRUD /admin/settings/fees
- [ ] Set registration fee default, tuition per prodi, dormitory fee
- [ ] Test: CRUD

---

## Feature 23: Settings - Staff Management

Manage consultant active status for assignment pool.

- [ ] `templates/admin/settings/staff.templ` - Staff list
- [ ] `handler/admin.go` - GET /admin/settings/staff, POST toggle active
- [ ] Assign supervisor role
- [ ] Test: Toggle, role assignment

---

## Feature 24: Dashboard - Consultant

Consultant sees their candidates and pending followups.

- [ ] `model/stats.go` - GetConsultantStats
- [ ] `templates/admin/dashboard_consultant.templ`
- [ ] `handler/admin.go` - GET /admin (role-based dashboard)
- [ ] Show: my candidates by status, overdue followups, recent interactions
- [ ] Test: Dashboard data

---

## Feature 25: Dashboard - Supervisor

Supervisor sees team performance and funnel.

- [ ] `model/stats.go` - GetTeamStats, GetFunnelStats
- [ ] `templates/admin/dashboard_supervisor.templ`
- [ ] Show: funnel (registered → committed → enrolled), consultant leaderboard
- [ ] Candidates stuck > 7 days without interaction
- [ ] Common obstacles this period
- [ ] Test: Dashboard data

---

## Feature 26: Reports - Funnel

Conversion funnel report.

- [ ] `model/stats.go` - GetFunnelByDateRange, GetFunnelByProdi
- [ ] `templates/admin/reports/funnel.templ`
- [ ] `handler/admin.go` - GET /admin/reports/funnel
- [ ] Filter by date range, prodi, campaign
- [ ] Test: Report with filters

---

## Feature 27: Reports - Consultant Performance

Individual consultant metrics.

- [ ] `model/stats.go` - GetConsultantPerformance
- [ ] `templates/admin/reports/consultants.templ`
- [ ] `handler/admin.go` - GET /admin/reports/consultants
- [ ] Show: candidates handled, success rate, avg days to commit, activity score
- [ ] Test: Performance report

---

## Feature 28: Reports - Campaign ROI

Campaign performance and conversion.

- [ ] `model/stats.go` - GetCampaignStats
- [ ] `templates/admin/reports/campaigns.templ`
- [ ] `handler/admin.go` - GET /admin/reports/campaigns
- [ ] Show: leads, commits, enrollments, conversion rate per campaign
- [ ] Test: Campaign report

---

## Feature 29: Reports - Referrer Leaderboard

Referrer performance and commissions.

- [ ] `model/stats.go` - GetReferrerStats
- [ ] `templates/admin/reports/referrers.templ`
- [ ] `handler/admin.go` - GET /admin/reports/referrers
- [ ] Show: referrals, enrollments, commission earned/paid
- [ ] Test: Referrer report

---

## Feature 30: CSV Export

Export data for external analysis.

- [ ] `handler/admin.go` - GET /admin/export/candidates, /admin/export/interactions
- [ ] Filter by date range, status, consultant
- [ ] Test: Export with filters

---

## Feature 31: WhatsApp Notifications

Send notifications at key events.

**Migrations:** 021

- [ ] `integration/whatsapp.go` - Send via API (Fonnte/similar)
- [ ] `model/notification.go` - Log sent messages
- [ ] Templates: registration_confirmed, followup_reminder, payment_reminder, enrolled
- [ ] Manual send from candidate detail
- [ ] Test: Send with mock API

---

## Success Criteria

- [ ] Candidate can register, see fee, upload proof
- [ ] Registration fee can be waived during promo campaigns
- [ ] Candidates auto-assigned to consultants (configurable algorithm)
- [ ] Consultants log each interaction with category/obstacle/remarks
- [ ] Supervisors can review interactions and provide suggestions
- [ ] Commitment generates tuition billing with installment options
- [ ] Payment tracking with verification
- [ ] Documents can be deferred (ijazah/transcript)
- [ ] Enrollment validates requirements, generates NIM
- [ ] Referrer commissions tracked and paid
- [ ] Campaign ROI trackable
- [ ] All HTMX interactions (no full page reloads)
- [ ] Mobile responsive
