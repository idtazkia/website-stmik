# Keamanan & Enkripsi

Bagian ini menjelaskan fitur keamanan yang diterapkan dalam sistem.

## Enkripsi Data Sensitif

Sistem mengenkripsi data personal (PII) menggunakan AES-256-GCM.

### Data yang Dienkripsi

| Data | Metode | Alasan |
|------|--------|--------|
| Email kandidat | Deterministic encryption | Agar bisa dicari (equality search) |
| Nomor telepon | Deterministic encryption | Agar bisa dicari |
| Nama user | AES-256-GCM | Proteksi data personal |
| Google ID | AES-256-GCM | Proteksi credential |

### Deterministic vs Randomized Encryption

| Metode | Kelebihan | Kekurangan |
|--------|-----------|------------|
| **Deterministic** | Bisa equality search di database | Pola data bisa terlihat |
| **Randomized (GCM)** | Lebih aman, tidak ada pola | Tidak bisa dicari langsung |

### Implikasi

- Data terenkripsi di database — aman meskipun database bocor
- Query langsung ke database tidak menampilkan data asli
- Backup database aman karena data terenkripsi
- Encryption key (`ENCRYPTION_KEY`) harus dijaga kerahasiaannya

---

## Autentikasi & Otorisasi

### Autentikasi Staff (Google OIDC)

Staff login melalui Google OAuth:
1. Redirect ke Google consent screen
2. Google memverifikasi identitas
3. Sistem menerima ID token dari Google
4. Sistem membuat JWT session token sendiri
5. JWT disimpan di HttpOnly cookie (7 hari)

### Autentikasi Kandidat (Email/Password)

Kandidat login dengan email dan password:
1. Submit form login
2. Sistem memverifikasi password dengan bcrypt
3. Sistem membuat JWT session token
4. JWT disimpan di HttpOnly cookie (7 hari)

### JWT Token

| Property | Nilai |
|----------|-------|
| Algoritma | HS256 |
| Masa berlaku | 168 jam (7 hari, konfigurasi via `JWT_EXPIRATION_HOURS`) |
| Penyimpanan | HttpOnly cookie |
| Isi (staff) | user ID, email, nama, role |
| Isi (kandidat) | candidate ID, email, flag isCandidate |

---

## CSRF Protection

Sistem menggunakan `http.CrossOriginProtection` dari Go 1.25 standard library.

### Cara Kerja

Proteksi berbasis header, bukan token:
- Memeriksa header `Sec-Fetch-Site` dan `Origin`
- Request cross-origin yang memodifikasi data (POST, PUT, DELETE) ditolak
- Tidak memerlukan CSRF token di form

### Konfigurasi

CSRF protection otomatis aktif untuk semua route yang memodifikasi data. Tidak perlu konfigurasi tambahan.

### Implikasi

- Form submit dari domain lain ditolak
- AJAX request dari domain lain ditolak
- Login via Google OAuth tetap berjalan (menggunakan redirect, bukan POST cross-origin)

---

## Role & Permission

### Matriks Permission

| Fitur | Admin | Supervisor | Consultant | Finance | Academic |
|-------|-------|------------|------------|---------|----------|
| Dashboard | Semua | Tim sendiri | Kandidat sendiri | Finance | Academic |
| Kandidat - Lihat | Semua | Tim sendiri | Assigned saja | Tidak | Tidak |
| Kandidat - Interaksi | Ya | Ya | Ya | Tidak | Tidak |
| Kandidat - Reassign | Ya | Tim sendiri | Tidak | Tidak | Tidak |
| Kandidat - Commitment | Ya | Ya | Ya | Tidak | Tidak |
| Kandidat - Enrollment | Ya | Ya | Tidak | Tidak | Ya |
| Kandidat - Mark Lost | Ya | Ya | Ya | Tidak | Tidak |
| Dokumen - Review | Ya | Ya | Tidak | Tidak | Tidak |
| Tagihan - Buat | Ya | Tidak | Tidak | Ya | Tidak |
| Pembayaran - Verifikasi | Ya | Tidak | Tidak | Ya | Tidak |
| Komisi - Approve | Ya | Tidak | Tidak | Ya | Tidak |
| Pengumuman | Ya | Tidak | Tidak | Tidak | Tidak |
| Pengaturan | Ya | Tidak | Tidak | Tidak | Tidak |
| Laporan | Ya | Ya | Tidak | Ya | Ya |
| User Management | Ya | Tidak | Tidak | Tidak | Tidak |

### Cookie Security

| Setting | Nilai | Alasan |
|---------|-------|--------|
| HttpOnly | true | Mencegah akses JavaScript (XSS) |
| Secure | true (production) | Hanya dikirim via HTTPS |
| SameSite | Lax | Mengizinkan redirect OAuth, blokir cross-site POST |
| Path | / | Berlaku untuk seluruh path |
| MaxAge | 7 hari | Sesuai JWT expiration |
