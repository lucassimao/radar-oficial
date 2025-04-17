package diarios

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

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
		exists, err := service.DiarioExists(ctx, 1, desc)
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

		diarios = append(diarios, &model.Diario{
			InstitutionID:  1,
			SourceURL:      fmt.Sprintf("https://%s.%s/%s", uploader.Bucket, os.Getenv("DO_SPACES_ENDPOINT"), objectPath),
			Description:    &desc,
			PublishedAt:    &publishedAt,
			LastModifiedAt: &lastModifiedAt,
		})
	}

	return diarios, nil
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
