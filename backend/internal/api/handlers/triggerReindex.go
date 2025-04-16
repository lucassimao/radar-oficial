package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ReindexHandler struct{}

func NewReindexHandler() *ReindexHandler {
	return &ReindexHandler{}
}

func (h *ReindexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := triggerReindex(r.Context(), "a4fc4135-1a22-11f0-bf8f-4e013e2ddde4", []string{})
	if err != nil {
		log.Printf("❌ Failed to trigger reindex: %v", err)
		http.Error(w, "Failed to trigger reindex.", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func triggerReindex(ctx context.Context, kbUUID string, dataSourceUUIDs []string) error {
	token := os.Getenv("DO_API_KEY")

	body := map[string]interface{}{
		"knowledge_base_uuid": kbUUID,
	}
	if len(dataSourceUUIDs) > 0 {
		body["data_source_uuids"] = dataSourceUUIDs
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.digitalocean.com/v2/gen-ai/indexing_jobs", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("❌ Reindex request failed with status: %s", resp.Status)
		return err
	}

	log.Println("✅ Reindex triggered successfully")
	return nil
}
