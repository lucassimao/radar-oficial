package diarios

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"radaroficial.app/internal/config"
	"radaroficial.app/internal/model"
	"radaroficial.app/internal/storage"
)

type diarioAPIResponse struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

var diarioURLBase = "https://www.diario.pi.gov.br"
var hrefRegexp = regexp.MustCompile(`href="(.+?\.pdf)"`)

// Institution IDs for database references
const (
	InstitutionIDGovernoPiaui    = 1 // ID for Governo do Estado do Piauí
	InstitutionIDMunicipiosPiaui = 2 // ID for Diário dos Municípios do Piauí
)

// FetchGovernoPiauiDiarios fetches diarios from the Governo do Piauí website, uploads them to storage, and inserts them into the database
func FetchGovernoPiauiDiarios(ctx context.Context, date time.Time, uploader *storage.SpacesUploader, service *DiarioService) ([]*model.Diario, error) {
	form := url.Values{}
	form.Set("draw", "3")
	form.Set("start", "0")
	form.Set("length", "50")
	form.Set("filter_data", date.Format("2006-01-02"))

	for i := 0; i <= 2; i++ {
		prefix := fmt.Sprintf("columns[%d]", i)
		form.Set(prefix+"[data]", fmt.Sprintf("%d", i))
		form.Set(prefix+"[searchable]", "true")
		form.Set(prefix+"[orderable]", "true")
		form.Set(prefix+"[search][value]", "")
		form.Set(prefix+"[search][regex]", "false")
	}
	form.Set("order[0][column]", "0")
	form.Set("order[0][dir]", "asc")
	form.Set("search[value]", "")
	form.Set("search[regex]", "false")
	form.Set("filter_numero", "")

	req, err := http.NewRequestWithContext(ctx, "POST", diarioURLBase+"/doe/Api/listardiarios.json", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", diarioURLBase+"/doe/")
	req.Header.Set("Origin", diarioURLBase)

	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var parsed diarioAPIResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	var diarios []*model.Diario
	var processedCount int

	for _, row := range parsed.Data {
		if len(row) < 4 {
			continue
		}

		match := hrefRegexp.FindStringSubmatch(row[0])
		if len(match) < 2 {
			continue
		}

		publishedAt, _ := time.Parse("02/01/2006", row[2])
		lastModifiedAt, _ := time.Parse("02/01/2006 15:04:05", row[3])
		desc := strings.TrimSpace(row[1])

		// Check if this diario already exists in our database
		exists, err := service.DiarioExists(ctx, InstitutionIDGovernoPiaui, desc)
		if err != nil {
			log.Printf("⚠️ Error checking if diario exists: %v", err)
		}

		if exists {
			log.Printf("✅ Skipping already downloaded diario: %s", desc)
			continue
		}

		// If we get here, this is a new diario that needs downloading
		log.Printf("📥 Downloading new diario: %s", desc)

		rawPDFPath := strings.ReplaceAll(match[1], "..", "")
		pdfURL := diarioURLBase + rawPDFPath

		resp, err := http.Get(pdfURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("❌ Failed to download PDF from %s: %v", pdfURL, err)
			continue
		}
		defer resp.Body.Close()

		pdfContent, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("❌ Failed to read PDF content: %v", err)
			continue
		}

		sanitized := sanitizeDescription(desc)
		// Construct object path: e.g., "2025/04/DOEPI_71_2025.pdf"
		filename := filepath.Base(rawPDFPath)
		objectPath := fmt.Sprintf("governo-pi/%d/%02d/%s_%s", date.Year(), date.Month(), sanitized, filename)

		err = uploader.UploadFile(ctx, objectPath, bytes.NewReader(pdfContent), int64(len(pdfContent)), "application/pdf")
		if err != nil {
			log.Printf("❌ Failed to upload PDF to storage: %v", err)
			continue
		}

		log.Printf("✅ Successfully uploaded %s", objectPath)

		// Create diario object
		diario := &model.Diario{
			InstitutionID:  InstitutionIDGovernoPiaui,
			SourceURL:      fmt.Sprintf("https://%s.%s/%s", uploader.Bucket, os.Getenv("DO_SPACES_ENDPOINT"), objectPath),
			Description:    &desc,
			PublishedAt:    &publishedAt,
			LastModifiedAt: &lastModifiedAt,
		}

		// Insert directly into database
		if err := service.Insert(ctx, diario); err != nil {
			log.Printf("⚠️ Failed to insert diário for %s: %v", diario.SourceURL, err)
		} else {
			log.Printf("✅ Inserted diário %s", diario.SourceURL)
			processedCount++
		}

		diarios = append(diarios, diario)
	}

	log.Printf("📊 Successfully processed %d diários from Governo do Piauí", processedCount)
	return diarios, nil
}

// CurrentEditionResponse represents the data structure returned by the Supabase API
type CurrentEditionResponse struct {
	Edition int    `json:"edicao"`
	Date    string `json:"data"`
	ID      int    `json:"id"`
}

