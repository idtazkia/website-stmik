package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Section represents a single manual section extracted from a markdown file.
type Section struct {
	ID          string
	Title       string
	File        string
	Screenshots []string
}

// SectionGroup is a collapsible group of sections in the sidebar TOC.
type SectionGroup struct {
	ID       string
	Title    string
	Icon     string
	Sections []Section
}

func getSectionGroups() []SectionGroup {
	return []SectionGroup{
		{
			ID:    "login-dashboard",
			Title: "Login & Dashboard",
			Icon:  iconDashboard,
			Sections: []Section{
				{ID: "login-admin", Title: "Login Admin (Google OIDC)", File: "01-login-dashboard.md"},
				{ID: "login-calon", Title: "Login Calon Mahasiswa", File: "01-login-dashboard.md"},
				{ID: "dashboard-admin", Title: "Dashboard Admin", File: "01-login-dashboard.md"},
				{ID: "dashboard-konsultan", Title: "Dashboard Konsultan", File: "01-login-dashboard.md"},
				{ID: "dashboard-supervisor", Title: "Dashboard Supervisor", File: "01-login-dashboard.md"},
			},
		},
		{
			ID:    "registrasi",
			Title: "Registrasi Calon Mahasiswa",
			Icon:  iconRegistration,
			Sections: []Section{
				{ID: "registrasi-akun", Title: "Pembuatan Akun (Step 1)", File: "02-registrasi.md"},
				{ID: "registrasi-data-diri", Title: "Data Diri (Step 2)", File: "02-registrasi.md"},
				{ID: "registrasi-pendidikan", Title: "Riwayat Pendidikan (Step 3)", File: "02-registrasi.md"},
				{ID: "registrasi-sumber", Title: "Sumber Informasi (Step 4)", File: "02-registrasi.md"},
				{ID: "verifikasi-email", Title: "Verifikasi Email", File: "02-registrasi.md"},
				{ID: "verifikasi-telepon", Title: "Verifikasi Telepon (WhatsApp)", File: "02-registrasi.md"},
			},
		},
		{
			ID:    "kandidat",
			Title: "Manajemen Kandidat",
			Icon:  iconCandidates,
			Sections: []Section{
				{ID: "daftar-kandidat", Title: "Daftar Kandidat", File: "03-manajemen-kandidat.md"},
				{ID: "detail-kandidat", Title: "Detail Kandidat", File: "03-manajemen-kandidat.md"},
				{ID: "assignment-konsultan", Title: "Assignment Konsultan", File: "03-manajemen-kandidat.md"},
				{ID: "reassign-kandidat", Title: "Reassign Kandidat", File: "03-manajemen-kandidat.md"},
				{ID: "status-kandidat", Title: "Status & Lifecycle Kandidat", File: "03-manajemen-kandidat.md"},
			},
		},
		{
			ID:    "interaksi",
			Title: "Pencatatan Interaksi",
			Icon:  iconInteraction,
			Sections: []Section{
				{ID: "catat-interaksi", Title: "Mencatat Interaksi", File: "04-interaksi.md"},
				{ID: "kategori-interaksi", Title: "Kategori Interaksi", File: "04-interaksi.md"},
				{ID: "hambatan-keberatan", Title: "Hambatan & Keberatan", File: "04-interaksi.md"},
				{ID: "follow-up", Title: "Follow-Up & Tindak Lanjut", File: "04-interaksi.md"},
				{ID: "commitment-enrollment", Title: "Commitment & Enrollment", File: "04-interaksi.md"},
				{ID: "mark-lost", Title: "Mark as Lost", File: "04-interaksi.md"},
			},
		},
		{
			ID:    "dokumen",
			Title: "Manajemen Dokumen",
			Icon:  iconDocuments,
			Sections: []Section{
				{ID: "upload-dokumen", Title: "Upload Dokumen (Portal)", File: "05-dokumen.md"},
				{ID: "review-dokumen", Title: "Review Dokumen (Admin)", File: "05-dokumen.md"},
				{ID: "tipe-dokumen", Title: "Pengaturan Tipe Dokumen", File: "05-dokumen.md"},
			},
		},
		{
			ID:    "keuangan",
			Title: "Keuangan & Pembayaran",
			Icon:  iconFinance,
			Sections: []Section{
				{ID: "struktur-biaya", Title: "Struktur Biaya", File: "06-keuangan.md"},
				{ID: "tagihan", Title: "Pembuatan Tagihan", File: "06-keuangan.md"},
				{ID: "upload-bukti", Title: "Upload Bukti Pembayaran (Portal)", File: "06-keuangan.md"},
				{ID: "verifikasi-pembayaran", Title: "Verifikasi Pembayaran (Admin)", File: "06-keuangan.md"},
			},
		},
		{
			ID:    "marketing",
			Title: "Marketing & Referral",
			Icon:  iconMarketing,
			Sections: []Section{
				{ID: "kampanye", Title: "Manajemen Kampanye", File: "07-marketing-referral.md"},
				{ID: "referrer", Title: "Manajemen Referrer", File: "07-marketing-referral.md"},
				{ID: "referral-claim", Title: "Referral Claim & Linking", File: "07-marketing-referral.md"},
				{ID: "komisi", Title: "Komisi & Reward", File: "07-marketing-referral.md"},
				{ID: "mgm-reward", Title: "Multi-Level Referral (MGM)", File: "07-marketing-referral.md"},
			},
		},
		{
			ID:    "portal",
			Title: "Portal Calon Mahasiswa",
			Icon:  iconPortal,
			Sections: []Section{
				{ID: "portal-dashboard", Title: "Dashboard Portal", File: "08-portal.md"},
				{ID: "portal-dokumen", Title: "Upload Dokumen", File: "08-portal.md"},
				{ID: "portal-pembayaran", Title: "Status Pembayaran", File: "08-portal.md"},
				{ID: "portal-pengumuman", Title: "Pengumuman", File: "08-portal.md"},
				{ID: "portal-referral", Title: "Program Referral", File: "08-portal.md"},
			},
		},
		{
			ID:    "laporan",
			Title: "Laporan & Analisis",
			Icon:  iconReports,
			Sections: []Section{
				{ID: "laporan-funnel", Title: "Analisis Funnel", File: "09-laporan.md"},
				{ID: "laporan-konsultan", Title: "Performa Konsultan", File: "09-laporan.md"},
				{ID: "laporan-kampanye", Title: "Efektivitas Kampanye", File: "09-laporan.md"},
				{ID: "laporan-referrer", Title: "Leaderboard Referrer", File: "09-laporan.md"},
			},
		},
		{
			ID:    "pengumuman",
			Title: "Pengumuman",
			Icon:  iconAnnouncements,
			Sections: []Section{
				{ID: "buat-pengumuman", Title: "Buat & Kelola Pengumuman", File: "10-pengumuman.md"},
				{ID: "publish-pengumuman", Title: "Publikasi Pengumuman", File: "10-pengumuman.md"},
				{ID: "target-pengumuman", Title: "Target Penerima", File: "10-pengumuman.md"},
			},
		},
		{
			ID:    "pengaturan",
			Title: "Pengaturan Sistem",
			Icon:  iconSettings,
			Sections: []Section{
				{ID: "pengaturan-user", Title: "Manajemen User & Role", File: "11-pengaturan.md"},
				{ID: "pengaturan-prodi", Title: "Program Studi", File: "11-pengaturan.md"},
				{ID: "pengaturan-kategori", Title: "Kategori Interaksi", File: "11-pengaturan.md"},
				{ID: "pengaturan-hambatan", Title: "Hambatan Pendaftaran", File: "11-pengaturan.md"},
				{ID: "pengaturan-biaya", Title: "Struktur Biaya", File: "11-pengaturan.md"},
				{ID: "pengaturan-assignment", Title: "Algoritma Assignment", File: "11-pengaturan.md"},
				{ID: "pengaturan-lost-reason", Title: "Alasan Lost", File: "11-pengaturan.md"},
			},
		},
		{
			ID:    "keamanan",
			Title: "Keamanan & Enkripsi",
			Icon:  iconSecurity,
			Sections: []Section{
				{ID: "enkripsi-data", Title: "Enkripsi Data Sensitif", File: "12-keamanan.md"},
				{ID: "autentikasi", Title: "Autentikasi & Otorisasi", File: "12-keamanan.md"},
				{ID: "csrf-protection", Title: "CSRF Protection", File: "12-keamanan.md"},
				{ID: "role-permission", Title: "Role & Permission", File: "12-keamanan.md"},
			},
		},
	}
}

