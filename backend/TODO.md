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

Phase 2.5 - Remaining Configuration (DONE):
  013_create_assignment_algorithms.sql  ✅
  014_create_document_types.sql         ✅
  015_create_lost_reasons.sql           ✅

Phase 3 - Candidates (DONE):
  016_create_candidates.sql             ✅ → prodis, campaigns, referrers, users, lost_reasons
  017_create_verification_tokens.sql    ✅ → candidates
  018_document_source_types.sql         ✅ (seed data)
  019_create_announcements.sql          ✅ → prodis
  020_create_interactions.sql           ✅ → candidates, users, categories, obstacles

Phase 4 - Transactions:
  021_create_billings.sql             → candidates, fee_types
  022_create_payments.sql             → billings, users
  023_create_documents.sql            → candidates, document_types
  024_create_commission_ledger.sql    → candidates, referrers
  025_create_notification_logs.sql    → candidates
  026_create_announcement_reads.sql   → announcements, candidates
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

## Feature 14: Settings - Lost Reasons ✅

Configure reasons for lost candidates.

**Migrations:** 015_create_lost_reasons ✅

- [x] `model/lost_reason.go` - CRUD, ListActive
- [x] `templates/admin/settings_lost_reasons.templ` - Lost reason list with HTMX
- [x] `handler/admin.go` - CRUD /admin/settings/lost-reasons
- [x] Fields: name, description, display_order, is_active
- [x] Seed: No response, Chose competitor, Financial issues, Not qualified, Bad timing, Location, Parents disagree, Other
- [x] Playwright E2E tests: CRUD, toggle status with database persistence

---

# Phase 3: Public Registration & Candidate Portal

Candidate-facing features.

---

## Feature 15: Candidate Registration ✅

Registration with password. Email/phone verification optional.

**Migrations:** 016_create_candidates ✅, 017_create_verification_tokens ✅, 018_document_source_types ✅

**Setup - Email (for OTP) - Optional:**
- [x] Provider: Resend (configured via RESEND_API_KEY, RESEND_FROM)
- [x] Integration created but optional - candidates can register without verification

**Setup - WhatsApp (for OTP) - Optional:**
- [x] Provider: Custom WhatsApp API (configured via WHATSAPP_API_URL, WHATSAPP_API_TOKEN)
- [x] Integration created but optional - candidates can register without verification

**Implementation:**
- [x] `model/candidate.go` - Create, FindByEmail, FindByPhone, FindByID
- [x] `model/candidate.go` - UpdatePersonalInfo, UpdateEducation, UpdateSourceTracking
- [x] `model/candidate.go` - Authenticate, SetEmailVerified, SetPhoneVerified
- [x] `model/candidate.go` - AssignConsultant (using active algorithm)
- [x] `model/verification_token.go` - CreateToken, VerifyToken, CleanupExpired
- [x] `integration/resend.go` - SendOTP via Resend (optional)
- [x] `integration/whatsapp.go` - SendOTP via WhatsApp (optional)
- [x] `templates/portal/registration.templ` - 4-step registration form
- [x] `handler/public.go` - Registration routes
- [x] `auth/session.go` - CreateCandidateToken for candidate sessions
- [x] Step 1: Account (email or phone required, password)
- [x] Step 2: Personal info (name, address, city, province)
- [x] Step 3: Education (high_school, graduation_year, prodi)
- [x] Step 4: Source tracking (source_type dropdown, source_detail if referral)
- [x] Source types seeded: instagram, google, tiktok, youtube, expo, school_visit, friend_family, teacher_alumni, walkin, other
- [x] If URL has `?ref=CODE`, auto-link to referrer
- [x] If URL has `?utm_campaign=X`, auto-link to campaign
- [x] Auto-assign to consultant (using active assignment algorithm)
- [x] Hash password with bcrypt
- [x] E2E Test: Full registration flow
- [x] E2E Test: Login by email and phone

---

## Feature 16: Candidate Login ✅

Candidate authenticates with email or phone + password.

- [x] `model/candidate.go` - Authenticate (by email or phone)
- [x] `templates/portal/login.templ` - Login form
- [x] `handler/public.go` - GET/POST /login, /logout
- [x] Email/phone verification not required (optional)
- [x] Cookie-based session (HttpOnly JWT)
- [x] Redirect to portal dashboard after login
- [x] E2E Test: Login by email, login by phone, session persistence

---

## Feature 17: Candidate Portal - Dashboard ✅

Overview of candidate status and actions.

