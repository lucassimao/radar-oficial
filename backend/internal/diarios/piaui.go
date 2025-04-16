package diarios

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"radaroficial.app/internal/model"
)

type diarioAPIResponse struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

var diarioURLBase = "https://www.diario.pi.gov.br"

var hrefRegexp = regexp.MustCompile(`href="(.+?\.pdf)"`)

func FetchGovernoPiauiDiarios(ctx context.Context, forDate time.Time) ([]*model.Diario, error) {
	form := url.Values{}
	form.Set("draw", "3")
	form.Set("start", "0")
	form.Set("length", "50")
	form.Set("filter_data", forDate.Format("2006-01-02")) // "2006-01-02" is the layout string to get YYYY-MM-DD

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

	req, err := http.NewRequestWithContext(ctx, "POST", diarioURLBase+"/doe/Api/listardiarios.json", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", diarioURLBase+"/doe/")
	req.Header.Set("Origin", diarioURLBase)

	client := &http.Client{Timeout: 15 * time.Second}
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

		// Extract .pdf path from HTML
		match := hrefRegexp.FindStringSubmatch(row[0])
		if len(match) < 2 {
			continue
		}
		pdfURL := diarioURLBase + strings.ReplaceAll(match[1], "..", "")

		// Parse published and last modified dates
		publishedAt, _ := time.Parse("02/01/2006", row[2])
		lastModifiedAt, _ := time.Parse("02/01/2006 15:04:05", row[3])

		desc := strings.TrimSpace(row[1])

		diarios = append(diarios, &model.Diario{
			InstitutionID:  1,
			SourceURL:      pdfURL,
			Description:    &desc,
			PublishedAt:    &publishedAt,
			LastModifiedAt: &lastModifiedAt,
		})
	}

	return diarios, nil
}
