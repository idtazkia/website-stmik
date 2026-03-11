# Registrasi Calon Mahasiswa

Proses registrasi calon mahasiswa baru terdiri dari 4 step yang harus diselesaikan secara berurutan, dilanjutkan verifikasi email dan telepon.

## Pembuatan Akun (Step 1)

Langkah pertama adalah membuat akun dengan email dan password.

### Langkah

1. Buka halaman `/register`
2. Isi form:

| Field | Keterangan |
|-------|------------|
| Email | Alamat email aktif (akan diverifikasi) |
| Nomor Telepon | Nomor WhatsApp aktif (akan diverifikasi) |
| Password | Minimal 8 karakter |
| Konfirmasi Password | Harus sama dengan password |

3. Klik **Lanjutkan**

![Registrasi Step 1](screenshots/portal/register-step1.png)

> **Catatan**: Email dan nomor telepon harus unik — tidak bisa digunakan untuk registrasi lebih dari satu kali.

---

## Data Diri (Step 2)

Setelah akun dibuat, calon mahasiswa mengisi data diri.

### Langkah

1. Isi form data personal:

| Field | Keterangan |
|-------|------------|
| Nama Lengkap | Sesuai KTP/identitas resmi |
| Tempat Lahir | Kota tempat lahir |
| Tanggal Lahir | Format: DD/MM/YYYY |
| Jenis Kelamin | Laki-laki / Perempuan |
| Alamat | Alamat lengkap saat ini |
| Program Studi | Pilihan prodi yang diminati |

2. Klik **Lanjutkan**

![Registrasi Step 2](screenshots/portal/register-step2.png)

---

## Riwayat Pendidikan (Step 3)

Data riwayat pendidikan terakhir calon mahasiswa.

### Langkah

1. Isi form pendidikan:

| Field | Keterangan |
|-------|------------|
| Jenjang | SMA/SMK/MA/Sederajat |
| Nama Sekolah | Nama institusi pendidikan terakhir |
| Tahun Lulus | Tahun kelulusan |
| Jurusan | Jurusan di sekolah |
| Nilai Rata-rata | Nilai rata-rata rapor/ijazah |

2. Klik **Lanjutkan**

![Registrasi Step 3](screenshots/portal/register-step3.png)

---

## Sumber Informasi (Step 4)

Data tentang bagaimana calon mahasiswa mengetahui STMIK Tazkia.

### Langkah

1. Isi form sumber informasi:

| Field | Keterangan |
|-------|------------|
| Sumber | Media sosial, website, teman, guru, pameran, dll |
| Kode Referral | Opsional — jika direferensikan oleh seseorang |
| Kampanye | Otomatis terisi dari UTM parameter jika ada |

2. Klik **Selesai**

![Registrasi Step 4](screenshots/portal/register-step4.png)

Setelah menyelesaikan step 4:
- Akun kandidat dibuat dengan status **registered**
- Kandidat otomatis di-assign ke konsultan berdasarkan algoritma assignment
- Redirect ke halaman portal

---

## Verifikasi Email

Verifikasi email menggunakan kode OTP yang dikirim ke alamat email.

### Langkah

1. Di portal, buka halaman **Verifikasi Email**
2. Klik **Kirim Kode OTP**
3. Buka email dan salin kode 6 digit
4. Masukkan kode OTP
5. Klik **Verifikasi**

![Verifikasi Email](screenshots/portal/verify-email.png)

Email dikirim melalui layanan Resend. Jika tidak menerima email, periksa folder spam.

---

## Verifikasi Telepon (WhatsApp)

Verifikasi nomor telepon menggunakan kode OTP yang dikirim via WhatsApp.

### Langkah

1. Di portal, klik **Verifikasi Telepon**
2. Klik **Kirim OTP via WhatsApp**
3. Buka WhatsApp dan salin kode 6 digit
4. Masukkan kode OTP
5. Klik **Verifikasi**

![Verifikasi Telepon](screenshots/portal/verify-phone.png)

> **Catatan**: Pastikan nomor WhatsApp aktif dan bisa menerima pesan. OTP dikirim melalui WhatsApp API.