// --- Markdown processing ---

// extractSectionContent extracts content for a section from a markdown file.
// If the section title matches an H2, return that H2's content.
// If no H2 matches, return intro content excluding H2s that have their own section definitions.
func extractSectionContent(markdown string, sectionTitle string, siblingTitles []string) string {
	lines := strings.Split(markdown, "\n")

	// Try to find matching H2
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			h2Title := strings.TrimPrefix(line, "## ")
			if titlesMatch(h2Title, sectionTitle) {
				return extractH2Content(lines, i)
			}
		}
	}

	// No H2 match — return intro content excluding sibling H2s
	return extractIntroContent(lines, siblingTitles)
}

func extractH2Content(lines []string, startIdx int) string {
	var result []string
	for i := startIdx + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			break
		}
		result = append(result, lines[i])
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

func extractIntroContent(lines []string, siblingTitles []string) string {
	var result []string
	skip := false
	pastH1 := false

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "## ") {
			pastH1 = true
			continue
		}
		if !pastH1 {
			continue
		}

		if strings.HasPrefix(line, "## ") {
			h2Title := strings.TrimPrefix(line, "## ")
			skip = false
			for _, sibling := range siblingTitles {
				if titlesMatch(h2Title, sibling) {
					skip = true
					break
				}
			}
			if !skip {
				result = append(result, line)
			}
			continue
		}

		if !skip {
			result = append(result, line)
		}
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

