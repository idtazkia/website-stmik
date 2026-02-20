---
title: "Lima Aplikasi Hasil Proyek Terintegrasi Mahasiswa Semester 3 STMIK Tazkia"
date: 2026-02-13
author: Tim Humas STMIK Tazkia
excerpt: Mahasiswa Sistem Informasi STMIK Tazkia mempresentasikan lima aplikasi hasil proyek terintegrasi semester 3 yang dikembangkan untuk klien nyata, mulai dari AI chatbot, manajemen RT/RW, hingga aplikasi perbankan digital.
images:
  - /images/news/DSC_8124.jpg
  - /images/news/DSC_8126.jpg
  - /images/news/DSC_8130.jpg
  - /images/news/DSC_8131.jpg
  - /images/news/DSC_8132.jpg
  - /images/news/DSC_8135.jpg
  - /images/news/DSC_8136.jpg
  - /images/news/DSC_8138.jpg
  - /images/news/DSC_8139.jpg
  - /images/news/DSC_8143.jpg
  - /images/news/DSC_8149.jpg
  - /images/news/DSC_8159.jpg
tags:
  - proyek-mahasiswa
  - project-based-learning
  - sistem-informasi
  - presentasi
---

Pada hari Jumat, 13 Februari 2026, mahasiswa semester 3 Program Studi [Sistem Informasi](/programs/information-systems) STMIK Tazkia mempresentasikan lima aplikasi yang mereka bangun selama satu semester di International Class Tazkia Sentul. Kelima aplikasi ini merupakan hasil dari **Proyek Terintegrasi** yang menggabungkan empat mata kuliah sekaligus:

- **MISI.304** — Analisis dan Perancangan Sistem Informasi
- **MISI.301** — Manajemen Basis Data Terintegrasi
- **MISI.302** — Sistem Informasi Berbasis Web
- **MISI.303** — Integrasi dan Deployment Sistem

Setiap kelompok bekerja dengan klien nyata, menganalisis kebutuhan, merancang sistem, mengimplementasikan aplikasi, hingga men-deploy ke server produksi. Seluruh aplikasi dapat diakses publik melalui internet.

## 1. Sapa Tazkia — AI Chatbot Kampus

