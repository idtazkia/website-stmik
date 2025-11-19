import { defineCollection, z } from 'astro:content';

const programs = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    slug: z.string(),
    description: z.string(),
    duration: z.string(),
    degree: z.string(),
  }),
});

const about = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    slug: z.string(),
  }),
});

const admissions = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    slug: z.string(),
  }),
});

export const collections = { programs, about, admissions };
