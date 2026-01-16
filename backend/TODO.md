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

REWARD_CONFIG (for external referrers)
├─ id, referrer_type (alumni, teacher, student, partner, staff)
├─ reward_type (cash, tuition_discount, merchandise)
├─ amount (fixed or percentage)
├─ is_percentage
├─ trigger_event (enrollment, commitment, registration)
├─ description
└─ is_active

MGM_REWARD_CONFIG (for enrolled student referrals)
├─ id, academic_year
├─ reward_type (tuition_discount, cash, merchandise)
├─ referrer_amount (reward for the referring student)
├─ referee_amount (discount for new candidate, optional)
├─ trigger_event (enrollment, commitment)
├─ description
└─ is_active

REFERRER
├─ id, name
├─ type (alumni, teacher, student, partner, staff)
├─ institution (helps match claims)
├─ phone, email, code (all optional)
├─ bank_name, bank_account, account_holder
├─ commission_override (nullable, overrides reward_config)
├─ payout_preference (monthly, per_enrollment)
└─ is_active

CANDIDATE
├─ id, created_at, updated_at
├─ name, email, phone, whatsapp
├─ password_hash
├─ email_verified_at, phone_verified_at
├─ address, city, province
├─ high_school, graduation_year
├─ prodi_id
├─ source_type, source_detail (referral claim text)
├─ campaign_id, referrer_id, referrer_verified_at
├─ referred_by_candidate_id (for member-get-member)
├─ status (registered → prospecting → committed → enrolled / lost)
├─ assigned_consultant_id
├─ registration_fee_paid_at
├─ lost_at, lost_reason_id
├─ enrolled_at, nim
└─ referral_code (generated on enrollment, for MGM)

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

VERIFICATION_TOKEN (email/phone OTP)
├─ id, candidate_id
├─ type (email, phone)
├─ token (6-digit OTP)
├─ expires_at
└─ verified_at

ANNOUNCEMENT
├─ id, created_at
├─ title, content
├─ target_status (null=all, or specific status)
├─ target_prodi_id (null=all)
├─ published_at
└─ is_active

ANNOUNCEMENT_READ
├─ id, announcement_id, candidate_id
└─ read_at
```

### Migration Plan (by dependency)

```
Phase 1 - Admin Foundation (DONE):
  001_create_users.sql                ✅
  002_create_prodis.sql               ✅
  003_create_fee_types.sql            ✅
  004_create_fee_structures.sql       ✅ → fee_types, prodis
  005_create_interaction_categories.sql ✅
  006_create_obstacles.sql            ✅
  007_seed_phase1_data.sql            ✅ → fee_types, categories, obstacles, prodis

Phase 2 - Configuration (DONE):
  008_create_campaigns.sql            ✅
  009_create_reward_configs.sql       ✅
  010_create_mgm_reward_configs.sql   ✅
  011_seed_reward_configs.sql         ✅ → reward_configs, mgm_reward_configs
  012_create_referrers.sql            ✅

Phase 2.5 - Remaining Configuration:
  013_create_assignment_algorithms.sql  ✅
  014_create_document_types.sql         ✅
  015_create_lost_reasons.sql

Phase 3 - Candidates:
  016_create_announcements.sql        → prodis
  017_create_candidates.sql           → prodis, campaigns, referrers, users, lost_reasons

Phase 4 - Transactions:
  018_create_billings.sql             → candidates, fee_types
  019_create_payments.sql             → billings, users
  020_create_interactions.sql         → candidates, users, categories, obstacles
  021_create_documents.sql            → candidates, document_types
  022_create_commission_ledger.sql    → candidates, referrers
  023_create_notification_logs.sql    → candidates
  024_create_verification_tokens.sql  → candidates
  025_create_announcement_reads.sql   → announcements, candidates