func titlesMatch(a, b string) bool {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	if a == b {
		return true
	}
	if strings.Contains(a, b) || strings.Contains(b, a) {
		return true
	}

	// Keyword overlap: all words (len>=4) from shorter must appear in longer
	shorter, longer := a, b
	if len(a) > len(b) {
		shorter, longer = b, a
	}
	words := strings.Fields(shorter)
	matched := 0
	total := 0
	for _, w := range words {
		if len(w) >= 4 {
			total++
			if strings.Contains(longer, w) {
				matched++
			}
		}
	}
	return total > 0 && matched == total
}

// convertMarkdownToHTML converts markdown to HTML.
// Handles: headings, bold, italic, code blocks, inline code, tables, images, links, lists, blockquotes, horizontal rules.
func convertMarkdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var out []string
	inCodeBlock := false
	inTable := false
	inList := false
	listType := "" // "ul" or "ol"

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Code blocks
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				out = append(out, "</code></pre>")
				inCodeBlock = false
			} else {
				lang := strings.TrimPrefix(line, "```")
				cls := ""
				if lang != "" {
					cls = fmt.Sprintf(` class="language-%s"`, html.EscapeString(lang))
				}
				out = append(out, fmt.Sprintf("<pre><code%s>", cls))
				inCodeBlock = true
			}
			continue
		}
		if inCodeBlock {
			out = append(out, html.EscapeString(line))
			continue
		}

		// Close table if needed
		if inTable && !strings.HasPrefix(strings.TrimSpace(line), "|") {
			out = append(out, "</tbody></table>")
			inTable = false
		}

		// Close list if needed
		if inList && !isListItem(line) && strings.TrimSpace(line) != "" {
			out = append(out, fmt.Sprintf("</%s>", listType))
			inList = false
		}

		trimmed := strings.TrimSpace(line)

		// Empty line
		if trimmed == "" {
			if inList {
				out = append(out, fmt.Sprintf("</%s>", listType))
				inList = false
			}
			continue
		}

		// Horizontal rule
		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			out = append(out, "<hr>")
			continue
		}

		// Headings
		if strings.HasPrefix(trimmed, "### ") {
			out = append(out, fmt.Sprintf("<h3>%s</h3>", inlineFormat(strings.TrimPrefix(trimmed, "### "))))
			continue
		}
		if strings.HasPrefix(trimmed, "#### ") {
			out = append(out, fmt.Sprintf("<h4>%s</h4>", inlineFormat(strings.TrimPrefix(trimmed, "#### "))))
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			out = append(out, fmt.Sprintf("<h2>%s</h2>", inlineFormat(strings.TrimPrefix(trimmed, "## "))))
			continue
		}

		// Blockquote
		if strings.HasPrefix(trimmed, "> ") {
			out = append(out, fmt.Sprintf("<blockquote>%s</blockquote>", inlineFormat(strings.TrimPrefix(trimmed, "> "))))
			continue
		}

		// Table
		if strings.HasPrefix(trimmed, "|") {
			if !inTable {
				out = append(out, "<table>")
				// Header row
				cells := parseTableRow(trimmed)
				out = append(out, "<thead><tr>")
				for _, c := range cells {
					out = append(out, fmt.Sprintf("<th>%s</th>", inlineFormat(c)))
				}
				out = append(out, "</tr></thead>")
				// Skip separator row
				if i+1 < len(lines) && isTableSeparator(lines[i+1]) {
					i++
				}
				out = append(out, "<tbody>")
				inTable = true
			} else {
				cells := parseTableRow(trimmed)
				out = append(out, "<tr>")
				for _, c := range cells {
					out = append(out, fmt.Sprintf("<td>%s</td>", inlineFormat(c)))
				}
				out = append(out, "</tr>")
			}
			continue
		}

		// Lists
		if isListItem(line) {
			if !inList {
				if isOrderedListItem(line) {
					listType = "ol"
				} else {
					listType = "ul"
				}
				out = append(out, fmt.Sprintf("<%s>", listType))
				inList = true
			}
			content := extractListContent(line)
			out = append(out, fmt.Sprintf("<li>%s</li>", inlineFormat(content)))
			continue
		}

		// Image
		if imgMatch := regexp.MustCompile(`^!\[([^\]]*)\]\(([^)]+)\)$`).FindStringSubmatch(trimmed); imgMatch != nil {
			out = append(out, fmt.Sprintf(`<figure><img src="%s" alt="%s"><figcaption>%s</figcaption></figure>`,
				html.EscapeString(imgMatch[2]), html.EscapeString(imgMatch[1]), html.EscapeString(imgMatch[1])))
			continue
		}

		// Paragraph
		out = append(out, fmt.Sprintf("<p>%s</p>", inlineFormat(trimmed)))
	}

	if inTable {
		out = append(out, "</tbody></table>")
	}
	if inList {
		out = append(out, fmt.Sprintf("</%s>", listType))
	}
	if inCodeBlock {
		out = append(out, "</code></pre>")
	}

	return strings.Join(out, "\n")
}

