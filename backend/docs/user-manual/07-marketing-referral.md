# Marketing & Referral

Bagian ini menjelaskan fitur pemasaran, pelacakan kampanye, dan program referral.

## Manajemen Kampanye

Kampanye (campaign) digunakan untuk melacak sumber kandidat dari berbagai channel pemasaran.

### Akses

Buka menu **Kampanye** di sidebar admin.

![Daftar Kampanye](screenshots/admin/campaigns-list.png)

### Cara Kerja

1. Setiap kampanye memiliki kode unik
2. Kode kampanye dipasang sebagai UTM parameter di link marketing
3. Saat kandidat mendaftar melalui link tersebut, kampanye otomatis tercatat
4. Laporan menampilkan efektivitas setiap kampanye

### Contoh UTM

```
https://stmik.tazkia.ac.id/register?utm_campaign=instagram-ads-2026
```

### Konfigurasi Kampanye

Buka menu **Pengaturan > Kampanye** untuk membuat dan mengelola kampanye.

| Field | Keterangan |
|-------|------------|
| Nama | Nama kampanye |
| Kode | Kode unik (digunakan di UTM) |
| Channel | Instagram, Facebook, Google Ads, Offline, dll |
| Status | Aktif / Nonaktif |

---

## Manajemen Referrer

Referrer adalah pihak yang mereferensikan calon mahasiswa (bisa mahasiswa aktif, alumni, guru, atau agen).

### Akses

Buka menu **Referrer** di sidebar admin.

![Daftar Referrer](screenshots/admin/referrers-list.png)

### Tipe Referrer

| Tipe | Deskripsi |
|------|-----------|
| Mahasiswa | Mahasiswa aktif yang mereferensikan teman |
| Alumni | Lulusan STMIK Tazkia |
| Guru/Sekolah | Guru BK atau sekolah mitra |
| Agen | Agen rekrutmen resmi |

### Kode Referral

Setiap referrer memiliki kode referral unik. Calon mahasiswa memasukkan kode ini saat registrasi (Step 4).

### Konfigurasi Referrer

Buka menu **Pengaturan > Referrer** untuk membuat dan mengelola referrer.

---

## Referral Claim & Linking

Ketika kandidat mendaftar dengan kode referral, sistem mencatat referral claim yang perlu diverifikasi.

### Proses

1. Kandidat memasukkan kode referral saat registrasi
2. Sistem mencatat claim referral
3. Admin memverifikasi dan me-link claim ke referrer yang benar

### Akses

Buka menu **Referral Claims** di sidebar admin.

![Referral Claims](screenshots/admin/referral-claims.png)

### Langkah Linking

1. Buka daftar referral claims
2. Klik claim yang belum di-link
3. Verifikasi kode referral dan referrer
4. Klik **Link** untuk menghubungkan claim ke referrer

---

## Komisi & Reward

Referrer mendapat komisi ketika kandidat yang direferensikan berhasil enrolled.

### Konfigurasi Reward

Buka menu **Pengaturan > Reward** untuk mengatur besaran komisi.

![Reward Config](screenshots/admin/settings-rewards.png)

| Field | Keterangan |
|-------|------------|
| Tipe Referrer | Mahasiswa, Alumni, Guru, Agen |
| Nominal | Besaran komisi per kandidat enrolled |
| Status | Aktif / Nonaktif |

### Proses Komisi

1. Kandidat di-refer oleh referrer
2. Kandidat berhasil enrolled
3. Komisi otomatis tercatat di commission ledger
4. Admin meng-approve komisi
5. Setelah dibayarkan, admin menandai sebagai **Paid**

### Aksi Batch

Admin dapat melakukan approve dan mark as paid secara batch:
- **Batch Approve** — approve beberapa komisi sekaligus
- **Batch Paid** — tandai beberapa komisi sebagai dibayar sekaligus

---

## Multi-Level Referral (MGM)

MGM (Member Get Member) adalah skema komisi berlapis dimana referrer mendapat reward tambahan ketika kandidat yang direferensikan juga mereferensikan kandidat lain.

### Cara Kerja

```
Referrer A → refer → Kandidat B (enrolled) → komisi level 1
Kandidat B → refer → Kandidat C (enrolled) → komisi level 2 untuk Referrer A
```

### Konfigurasi MGM

Buka menu **Pengaturan > Reward** tab **MGM**.

| Field | Keterangan |
|-------|------------|
| Level | Kedalaman level (1, 2, 3, ...) |
| Nominal | Besaran komisi per level |
| Status | Aktif / Nonaktif |
