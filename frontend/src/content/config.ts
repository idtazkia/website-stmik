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

const lecturers = defineCollection({
  type: 'content',
  schema: z.object({
    name: z.string(),
    title: z.string(), // e.g., "Dosen Tetap", "Dosen Luar Biasa"
    position: z.string().optional(), // e.g., "Ketua Program Studi"
    expertise: z.array(z.string()),
    education: z.array(z.object({
      degree: z.string(),
      institution: z.string(),
      year: z.string().optional(),
    })),
    email: z.string().optional(),
    phone: z.string().optional(),
    website: z.string().optional(),
    github: z.string().optional(),
    linkedin: z.string().optional(),
    youtube: z.string().optional(),
    photo: z.string().optional(),
    order: z.number().default(999), // For sorting
  }),
});

const news = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    date: z.date(),
    author: z.string().optional(),
    excerpt: z.string(),
    images: z.array(z.string()).optional(),
    tags: z.array(z.string()).optional(),
  }),
});

export const collections = { programs, about, admissions, lecturers, news };