func inlineFormat(s string) string {
	// Inline code
	s = regexp.MustCompile("`([^`]+)`").ReplaceAllString(s, "<code>$1</code>")
	// Bold
	s = regexp.MustCompile(`\*\*([^*]+)\*\*`).ReplaceAllString(s, "<strong>$1</strong>")
	// Italic
	s = regexp.MustCompile(`\*([^*]+)\*`).ReplaceAllString(s, "<em>$1</em>")
	// Links
	s = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`).ReplaceAllString(s, `<a href="$2">$1</a>`)
	// Images (inline)
	s = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`).ReplaceAllStringFunc(s, func(m string) string {
		parts := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`).FindStringSubmatch(m)
		return fmt.Sprintf(`<img src="%s" alt="%s">`, html.EscapeString(parts[2]), html.EscapeString(parts[1]))
	})
	return s
}

func parseTableRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.Trim(line, "|")
	parts := strings.Split(line, "|")
	var cells []string
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	return cells
}

func isTableSeparator(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "|") && strings.Contains(trimmed, "---")
}

func isListItem(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || isOrderedListItem(line)
}

func isOrderedListItem(line string) bool {
	trimmed := strings.TrimSpace(line)
	return regexp.MustCompile(`^\d+\. `).MatchString(trimmed)
}

func extractListContent(line string) string {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "- ") {
		return strings.TrimPrefix(trimmed, "- ")
	}
	if strings.HasPrefix(trimmed, "* ") {
		return strings.TrimPrefix(trimmed, "* ")
	}
	return regexp.MustCompile(`^\d+\. `).ReplaceAllString(trimmed, "")
}

// --- HTML Generation ---

func generateHTML(groups []SectionGroup, manualDir string) string {
	var toc strings.Builder
	var content strings.Builder

	groupIdx := 0
	for _, g := range groups {
		groupIdx++
		toc.WriteString(fmt.Sprintf(`<div class="toc-group">
			<button class="toc-group-header" onclick="toggleGroup('%s')">
				<span class="toc-group-icon">%s</span>
				<span class="toc-group-title">%d. %s</span>
				<span class="toc-chevron" id="chevron-%s">&#9654;</span>
			</button>
			<div class="toc-group-items" id="group-%s">`, g.ID, g.Icon, groupIdx, g.Title, g.ID, g.ID))

		// Collect sibling titles for this file
		fileSiblings := map[string][]string{}
		for _, s := range g.Sections {
			fileSiblings[s.File] = append(fileSiblings[s.File], s.Title)
		}

		sectionIdx := 0
		for _, s := range g.Sections {
			sectionIdx++
			num := fmt.Sprintf("%d.%d", groupIdx, sectionIdx)

			toc.WriteString(fmt.Sprintf(`<a class="toc-item" href="#%s" onclick="navigateTo('%s','%s')">
				<span class="toc-num">%s</span> %s
			</a>`, s.ID, s.ID, g.ID, num, s.Title))

			// Read markdown and extract section content
			mdPath := filepath.Join(manualDir, s.File)
			mdBytes, err := os.ReadFile(mdPath)
			var htmlContent string
			if err != nil {
				htmlContent = fmt.Sprintf(`<p class="placeholder">Konten belum tersedia. File: <code>%s</code></p>`, s.File)
			} else {
				siblings := fileSiblings[s.File]
				extracted := extractSectionContent(string(mdBytes), s.Title, siblings)
				if extracted == "" {
					htmlContent = fmt.Sprintf(`<p class="placeholder">Section "%s" belum ditulis di <code>%s</code></p>`, s.Title, s.File)
				} else {
					htmlContent = convertMarkdownToHTML(extracted)
				}
			}

			// Screenshot gallery
			screenshotHTML := ""
			if len(s.Screenshots) > 0 {
				var gallery []string
				gallery = append(gallery, `<div class="screenshot-gallery">`)
				for _, ssID := range s.Screenshots {
					ssPath := fmt.Sprintf("screenshots/%s.png", ssID)
					fullPath := filepath.Join(manualDir, ssPath)
					if _, err := os.Stat(fullPath); err == nil {
						gallery = append(gallery, fmt.Sprintf(
							`<figure class="screenshot"><a href="%s" target="_blank"><img src="%s" alt="%s"></a><figcaption>%s</figcaption></figure>`,
							ssPath, ssPath, ssID, ssID))
					} else {
						gallery = append(gallery, fmt.Sprintf(
							`<figure class="screenshot placeholder-screenshot"><div class="screenshot-placeholder">📷 %s</div></figure>`, ssID))
					}
				}
				gallery = append(gallery, `</div>`)
				screenshotHTML = strings.Join(gallery, "\n")
			}

			content.WriteString(fmt.Sprintf(`<section id="%s" class="manual-section">
				<div class="section-header">
					<span class="section-num">%s</span>
					<h2>%s</h2>
				</div>
				<div class="section-content">
					%s
					%s
				</div>
			</section>`, s.ID, num, s.Title, htmlContent, screenshotHTML))
		}

		toc.WriteString(`</div></div>`)
	}

	buildDate := time.Now().Format("2 January 2006")

	return fmt.Sprintf(htmlTemplate, buildDate, toc.String(), content.String(), buildDate)
}

// GenerateManual reads markdown sources and produces a single-page HTML manual.
func GenerateManual(manualDir, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy screenshots directory if exists
	screenshotSrc := filepath.Join(manualDir, "screenshots")
	screenshotDst := filepath.Join(outputDir, "screenshots")
	if info, err := os.Stat(screenshotSrc); err == nil && info.IsDir() {
		copyDir(screenshotSrc, screenshotDst)
	}

	groups := getSectionGroups()
	htmlOut := generateHTML(groups, manualDir)

	outputPath := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(outputPath, []byte(htmlOut), 0o644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	log.Printf("User manual generated: %s", outputPath)

	// Print section statistics
	totalSections := 0
	filledSections := 0
	for _, g := range groups {
		for _, s := range g.Sections {
			totalSections++
			mdPath := filepath.Join(manualDir, s.File)
			if mdBytes, err := os.ReadFile(mdPath); err == nil {
				fileSiblings := map[string][]string{}
				for _, ss := range g.Sections {
					fileSiblings[ss.File] = append(fileSiblings[ss.File], ss.Title)
				}
				extracted := extractSectionContent(string(mdBytes), s.Title, fileSiblings[s.File])
				if extracted != "" {
					filledSections++
				}
			}
		}
	}
	log.Printf("Sections: %d/%d filled", filledSections, totalSections)
	return nil
}

func copyDir(src, dst string) {
	os.MkdirAll(dst, 0o755)
	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			copyDir(srcPath, dstPath)
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				continue
			}
			os.WriteFile(dstPath, data, 0o644)
		}
	}
}

// --- SVG Icons ---

const (
	iconDashboard    = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>`
	iconRegistration = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>`
	iconCandidates   = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`
	iconInteraction  = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>`
	iconDocuments    = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>`
	iconFinance      = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>`
	iconMarketing    = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>`
	iconPortal       = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>`
	iconReports      = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>`
	iconAnnouncements = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>`
	iconSettings     = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>`
	iconSecurity     = `<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>`
)

// --- HTML Template ---

var htmlTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>STMIK Tazkia - Panduan Pengguna Sistem Penerimaan Mahasiswa Baru</title>
<style>
:root {
  --primary-600: #194189;
  --primary-700: #143370;
  --primary-50: #eef2f9;
  --secondary-500: #EE7B1D;
  --secondary-50: #fef6ee;
  --gray-50: #f9fafb;
  --gray-100: #f3f4f6;
  --gray-200: #e5e7eb;
  --gray-300: #d1d5db;
  --gray-500: #6b7280;
  --gray-600: #4b5563;
  --gray-700: #374151;
  --gray-800: #1f2937;
  --gray-900: #111827;
  --radius: 8px;
}

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  color: var(--gray-800);
  line-height: 1.7;
  background: var(--gray-50);
}

