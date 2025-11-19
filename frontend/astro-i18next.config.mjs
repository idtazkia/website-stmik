/** @type {import('astro-i18next').AstroI18nextConfig} */
export default {
  defaultLocale: 'id',
  locales: ['id', 'en'],
  showDefaultLocale: false,
  routes: {
    en: {
      about: 'about',
      programs: 'programs',
      admissions: 'admissions',
      contact: 'contact',
      login: 'login',
      register: 'register',
      dashboard: 'dashboard',
      apply: 'apply',
      applications: 'applications',
    },
    id: {
      about: 'tentang',
      programs: 'program',
      admissions: 'pendaftaran',
      contact: 'kontak',
      login: 'masuk',
      register: 'daftar',
      dashboard: 'dasbor',
      apply: 'lamar',
      applications: 'lamaran',
    },
  },
};