- [x] `templates/portal/dashboard.templ` - Status summary with checklist
- [x] `handler/portal.go` - GET /portal/dashboard
- [x] `handler/middleware.go` - RequireCandidateAuth middleware
- [x] Show: status badge, assigned consultant contact
- [x] Show: checklist (verification, personal info, education, documents, payment)
- [x] Show: recent announcements
- [x] Quick links: upload documents, view payments, announcements
- [x] E2E Test: Dashboard displays correct info

---

## Feature 18: Candidate Portal - Documents

Candidate uploads and tracks documents.

**Migrations:** 023_create_documents (pending)

**Setup - File Storage:**
- [ ] Choose storage provider (Cloudflare R2, AWS S3, local disk, etc.)
- [ ] Create bucket/container for documents
- [ ] Configure CORS if using object storage
- [ ] Set up access credentials and permissions
- [ ] Store credentials: STORAGE_TYPE, STORAGE_BUCKET, STORAGE_ACCESS_KEY, STORAGE_SECRET_KEY, STORAGE_ENDPOINT

**Implementation:**
- [x] `templates/portal/documents.templ` - Upload form with status (UI mockup done)
- [ ] `storage/storage.go` - Upload, Download, Delete interface
- [ ] `storage/r2.go` or `storage/s3.go` - Provider implementation
- [ ] `model/document.go` - Upload, ListByCandidate
- [ ] `handler/portal.go` - GET/POST /portal/documents (wire to real data)
- [ ] List: document type, status (pending/approved/rejected), rejection reason
- [ ] Upload: file picker with type/size validation
- [ ] Re-upload rejected documents
- [ ] Show deferrable documents with note
- [ ] Test: Upload, re-upload, status display

---

## Feature 19: Candidate Portal - Payments

Candidate views billing and uploads payment proof.

**Migrations:** 021_create_billings, 022_create_payments (pending)

- [x] `templates/portal/payments.templ` - Billing list with installments (UI mockup done)
- [ ] `model/billing.go` - Create, FindByCandidate
- [ ] `model/payment.go` - UploadProof
- [ ] `handler/portal.go` - GET /portal/payments, POST /portal/payments/{id}/proof (wire to real data)
- [ ] List all billings: registration, tuition, dormitory
- [ ] Show per billing: total, paid, remaining, installments
- [ ] Show per installment: amount, due date, status, proof
- [ ] Upload payment proof (image)
- [ ] Registration fee: generated on registration
- [ ] Tuition/dormitory: generated on commitment (by admin)
- [ ] Test: Billing display, proof upload

---

## Feature 20: Candidate Portal - Announcements ✅

Candidate receives targeted announcements.

**Migrations:** 019_create_announcements ✅, 026_create_announcement_reads ✅

- [x] `templates/portal/announcements.templ` - Announcement list with HTMX mark-read
- [x] `model/announcement.go` - ListForCandidate, MarkRead, CountUnread
- [x] `handler/portal.go` - GET /portal/announcements, POST mark-read (wired to real data)
- [x] Filter by: target_status, target_prodi (or null for all)
- [x] Show: title, preview, published date, read status
- [x] Detail view with full content
- [x] Mark as read via HTMX
- [x] E2E Test: Covered in portal-dashboard.spec.ts

---

## Feature 21: Announcement Management (Admin) ✅

Admin creates and targets announcements.

**Migrations:** 019_create_announcements ✅

- [x] `model/announcement.go` - CRUD, Publish, Unpublish
- [x] `templates/admin/settings_announcements.templ` - Announcement list with HTMX CRUD
- [x] `handler/admin_announcements.go` - Full CRUD handlers
- [x] Routes wired in `handler/admin.go`
- [x] Fields: title, content, target_status, target_prodi
- [x] Publish/unpublish toggle
- [x] E2E Test: CRUD, targeting (announcements.spec.ts)

---

## Feature 22: Member Get Member Referral

Enrolled students refer new candidates.

- [x] `templates/portal/referral.templ` - Referral dashboard (UI mockup done)
- [ ] `model/candidate.go` - GenerateReferralCode, FindByReferralCode
- [ ] `handler/portal.go` - GET /portal/referral (wire to real data)
- [ ] Generate unique referral_code on enrollment
- [ ] Show: shareable link with code
- [ ] List: referred candidates with status
- [ ] Show: reward status (pending/earned based on enrollment)
- [x] Registration form: `?ref=CODE` support implemented
- [ ] Test: Code generation, referral tracking

---

# Phase 4: CRM Operations

Day-to-day sales operations.

---

## Feature 23: Candidate List & Filters ✅

Admin/consultant views candidates.

