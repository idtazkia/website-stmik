# Keuangan & Pembayaran

Bagian ini menjelaskan pengelolaan biaya pendaftaran, tagihan, dan verifikasi pembayaran.

## Struktur Biaya

Biaya pendaftaran ditentukan berdasarkan program studi dan jenis biaya.

### Jenis Biaya

| Jenis | Deskripsi |
|-------|-----------|
| Biaya Pendaftaran | Biaya satu kali saat mendaftar |
| Biaya Ujian | Biaya tes masuk |
| SPP Semester 1 | Biaya kuliah semester pertama |
| Biaya Lainnya | Biaya seragam, almamater, dll |

### Konfigurasi

Buka menu **Pengaturan > Struktur Biaya** untuk mengelola daftar biaya per program studi.

![Struktur Biaya](screenshots/admin/settings-fees.png)

Setiap struktur biaya terkait dengan:
- Program studi tertentu
- Jenis biaya (fee type)
- Nominal
- Semester (jika applicable)

---

## Pembuatan Tagihan

Tim finance membuat tagihan untuk kandidat yang sudah committed.

### Akses

Buka menu **Finance > Tagihan** di sidebar admin.

![Daftar Tagihan](screenshots/admin/finance-billings.png)

### Langkah Buat Tagihan

1. Klik **Buat Tagihan**
2. Pilih kandidat
3. Pilih jenis biaya yang akan ditagihkan
4. Nominal otomatis terisi dari struktur biaya
5. Klik **Simpan**

### Status Tagihan

| Status | Deskripsi |
|--------|-----------|
| **Pending** | Tagihan dibuat, menunggu pembayaran |
| **Paid** | Pembayaran sudah diverifikasi |

---

## Upload Bukti Pembayaran (Portal)

Kandidat mengupload bukti transfer pembayaran melalui portal.

### Langkah

1. Login ke portal
2. Buka menu **Pembayaran**
3. Lihat daftar tagihan yang belum dibayar
4. Klik **Upload Bukti** pada tagihan yang akan dibayar
5. Pilih file bukti transfer (gambar atau PDF)
6. Klik **Upload**

![Upload Bukti](screenshots/portal/payments.png)

### Setelah Upload

- Status pembayaran berubah menjadi **Menunggu Verifikasi**
- Tim finance mendapat notifikasi untuk melakukan verifikasi

---

## Verifikasi Pembayaran (Admin)

Tim finance memverifikasi bukti pembayaran yang diupload kandidat.

### Akses

Buka menu **Finance > Pembayaran** atau klik notifikasi pembayaran baru.

![Verifikasi Pembayaran](screenshots/admin/finance-payments.png)

### Langkah Verifikasi

1. Buka daftar pembayaran pending
2. Klik bukti pembayaran untuk melihat preview
3. Cocokkan dengan data tagihan (nominal, tanggal, pengirim)
4. Klik **Approve** jika valid
5. Klik **Reject** jika tidak valid — isi alasan penolakan

### Notifikasi

| Aksi | Notifikasi |
|------|------------|
| Approve | Email konfirmasi pembayaran ke kandidat |
| Reject | Email penolakan dengan alasan ke kandidat |

### Siapa yang Bisa Verifikasi?

| Role | Akses |
|------|-------|
| Finance | Bisa approve dan reject |
| Admin | Bisa approve dan reject |
| Role lain | Tidak bisa |