| | |
|---|---|
| **Klien** | STMIK Tazkia |
| **URL** | [sapa.tazkia.ac.id](https://sapa.tazkia.ac.id) |
| **Tech Stack** | Node.js, React, MySQL, Qdrant (vector DB), Redis |
| **GitHub** | [Code](https://github.com/IchsanJunaedi/Sapa-Tazkia) · [Docs](https://github.com/IchsanJunaedi/sapa-tazkia-documents) |

![Presentasi Sapa Tazkia](/images/news/DSC_8124.jpg)

Sapa Tazkia adalah chatbot berbasis Retrieval-Augmented Generation (RAG) yang menjawab pertanyaan seputar kampus menggunakan dokumen resmi sebagai sumber data. Sistem melakukan embedding terhadap pertanyaan pengguna, mencari dokumen relevan melalui Qdrant, lalu mengirimkan konteks ke LLM untuk menghasilkan jawaban yang akurat.

Fitur utama:
- Chatbot RAG dengan sumber data dokumen kampus
- Autentikasi mahasiswa terintegrasi database akademik
- Akses data nilai dan IPK
- Download transkrip otomatis dalam format PDF
- Rate limiting dan chat history

| Nama | Peran |
|------|-------|
| Muhammad Ichsan Junaedi | PM / Developer |
| Rahmawati | Business Analyst |
| Rafly Ariel Hidayat | Backend Developer |
| Nabila Nurul Haq | UI/UX Designer |

![Demo Sapa Tazkia](/images/news/DSC_8126.jpg)

## 2. Lingkar Warga — Manajemen RT/RW

| | |
|---|---|
| **Klien** | Ketua RW/RT |
| **URL** | [rtrw.demo.tazkia.ac.id](https://rtrw.demo.tazkia.ac.id) |
| **Tech Stack** | Express.js, Flutter, PostgreSQL, Docker |
| **GitHub** | [Code](https://github.com/azmttqi/Aplikasi-Manajemen-RT-RW) · [Docs](https://github.com/azmttqi/Project-Manajemen-RT-RW-Document) |

![Presentasi Lingkar Warga](/images/news/DSC_8130.jpg)

Lingkar Warga adalah platform administrasi lingkungan digital yang menggantikan proses manual pendataan warga, verifikasi akun, dan pelaporan. Aplikasi ini dibangun dengan Flutter sehingga dapat berjalan di web maupun mobile dari satu codebase. Backend dikontainerisasi dengan Docker Compose bersama database PostgreSQL.

Fitur utama:
- Pendataan warga digital
- Sistem verifikasi akun dengan email SMTP
- Pelaporan mandiri oleh warga
- Manajemen data RT/RW multi-level
- Containerized deployment dengan Docker Compose

| Nama | Peran |
|------|-------|
| Azmi Ittaqi Hammami | PM / System Analyst / Backend |
| Amanda Wijayanti | UI/UX / Frontend Developer |
| Muhammad Nabil Thoriq | Test Engineer |

![Demo Lingkar Warga](/images/news/DSC_8131.jpg)

## 3. Aplikasi Mutaba'ah — Monitoring Ibadah Mahasiswa

| | |
|---|---|
| **Klien** | STMIK Tazkia |
| **URL** | [ibadah.demo.tazkia.ac.id](https://ibadah.demo.tazkia.ac.id) |
| **Tech Stack** | Express.js, MongoDB Atlas, React, Chart.js |
| **GitHub** | [Code](https://github.com/runaisyah1337/ProjectWebsite_Mutabaah-Mahasiswa) · [Docs](https://github.com/PatoyUhuy/Mutaba-ah_Mahasiswa_Docs) |

![Presentasi Aplikasi Mutaba'ah](/images/news/DSC_8136.jpg)

Sistem monitoring ibadah harian (mutaba'ah) yang mendigitalisasi proses evaluasi ibadah pekanan mahasiswa. Mahasiswa mengisi 9 indikator ibadah setiap minggu, lalu sistem menghitung skor secara otomatis. Data dikunci setiap Minggu pukul 23:59 WIB untuk menjaga integritas.

Fitur utama:
- Form pengisian mandiri 9 indikator ibadah
- Automated scoring (skor 1-3 per indikator, maks 27/minggu)
- Dashboard visualisasi dengan bar chart dan line chart
- Peran mahasiswa, pembina, dan admin
- Smart search lintas kelompok bimbingan
- Auto-locking data per periode

| Nama | Peran |
|------|-------|
| Abdurrahman Fathi M. | Project Manager |
| Destri | Business Analyst |
| Mutiara | UI/UX Designer |
| Aisyah | Quality Assurance |

![Demo Aplikasi Mutaba'ah](/images/news/DSC_8138.jpg)

## 4. Manajemen Yayasan Sahabat Qur'an

| | |
|---|---|
| **Klien** | Yayasan Sahabat Qur'an |
| **URL** | [akademik.sahabatquran.com](https://akademik.sahabatquran.com) |
| **Tech Stack** | HTML, Tailwind CSS, JavaScript, PostgreSQL, Docker |
| **GitHub** | [Code](https://github.com/nisa262006/Manajemen_YSQ) · [Docs](https://github.com/rizkasugiarto/ysq-project-documents) |

![Presentasi Manajemen YSQ](/images/news/DSC_8132.jpg)

Aplikasi manajemen operasional untuk lembaga pendidikan Al-Qur'an yang sebelumnya mengelola seluruh proses secara manual menggunakan Google Forms dan spreadsheet. Sistem ini mengintegrasikan pendaftaran santri, penjadwalan kelas, pencatatan kehadiran, hingga pengelolaan keuangan dalam satu platform.

Fitur utama:
- Registrasi dan verifikasi santri baru
- Manajemen kelas dan penempatan santri
- Penjadwalan pelajaran
- Absensi santri dan pengajar
- Pencatatan nilai dan rapor digital (export PDF)
- Manajemen tagihan, verifikasi pembayaran, dan laporan keuangan
- Perhitungan honor pengajar

| Nama | Peran |
|------|-------|
| Rizka | PM / Quality Documentation |
| Fikri | Backend / Database Developer |
| Jingga | Frontend / UI/UX Developer |
| Rahma Fitria | Developer |

![Demo Manajemen YSQ](/images/news/DSC_8135.jpg)

## 5. MiniBank — Aplikasi Keuangan Digital

| | |
|---|---|
| **Klien** | Melvina Lubis |
| **URL** | [dev.minibank.tazkia.ac.id](https://dev.minibank.tazkia.ac.id) |
| **Tech Stack** | Spring Boot, Thymeleaf, PostgreSQL |
| **GitHub** | [Code](https://github.com/nflFauzan/project-minibank) · [Docs](https://github.com/MuhammadAkmalSyarif/minibank-docs) |

![Presentasi MiniBank](/images/news/DSC_8139.jpg)

MiniBank adalah simulasi aplikasi perbankan digital yang menyediakan layanan transaksi dasar. Sistem dikembangkan dengan Spring Boot di backend dan Thymeleaf untuk server-side rendering, menggunakan PostgreSQL sebagai database transaksional.

Fitur utama:
- Onboarding nasabah (KYC digital)
- Manajemen rekening dan saldo real-time
- Transfer dana antar-rekening
- Pembayaran tagihan
- Riwayat transaksi dan laporan mutasi (PDF/CSV)
- Autentikasi multi-faktor

| Nama | Peran |
|------|-------|
| Muhammad Naufal Fauzan | Developer |
| Muhammad Akmal Syarif | Developer |
| Siti Tahtia Ainun Zahra | Developer |
| Rackisha Dhia Ezelly Lathief | Developer |

![Demo MiniBank](/images/news/DSC_8143.jpg)

![Foto bersama peserta presentasi](/images/news/DSC_8149.jpg)

![Foto bersama seluruh mahasiswa dan dosen](/images/news/DSC_8159.jpg)

## Proyek Terintegrasi: Belajar dari Klien Nyata

Model [project-based learning](/about) di STMIK Tazkia mengharuskan mahasiswa mengerjakan proyek untuk klien nyata sejak semester awal. Proyek terintegrasi semester 3 ini menggabungkan empat mata kuliah dalam satu deliverable, sehingga mahasiswa mengalami siklus pengembangan software secara utuh: dari analisis kebutuhan, perancangan database, implementasi web, hingga deployment ke server produksi.

Kelima proyek di atas menunjukkan variasi tech stack yang luas — dari React + Node.js, Flutter + Express, hingga Spring Boot + Thymeleaf — sesuai dengan kebutuhan masing-masing klien. Seluruh aplikasi telah di-deploy dan dapat diakses publik, bukan sekadar tugas yang berhenti di localhost.

Tertarik kuliah IT dengan pendekatan project-based learning? STMIK Tazkia membuka pendaftaran untuk program [Sistem Informasi](/programs/information-systems) dan [Teknik Informatika](/programs/computer-engineering).

[Daftar sekarang](/admissions) atau [hubungi kami](/contact) untuk informasi lebih lanjut.