- [x] `templates/admin/candidates.templ` - Table with filters and HTMX
- [x] `model/candidate.go` - ListCandidates with filters, pagination, role-based visibility
- [x] `handler/admin_candidates.go` - GET /admin/candidates (wired to real data)
- [x] Filters: status, assigned consultant, prodi, campaign, source_type, search
- [x] Sort: newest (default), name, status
- [x] Consultant sees only their candidates, supervisor sees team, admin sees all
- [x] HTMX: Filter without reload
- [x] E2E Test: Filter combinations, role-based visibility (candidates.spec.ts)

---

## Feature 24: Candidate Detail & Timeline ✅

View candidate info and history.

- [x] `templates/admin/candidate_detail.templ` - Info + timeline with HTMX
- [x] `handler/admin_candidates.go` - GET /admin/candidates/{id} (wired to real data)
- [x] Show: personal info, prodi, source, campaign/referrer, status, assigned consultant
- [x] Timeline: interactions with consultant name and timestamps
- [x] Quick actions: log interaction button
- [x] E2E Test: Detail view (candidates.spec.ts)

---

## Feature 25: Interaction Logging ✅

Consultants log each contact.

**Migrations:** 020_create_interactions ✅

- [x] `templates/admin/interaction_form.templ` - Full interaction form page
- [x] `model/interaction.go` - CreateInteraction, ListInteractionsByCandidate
- [x] `handler/admin_interactions.go` - GET/POST /admin/candidates/{id}/interaction
- [x] Fields: channel, category, obstacle (optional), remarks, next_followup_date
- [x] Channels: call, whatsapp, email, campus_visit, home_visit
- [x] E2E Test: Create interaction, verify in timeline (candidates-crud.spec.ts)

---

## Feature 26: Supervisor Suggestions

Supervisor reviews and provides guidance.

- [x] `templates/admin/consultant_dashboard.templ` - Suggestions section (UI mockup done)
- [ ] `templates/admin/candidate_detail.templ` - Suggestion field in timeline
- [ ] `handler/admin.go` - POST /admin/interactions/{id}/suggestion
- [ ] Consultant sees suggestion, marks as read
- [ ] Notification badge for unread suggestions
- [ ] Test: Add suggestion, mark as read

---

## Feature 27: Consultant Assignment

Manual reassignment of candidates.

- [ ] `model/candidate.go` - Assign, GetAssignmentStats
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/assign
- [ ] Supervisor/admin can reassign candidates
- [ ] Show consultant workload in dropdown
- [ ] Log assignment change in timeline
- [ ] Test: Reassignment, workload display

---

## Feature 28: Referral Claim Verification

Link referral claims to referrers.

- [x] `templates/admin/referral_claims.templ` - List unverified claims (UI mockup done)
- [ ] `handler/admin.go` - GET /admin/referral-claims, POST link (wire to real data)
- [ ] Show candidates with source_detail (referral claim) but no referrer_id
- [ ] Search existing referrers by name/institution
- [ ] Actions: link to existing referrer, create new referrer then link, mark as invalid
- [ ] Test: Claim verification flow

---

# Phase 5: Commitment & Enrollment

Convert candidates to students.

---

## Feature 29: Commitment & Tuition Billing

Generate billing when candidate commits.

**Migrations:** 021_create_billings, 022_create_payments (pending)

- [ ] `model/candidate.go` - Commit (change status)
- [ ] `model/billing.go` - CreateTuitionBilling, CreateDormitoryBilling
- [ ] `templates/admin/commitment_form.templ` - Commitment modal
- [ ] `handler/admin.go` - POST /admin/candidates/{id}/commit
- [ ] Select: tuition installments (1x), dormitory (1x, 2x, or 10x)
- [ ] Generate billing records with due dates
- [ ] Change status: prospecting → committed
- [ ] Test: Commit with various installment options

---

## Feature 30: Payment Tracking

Track and verify installment payments.

**Migrations:** 022_create_payments (pending)

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

## Feature 31: Document Review

Admin reviews uploaded documents.

**Migrations:** 023_create_documents (pending)

- [ ] `model/document.go` - UpdateStatus
- [ ] `templates/admin/document_review.templ` - Review modal
- [ ] `handler/admin.go` - GET /admin/candidates/{id}/documents, POST approve/reject
- [ ] View document, approve or reject with reason
- [ ] Notify candidate of rejection (for re-upload)
- [ ] Test: Review flow, rejection

---

## Feature 32: Enrollment

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

## Feature 33: Lost Candidate

Mark candidate as lost.

**Migrations:** 015_create_lost_reasons ✅

