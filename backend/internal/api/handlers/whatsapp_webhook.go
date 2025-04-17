package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// WhatsAppWebhookHandler handles incoming webhook requests from WhatsApp
type WhatsAppWebhookHandler struct {}

// NewWhatsAppWebhookHandler creates a new WhatsAppWebhookHandler
func NewWhatsAppWebhookHandler() *WhatsAppWebhookHandler {
	return &WhatsAppWebhookHandler{}
}

// WhatsAppMessage represents a simplified structure of an incoming WhatsApp message
type WhatsAppMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type string `json:"type"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// ServeHTTP handles both GET (verification) and POST (message webhook) requests
func (h *WhatsAppWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Print request details for debugging
	log.Printf("📝 WhatsApp Webhook Request: Method=%s, URL=%s, RemoteAddr=%s, Headers=%v", 
		r.Method, r.URL.String(), r.RemoteAddr, r.Header)
	
	switch r.Method {
	case http.MethodGet:
		// Handle verification request from WhatsApp
		h.handleVerification(w, r)
	case http.MethodPost:
		// Handle incoming messages
		h.handleWebhook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleVerification verifies the webhook with WhatsApp
func (h *WhatsAppWebhookHandler) handleVerification(w http.ResponseWriter, r *http.Request) {
	// WhatsApp sends a challenge parameter that we need to echo back
	challenge := r.URL.Query().Get("hub.challenge")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")

	// Replace with your actual verify token (should match what you set in WhatsApp)
	verifyToken := os.Getenv("WHATSAPP_WEBHOOK_TOKEN")

	if mode == "subscribe" && token == verifyToken {
		log.Println("✅ WhatsApp webhook verified")
		w.Write([]byte(challenge))
	} else {
		log.Println("❌ WhatsApp webhook verification failed")
		http.Error(w, "Verification failed", http.StatusForbidden)
	}
}

// handleWebhook processes incoming messages from WhatsApp
func (h *WhatsAppWebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading request body: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	
	// Log the raw request body for debugging
	log.Printf("📦 Raw WhatsApp webhook payload: %s", string(body))

	var message WhatsAppMessage
	if err := json.Unmarshal(body, &message); err != nil {
		log.Printf("❌ Error parsing webhook: %v", err)
		http.Error(w, "Error parsing webhook", http.StatusBadRequest)
		return
	}

	// Process the message
	log.Printf("✅ Received WhatsApp webhook")
	
	// Extract the sender info and message text (if present)
	for _, entry := range message.Entry {
		for _, change := range entry.Changes {
			if change.Field == "messages" {
				for _, msg := range change.Value.Messages {
					if msg.Type == "text" {
						senderID := msg.From
						messageText := msg.Text.Body
						
						log.Printf("📱 Message from %s: %s", senderID, messageText)
						
						// Send a simple response
						responseText := "Obrigado por sua mensagem! Estamos processando sua solicitação."
						h.sendWhatsAppMessage(senderID, responseText)
					}
				}
			}
		}
	}

	// Acknowledge receipt
	w.WriteHeader(http.StatusOK)
}

// sendWhatsAppMessage sends a text message to a WhatsApp user
func (h *WhatsAppWebhookHandler) sendWhatsAppMessage(recipientID, message string) error {
	token := os.Getenv("WHATSAPP_TOKEN")
	if token == "" {
		return fmt.Errorf("WHATSAPP_TOKEN environment variable not set")
	}
	
	phoneNumberID := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	if phoneNumberID == "" {
		return fmt.Errorf("WHATSAPP_PHONE_NUMBER_ID environment variable not set")
	}
	
	url := fmt.Sprintf("https://graph.facebook.com/v17.0/%s/messages", phoneNumberID)
	
	// Construct the request payload
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                recipientID,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Error marshaling message payload: %v", err)
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("❌ Error creating request: %v", err)
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error sending message: %v", err)
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("❌ WhatsApp API error (status %d): %s", resp.StatusCode, string(responseBody))
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}
	
	log.Printf("✅ Message sent successfully to %s", recipientID)
	return nil
}