// FetchDiarioDosMunicipiosPiaui fetches the latest edition of Diário dos Municípios using go-rod
// and directly inserts it into the database
func FetchDiarioDosMunicipiosPiaui(ctx context.Context, uploader *storage.SpacesUploader, service *DiarioService) ([]*model.Diario, error) {
	log.Printf("📥 Fetching Diário dos Municípios using go-rod...")

	// Launch a new browser with headless mode and no-sandbox for Linux compatibility
	l := launcher.New().
		Headless(true).
		Set("no-sandbox", "").
		Set("disable-setuid-sandbox", "")

	if config.Env() == "production" {
		l.Bin("/usr/bin/google-chrome")
	}
	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// Create a new page and navigate to the site
	page := browser.MustPage("https://www.diarioficialdosmunicipios.org/edicao_atual.html")

	// Wait for the page to load
	page.MustWaitLoad()
	log.Printf("✅ Page loaded successfully")

	// Try to find the edition information with different strategies
	var editionText string

	// First attempt: Look for element with class .title-edi
	titleElement, err := page.Element("span#newDesc > h1")
	if err == nil {
		editionText = titleElement.MustText()
		log.Printf("✅ Found edition text using span#newDesc > h1: %s", editionText)
	}

	// Try different regex patterns to match edition info
	var matches []string

	// Try with comma pattern first: "Edição 5302, 16/04/2025"
	editionRegex := regexp.MustCompile(`Edição\s+(\d+),\s+(\d{2}/\d{2}/\d{4})`)
	matches = editionRegex.FindStringSubmatch(editionText)

	log.Printf("📄 Parsed edition info: %v", matches)

	if len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse edition information from text: %s", editionText)
	}

	// Extract the edition number and date
	editionNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse edition number: %w", err)
	}

	dateStr := matches[2]

	// Parse the publication date
	publishDate, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		log.Printf("⚠️ Error parsing publication date %s, using current date", dateStr)
		publishDate = time.Now()
	}

	// Create a descriptive title
	description := fmt.Sprintf("Edição %d (%s)", editionNumber, dateStr)

	// Check if this diario already exists in our database
	exists, err := service.DiarioExists(ctx, InstitutionIDMunicipiosPiaui, description)
	if err != nil {
		log.Printf("⚠️ Error checking if diario exists: %v", err)
	}

	if exists {
		log.Printf("✅ Skipping already downloaded diario: %s", description)
		return nil, nil
	}

	var pdfURL string

	// Try to find a link containing "Baixar Edição"
	downloadLinks, err := page.Elements("a")
	if err == nil {
		for _, link := range downloadLinks {
			text, err := link.Text()
			if err == nil && strings.Contains(text, "Baixar Edição") {
				pdfURL = link.MustProperty("href").String()
				log.Printf("✅ Found download link with text 'Baixar Edição': %s", pdfURL)
				break
			}
		}
	}

	if pdfURL == "" {
		return nil, fmt.Errorf("could not find any download link for the PDF")
	}

	log.Printf("📄 Found PDF URL: %s", pdfURL)
	// kill the browser to save memory
	browser.MustClose()

	// Get the PDF file directly using a standard HTTP request instead of the browser
	client := &http.Client{
		Timeout: 1 * time.Hour,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", pdfURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://www.diarioficialdosmunicipios.org/edicao_atual.html")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download PDF, HTTP status: %d", resp.StatusCode)
	}

	// Read the PDF content
	pdfContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF content: %w", err)
	}

	log.Printf("📥 Downloaded PDF: %d bytes", len(pdfContent))

	// Format upload path
	objectPath := fmt.Sprintf("municipios-pi/%d/%02d/edicao_%d_%s.pdf",
		publishDate.Year(),
		publishDate.Month(),
		editionNumber,
		publishDate.Format("2006-01-02"))

	// Upload to storage
	err = uploader.UploadFile(ctx, objectPath, bytes.NewReader(pdfContent), int64(len(pdfContent)), "application/pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to upload PDF: %w", err)
	}

	log.Printf("✅ Successfully uploaded %s", objectPath)

	// Create the Diario object
	lastModifiedAt := time.Now()

	diario := &model.Diario{
		InstitutionID:  InstitutionIDMunicipiosPiaui, // ID for Diário dos Municípios
		SourceURL:      fmt.Sprintf("https://%s.%s/%s", uploader.Bucket, os.Getenv("DO_SPACES_ENDPOINT"), objectPath),
		Description:    &description,
		PublishedAt:    &publishDate,
		LastModifiedAt: &lastModifiedAt,
	}

	// Insert directly into database
	if err := service.Insert(ctx, diario); err != nil {
		log.Printf("⚠️ Failed to insert diário for %s: %v", diario.SourceURL, err)
	} else {
		log.Printf("✅ Inserted diário %s", diario.SourceURL)
	}

	return []*model.Diario{diario}, nil
}

// sanitizeDescription returns a safe string for filenames
func sanitizeDescription(s string) string {
	// Convert to ASCII (basic filter only — no accent removal here)
	var b strings.Builder
	for _, r := range s {
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_') {
			b.WriteRune(r)
		}
	}

	// Replace spaces with underscores
	safe := strings.ReplaceAll(b.String(), " ", "_")

	// Collapse multiple underscores
	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_")

	return safe
}
