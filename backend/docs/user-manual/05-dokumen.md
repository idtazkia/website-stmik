# Manajemen Dokumen

Bagian ini menjelaskan proses upload dan review dokumen persyaratan pendaftaran.

## Upload Dokumen (Portal)

Calon mahasiswa mengupload dokumen persyaratan melalui portal.

### Langkah

1. Login ke portal
2. Buka menu **Dokumen**
3. Pilih tipe dokumen yang akan diupload
4. Klik **Pilih File** dan pilih file dari perangkat
5. Klik **Upload**

![Upload Dokumen](screenshots/portal/documents.png)

### Ketentuan Upload

| Ketentuan | Nilai |
|-----------|-------|
| Format file | PDF, JPG, PNG |
| Ukuran maksimal | 5 MB per file |
| Lokasi penyimpanan | Server (direktori upload) |

### Status Dokumen

| Status | Deskripsi |
|--------|-----------|
| **Pending** | Baru diupload, menunggu review |
| **Approved** | Dokumen diterima |
| **Rejected** | Dokumen ditolak, perlu upload ulang |

Kandidat mendapat notifikasi email saat status dokumen berubah.

---

## Review Dokumen (Admin)

Admin mereview dokumen yang diupload oleh kandidat.

### Akses

Buka menu **Dokumen** di sidebar admin.

![Review Dokumen](screenshots/admin/documents-review.png)

### Langkah Review

1. Buka halaman **Dokumen**
2. Pilih dokumen dengan status **Pending**
3. Klik dokumen untuk melihat preview
4. Klik **Approve** jika dokumen valid
5. Klik **Reject** jika dokumen tidak valid — isi alasan penolakan

### Notifikasi

Setelah approve atau reject:
- Email notifikasi dikirim ke kandidat
- Status dokumen terupdate di portal kandidat
- Jika ditolak, kandidat dapat mengupload ulang

---

## Pengaturan Tipe Dokumen

Admin dapat mengkonfigurasi tipe dokumen yang wajib diupload.

### Akses

Buka menu **Pengaturan > Tipe Dokumen**.

![Tipe Dokumen](screenshots/admin/settings-document-types.png)

### Contoh Tipe Dokumen

| Tipe Dokumen | Wajib | Deskripsi |
|--------------|-------|-----------|
| KTP/Kartu Identitas | Ya | Scan KTP atau kartu identitas |
| Ijazah | Ya | Scan ijazah terakhir |
| Transkrip Nilai | Ya | Scan transkrip atau rapor |
| Pas Foto | Ya | Foto 3x4 latar belakang merah |
| Surat Keterangan Sehat | Tidak | Dari dokter/puskesmas |

### Tambah Tipe Dokumen

1. Klik **Tipe Dokumen Baru**
2. Isi nama tipe dokumen
3. Tentukan apakah wajib atau opsional
4. Klik **Simpan**
