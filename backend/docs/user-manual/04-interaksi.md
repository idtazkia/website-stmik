# Pencatatan Interaksi

Bagian ini menjelaskan cara mencatat komunikasi dengan kandidat. Interaksi adalah catatan setiap kontak antara konsultan dan calon mahasiswa.

## Mencatat Interaksi

### Langkah

1. Buka detail kandidat
2. Klik **Catat Interaksi** atau buka tab **Interaksi**
3. Isi form:

| Field | Keterangan |
|-------|------------|
| Tipe Kontak | Telepon, WhatsApp, Email, Kunjungan Kampus, Kunjungan Rumah |
| Kategori | Positif, Negatif, Netral |
| Catatan | Ringkasan pembicaraan |
| Hambatan | Pilih hambatan jika ada (opsional) |
| Tanggal Follow-Up | Tanggal untuk menghubungi kembali (opsional) |
| Saran Supervisor | Catatan saran dari supervisor (opsional) |

4. Klik **Simpan**

### Efek Pencatatan Interaksi

- Jika ini interaksi pertama, status kandidat otomatis berubah dari **registered** ke **prospecting**
- Tanggal follow-up akan muncul di dashboard konsultan
- Riwayat interaksi tercatat di detail kandidat

---

## Kategori Interaksi

Setiap interaksi dikategorikan untuk mengukur sentimen komunikasi.

| Kategori | Deskripsi | Contoh |
|----------|-----------|--------|
| **Positif** | Kandidat menunjukkan minat | "Tertarik, minta info lebih lanjut tentang beasiswa" |
| **Negatif** | Kandidat menunjukkan keberatan | "Masih ragu karena biaya" |
| **Netral** | Tidak ada indikasi jelas | "Belum bisa dihubungi, coba lagi besok" |

Kategori interaksi dapat dikonfigurasi di menu **Pengaturan > Kategori Interaksi**.

---

## Hambatan & Keberatan

Hambatan (obstacles) adalah alasan spesifik yang menghalangi kandidat untuk mendaftar.

### Contoh Hambatan

| Hambatan | Deskripsi |
|----------|-----------|
| Biaya | Kendala finansial |
| Lokasi | Jarak terlalu jauh |
| Waktu | Bentrok dengan aktivitas lain |
| Orang Tua | Belum mendapat izin orang tua |
| Kompetitor | Mempertimbangkan kampus lain |

### Cara Mencatat

Saat mengisi form interaksi, pilih hambatan dari dropdown. Satu interaksi bisa memiliki satu hambatan.

Daftar hambatan dapat dikonfigurasi di menu **Pengaturan > Hambatan**.

---

## Follow-Up & Tindak Lanjut

Setiap interaksi dapat memiliki tanggal follow-up untuk mengingatkan konsultan menghubungi kandidat kembali.

### Cara Kerja

1. Saat mencatat interaksi, isi **Tanggal Follow-Up**
2. Pada tanggal tersebut, kandidat muncul di widget **Follow-Up Hari Ini** di dashboard konsultan
3. Setelah menghubungi, catat interaksi baru

### Tips Follow-Up

- Kandidat yang belum di-follow-up akan muncul dengan highlight di dashboard
- Follow-up yang terlewat tetap muncul sampai interaksi baru dicatat
- Supervisor dapat melihat follow-up terlewat di dashboardnya

---

## Commitment & Enrollment

### Commitment

Ketika kandidat menyatakan komitmen untuk mendaftar:

1. Buka detail kandidat
2. Klik **Commitment**
3. Status berubah dari **prospecting** ke **committed**

Setelah committed, tim finance dapat membuat tagihan untuk kandidat.

### Enrollment

Ketika kandidat sudah menyelesaikan seluruh persyaratan dan pembayaran:

1. Buka detail kandidat
2. Klik **Enroll**
3. Status berubah dari **committed** ke **enrolled**

---

## Mark as Lost

Ketika kandidat tidak melanjutkan proses pendaftaran:

1. Buka detail kandidat
2. Klik **Mark as Lost**
3. Pilih **Alasan Lost** dari dropdown
4. Tambahkan catatan (opsional)
5. Klik **Simpan**

### Alasan Lost

| Alasan | Deskripsi |
|--------|-----------|
| Tidak Merespon | Kandidat tidak bisa dihubungi |
| Memilih Kampus Lain | Kandidat mendaftar ke institusi lain |
| Kendala Biaya | Tidak mampu membayar biaya kuliah |
| Alasan Pribadi | Alasan personal lainnya |

Daftar alasan lost dapat dikonfigurasi di menu **Pengaturan > Alasan Lost**.
