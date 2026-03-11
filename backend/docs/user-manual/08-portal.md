# Portal Calon Mahasiswa

Bagian ini menjelaskan fitur-fitur yang tersedia di portal self-service untuk calon mahasiswa.

## Dashboard Portal

Setelah login, calon mahasiswa melihat dashboard yang menampilkan ringkasan status pendaftarannya.

![Dashboard Portal](screenshots/portal/dashboard.png)

### Informasi yang Ditampilkan

| Widget | Deskripsi |
|--------|-----------|
| Status Pendaftaran | Status saat ini (registered, prospecting, committed, enrolled) |
| Dokumen | Jumlah dokumen yang sudah dan belum diupload |
| Pembayaran | Ringkasan tagihan dan status pembayaran |
| Pengumuman | Pengumuman terbaru dari kampus |
| Verifikasi | Status verifikasi email dan telepon |

---

## Upload Dokumen

Calon mahasiswa mengupload dokumen persyaratan melalui portal.

### Akses

Klik menu **Dokumen** di sidebar portal.

![Dokumen Portal](screenshots/portal/documents.png)

### Informasi yang Ditampilkan

Setiap tipe dokumen menampilkan:
- Nama tipe dokumen
- Status: Belum Upload / Pending / Approved / Rejected
- Tanggal upload (jika sudah)
- Alasan penolakan (jika ditolak)

### Upload Ulang

Jika dokumen ditolak:
1. Baca alasan penolakan
2. Siapkan dokumen yang benar
3. Klik **Upload Ulang**
4. Pilih file baru
5. Klik **Upload**

---

## Status Pembayaran

Halaman pembayaran menampilkan daftar tagihan dan status pembayaran.

### Akses

Klik menu **Pembayaran** di sidebar portal.

![Pembayaran Portal](screenshots/portal/payments.png)

### Informasi yang Ditampilkan

| Kolom | Deskripsi |
|-------|-----------|
| Jenis Biaya | Nama biaya (pendaftaran, SPP, dll) |
| Nominal | Jumlah yang harus dibayar |
| Status | Belum Bayar / Menunggu Verifikasi / Lunas / Ditolak |
| Aksi | Upload bukti pembayaran |

### Upload Bukti Pembayaran

1. Klik **Upload Bukti** pada tagihan yang akan dibayar
2. Pilih file bukti transfer
3. Klik **Upload**
4. Status berubah menjadi **Menunggu Verifikasi**

---

## Pengumuman

Halaman pengumuman menampilkan informasi dan pengumuman dari kampus.

### Akses

Klik menu **Pengumuman** di sidebar portal.

![Pengumuman Portal](screenshots/portal/announcements.png)

### Fitur

- Daftar pengumuman yang ditargetkan untuk kandidat
- Pengumuman baru ditandai dengan badge
- Klik pengumuman untuk membaca detail
- Klik **Tandai Sudah Dibaca** setelah membaca

---

## Program Referral

Kandidat yang sudah enrolled dapat ikut program referral untuk mereferensikan calon mahasiswa baru.

### Akses

Klik menu **Referral** di sidebar portal.

![Referral Portal](screenshots/portal/referral.png)

### Informasi yang Ditampilkan

| Info | Deskripsi |
|------|-----------|
| Kode Referral | Kode unik kandidat untuk dibagikan |
| Jumlah Referensi | Berapa orang yang sudah mendaftar dengan kode ini |
| Status Referensi | Status masing-masing orang yang direferensikan |
| Komisi | Besaran komisi yang akan didapat jika berhasil enrolled |
