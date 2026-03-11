# Pengaturan Sistem

Bagian ini menjelaskan konfigurasi yang tersedia di menu Pengaturan untuk admin.

## Manajemen User & Role

Admin mengelola akun staff yang memiliki akses ke panel admin.

### Akses

Buka menu **Pengaturan > User** di sidebar admin.

![Daftar User](screenshots/admin/settings-users.png)

### Daftar Role

| Role | Deskripsi | Akses |
|------|-----------|-------|
| **Admin** | Administrator sistem | Full access |
| **Supervisor** | Pengawas konsultan | Dashboard supervisor, reassign, laporan |
| **Consultant** | Konsultan penerimaan | Kandidat yang di-assign, interaksi |
| **Finance** | Tim keuangan | Tagihan, verifikasi pembayaran |
| **Academic** | Tim akademik | Data akademik, enrollment |

### Tambah User

1. Klik **User Baru**
2. Isi email Google institusi (`@tazkia.ac.id`)
3. Pilih role
4. Pilih supervisor (jika role = Consultant)
5. Klik **Simpan**

### Ubah Role

1. Klik user di daftar
2. Ubah role dari dropdown
3. Klik **Simpan**

### Nonaktifkan User

1. Klik user di daftar
2. Klik toggle **Aktif**
3. User yang dinonaktifkan tidak bisa login

---

## Program Studi

Konfigurasi program studi yang tersedia untuk calon mahasiswa.

### Akses

Buka menu **Pengaturan > Program Studi**.

![Program Studi](screenshots/admin/settings-programs.png)

### Tambah Program Studi

1. Klik **Program Studi Baru**
2. Isi form:

| Field | Keterangan |
|-------|------------|
| Nama | Nama program studi |
| Jenjang | S1 / D3 / D4 |
| Kode | Kode unik prodi |
| Status | Aktif / Nonaktif |

3. Klik **Simpan**

---

## Kategori Interaksi

Konfigurasi kategori yang digunakan saat mencatat interaksi dengan kandidat.

### Akses

Buka menu **Pengaturan > Kategori Interaksi**.

### Kategori Default

| Kategori | Sentimen | Deskripsi |
|----------|----------|-----------|
| Tertarik | Positif | Kandidat menunjukkan minat |
| Ragu-ragu | Netral | Kandidat masih mempertimbangkan |
| Menolak | Negatif | Kandidat menolak |
| Tidak Bisa Dihubungi | Netral | Tidak ada respons |

### Tambah Kategori

1. Klik **Kategori Baru**
2. Isi nama kategori dan sentimen (positif/netral/negatif)
3. Klik **Simpan**

---

## Hambatan Pendaftaran

Konfigurasi daftar hambatan (obstacles) yang dapat dipilih saat mencatat interaksi.

### Akses

Buka menu **Pengaturan > Hambatan**.

### Contoh Hambatan

| Hambatan | Deskripsi |
|----------|-----------|
| Biaya | Kendala finansial |
| Lokasi | Jarak dari kampus |
| Waktu | Jadwal bentrok |
| Orang Tua | Izin orang tua |
| Kompetitor | Mempertimbangkan kampus lain |

---

## Struktur Biaya

Konfigurasi biaya pendaftaran per program studi.

### Akses

Buka menu **Pengaturan > Biaya**.

![Struktur Biaya](screenshots/admin/settings-fees.png)

### Konfigurasi

Setiap entri struktur biaya terdiri dari:

| Field | Keterangan |
|-------|------------|
| Program Studi | Prodi yang terkait |
| Jenis Biaya | Tipe biaya (pendaftaran, SPP, dll) |
| Nominal | Jumlah dalam rupiah |
| Semester | Semester ke-berapa (untuk SPP) |

---

## Algoritma Assignment

Konfigurasi cara sistem mendistribusikan kandidat baru ke konsultan.

### Akses

Buka menu **Pengaturan > Assignment**.

![Assignment](screenshots/admin/settings-assignment.png)

### Algoritma yang Tersedia

| Algoritma | Cara Kerja |
|-----------|------------|
| **Round-Robin** | Giliran secara berurutan, setiap konsultan mendapat 1 kandidat sebelum beralih ke berikutnya |
| **Load-Based** | Prioritaskan konsultan dengan jumlah kandidat aktif paling sedikit |

### Cara Mengubah

1. Pilih algoritma yang diinginkan
2. Klik **Simpan**
3. Algoritma baru berlaku untuk kandidat yang mendaftar selanjutnya

---

## Alasan Lost

Konfigurasi alasan-alasan mengapa kandidat tidak melanjutkan pendaftaran.

### Akses

Buka menu **Pengaturan > Alasan Lost**.

### Contoh Alasan

| Alasan | Deskripsi |
|--------|-----------|
| Tidak Merespon | Kandidat tidak bisa dihubungi dalam waktu lama |
| Memilih Kampus Lain | Kandidat mendaftar ke institusi lain |
| Kendala Biaya | Tidak mampu membayar |
| Alasan Pribadi | Alasan personal |
| Tidak Memenuhi Syarat | Tidak lolos persyaratan |

### Tambah Alasan

1. Klik **Alasan Baru**
2. Isi nama alasan
3. Klik **Simpan**
