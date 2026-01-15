// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';

import mdx from '@astrojs/mdx';
import sitemap from '@astrojs/sitemap';

// https://astro.build/config
export default defineConfig({
  site: 'https://stmik.tazkia.ac.id',
  i18n: {
    defaultLocale: 'id',
    locales: ['id', 'en'],
  },
  vite: {
    plugins: [tailwindcss()]
  },

  integrations: [mdx(), sitemap()]
});