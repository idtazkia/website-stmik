# API Documentation - STMIK Tazkia Admission System

Complete REST API and web routes documentation.

## Table of Contents

- [Public API](#public-api)
- [Portal Routes (Registrants)](#portal-routes-registrants)
- [Admin Routes (Staff)](#admin-routes-staff)

## Public API

REST endpoints for external integrations (landing page, payment system).

### Endpoints

```
POST /api/prospects              # Create prospect (from landing page)
GET  /api/health                 # Health check
GET  /api/referrers/{code}       # Validate referral code (for landing page)
```

### POST /api/prospects

Create a new prospect from the landing page registration form.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "whatsapp": "08123456789",
  "intake_id": 1,
  "referral_code": "REF123",
  "utm_source": "google",
  "utm_medium": "cpc",
  "utm_campaign": "intake_2025_ganjil",
  "utm_term": "kuliah it jakarta",
  "utm_content": "banner_v1",
  "landing_page": "https://stmik.tazkia.ac.id/promo",
  "device_type": "mobile"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "whatsapp": "08123456789",
  "status": "new",
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Field Descriptions:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Prospect full name |
| `email` | string | No | Email address (optional) |
| `whatsapp` | string | Yes | WhatsApp number (Indonesian format) |
| `intake_id` | integer | Yes | Target intake period ID |
| `referral_code` | string | No | Referrer's unique code |
| `utm_source` | string | No | Traffic source (google, facebook, instagram) |
| `utm_medium` | string | No | Marketing medium (cpc, social, email) |
| `utm_campaign` | string | No | Campaign name |
| `utm_term` | string | No | Paid search keywords |
| `utm_content` | string | No | Ad variation identifier |
| `landing_page` | string | No | Landing page URL |
| `device_type` | string | No | Device type (mobile, desktop) |

### GET /api/health

Health check endpoint for monitoring.

**Response:** `200 OK`
```json
{
  "status": "ok",
  "version": {
    "commit": "abc123",
    "branch": "main",
    "build_time": "2025-01-15T10:00:00Z"
  }
}
```

### GET /api/referrers/{code}

Validate referral code exists and is active.

**Response:** `200 OK`
```json
{
  "valid": true,
  "referrer_name": "Jane Doe",
  "referrer_type": "student"
}
```

**Error Response:** `404 Not Found`
```json
{
  "valid": false,
  "error": "Referral code not found"
}
```

## Portal Routes (Registrants)

Web routes for prospective students to register, login, and complete applications.

### Authentication

```
GET  /portal/login               # Login page
POST /portal/login               # Login submit
GET  /portal/register            # Register page
POST /portal/register            # Register submit
GET  /portal/auth/google         # Google OAuth
GET  /portal/auth/google/callback
POST /portal/logout
```

### Application Management

```
GET  /portal                     # Dashboard (status overview)
GET  /portal/application         # Application form
POST /portal/application         # Create/update application (HTMX)
GET  /portal/documents           # Document upload page
POST /portal/documents           # Upload document (HTMX)
DELETE /portal/documents/{id}    # Remove document (HTMX)
```

### Cancellation

```
GET  /portal/cancel              # Cancel confirmation page
POST /portal/cancel              # Submit cancellation
```

## Admin Routes (Staff)

Web routes for marketing staff to manage prospects, applications, and reports.

### Authentication

```
GET  /admin/login                # Login (Google only)
GET  /admin/auth/google
GET  /admin/auth/google/callback
POST /admin/logout
```

### Dashboard

```
GET  /admin                      # Dashboard overview
```

### Prospect Management

```
GET  /admin/prospects                        # Prospect list (filterable)
GET  /admin/prospects/{id}                   # Prospect detail
POST /admin/prospects/{id}/assign            # Assign to staff (HTMX)
POST /admin/prospects/{id}/status            # Update status (HTMX)
POST /admin/prospects/{id}/whatsapp          # Send WhatsApp (HTMX)
POST /admin/prospects/{id}/cancel            # Mark cancelled (HTMX)
```

### Application Management

```
GET  /admin/applications                     # Application list
GET  /admin/applications/{id}                # Application detail
GET  /admin/applications/{id}/documents/{docId}/review  # Review modal
POST /admin/applications/{id}/documents/{docId}/review  # Submit review (HTMX)
POST /admin/applications/{id}/approve        # Approve application (HTMX)
POST /admin/applications/{id}/cancel         # Mark cancelled (HTMX)
```

### Referrer Management

```
GET  /admin/referrers            # Referrer list
GET  /admin/referrers/{id}       # Referrer detail + stats
POST /admin/referrers            # Create referrer (HTMX)
PUT  /admin/referrers/{id}       # Update referrer (HTMX)
POST /admin/referrers/{id}/toggle        # Toggle active (HTMX)
```

### Campaign Management

```
GET  /admin/campaigns            # Campaign list
GET  /admin/campaigns/{id}       # Campaign detail + stats
POST /admin/campaigns            # Create campaign (HTMX)
PUT  /admin/campaigns/{id}       # Update campaign (HTMX)
```

### Settings

```
GET  /admin/settings             # Settings overview
GET  /admin/settings/intakes     # Intake management
POST /admin/settings/intakes     # Create intake (HTMX)
PUT  /admin/settings/intakes/{id}        # Update intake (HTMX)
GET  /admin/settings/tracks      # Track management
GET  /admin/settings/reasons     # Cancel reasons
GET  /admin/settings/checklists  # Document checklists
GET  /admin/settings/staff       # Staff management
POST /admin/settings/staff/{id}/toggle   # Toggle active (HTMX)
```

### Reports

```
GET  /admin/reports              # Reports page
GET  /admin/reports/funnel       # Funnel data (HTMX)
GET  /admin/reports/sources      # Conversion by source (HTMX)
GET  /admin/reports/campaigns    # Campaign performance (HTMX)
GET  /admin/reports/referrers    # Referrer leaderboard (HTMX)
GET  /admin/reports/export       # CSV export
```

## HTMX Integration

Many admin routes use HTMX for dynamic updates without page reloads. These endpoints return HTML fragments instead of full pages.

**Example HTMX Request:**
```html
<button hx-post="/admin/prospects/123/assign"
        hx-target="#assignment-status"
        hx-swap="innerHTML">
    Assign to Staff
</button>
```

**Response:** HTML fragment
```html
<div id="assignment-status">
    Assigned to: John Doe
</div>
```

## Authentication

### Portal Authentication

- **Email/Password**: Traditional login with bcrypt password hashing
- **Google OAuth**: OIDC flow for quick registration

### Admin Authentication

- **Google OAuth Only**: Staff must use `@tazkia.ac.id` email domain
- Session managed via JWT tokens in HttpOnly cookies

## Error Handling

All endpoints return consistent error responses:

**4xx Client Errors:**
```json
{
  "error": "Validation failed",
  "details": {
    "email": "Invalid email format"
  }
}
```

**5xx Server Errors:**
```json
{
  "error": "Internal server error",
  "request_id": "abc-123"
}
```

## Rate Limiting

- Public API: 100 requests/minute per IP
- Portal routes: 60 requests/minute per user
- Admin routes: 120 requests/minute per user

Rate limit headers included in all responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642234567
```
