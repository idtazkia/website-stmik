# Login & Dashboard

Bagian ini menjelaskan cara masuk ke sistem dan navigasi dashboard untuk setiap role.

## Login Admin (Google OIDC)

Staff admin, supervisor, konsultan, finance, dan academic login menggunakan akun Google institusi.

### Langkah Login

1. Buka halaman `/admin/login`
2. Klik **Login dengan Google**
3. Pilih akun Google dengan domain `@tazkia.ac.id`
4. Sistem memverifikasi domain email — hanya `@tazkia.ac.id` yang diizinkan
5. Setelah berhasil, redirect ke dashboard sesuai role

![Login Admin](screenshots/admin/login.png)

### Validasi Domain

| Domain | Akses |
|--------|-------|
| `@tazkia.ac.id` | Diizinkan — role ditentukan oleh admin |
| Domain lain | Ditolak — hanya domain institusi |

> **Catatan**: Akun Google baru yang belum terdaftar di sistem akan ditolak. Admin harus menambahkan user terlebih dahulu di menu Pengaturan > User.

---

## Login Calon Mahasiswa

Calon mahasiswa login menggunakan email/password yang dibuat saat registrasi.

### Langkah Login

1. Buka halaman `/login`
2. Masukkan email dan password
3. Klik **Login**
4. Redirect ke portal calon mahasiswa

![Login Portal](screenshots/portal/login.png)

### Lupa Password

Saat ini belum tersedia fitur reset password mandiri. Hubungi admin untuk reset password.

---

## Dashboard Admin

Dashboard admin menampilkan ringkasan data penerimaan mahasiswa baru.

![Dashboard Admin](screenshots/admin/dashboard.png)

### Informasi yang Ditampilkan

| Widget | Deskripsi |
|--------|-----------|
| Total Kandidat | Jumlah seluruh calon mahasiswa terdaftar |
| Status Funnel | Distribusi kandidat per status (registered, prospecting, committed, enrolled, lost) |
| Kandidat Baru | Kandidat yang baru mendaftar hari ini |
| Dokumen Pending | Dokumen yang menunggu review |
| Pembayaran Pending | Bukti pembayaran yang menunggu verifikasi |

---

## Dashboard Konsultan

Dashboard konsultan menampilkan data kandidat yang di-assign kepadanya.

![Dashboard Konsultan](screenshots/admin/my-dashboard.png)

### Informasi yang Ditampilkan

| Widget | Deskripsi |
|--------|-----------|
| Kandidat Saya | Jumlah kandidat yang di-assign |
| Follow-Up Hari Ini | Kandidat yang perlu di-follow-up hari ini |
| Interaksi Terakhir | Riwayat interaksi terbaru |
| Status Pipeline | Distribusi kandidat per status |

### Aksi Cepat

- Klik kandidat untuk melihat detail
- Klik **Catat Interaksi** untuk mencatat komunikasi baru
- Filter berdasarkan status atau tanggal follow-up

---

## Dashboard Supervisor

Dashboard supervisor menampilkan ringkasan performa seluruh konsultan di bawah supervisinya.

![Dashboard Supervisor](screenshots/admin/supervisor-dashboard.png)

### Informasi yang Ditampilkan

| Widget | Deskripsi |
|--------|-----------|
| Overview Tim | Total kandidat dan distribusi status per konsultan |
| Konsultan Aktif | Daftar konsultan dengan jumlah kandidat masing-masing |
| Follow-Up Terlewat | Kandidat yang melewati tanggal follow-up tanpa interaksi |
| Conversion Rate | Rasio enrolled/total per konsultan |