/* Header */
.header {
  position: sticky; top: 0; z-index: 100;
  background: linear-gradient(135deg, var(--primary-600), var(--primary-700));
  color: white; padding: 16px 24px;
  display: flex; justify-content: space-between; align-items: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}
.header h1 { font-size: 18px; font-weight: 600; }
.header .meta { font-size: 13px; opacity: 0.85; }
.header button {
  background: rgba(255,255,255,0.15); border: 1px solid rgba(255,255,255,0.3);
  color: white; padding: 6px 16px; border-radius: var(--radius);
  cursor: pointer; font-size: 13px;
}
.header button:hover { background: rgba(255,255,255,0.25); }

/* Layout */
.layout { display: flex; max-width: 1400px; margin: 0 auto; }

/* Sidebar */
.sidebar {
  width: 320px; min-width: 320px;
  position: sticky; top: 64px; height: calc(100vh - 64px);
  overflow-y: auto; padding: 16px;
  border-right: 1px solid var(--gray-200);
  background: white;
}
.toc-group { margin-bottom: 4px; }
.toc-group-header {
  display: flex; align-items: center; gap: 8px;
  width: 100%%; padding: 8px 12px; border: none;
  background: none; cursor: pointer; border-radius: var(--radius);
  font-size: 14px; font-weight: 600; color: var(--gray-700);
  text-align: left;
}
.toc-group-header:hover { background: var(--gray-100); }
.toc-group-icon { display: flex; color: var(--primary-600); }
.toc-chevron {
  margin-left: auto; font-size: 10px; transition: transform 0.2s;
  color: var(--gray-500);
}
.toc-chevron.open { transform: rotate(90deg); }
.toc-group-items {
  max-height: 0; overflow: hidden; transition: max-height 0.3s ease;
}
.toc-group-items.open { max-height: 800px; }
.toc-item {
  display: block; padding: 5px 12px 5px 42px;
  text-decoration: none; color: var(--gray-600);
  font-size: 13px; border-radius: var(--radius);
}
.toc-item:hover { background: var(--primary-50); color: var(--primary-600); }
.toc-item.active { background: var(--primary-50); color: var(--primary-600); font-weight: 600; }
.toc-num { color: var(--gray-500); font-size: 12px; min-width: 28px; display: inline-block; }

