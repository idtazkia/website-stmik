# Manajemen Kandidat

Bagian ini menjelaskan cara mengelola data calon mahasiswa (kandidat) di panel admin.

## Daftar Kandidat

Halaman daftar kandidat menampilkan seluruh calon mahasiswa yang terdaftar.

### Akses

Buka menu **Kandidat** di sidebar admin.

![Daftar Kandidat](screenshots/admin/candidates-list.png)

### Filter dan Pencarian

| Filter | Deskripsi |
|--------|-----------|
| Status | Filter berdasarkan status: registered, prospecting, committed, enrolled, lost |
| Konsultan | Filter berdasarkan konsultan yang di-assign |
| Program Studi | Filter berdasarkan prodi pilihan |
| Tanggal | Filter berdasarkan tanggal registrasi |
| Pencarian | Cari berdasarkan nama atau email |

### Informasi yang Ditampilkan

Setiap baris menampilkan:
- Nama kandidat
- Email dan nomor telepon
- Program studi pilihan
- Status saat ini
- Konsultan yang di-assign
- Tanggal registrasi

---

## Detail Kandidat

Halaman detail menampilkan informasi lengkap seorang kandidat.

### Akses

Klik nama kandidat di daftar.

![Detail Kandidat](screenshots/admin/candidate-detail.png)

### Tab Informasi

| Tab | Isi |
|-----|-----|
| Data Diri | Nama, email, telepon, alamat, tanggal lahir |
| Pendidikan | Riwayat pendidikan terakhir |
| Dokumen | Dokumen yang sudah diupload beserta status review |
| Interaksi | Riwayat seluruh interaksi dengan kandidat |
| Pembayaran | Status tagihan dan bukti pembayaran |

---

## Assignment Konsultan

Kandidat baru otomatis di-assign ke konsultan berdasarkan algoritma yang dikonfigurasi.

### Algoritma Assignment

| Algoritma | Cara Kerja |
|-----------|------------|
| Round-Robin | Distribusi merata secara bergiliran |
| Load-Based | Assign ke konsultan dengan jumlah kandidat aktif paling sedikit |

Konfigurasi algoritma dilakukan di menu **Pengaturan > Assignment**.

### Assignment Otomatis

Saat kandidat menyelesaikan registrasi (step 4), sistem otomatis:
1. Memilih konsultan berdasarkan algoritma aktif
2. Assign kandidat ke konsultan tersebut
3. Status kandidat menjadi **registered**

---

## Reassign Kandidat

Supervisor atau admin dapat memindahkan kandidat ke konsultan lain.

### Langkah

1. Buka detail kandidat
2. Klik **Reassign**
3. Pilih konsultan baru dari dropdown
4. Isi alasan reassign (opsional)
5. Klik **Simpan**

![Reassign Kandidat](screenshots/admin/candidate-reassign.png)

### Siapa yang Bisa Reassign?

| Role | Akses |
|------|-------|
| Admin | Bisa reassign semua kandidat |
| Supervisor | Bisa reassign kandidat di bawah supervisinya |
| Konsultan | Tidak bisa reassign |

---

## Status & Lifecycle Kandidat

Setiap kandidat memiliki status yang mencerminkan posisinya dalam funnel penerimaan.

### Status Flow

```
registered → prospecting → committed → enrolled
                ↓
              lost
```

### Deskripsi Status

| Status | Deskripsi | Trigger |
|--------|-----------|---------|
| **registered** | Baru mendaftar, belum ada interaksi | Selesai registrasi step 4 |
| **prospecting** | Sedang dalam proses komunikasi | Interaksi pertama dicatat |
| **committed** | Kandidat menyatakan komitmen untuk mendaftar | Klik tombol **Commitment** |
| **enrolled** | Sudah terdaftar resmi sebagai mahasiswa | Klik tombol **Enroll** |
| **lost** | Tidak melanjutkan pendaftaran | Klik tombol **Mark as Lost** |

### Transisi Status yang Diizinkan

| Dari | Ke | Aksi |
|------|----|------|
| registered | prospecting | Otomatis saat interaksi pertama dicatat |
| prospecting | committed | Klik **Commitment** di detail kandidat |
| prospecting | lost | Klik **Mark as Lost** di detail kandidat |
| committed | enrolled | Klik **Enroll** di detail kandidat |
| committed | lost | Klik **Mark as Lost** di detail kandidat |
