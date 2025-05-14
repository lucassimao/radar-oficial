package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/chat"
	"radaroficial.app/internal/diarios"
)

// ChatHandler handles incoming webhook requests from WhatsApp
type ChatHandler struct {
	diarioService *diarios.DiarioService
	chatService   *chat.ChatService
	db            *pgxpool.Pool
}

func NewChatHandler(db *pgxpool.Pool) *ChatHandler {

	return &ChatHandler{
		diarioService: diarios.NewInstitutionService(db),
		chatService:   chat.NewChatService(),
		db:            db,
	}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.chatCompletion(w, r)
}

func (h *ChatHandler) chatCompletion(w http.ResponseWriter, r *http.Request) {

	queryValues := r.URL.Query()

	// if somehow the state is not present, trigger the select-diario-state tool
	if !queryValues.Has("state") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(selectDiarioStateToolCall())
		return
	}

	state := queryValues.Get("state")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading request body: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	var message MessageSet
	if err := json.Unmarshal(body, &message); err != nil {
		http.Error(w, "Error parsing chat request", http.StatusBadRequest)
		return
	}

	lastMessage := message.Messages[len(message.Messages)-1]

	agentResponse, err := h.sendMessageToAIAgent(state, lastMessage.Content[0].Text)

	if err != nil {
		log.Printf("❌ Failed to process chat completion: %v", err)
		http.Error(w, "Failed to process chat completion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"text": agentResponse,
	})

}

// sendMessageToAIAgent sends a message to the Digital Ocean AI agent and returns the response
func (h *ChatHandler) sendMessageToAIAgent(state, message string) (string, error) {

	url, agentAccessKey, err := getAgentSettings(state)
	if err != nil {
		return "", err
	}

	// Construct the request payload
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	payload := struct {
		Messages              []Message `json:"messages"`
		Stream                bool      `json:"stream"`
		IncludeFunctionsInfo  bool      `json:"include_functions_info"`
		IncludeRetrievalInfo  bool      `json:"include_retrieval_info"`
		IncludeGuardrailsInfo bool      `json:"include_guardrails_info"`
	}{
		Messages: []Message{
			{
				Role:    "user",
				Content: message,
			},
		},
		Stream:                false,
		IncludeFunctionsInfo:  false,
		IncludeRetrievalInfo:  false,
		IncludeGuardrailsInfo: false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+agentAccessKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("❌ AI agent API error (status %d): %s", resp.StatusCode, string(responseBody))
		return "", fmt.Errorf("AI agent API error: %d", resp.StatusCode)
	}

	// Parse the response
	type AIResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	var aiResponse AIResponse
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(responseBody, &aiResponse); err != nil {
		return "", err
	}

	// Extract the response message
	if len(aiResponse.Choices) == 0 {
		return "", fmt.Errorf("no response from AI agent")
	}

	// log.Printf("✅ Got response from AI agent: %s", aiResponse.Choices[0].Message.Content)
	return aiResponse.Choices[0].Message.Content, nil
}

type MessageSet struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ID          string        `json:"id"`
	CreatedAt   string        `json:"createdAt"` // Use time.Time if you want to parse the timestamp
	Role        string        `json:"role"`
	Content     []Content     `json:"content"`
	Attachments []interface{} `json:"attachments,omitempty"` // Or define a struct if needed
	Metadata    Metadata      `json:"metadata"`
	Status      *Status       `json:"status,omitempty"` // Present only for assistant messages
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Metadata struct {
	Custom              map[string]interface{} `json:"custom"`
	UnstableAnnotations []interface{}          `json:"unstable_annotations,omitempty"`
	UnstableData        []interface{}          `json:"unstable_data,omitempty"`
	Steps               []interface{}          `json:"steps,omitempty"`
}

type Status struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Error  string `json:"error"`
}

// Get endpoint and access key from environment variables
func getAgentSettings(state string) (url, accessKey string, err error) {

	var agentEndpoint string

	switch state {
	case "PI":
		agentEndpoint = os.Getenv("DO_AGENT_PIAUI_URL")
		if agentEndpoint == "" {
			return "", "", fmt.Errorf("DO_AGENT_PIAUI_URL environment variable not set")
		}

		accessKey = os.Getenv("DO_AGENT_PIAUI_ACCESS_KEY")
		if accessKey == "" {
			return "", "", fmt.Errorf("DO_AGENT_PIAUI_ACCESS_KEY environment variable not set")
		}
	default:
		return "", "", fmt.Errorf("Invalid diario state: %s", state)

	}

	url = fmt.Sprintf("%s/api/v1/chat/completions", agentEndpoint)
	return url, accessKey, nil

}

func selectDiarioStateToolCall() map[string]any {
	return map[string]any{
		"content": []any{
			map[string]any{
				"type":       "tool-call",
				"toolName":   "select-diario-state",
				"toolCallId": fmt.Sprintf("%d", time.Now().UnixMilli()),
				"argsText":   "",
				"args":       map[string]any{},
			},
		},
	}
}