- [x] `model/lost_reason.go` - ListActive (already implemented)
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

## Feature 34: Commission Tracking

Auto-create and track commissions.

**Migrations:** 024_create_commission_ledger (pending)

- [x] `templates/admin/commissions.templ` - Commission list (UI mockup done)
- [ ] `model/commission.go` - Create, ListByReferrer, ListPending
- [ ] Auto-create commission when referred candidate enrolls
- [ ] Amount from referrer.commission_per_enrollment
- [ ] Status: pending → approved → paid
- [ ] Test: Auto-creation on enrollment

---

## Feature 35: Commission Payout

Approve and pay commissions.

- [x] `templates/admin/commissions.templ` - Commission list (UI mockup done)
- [ ] `handler/admin.go` - GET /admin/commissions, POST approve, POST mark-paid (wire to real data)
- [ ] Filter by: referrer, status, date range
- [ ] Batch approve, batch mark as paid
- [ ] Export for bank transfer
- [ ] Test: Approval flow, batch operations

---

# Phase 7: Reporting & Analytics

Insights for decision making.

---

## Feature 36: Dashboard - Consultant ✅

Consultant's daily view.

- [x] `templates/admin/consultant_dashboard.templ` - Consultant personal dashboard (UI mockup done)
- [x] `handler/admin_mockups.go` - GET /admin/my-dashboard (mockup handler)
- [x] Show: my candidates by status, overdue followups, today's tasks
- [x] Quick access to pending followups
- [x] Supervisor suggestions section
- [ ] `model/stats.go` - GetConsultantStats
- [ ] Wire to real data
- [ ] Test: Dashboard data accuracy

---

## Feature 37: Dashboard - Supervisor

Supervisor's team view.

- [x] `templates/admin/dashboard.templ` - Admin/Supervisor dashboard (UI mockup done)
- [ ] `model/stats.go` - GetTeamStats, GetFunnelStats
- [ ] Show: team funnel, consultant leaderboard, stuck candidates (> 7 days no interaction)
- [ ] Common obstacles this period
- [ ] Test: Dashboard data accuracy

---

## Feature 38: Reports - Funnel

Conversion funnel analysis.

- [x] `templates/admin/reports_funnel.templ` - Funnel report (UI mockup done)
- [ ] `model/stats.go` - GetFunnelByDateRange, GetFunnelByProdi
- [ ] `handler/admin.go` - GET /admin/reports/funnel (wire to real data)
- [ ] Filter by: date range, prodi, campaign
- [ ] Show: registered → prospecting → committed → enrolled, with conversion rates
- [ ] Test: Report accuracy

---

## Feature 39: Reports - Consultant Performance ✅

Individual performance metrics.

- [x] `templates/admin/consultant_report.templ` - Consultant performance report (UI mockup done)
- [ ] `model/stats.go` - GetConsultantPerformance
- [ ] `handler/admin.go` - GET /admin/reports/consultants (wire to real data)
- [ ] Metrics: candidates handled, success rate, avg days to commit, interaction frequency
- [ ] Ranking by success rate
- [ ] Test: Metrics calculation

---

## Feature 40: Reports - Campaign ROI

Campaign effectiveness.

- [x] `templates/admin/reports_campaigns.templ` - Campaign report (UI mockup done)
- [ ] `model/stats.go` - GetCampaignStats
- [ ] `handler/admin.go` - GET /admin/reports/campaigns (wire to real data)
- [ ] Show: leads, commits, enrollments, conversion rate per campaign
- [ ] Cost per enrollment (if cost data available)
- [ ] Test: Report accuracy

---

## Feature 41: Reports - Referrer Leaderboard

Referrer performance.

- [x] `templates/admin/reports_referrers.templ` - Referrer report (UI mockup done)
- [ ] `model/stats.go` - GetReferrerStats
- [ ] `handler/admin.go` - GET /admin/reports/referrers (wire to real data)
- [ ] Show: referrals, enrollments, conversion rate, commission earned/paid
- [ ] Test: Report accuracy

---

## Feature 42: CSV Export

Export data for external analysis.

- [ ] `handler/admin.go` - GET /admin/export/candidates, /admin/export/interactions
- [ ] Filter by: date range, status, consultant, campaign
- [ ] Include all relevant fields
- [ ] Test: Export with filters

---

# Phase 8: Notifications

Communication automation.

---

## Feature 43: WhatsApp Notifications

Send notifications at key events.

**Migrations:** 025_create_notification_logs (pending)

- [x] `integration/whatsapp.go` - SendOTP via WhatsApp API (implemented for OTP)
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