/* Main content */
.main { flex: 1; padding: 24px; max-width: 900px; }

.manual-section { margin-bottom: 32px; background: white; border-radius: var(--radius); box-shadow: 0 1px 3px rgba(0,0,0,0.08); overflow: hidden; }
.section-header {
  background: linear-gradient(135deg, var(--primary-600), var(--primary-700));
  color: white; padding: 16px 24px;
  display: flex; align-items: center; gap: 12px;
}
.section-num {
  background: rgba(255,255,255,0.2); padding: 2px 10px;
  border-radius: 12px; font-size: 13px; font-weight: 600;
}
.section-header h2 { font-size: 18px; font-weight: 600; margin: 0; }
.section-content { padding: 24px; }
.section-content h2 { font-size: 18px; margin: 24px 0 12px; color: var(--primary-700); border-bottom: 2px solid var(--primary-50); padding-bottom: 4px; }
.section-content h3 { font-size: 16px; margin: 20px 0 8px; color: var(--gray-800); }
.section-content h4 { font-size: 14px; margin: 16px 0 6px; color: var(--gray-700); }
.section-content p { margin: 8px 0; }
.section-content ul, .section-content ol { margin: 8px 0; padding-left: 24px; }
.section-content li { margin: 4px 0; }
.section-content blockquote {
  border-left: 4px solid var(--secondary-500);
  background: var(--secondary-50); padding: 12px 16px;
  margin: 12px 0; border-radius: 0 var(--radius) var(--radius) 0;
}
.section-content table { width: 100%%; border-collapse: collapse; margin: 12px 0; font-size: 14px; }
.section-content th { background: var(--primary-50); text-align: left; padding: 8px 12px; border: 1px solid var(--gray-200); font-weight: 600; }
.section-content td { padding: 8px 12px; border: 1px solid var(--gray-200); }
.section-content tr:nth-child(even) { background: var(--gray-50); }
.section-content code { background: var(--gray-100); padding: 2px 6px; border-radius: 4px; font-size: 13px; }
.section-content pre { background: var(--gray-900); color: #e5e7eb; padding: 16px; border-radius: var(--radius); overflow-x: auto; margin: 12px 0; }
.section-content pre code { background: none; padding: 0; color: inherit; }
.section-content img { max-width: 100%%; border-radius: var(--radius); border: 1px solid var(--gray-200); margin: 8px 0; }
.section-content figure { margin: 16px 0; text-align: center; }
.section-content figcaption { font-size: 13px; color: var(--gray-500); margin-top: 4px; }
.section-content hr { border: none; border-top: 1px solid var(--gray-200); margin: 20px 0; }

.placeholder { color: var(--gray-500); font-style: italic; background: var(--gray-50); padding: 24px; border-radius: var(--radius); text-align: center; }

/* Screenshots */
.screenshot-gallery { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 16px; margin-top: 16px; }
.screenshot { margin: 0; }
.screenshot img { width: 100%%; cursor: pointer; transition: transform 0.2s; }
.screenshot img:hover { transform: scale(1.02); }
.screenshot-placeholder { background: var(--gray-100); padding: 40px; text-align: center; border-radius: var(--radius); color: var(--gray-500); }

/* Footer */
.footer { text-align: center; padding: 24px; color: var(--gray-500); font-size: 13px; border-top: 1px solid var(--gray-200); margin-top: 40px; }

/* Print */
@media print {
  .sidebar, .header button { display: none; }
  .header { position: static; }
  .layout { display: block; }
  .manual-section { break-inside: avoid; box-shadow: none; border: 1px solid var(--gray-300); }
}

/* Responsive */
@media (max-width: 900px) {
  .sidebar { display: none; }
  .main { padding: 16px; }
}
</style>
</head>
<body>

<div class="header">
  <div>
    <h1>Panduan Pengguna — Sistem Penerimaan Mahasiswa Baru STMIK Tazkia</h1>
    <div class="meta">Dibuat: %s</div>
  </div>
  <button onclick="window.print()">🖨 Cetak</button>
</div>

<div class="layout">
  <nav class="sidebar">%s</nav>
  <main class="main">%s</main>
</div>

<div class="footer">
  &copy; STMIK Tazkia — Panduan Pengguna Sistem Penerimaan Mahasiswa Baru — %s
</div>

<script>
function toggleGroup(id) {
  const items = document.getElementById('group-' + id);
  const chevron = document.getElementById('chevron-' + id);
  items.classList.toggle('open');
  chevron.classList.toggle('open');
}

function navigateTo(sectionId, groupId) {
  // Expand the group
  const items = document.getElementById('group-' + groupId);
  const chevron = document.getElementById('chevron-' + groupId);
  items.classList.add('open');
  chevron.classList.add('open');

  // Update active state
  document.querySelectorAll('.toc-item').forEach(el => el.classList.remove('active'));
  event.currentTarget.classList.add('active');

  // Scroll to section
  const section = document.getElementById(sectionId);
  if (section) {
    section.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }

  // Update URL hash
  history.replaceState(null, null, '#' + sectionId);
}

// Handle page load with hash
window.addEventListener('load', function() {
  const hash = window.location.hash.slice(1);
  if (hash) {
    const section = document.getElementById(hash);
    if (section) {
      section.scrollIntoView({ block: 'start' });
      // Find and expand the right group
      document.querySelectorAll('.toc-item').forEach(el => {
        if (el.getAttribute('href') === '#' + hash) {
          el.classList.add('active');
        }
      });
    }
  }
  // Expand first group by default
  const firstGroup = document.querySelector('.toc-group-items');
  const firstChevron = document.querySelector('.toc-chevron');
  if (firstGroup) { firstGroup.classList.add('open'); }
  if (firstChevron) { firstChevron.classList.add('open'); }
});
</script>

</body>
</html>`
