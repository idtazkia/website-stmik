# Search Engine Indexing - Under Construction Mode

## Current Status: 🚫 INDEXING DISABLED

The site is currently configured to **prevent all search engine indexing** while under construction. This prevents incomplete or placeholder content from appearing in search results.

---

## How Indexing is Currently Blocked

### 1. robots.txt (Primary Block)
**File:** `frontend/public/robots.txt`

```
User-agent: *
Disallow: /
```

This tells all search engine crawlers (Google, Bing, etc.) not to crawl or index any pages on the site.

### 2. Meta Robots Tag (Secondary Block)
**File:** `frontend/src/layouts/BaseLayout.astro`

```html
<meta name="robots" content="noindex, nofollow" />
```

This is added to the `<head>` of every page via the `noIndex = true` default in BaseLayout.

**Both methods work together** to ensure maximum protection against indexing.

---

## When to Re-Enable Indexing

Enable search engine indexing when:
- ✅ All pages have accurate, final content (no placeholders)
- ✅ All images and assets are high-quality and final
- ✅ All links are working correctly
- ✅ Content has been proofread and approved
- ✅ Site is ready for public visibility
- ✅ Legal pages (Privacy Policy, Terms) are published
- ✅ Contact information is accurate and current

---

## How to Re-Enable Indexing (When Ready to Launch)

### Step 1: Update robots.txt

**File:** `frontend/public/robots.txt`

Replace current content with:

```
# STMIK Tazkia - Allow Search Engine Indexing
User-agent: *
Allow: /

# Disallow authenticated pages
Disallow: /dashboard
Disallow: /apply
Disallow: /applications
Disallow: /admin

# Sitemap location
Sitemap: https://stmik.tazkia.ac.id/sitemap-index.xml
```

### Step 2: Update BaseLayout Default

**File:** `frontend/src/layouts/BaseLayout.astro`

Change line 14 from:
```typescript
const { title, description, noIndex = true } = Astro.props;
```

To:
```typescript
const { title, description, noIndex = false } = Astro.props;
```

Also remove the "SITE UNDER CONSTRUCTION" comment (lines 11-13).

### Step 3: Rebuild and Deploy

```bash
cd frontend
npm run build
# Deploy to production
```

### Step 4: Verify Changes

After deployment, check:

1. **robots.txt is accessible:**
   - Visit: https://stmik.tazkia.ac.id/robots.txt
   - Verify it shows "Allow: /" instead of "Disallow: /"

2. **Meta robots tag is removed:**
   - Visit any page
   - View page source (Ctrl+U or Cmd+Option+U)
   - Search for `<meta name="robots"`
   - Should NOT see `content="noindex, nofollow"`

3. **Sitemap is accessible:**
   - Visit: https://stmik.tazkia.ac.id/sitemap-index.xml
   - Should show XML sitemap

### Step 5: Submit Sitemap to Search Engines

**Google Search Console:**
1. Go to https://search.google.com/search-console
2. Add property: stmik.tazkia.ac.id
3. Verify ownership (DNS or HTML file method)
4. Submit sitemap: https://stmik.tazkia.ac.id/sitemap-index.xml

**Bing Webmaster Tools:**
1. Go to https://www.bing.com/webmasters
2. Add site: stmik.tazkia.ac.id
3. Verify ownership
4. Submit sitemap

---

## Testing Indexing Status

### Check if Site is Blocked (Current State)

**Method 1: Google Search**
```
site:stmik.tazkia.ac.id
```
Should return: "No results found" or very few results

**Method 2: Google Search Console**
- Check "Coverage" report
- Should show pages as "Excluded by 'noindex' tag"

**Method 3: robots.txt Tester**
- Google Search Console → Crawl → robots.txt Tester
- Test URL: https://stmik.tazkia.ac.id/
- Should show: "Blocked"

### Check if Site is Indexed (After Launch)

**Method 1: Google Search**
```
site:stmik.tazkia.ac.id
```
Should return: Multiple pages listed

**Method 2: Google Search Console**
- Check "Coverage" report
- Should show pages as "Valid" and indexed

**Method 3: URL Inspection Tool**
- Google Search Console → URL Inspection
- Enter any page URL
- Should show: "URL is on Google"

---

## Gradual Indexing (Alternative Approach)

If you want to allow indexing for some pages but not others:

### Option 1: Per-Page Control

In individual page files, override the noIndex prop:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
---

<!-- Allow this specific page to be indexed -->
<BaseLayout noIndex={false} title="About Us">
  ...
</BaseLayout>
```

### Option 2: Selective robots.txt

Update robots.txt to allow specific sections:

```
User-agent: *
Disallow: /
Allow: /about
Allow: /programs
Allow: /lecturers
```

---

## Common Issues & Solutions

### Issue: Pages still showing in Google after blocking

**Cause:** Google has already indexed the pages
**Solution:**
1. Update robots.txt and meta tags as above
2. In Google Search Console, use "Removals" tool
3. Request temporary removal of URLs
4. Wait 24-48 hours for Google to re-crawl

### Issue: Changes not taking effect

**Cause:** Browser or CDN caching
**Solution:**
1. Clear browser cache (Ctrl+Shift+R or Cmd+Shift+R)
2. Clear Cloudflare cache (if using Cloudflare)
3. Wait 24 hours for search engines to re-crawl

### Issue: Sitemap shows 404

**Cause:** Sitemap integration not working
**Solution:**
1. Verify `@astrojs/sitemap` is installed: `npm list @astrojs/sitemap`
2. Check `astro.config.mjs` includes `sitemap()` in integrations
3. Rebuild: `npm run build`

---

## Monitoring Indexing Status

### Weekly Checks (While Blocked)
- [ ] Verify robots.txt still shows "Disallow: /"
- [ ] Check Google Search: `site:stmik.tazkia.ac.id` returns no results
- [ ] Review any accidental backlinks or references

### Weekly Checks (After Launch)
- [ ] Monitor Google Search Console for crawl errors
- [ ] Check indexing status of new pages
- [ ] Review search performance (clicks, impressions, CTR)
- [ ] Fix any "Excluded" pages in Coverage report

---

## Important Notes

1. **robots.txt is a request, not enforcement:** Some bad bots may ignore it. The meta robots tag provides additional protection.

2. **Existing indexed pages:** If Google already indexed pages before blocking, they may remain in search results for a while. Use Google Search Console's Removal tool to expedite.

3. **Development vs Production:** If you have separate dev and production sites, keep dev blocked permanently and only enable indexing on production.

4. **Staging environments:** Always block staging/preview URLs:
   ```
   # In robots.txt
   User-agent: *
   Disallow: /
   ```

5. **Sitemap:** The sitemap is still generated even when blocking. This is fine and useful for internal testing.

---

## Checklist: Launch Day Indexing Tasks

- [ ] Update `robots.txt` to allow crawling
- [ ] Change `noIndex` default to `false` in BaseLayout
- [ ] Rebuild site: `npm run build`
- [ ] Deploy to production
- [ ] Verify robots.txt is accessible and shows "Allow: /"
- [ ] Verify meta robots tags are removed (view source)
- [ ] Submit sitemap to Google Search Console
- [ ] Submit sitemap to Bing Webmaster Tools
- [ ] Request indexing of homepage via URL Inspection tool
- [ ] Monitor Google Search Console for 1-2 weeks
- [ ] Check `site:` search after 1 week

---

**Last Updated:** 2025-01-19
**Status:** Indexing Disabled ✅
**Ready to Launch:** Change files as instructed above