```

---

# Phase 0: UI Mockup

Clickable prototype untuk validasi dengan stakeholder sebelum implementasi.

---

## Feature 3: UI Mockup (No Database) ✅

Semua halaman dengan data hardcoded untuk demo dan validasi UI/UX.

**Tujuan:**
- Validasi alur dan tampilan dengan tim sales/marketing
- Dapat diklik dan dinavigasi seperti aplikasi asli
- Tidak perlu database atau backend logic
- Mendapat commitment stakeholder sebelum implementasi

**Admin Pages:**
- [x] Login page (mock, langsung redirect)
- [x] Dashboard konsultan (statistik hardcoded)
- [x] Dashboard supervisor (statistik tim hardcoded)
- [x] Daftar kandidat (tabel dengan filter, data dummy)
- [x] Detail kandidat (info + timeline interaksi dummy)
- [x] Form log interaksi (modal)
- [x] Settings: User management (list dummy users)
- [x] Settings: Prodi (list dummy)
- [x] Settings: Fee structure (matrix dummy)
- [x] Settings: Reward config (list dummy)
- [x] Settings: Kategori & hambatan (list dummy)
- [x] Kampanye management (list + form)
- [x] Referrer management (list + form)
- [x] Referral claims (list unverified)
- [x] Komisi/commission (list + approve)
- [x] Laporan funnel (chart dummy)
- [x] Laporan performa konsultan (table dummy)

**Portal Kandidat Pages:**
- [x] Landing/register page
- [x] Login page
- [x] Dashboard kandidat (status, checklist)
- [x] Upload dokumen (list + upload form)
- [x] Pembayaran (list tagihan + upload bukti)
- [x] Pengumuman (list + detail)
- [x] Referral MGM (kode + list referred)

**Public Pages:**
- [x] Form pendaftaran multi-step (6 langkah)
- [x] Halaman verifikasi OTP (mock)
- [x] Halaman sukses registrasi

**Navigation:**
- [x] Admin sidebar dengan semua menu
- [x] Portal sidebar dengan semua menu
- [x] Responsive mobile view

**Demo Data:**
- [x] 10 dummy candidates dengan berbagai status
- [x] 3 dummy consultants
- [x] 5 dummy interactions per candidate
- [x] 2 dummy prodis
- [x] Sample fee structure
- [x] Sample campaigns & referrers

---

# Phase 1: Admin Foundation

Must complete before opening registration.

---

## Feature 4: Staff Login (Google OAuth) ✅

All staff (admin, supervisor, consultant) login with domain-restricted Google.

**Migrations:** 001_create_users

**Setup:**
- [x] Create Google Cloud project (or use existing)
- [x] Enable Google+ API / People API
- [x] Configure OAuth consent screen (internal for workspace domain)
- [x] Create OAuth 2.0 credentials (Web application)
- [x] Add authorized redirect URIs: `{BASE_URL}/admin/auth/google/callback`
- [x] Store credentials: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET env vars
- [x] Set STAFF_EMAIL_DOMAIN env var (e.g., tazkia.ac.id)

**Implementation:**
- [x] `model/user.go` - Create, FindByEmail, FindByGoogleID
- [x] `auth/google.go` - OAuth flow (GetAuthURL, ExchangeCode, GetUserInfo)
- [x] `handler/admin_auth.go` - GET /admin/login, /admin/auth/google, callback
- [x] Domain check (validate email ends with STAFF_EMAIL_DOMAIN)
- [x] Auto-create user with role=consultant if valid domain
- [x] `handler/middleware.go` - RequireAuth, RequireRole
- [x] Cookie-based session (HttpOnly JWT, SameSite=Lax)
- [x] User widget in sidebar showing logged-in user name and role
- [x] Logout functionality
- [ ] Test: Login with valid/invalid domain

---

## Feature 5: Staff Management ✅

Admin manages staff accounts (admin, supervisor, consultant).

**Migrations:** 001_create_users ✅

- [x] `model/user.go` - Create, FindByEmail, FindByGoogleID (basic)
- [x] `model/user.go` - List, UpdateRole, ToggleActive, SetSupervisor
- [x] `templates/admin/settings/users.templ` - User list with HTMX inline editing
- [x] Wire `handler/admin.go` to real data - GET /admin/settings/users, POST update role/active/supervisor
- [x] Roles: admin (full access), supervisor (team + suggestions), consultant (own candidates)
- [x] Assign supervisor to consultants (hierarchy)
- [x] Toggle active status (for assignment pool)
- [x] Playwright E2E tests for role changes, status toggle, supervisor assignment

---

## Feature 6: Settings - Prodi Management ✅

Admin configures available programs.

**Migrations:** 002_create_prodis ✅

- [x] `model/prodi.go` - CRUD, ListActive (model layer done)
- [x] `templates/admin/settings/prodis.templ` - Prodi list with HTMX forms
- [x] Wire `handler/admin.go` to real data - CRUD /admin/settings/prodis
- [x] Fields: name, code, degree (S1/D3), is_active
- [x] HTMX: Inline edit without reload (toggle status, edit modal)
- [x] Playwright E2E tests: Add, edit, toggle status with database persistence

---

## Feature 7: Settings - Fee Structure ✅

Admin configures fees per prodi and academic year.

**Migrations:** 003_create_fee_types ✅, 004_create_fee_structures ✅

- [x] `model/fee_type.go` - List (seeded: registration, tuition, dormitory)
- [x] `model/fee_structure.go` - CRUD, FindByTypeAndProdi (model layer done)
- [x] `templates/admin/settings/fees.templ` - Fee matrix with HTMX forms
- [x] Wire `handler/admin.go` to real data - CRUD /admin/settings/fees
- [x] Set: registration fee (global), tuition per prodi, dormitory (global)
- [ ] Installment options per fee type (deferred to Phase 5)
- [x] Playwright E2E tests: Add, edit, toggle status with database persistence

---

## Feature 8: Settings - Categories & Obstacles ✅

Supervisor manages interaction categories and obstacles.

**Migrations:** 005_create_interaction_categories ✅, 006_create_obstacles ✅, 007_seed_phase1_data ✅

- [x] `model/interaction_category.go` - CRUD, ListActive (model layer done)
- [x] `model/obstacle.go` - CRUD, ListActive (model layer done)
- [x] `templates/admin/settings/categories.templ` - Category and Obstacle list with HTMX forms
- [x] Wire `handler/admin.go` to real data - CRUD for categories, obstacles
- [x] Seed default categories: Tertarik, Mempertimbangkan, Ragu-ragu, Dingin, Tidak bisa dihubungi
- [x] Seed default obstacles: Biaya terlalu mahal, Lokasi jauh, Orang tua belum setuju, Waktu belum tepat, Memilih kampus lain
- [x] Playwright E2E tests: Add, edit, toggle status for both categories and obstacles

---

# Phase 2: Configuration

Setup before opening registration.

---

## Feature 9: Campaign Management ✅

Admin manages campaigns with promo pricing.

**Migrations:** 008_create_campaigns ✅

- [x] `model/campaign.go` - CRUD, FindActive, FindByCode
- [x] `templates/admin/settings_campaigns.templ` - Campaign settings page with HTMX CRUD
- [x] `handler/admin.go` - CRUD /admin/settings/campaigns
- [x] Fields: name, type, channel, dates, registration_fee_override
- [x] Fee override: fixed amount (nullable, overrides default registration fee)
- [ ] Generate UTM-compatible tracking code (deferred enhancement)
- [x] Playwright E2E tests: CRUD with database persistence

---

## Feature 10: Reward Configuration ✅

Configure default rewards by referrer type and MGM.

**Migrations:** 009_create_reward_configs ✅, 010_create_mgm_reward_configs ✅, 011_seed_reward_data ✅

- [x] `model/reward_config.go` - CRUD, FindByType
- [x] `model/mgm_reward_config.go` - CRUD, FindActive (FindByYear)
- [x] `templates/admin/settings_rewards.templ` - Reward config list with HTMX CRUD
- [x] `handler/admin.go` - CRUD /admin/settings/rewards, /admin/settings/mgm-rewards
- [x] External referrer rewards by type: alumni, teacher, student, partner, staff
- [x] Fields: reward_type (cash, tuition_discount, merchandise), amount, is_percentage, trigger_event
- [x] MGM rewards: referrer_amount (for enrolled student), referee_amount (for new candidate)
- [x] Seed defaults: alumni Rp500k, teacher Rp750k, student Rp300k, partner Rp1M, staff Rp250k, MGM Rp200k + 10% tuition
- [x] Playwright E2E tests: CRUD with database persistence

---

## Feature 11: Referrer Management ✅

Admin manages referrers for commission tracking.

**Migrations:** 012_create_referrers ✅

- [x] `model/referrer.go` - CRUD, GenerateCode, FindByCode, SearchByName
- [x] `templates/admin/settings_referrers.templ` - Referrer list with HTMX forms
- [x] `handler/admin.go` - CRUD /admin/settings/referrers
- [x] Fields: name, type, institution, contact, bank details
- [x] commission_override: optional, overrides reward_config default
- [x] payout_preference: monthly or per_enrollment
- [ ] Generate optional referral code (for partners who want trackable links)
- [x] Playwright E2E tests: CRUD, status toggle with database persistence

---

## Feature 12: Settings - Assignment Algorithm ✅

Configure consultant assignment algorithm.

**Migrations:** 013_create_assignment_algorithms ✅

- [x] `model/assignment_algorithm.go` - List, SetActive, FindActive
- [x] `templates/admin/settings_assignment.templ` - Algorithm selection with HTMX
- [x] `handler/admin.go` - GET /admin/settings/assignment, POST /admin/settings/assignment/{id}/activate
- [x] Algorithms: round_robin, load_balanced, performance_weighted, activity_based (seeded)
- [x] Only one algorithm active at a time (unique index constraint)
- [x] Playwright E2E tests: Display, switch active algorithm

---

## Feature 13: Settings - Document Types ✅

Configure required documents.

**Migrations:** 014_create_document_types ✅

- [x] `model/document_type.go` - CRUD, ListActive
- [x] `templates/admin/settings_documents.templ` - Document type list with HTMX
- [x] `handler/admin.go` - CRUD /admin/settings/document-types
- [x] Fields: name, code, description, is_required, can_defer, max_file_size_mb, display_order
- [x] Seed: KTP (required), Photo (required), Ijazah (required, can_defer), Transcript (required, can_defer)
- [x] Playwright E2E tests: CRUD, toggle status with database persistence

---

# Phase 3: Public Registration & Candidate Portal

Candidate-facing features.

---

## Feature 14: Candidate Registration

Registration with password and email/phone verification.

**Migrations:** 018, 025

**Setup - Email (for OTP):**
- [ ] Choose email provider (SMTP, SendGrid, AWS SES, Resend, etc.)
- [ ] Create account and get API key/credentials
- [ ] Configure sender email and domain verification
- [ ] Store credentials: SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS (or provider-specific)

**Setup - WhatsApp (for OTP):**
- [ ] Choose WhatsApp API provider (Fonnte, Wablas, Twilio, etc.)
- [ ] Create account and get API key
- [ ] Register sender phone number
- [ ] Store credentials: WHATSAPP_API_URL, WHATSAPP_API_KEY

**Implementation:**
- [ ] `model/candidate.go` - Create, FindByEmail, FindByPhone
- [ ] `model/verification_token.go` - Create, Verify, Cleanup
- [ ] `integration/email.go` - SendOTP via email
- [ ] `integration/whatsapp.go` - SendOTP via WhatsApp
- [ ] `templates/public/register.templ` - Multi-step registration form
- [ ] `handler/public.go` - GET/POST /register, /verify-email, /verify-phone
- [ ] Step 1: Account (email, password, phone)
- [ ] Step 2: Email verification (6-digit OTP via email)
- [ ] Step 3: Phone verification (6-digit OTP via WhatsApp)
- [ ] Step 4: Personal info (name, address, city, province)
- [ ] Step 5: Education (high_school, graduation_year, prodi)
- [ ] Step 6: Source tracking (source_type dropdown, source_detail if referral)
- [ ] Source types: instagram, google, tiktok, youtube, expo, school_visit, friend_family, teacher_alumni, walkin, other
- [ ] If URL has `?ref=CODE`, auto-link to referrer or referred_by_candidate
- [ ] If URL has `?utm_campaign=X`, auto-link to campaign
- [ ] Auto-assign to consultant (using active algorithm)
- [ ] Hash password with bcrypt
- [ ] Test: Registration flow, OTP verification

---

## Feature 15: Candidate Login

Candidate authenticates with email + password.

- [ ] `model/candidate.go` - Authenticate
- [ ] `templates/public/login.templ` - Login form
- [ ] `handler/public.go` - GET/POST /login, /logout
- [ ] Validate email is verified
- [ ] Cookie-based session (HttpOnly JWT, 30-day expiry)
- [ ] Redirect to portal dashboard after login
- [ ] Test: Login, session persistence

---

## Feature 16: Candidate Portal - Dashboard

Overview of candidate status and actions.

- [ ] `templates/portal/dashboard.templ` - Status summary
- [ ] `handler/portal.go` - GET /portal
- [ ] `handler/middleware.go` - RequireCandidateAuth
- [ ] Show: status badge, assigned consultant contact, registration fee status
- [ ] Show: document checklist (uploaded/pending/rejected)
- [ ] Show: unread announcements count
- [ ] Quick links: upload documents, view payments, announcements
- [ ] Test: Dashboard displays correct info

---

## Feature 17: Candidate Portal - Documents

Candidate uploads and tracks documents.

**Migrations:** 022

**Setup - File Storage:**
- [ ] Choose storage provider (Cloudflare R2, AWS S3, local disk, etc.)
- [ ] Create bucket/container for documents
- [ ] Configure CORS if using object storage
- [ ] Set up access credentials and permissions
- [ ] Store credentials: STORAGE_TYPE, STORAGE_BUCKET, STORAGE_ACCESS_KEY, STORAGE_SECRET_KEY, STORAGE_ENDPOINT

**Implementation:**
- [ ] `storage/storage.go` - Upload, Download, Delete interface
- [ ] `storage/r2.go` or `storage/s3.go` - Provider implementation
- [ ] `model/document.go` - Upload, ListByCandidate
- [ ] `templates/portal/documents.templ` - Upload form with status
- [ ] `handler/portal.go` - GET/POST /portal/documents
- [ ] List: document type, status (pending/approved/rejected), rejection reason
- [ ] Upload: file picker with type/size validation
- [ ] Re-upload rejected documents
- [ ] Show deferrable documents with note
- [ ] Test: Upload, re-upload, status display

---

## Feature 18: Candidate Portal - Payments

Candidate views billing and uploads payment proof.

**Migrations:** 019, 020

- [ ] `model/billing.go` - Create, FindByCandidate
- [ ] `model/payment.go` - UploadProof
- [ ] `templates/portal/payments.templ` - Billing list with installments
- [ ] `handler/portal.go` - GET /portal/payments, POST /portal/payments/{id}/proof
- [ ] List all billings: registration, tuition, dormitory
- [ ] Show per billing: total, paid, remaining, installments
- [ ] Show per installment: amount, due date, status, proof
- [ ] Upload payment proof (image)
- [ ] Registration fee: generated on registration
- [ ] Tuition/dormitory: generated on commitment (by admin)
- [ ] Test: Billing display, proof upload

---

## Feature 19: Candidate Portal - Announcements

Candidate receives targeted announcements.

**Migrations:** 017, 026

- [ ] `model/announcement.go` - ListForCandidate, MarkRead
- [ ] `templates/portal/announcements.templ` - Announcement list
- [ ] `handler/portal.go` - GET /portal/announcements, POST mark-read
- [ ] Filter by: target_status, target_prodi (or null for all)
- [ ] Show: title, preview, published date, read status
- [ ] Detail view with full content
- [ ] Mark as read on open
- [ ] Test: Filtering, read status

---

## Feature 20: Announcement Management (Admin)

Admin creates and targets announcements.

- [ ] `model/announcement.go` - CRUD
- [ ] `templates/admin/announcements.templ` - Announcement list
- [ ] `templates/admin/announcement_form.templ` - Create/edit form
- [ ] `handler/admin.go` - CRUD /admin/announcements
- [ ] Fields: title, content (markdown), target_status, target_prodi, published_at
- [ ] Schedule future publish
- [ ] Preview before publish
- [ ] Test: CRUD, targeting

---

## Feature 21: Member Get Member Referral

Enrolled students refer new candidates.

- [ ] `model/candidate.go` - GenerateReferralCode, FindByReferralCode
- [ ] `templates/portal/referral.templ` - Referral dashboard
- [ ] `handler/portal.go` - GET /portal/referral
- [ ] Generate unique referral_code on enrollment
- [ ] Show: shareable link with code
- [ ] List: referred candidates with status
- [ ] Show: reward status (pending/earned based on enrollment)
- [ ] Registration form: if `?ref=CODE` matches candidate, set referred_by_candidate_id
- [ ] Test: Code generation, referral tracking

---

# Phase 4: CRM Operations

Day-to-day sales operations.

---

## Feature 22: Candidate List & Filters

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

## Feature 23: Candidate Detail & Timeline

View candidate info and history.

- [ ] `templates/admin/candidate_detail.templ` - Info + timeline
- [ ] `handler/admin.go` - GET /admin/candidates/{id}
- [ ] Show: personal info, prodi, source, campaign/referrer, status, assigned consultant
- [ ] Timeline: interactions, payments, documents, status changes
- [ ] Quick actions: log interaction, reassign, change status
- [ ] Test: Detail view, timeline ordering

---

## Feature 24: Interaction Logging

Consultants log each contact.

**Migrations:** 021

- [ ] `model/interaction.go` - Create, ListByCandidate, ListByConsultant
- [ ] `templates/admin/interaction_form.templ` - Log interaction modal
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/interactions
- [ ] Fields: channel, category, obstacle (optional), remarks, next_followup_date
- [ ] Channels: call, whatsapp, email, campus_visit, home_visit
- [ ] Auto-update candidate last_interaction_at
- [ ] Test: Create interaction, list

---

## Feature 25: Supervisor Suggestions

Supervisor reviews and provides guidance.

- [ ] `templates/admin/candidate_detail.templ` - Suggestion field in timeline
- [ ] `handler/admin.go` - POST /admin/interactions/{id}/suggestion
- [ ] Consultant sees suggestion, marks as read
- [ ] Notification badge for unread suggestions
- [ ] Test: Add suggestion, mark as read

---

## Feature 26: Consultant Assignment

Manual reassignment of candidates.

- [ ] `model/candidate.go` - Assign, GetAssignmentStats
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/assign
- [ ] Supervisor/admin can reassign candidates
- [ ] Show consultant workload in dropdown
- [ ] Log assignment change in timeline
- [ ] Test: Reassignment, workload display

---

## Feature 27: Referral Claim Verification

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

## Feature 28: Commitment & Tuition Billing

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

## Feature 29: Payment Tracking

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

## Feature 30: Document Review

Admin reviews uploaded documents.

- [ ] `model/document.go` - UpdateStatus
- [ ] `templates/admin/document_review.templ` - Review modal
- [ ] `handler/admin.go` - GET /admin/candidates/{id}/documents, POST approve/reject
- [ ] View document, approve or reject with reason
- [ ] Notify candidate of rejection (for re-upload)
- [ ] Test: Review flow, rejection

---

## Feature 31: Enrollment

Mark candidate as enrolled.

- [ ] `model/candidate.go` - Enroll, GenerateNIM
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/enroll
- [ ] Validation:
  - Registration fee: paid
  - Tuition: at least 1st installment paid
  - Documents: KTP + photo approved (ijazah/transcript can be pending)
- [ ] Generate NIM: YYYY + PRODI_CODE + SEQUENCE (e.g., 2026SI001)
- [ ] Generate referral_code for member-get-member (e.g., MGM-2026SI001)
- [ ] Change status: committed → enrolled
- [ ] Trigger commission creation if referred (by referrer or by enrolled candidate)
- [ ] Test: Enrollment validation, NIM generation, referral code

---

## Feature 32: Lost Candidate

Mark candidate as lost.

**Migrations:** 011, 027 (seed)

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

## Feature 33: Commission Tracking

Auto-create and track commissions.

**Migrations:** 023

- [ ] `model/commission.go` - Create, ListByReferrer, ListPending
- [ ] Auto-create commission when referred candidate enrolls
- [ ] Amount from referrer.commission_per_enrollment
- [ ] Status: pending → approved → paid
- [ ] Test: Auto-creation on enrollment

---

## Feature 34: Commission Payout

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

## Feature 35: Dashboard - Consultant

Consultant's daily view.

- [ ] `model/stats.go` - GetConsultantStats
- [ ] `templates/admin/dashboard_consultant.templ`
- [ ] `handler/admin.go` - GET /admin (role-based)
- [ ] Show: my candidates by status, overdue followups, today's tasks
- [ ] Quick access to pending followups
- [ ] Test: Dashboard data accuracy

---

## Feature 36: Dashboard - Supervisor

Supervisor's team view.

- [ ] `model/stats.go` - GetTeamStats, GetFunnelStats
- [ ] `templates/admin/dashboard_supervisor.templ`
- [ ] Show: team funnel, consultant leaderboard, stuck candidates (> 7 days no interaction)
- [ ] Common obstacles this period
- [ ] Test: Dashboard data accuracy

---

## Feature 37: Reports - Funnel

Conversion funnel analysis.

- [ ] `model/stats.go` - GetFunnelByDateRange, GetFunnelByProdi
- [ ] `templates/admin/reports/funnel.templ`
- [ ] `handler/admin.go` - GET /admin/reports/funnel
- [ ] Filter by: date range, prodi, campaign
- [ ] Show: registered → prospecting → committed → enrolled, with conversion rates
- [ ] Test: Report accuracy

---

## Feature 38: Reports - Consultant Performance

Individual performance metrics.

- [ ] `model/stats.go` - GetConsultantPerformance
- [ ] `templates/admin/reports/consultants.templ`
- [ ] `handler/admin.go` - GET /admin/reports/consultants
- [ ] Metrics: candidates handled, success rate, avg days to commit, interaction frequency
- [ ] Ranking by success rate
- [ ] Test: Metrics calculation

---

## Feature 39: Reports - Campaign ROI

Campaign effectiveness.

- [ ] `model/stats.go` - GetCampaignStats
- [ ] `templates/admin/reports/campaigns.templ`
- [ ] `handler/admin.go` - GET /admin/reports/campaigns
- [ ] Show: leads, commits, enrollments, conversion rate per campaign
- [ ] Cost per enrollment (if cost data available)
- [ ] Test: Report accuracy

---

## Feature 40: Reports - Referrer Leaderboard

Referrer performance.

- [ ] `model/stats.go` - GetReferrerStats
- [ ] `templates/admin/reports/referrers.templ`
- [ ] `handler/admin.go` - GET /admin/reports/referrers
- [ ] Show: referrals, enrollments, conversion rate, commission earned/paid
- [ ] Test: Report accuracy

---

## Feature 41: CSV Export

Export data for external analysis.

- [ ] `handler/admin.go` - GET /admin/export/candidates, /admin/export/interactions
- [ ] Filter by: date range, status, consultant, campaign
- [ ] Include all relevant fields
- [ ] Test: Export with filters

---

# Phase 8: Notifications

Communication automation.

---

## Feature 42: WhatsApp Notifications

Send notifications at key events.

**Migrations:** 024

- [ ] `integration/whatsapp.go` - Send via API (Fonnte/similar)
- [ ] `model/notification.go` - Create, ListByCandidate
- [ ] Templates: registration_confirmed, payment_reminder, document_reminder, enrolled
- [ ] Manual send from candidate detail
- [ ] Log all sent messages
- [ ] Test: Send with mock API

---

## Success Criteria

### Admin/Staff
- [ ] Staff (admin, supervisor, consultant) can login with Google OAuth
- [ ] Admin can configure prodis/fees/campaigns before opening registration
- [ ] Consultants log interactions with category/obstacle/remarks
- [ ] Supervisors provide suggestions on interactions
- [ ] Commitment generates tuition billing with installments
- [ ] Payment tracking with proof verification
- [ ] Document review with approve/reject
- [ ] Enrollment validates requirements, generates NIM
- [ ] All admin interactions use HTMX (no full page reloads)

### Candidate Portal
- [ ] Candidate registers with email/phone verification
- [ ] Candidate logs in with password
- [ ] Candidate views dashboard with status summary
- [ ] Candidate uploads/re-uploads documents
- [ ] Candidate views billing and uploads payment proof
- [ ] Candidate receives targeted announcements
- [ ] Enrolled candidate can refer new candidates (member-get-member)

### Business
- [ ] Registration fee waived during promo campaigns
- [ ] Candidates auto-assigned to consultants
- [ ] Documents can be deferred (ijazah/transcript)
- [ ] Referrer commissions auto-created and tracked
- [ ] Campaign ROI trackable via reports
- [ ] Mobile responsive
