# Features - STMIK Tazkia Admission System

Comprehensive feature documentation with implementation status.

## Table of Contents

- [Sales Funnel](#sales-funnel)
- [Implementation Progress](#implementation-progress)
- [Feature Details](#feature-details)
- [Reference Data](#reference-data)

---

## Sales Funnel

### Status Flow

```
PROSPECT                           APPLICATION
────────                           ───────────
   │                                    │
   ▼                                    ▼
┌─────┐    ┌───────────┐    ┌─────────────────┐    ┌──────────┐    ┌──────────┐
│ New │───▶│ Contacted │───▶│ Applicant       │───▶│ Approved │───▶│ Enrolled │
└─────┘    └───────────┘    │ (uploads docs)  │    │ (payment)│    │ (paid)   │
   │            │           └─────────────────┘    └──────────┘    └──────────┘
   │            │                   │                   │
   └────────────┴───────────────────┴───────────────────┘
                                    ▼
                              Cancelled
                         (by registrant)
```

### Candidate Status (Simplified)

| Status | Description | Next Status |
|--------|-------------|-------------|
| `registered` | Just submitted form | `prospecting` |
| `prospecting` | Consultant following up | `committed`, `lost` |
| `committed` | Paid registration, uploading docs | `enrolled`, `lost` |
| `enrolled` | Fully paid, NIM assigned | - |
| `lost` | Dropped out | - |

---

## Implementation Progress

**Overall Progress:** 70% (29/42 features complete)

| Phase | Features | Status | Progress |
|-------|----------|--------|----------|
| **Phase 0: UI Mockup** | Feature 3 | ✅ Complete | 1/1 (100%) |
| **Phase 1: Admin Foundation** | Features 4-8 | ✅ Complete | 5/5 (100%) |
| **Phase 2: Configuration** | Features 9-14 | ✅ Complete | 6/6 (100%) |
| **Phase 3: Registration & Portal** | Features 15-22 | 🔄 In Progress | 7/8 (87%) |
| **Phase 4: CRM Operations** | Features 23-28 | ✅ Complete | 6/6 (100%) |
| **Phase 5: Commitment & Enrollment** | Features 29-33 | 🔄 In Progress | 4/5 (80%) |
| **Phase 6: Commissions** | Features 34-35 | ✅ Complete | 2/2 (100%) |
| **Phase 7: Reporting** | Features 36-42 | 🔄 In Progress | 1/7 (14%) |
| **Phase 8: Notifications** | Feature 43 | ⏳ Pending | 0/1 (0%) |

---

## Feature Details

### Phase 0: UI Mockup

#### Feature 3: UI Mockup (No Database) ✅

**Status:** Complete | **Priority:** High | **Phase:** 0

Clickable prototype untuk validasi dengan stakeholder sebelum implementasi.

**Completed:**
- [x] Admin pages (login, dashboard, candidates, interactions, settings)
- [x] Portal pages (register, login, dashboard, documents, payments, announcements)
- [x] Public registration form (multi-step)
- [x] Navigation and responsive design
- [x] Demo data for all pages

---

### Phase 1: Admin Foundation

#### Feature 4: Staff Login (Google OAuth) ✅

**Status:** Complete | **Priority:** Critical | **Phase:** 1

All staff login with domain-restricted Google OAuth.

**Completed:**
- [x] Google Cloud project setup
- [x] OAuth consent screen configuration
- [x] Model: User CRUD, FindByEmail, FindByGoogleID
- [x] Auth: Google OAuth flow
- [x] Handler: Login routes and callbacks
- [x] Domain check (STAFF_EMAIL_DOMAIN)
- [x] Auto-create user with role=consultant
- [x] Middleware: RequireAuth, RequireRole
- [x] HttpOnly JWT cookies

**Pending:**
- [ ] Test: Login with valid/invalid domain

---

#### Feature 5: Staff Management ✅

**Status:** Complete | **Priority:** High | **Phase:** 1

Admin manages staff accounts (admin, supervisor, consultant).

**Completed:**
- [x] Model: User CRUD, List, UpdateRole, ToggleActive, SetSupervisor
- [x] Templates: User list with HTMX inline editing
- [x] Handler: User management routes
- [x] Roles: admin, supervisor, consultant
- [x] Assign supervisor to consultants
- [x] Toggle active status
- [x] E2E tests

---

#### Feature 6: Settings - Prodi Management ✅

**Status:** Complete | **Priority:** High | **Phase:** 1

Admin configures available programs.

**Completed:**
- [x] Model: Prodi CRUD, ListActive
- [x] Templates: Prodi list with HTMX forms
- [x] Handler: CRUD routes
- [x] Fields: name, code, degree (S1/D3), is_active
- [x] E2E tests

---

#### Feature 7: Settings - Fee Structure ✅

**Status:** Complete | **Priority:** High | **Phase:** 1

Admin configures fees per prodi and academic year.

**Completed:**
- [x] Model: FeeType List, FeeStructure CRUD
- [x] Templates: Fee matrix with HTMX forms
- [x] Handler: CRUD routes
- [x] Seed: registration, tuition, dormitory
- [x] E2E tests

**Pending:**
- [ ] Installment options per fee type

---

#### Feature 8: Settings - Categories & Obstacles ✅

**Status:** Complete | **Priority:** High | **Phase:** 1

Supervisor manages interaction categories and obstacles.

**Completed:**
- [x] Model: Category and Obstacle CRUD
- [x] Templates: Lists with HTMX forms
- [x] Handler: CRUD routes
- [x] Seed: default categories and obstacles
- [x] E2E tests

---

### Phase 2: Configuration

#### Feature 9: Campaign Management ✅

**Status:** Complete | **Priority:** High | **Phase:** 2

Admin manages campaigns with promo pricing.

**Completed:**
- [x] Model: Campaign CRUD, FindActive, FindByCode
- [x] Templates: Campaign settings with HTMX CRUD
- [x] Handler: CRUD routes
- [x] Fee override support
- [x] E2E tests

**Pending:**
- [ ] UTM tracking code generation

---

#### Feature 10: Reward Configuration ✅

**Status:** Complete | **Priority:** High | **Phase:** 2

Configure default rewards by referrer type and MGM.

**Completed:**
- [x] Model: RewardConfig, MgmRewardConfig CRUD
- [x] Templates: Reward config list with HTMX
- [x] Handler: CRUD routes
- [x] Seed: default rewards by type
- [x] E2E tests

---

#### Feature 11: Referrer Management ✅

**Status:** Complete | **Priority:** High | **Phase:** 2

Admin manages referrers for commission tracking.

**Completed:**
- [x] Model: Referrer CRUD, GenerateCode, FindByCode
- [x] Templates: Referrer list with HTMX forms
- [x] Handler: CRUD routes
- [x] Commission override support
- [x] E2E tests

**Pending:**
- [ ] Optional referral code generation for partners

---

#### Feature 12: Settings - Assignment Algorithm ✅

**Status:** Complete | **Priority:** Medium | **Phase:** 2

Configure consultant assignment algorithm.

**Completed:**
- [x] Model: Algorithm List, SetActive, FindActive
- [x] Templates: Algorithm selection with HTMX
- [x] Handler: Selection routes
- [x] Seed: round_robin, load_balanced, performance_weighted, activity_based
- [x] E2E tests

---

#### Feature 13: Settings - Document Types ✅

**Status:** Complete | **Priority:** High | **Phase:** 2

Configure required documents.

**Completed:**
- [x] Model: DocumentType CRUD, ListActive
- [x] Templates: Document type list with HTMX
- [x] Handler: CRUD routes
- [x] Seed: KTP, Photo, Ijazah, Transcript
- [x] E2E tests

---

#### Feature 14: Settings - Lost Reasons ✅

**Status:** Complete | **Priority:** Medium | **Phase:** 2

Configure reasons for lost candidates.

**Completed:**
- [x] Model: LostReason CRUD, ListActive
- [x] Templates: Lost reason list with HTMX
- [x] Handler: CRUD routes
- [x] Seed: No response, Chose competitor, Financial, etc.
- [x] E2E tests

---

### Phase 3: Public Registration & Candidate Portal

#### Feature 15: Candidate Registration ✅

**Status:** Complete | **Priority:** Critical | **Phase:** 3

Registration with password. Email/phone verification optional.

**Completed:**
- [x] Model: Candidate CRUD, verification methods
- [x] Integration: Resend and WhatsApp (optional)
- [x] Templates: 4-step registration form
- [x] Handler: Registration routes
- [x] Auto-assign to consultant
- [x] Referral and campaign tracking
- [x] E2E tests

---

#### Feature 16: Candidate Login ✅

**Status:** Complete | **Priority:** Critical | **Phase:** 3

Candidate authenticates with email or phone + password.

**Completed:**
- [x] Model: Authenticate by email or phone
- [x] Templates: Login form
- [x] Handler: Login/logout routes
- [x] Cookie-based session
- [x] E2E tests

---

#### Feature 17: Candidate Portal - Dashboard ✅

**Status:** Complete | **Priority:** High | **Phase:** 3

Overview of candidate status and actions.

**Completed:**
- [x] Templates: Dashboard with checklist
- [x] Handler: Dashboard route
- [x] Middleware: RequireCandidateAuth
- [x] Status badge and checklist
- [x] Recent announcements
- [x] E2E tests

---

#### Feature 18: Candidate Portal - Documents ✅

**Status:** Complete | **Priority:** High | **Phase:** 3

Candidate uploads and tracks documents.

**Completed:**
- [x] Storage: Local filesystem implementation
- [x] Model: Document CRUD
- [x] Templates: Upload form with status
- [x] Handler: Upload and list routes
- [x] Re-upload rejected documents
- [x] Admin review interface
- [x] E2E tests

---

#### Feature 19: Candidate Portal - Payments ✅

**Status:** Complete | **Priority:** High | **Phase:** 3

Candidate views billing and uploads payment proof.

**Completed:**
- [x] Model: Billing and Payment CRUD
- [x] Templates: Billing list with upload modal
- [x] Handler: Payment routes
- [x] Upload payment proof
- [x] E2E tests

**Pending:**
- [ ] Registration fee generation on registration
- [ ] Tuition/dormitory generation on commitment
- [ ] Admin payment review

---

#### Feature 20: Candidate Portal - Announcements ✅

**Status:** Complete | **Priority:** Medium | **Phase:** 3

Candidate receives targeted announcements.

**Completed:**
- [x] Model: Announcement ListForCandidate, MarkRead
- [x] Templates: Announcement list with HTMX
- [x] Handler: Announcement routes
- [x] Filter by status and prodi
- [x] Mark as read
- [x] E2E tests

---

#### Feature 21: Announcement Management (Admin) ✅

**Status:** Complete | **Priority:** Medium | **Phase:** 3

Admin creates and targets announcements.

**Completed:**
- [x] Model: Announcement CRUD, Publish
- [x] Templates: Announcement list with HTMX
- [x] Handler: CRUD routes
- [x] Targeting by status and prodi
- [x] Publish/unpublish toggle
- [x] E2E tests

---

#### Feature 22: Member Get Member Referral 🔄

**Status:** In Progress | **Priority:** Medium | **Phase:** 3

Enrolled students refer new candidates.

**Completed:**
- [x] Templates: Referral dashboard (UI mockup)
- [x] Registration form: `?ref=CODE` support

**Pending:**
- [ ] Model: GenerateReferralCode
- [ ] Handler: Wire to real data
- [ ] Referral tracking
- [ ] Test: Code generation, referral tracking

---

### Phase 4: CRM Operations

#### Feature 23: Candidate List & Filters ✅

**Status:** Complete | **Priority:** High | **Phase:** 4

Admin/consultant views candidates.

**Completed:**
- [x] Model: ListCandidates with filters
- [x] Templates: Table with HTMX filters
- [x] Handler: Candidate list route
- [x] Role-based visibility
- [x] E2E tests

---

#### Feature 24: Candidate Detail & Timeline ✅

**Status:** Complete | **Priority:** High | **Phase:** 4

View candidate info and history.

**Completed:**
- [x] Templates: Detail view with timeline
- [x] Handler: Candidate detail route
- [x] Interaction timeline
- [x] Quick actions
- [x] E2E tests

---

#### Feature 25: Interaction Logging ✅

**Status:** Complete | **Priority:** High | **Phase:** 4

Consultants log each contact.

**Completed:**
- [x] Model: Interaction CRUD
- [x] Templates: Interaction form
- [x] Handler: Interaction routes
- [x] Channels: call, whatsapp, email, campus_visit, home_visit
- [x] E2E tests

---

#### Feature 26: Supervisor Suggestions ✅

**Status:** Complete | **Priority:** High | **Phase:** 4

Supervisor reviews and provides guidance.

**Completed:**
- [x] Templates: Suggestion field in timeline
- [x] Handler: Suggestion routes
- [x] Consultant sees suggestion
- [x] Notification badge for unread
- [x] E2E tests

---

#### Feature 27: Consultant Assignment ✅

**Status:** Complete | **Priority:** High | **Phase:** 4

Manual reassignment of candidates.

**Completed:**
- [x] Model: ReassignCandidate, ListConsultantsWithWorkload
- [x] Handler: Reassignment routes
- [x] HTMX modal
- [x] Role-based access control
- [x] E2E tests

---

#### Feature 28: Referral Claim Verification ✅

**Status:** Complete | **Priority:** Medium | **Phase:** 4

Link referral claims to referrers.

**Completed:**
- [x] Templates: Unverified claims list
- [x] Handler: Verification routes
- [x] Search referrers
- [x] Link to existing referrer
- [x] E2E tests

---

### Phase 5: Commitment & Enrollment

#### Feature 29: Commitment & Tuition Billing ✅

**Status:** Complete | **Priority:** High | **Phase:** 5

Generate billing when candidate commits.

**Completed:**
- [x] Model: Commit, CreateTuitionBilling, CreateDormitoryBilling
- [x] Templates: Commitment modal
- [x] Handler: Commitment route
- [x] Generate billing records
- [x] Status change: prospecting → committed

---

#### Feature 30: Payment Tracking ⏳

**Status:** Pending | **Priority:** High | **Phase:** 5

Track and verify installment payments.

**Pending:**
- [ ] Model: Payment verification methods
- [ ] Templates: Payment list and verification
- [ ] Handler: Payment verification routes
- [ ] Overdue highlighting
- [ ] Test: Payment lifecycle

---

#### Feature 31: Document Review ✅

**Status:** Complete | **Priority:** High | **Phase:** 5

Admin reviews uploaded documents.

**Completed:**
- [x] Model: ApproveDocument, RejectDocument
- [x] Templates: Review list with modals
- [x] Handler: Review routes
- [x] Email notifications
- [x] E2E tests

---

#### Feature 32: Enrollment ✅

**Status:** Complete | **Priority:** Critical | **Phase:** 5

Mark candidate as enrolled.

**Completed:**
- [x] Model: Enroll, GenerateNIM, ValidateEnrollment
- [x] Templates: Enrollment modal with checklist
- [x] Handler: Enrollment route
- [x] Validation: fees paid, documents approved
- [x] Generate NIM and referral_code
- [x] Trigger commission creation

---

#### Feature 33: Lost Candidate ✅

**Status:** Complete | **Priority:** Medium | **Phase:** 5

Mark candidate as lost.

**Completed:**
- [x] Model: MarkCandidateLost
- [x] Templates: Lost modal with reason
- [x] Handler: Lost route
- [x] Record lost_at and reason
- [x] E2E tests

---

### Phase 6: Commissions

#### Feature 34: Commission Tracking ✅

**Status:** Complete | **Priority:** High | **Phase:** 6

Auto-create and track commissions.

**Completed:**
- [x] Model: Commission CRUD, GetStats
- [x] Templates: Commission list
- [x] Handler: Commission routes
- [x] Auto-create on status change
- [x] Status: pending → approved → paid

---

#### Feature 35: Commission Payout ✅

**Status:** Complete | **Priority:** High | **Phase:** 6

Approve and pay commissions.

**Completed:**
- [x] Templates: Commission list with actions
- [x] Handler: Approve, mark paid routes
- [x] Batch approve and mark paid
- [x] CSV export for bank transfer
- [x] E2E tests

---

### Phase 7: Reporting & Analytics

#### Feature 36: Dashboard - Consultant ✅

**Status:** Complete | **Priority:** High | **Phase:** 7

Consultant's daily view.

**Completed:**
- [x] Templates: Consultant dashboard
- [x] Handler: Dashboard route with real data
- [x] Model: GetConsultantStats, GetOverdueCandidates
- [x] Quick access to pending followups
- [x] Supervisor suggestions section

---

#### Feature 37: Dashboard - Supervisor ⏳

**Status:** Pending | **Priority:** Medium | **Phase:** 7

Supervisor's team view.

**Pending:**
- [ ] Model: GetTeamStats, GetFunnelStats
- [ ] Team funnel visualization
- [ ] Consultant leaderboard
- [ ] Stuck candidates detection
- [ ] Test: Dashboard data accuracy

---

#### Feature 38: Reports - Funnel ⏳

**Status:** Pending | **Priority:** Medium | **Phase:** 7

Conversion funnel analysis.

**Pending:**
- [ ] Model: GetFunnelByDateRange, GetFunnelByProdi
- [ ] Templates: Wire to real data
- [ ] Filter by date, prodi, campaign
- [ ] Conversion rates
- [ ] Test: Report accuracy

---

#### Feature 39: Reports - Consultant Performance ⏳

**Status:** Pending | **Priority:** Medium | **Phase:** 7

Individual performance metrics.

**Pending:**
- [ ] Model: GetConsultantPerformance
- [ ] Templates: Wire to real data
- [ ] Metrics: success rate, avg days to commit
- [ ] Ranking
- [ ] Test: Metrics calculation

---

#### Feature 40: Reports - Campaign ROI ⏳

**Status:** Pending | **Priority:** Medium | **Phase:** 7

Campaign effectiveness.

**Pending:**
- [ ] Model: GetCampaignStats
- [ ] Templates: Wire to real data
- [ ] Conversion rate per campaign
- [ ] Cost per enrollment
- [ ] Test: Report accuracy

---

#### Feature 41: Reports - Referrer Leaderboard ⏳

**Status:** Pending | **Priority:** Low | **Phase:** 7

Referrer performance.

**Pending:**
- [ ] Model: GetReferrerStats
- [ ] Templates: Wire to real data
- [ ] Referrals, enrollments, conversion rate
- [ ] Test: Report accuracy

---

#### Feature 42: CSV Export ⏳

**Status:** Pending | **Priority:** Low | **Phase:** 7

Export data for external analysis.

**Pending:**
- [ ] Handler: Export routes
- [ ] Filter by date, status, consultant
- [ ] Test: Export with filters

---

### Phase 8: Notifications

#### Feature 43: WhatsApp Notifications ⏳

**Status:** Pending | **Priority:** Medium | **Phase:** 8

Send notifications at key events.

**Completed:**
- [x] Integration: WhatsApp API client (for OTP)

**Pending:**
- [ ] Model: Notification log
- [ ] Templates: Event-based notifications
- [ ] Manual send from candidate detail
- [ ] Log all messages
- [ ] Test: Send with mock API

---

## Reference Data

### Cancel Reasons

| Code | Label | Applies To |
|------|-------|------------|
| `no_response` | Tidak merespon | prospect |
| `chose_other` | Memilih kampus lain | both |
| `financial` | Kendala biaya | both |
| `changed_mind` | Berubah pikiran | both |
| `age_requirement` | Tidak memenuhi syarat usia | prospect |
| `invalid_document` | Dokumen tidak valid | application |
| `other` | Lainnya | both |

### Document Checklists

**KTP (Kartu Tanda Penduduk):**
- Nama jelas terbaca
- NIK lengkap (16 digit)
- Foto tidak buram
- Masih berlaku

**Ijazah (Certificate/Diploma):**
- Nama sesuai KTP
- Tahun lulus terlihat
- Stempel sekolah ada
- Tanda tangan kepala sekolah

### WhatsApp Templates

| Template | Trigger | Variables |
|----------|---------|-----------|
| `welcome` | New prospect | `name` |
| `followup_3d` | 3 days no activity | `name` |
| `document_reminder` | 7 days incomplete | `name`, `missing_docs` |
| `revision_request` | Doc needs revision | `name`, `doc_type`, `remarks` |
| `approved` | Application approved | `name`, `program`, `va_number` |
| `enrolled` | Payment confirmed | `name`, `program` |

---

## Success Criteria

### Admin/Staff ✅
- [x] Staff can login with Google OAuth
- [x] Admin can configure prodis/fees/campaigns
- [x] Consultants log interactions
- [x] Supervisors provide suggestions
- [x] Commitment generates billing
- [ ] Payment tracking with verification
- [x] Document review with approve/reject
- [x] Enrollment validates requirements, generates NIM
- [x] All admin interactions use HTMX

### Candidate Portal ✅
- [x] Candidate registers with email/phone
- [x] Candidate logs in with password
- [x] Candidate views dashboard
- [x] Candidate uploads documents
- [x] Candidate views billing and uploads proof
- [x] Candidate receives announcements
- [ ] Enrolled candidate can refer (MGM pending)

### Business 🔄
- [x] Registration fee waived during promo
- [x] Candidates auto-assigned to consultants
- [x] Documents can be deferred
- [x] Referrer commissions tracked
- [ ] Campaign ROI reports
- [x] Mobile responsive
