package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/whatsapp"
)

// WhatsAppWebhookHandler handles incoming webhook requests from WhatsApp
type WhatsAppWebhookHandler struct {
	whatsappService *whatsapp.WhatsAppService
	db              *pgxpool.Pool
}

// NewWhatsAppWebhookHandler creates a new WhatsAppWebhookHandler
func NewWhatsAppWebhookHandler(db *pgxpool.Pool) (*WhatsAppWebhookHandler, error) {
	whatsappService, err := whatsapp.NewWhatsAppService(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp service: %w", err)
	}

	return &WhatsAppWebhookHandler{
		whatsappService: whatsappService,
		db:              db,
	}, nil
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
					} `json:"text,omitempty"`
					Interactive struct {
						ButtonReply struct {
							ID    string `json:"id"`
							Title string `json:"title"`
						} `json:"button_reply,omitempty"`
						ListReply struct {
							ID          string `json:"id"`
							Title       string `json:"title"`
							Description string `json:"description,omitempty"`
						} `json:"list_reply,omitempty"`
					} `json:"interactive,omitempty"`
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
	// Use the request's context for database operations
	ctx := r.Context()
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

	// Extract the sender info and messages (if present)
	for _, entry := range message.Entry {
		for _, change := range entry.Changes {
			if change.Field == "messages" {
				// Collect user info from contacts section
				var userName string
				if len(change.Value.Contacts) > 0 {
					userName = change.Value.Contacts[0].Profile.Name
				}

				// Collect all message texts from this change
				var allMessages []string
				var senderID string
				var isFirstMessage bool = false

				// Extract sender ID first to check user session
				if len(change.Value.Messages) > 0 {
					senderID = change.Value.Messages[0].From

					// Check if this is a first-time interaction by checking user session state
					userState, err := h.whatsappService.GetUserState(ctx, senderID)
					if err != nil || userState == "" {
						// User has no state or session, consider it a first message
						isFirstMessage = true
						log.Printf("📱 New user or user without state: %s", senderID)
					} else {
						log.Printf("📱 Returning user with state: %s, state: %s", senderID, userState)
					}
				}

				for _, msg := range change.Value.Messages {
					if senderID == "" {
						senderID = msg.From
					}

					// Handle different message types
					switch msg.Type {
					case "text":
						messageText := msg.Text.Body
						log.Printf("📱 Text message from %s: %s", senderID, messageText)
						allMessages = append(allMessages, messageText)

						// Check if this is potentially a first-time message
						lowerText := strings.ToLower(messageText)
						if isFirstMessage && (lowerText == "oi" || lowerText == "olá" || lowerText == "ola" ||
							lowerText == "hi" || lowerText == "hello" || strings.Contains(lowerText, "bom dia") ||
							strings.Contains(lowerText, "boa tarde") || strings.Contains(lowerText, "boa noite")) {
							// Send welcome message for first-time or greeting messages
							if err := h.whatsappService.SendWelcomeMessage(ctx, senderID, userName); err != nil {
								log.Printf("❌ Error sending welcome message: %v", err)
							}
							// Then send state selection
							if err := h.whatsappService.SendStateSelectionList(ctx, senderID); err != nil {
								log.Printf("❌ Error sending state selection: %v", err)
							}
							isFirstMessage = false
							// Skip AI processing since we're sending welcome message
							allMessages = nil
							break
						}

					case "interactive":
						// Handle interactive message responses
						if msg.Interactive.ListReply.ID != "" {
							selection := msg.Interactive.ListReply.ID
							log.Printf("📱 Interactive list selection from %s: %s", senderID, selection)

							if selection == "piaui" {
								// Update user state in database
								if err := h.whatsappService.UpdateUserState(ctx, senderID, "piaui"); err != nil {
									log.Printf("❌ Error updating user state: %v", err)
								}
								responseText := "Você selecionou o *Piauí*. Agora você pode me perguntar sobre qualquer publicação."
								if err := h.whatsappService.SendTextMessage(senderID, responseText); err != nil {
									log.Printf("❌ Error sending response message: %v", err)
								}
							} else if selection == "coming_soon" {
								responseText := "Estamos trabalhando para adicionar mais estados em breve."
								if err := h.whatsappService.SendTextMessage(senderID, responseText); err != nil {
									log.Printf("❌ Error sending response message: %v", err)
								}
							}

							// Skip AI processing for interactive responses
							allMessages = nil
						} else if msg.Interactive.ButtonReply.ID != "" {
							button := msg.Interactive.ButtonReply.ID
							log.Printf("📱 Interactive button selection from %s: %s", senderID, button)
							// Handle button replies if needed in the future
						}
					}
				}

				// Process collected messages with AI if we have any
				if len(allMessages) > 0 && senderID != "" {
					// Combine all messages into a single string
					combinedMessage := strings.Join(allMessages, "\n")

					// Send the combined message to the AI agent
					agentResponse, err := h.sendMessageToAIAgent(combinedMessage)
					if err != nil {
						log.Printf("❌ Error sending message to AI agent: %v", err)
						responseText := "Desculpe, estamos com dificuldades técnicas. Tente novamente mais tarde."
						h.whatsappService.SendTextMessage(senderID, responseText)
					} else {
						// Send the AI response back to the user
						h.whatsappService.SendTextMessage(senderID, agentResponse)
					}
				}
			}
		}
	}

	// Acknowledge receipt
	w.WriteHeader(http.StatusOK)
}

// sendMessageToAIAgent sends a message to the Digital Ocean AI agent and returns the response
func (h *WhatsAppWebhookHandler) sendMessageToAIAgent(message string) (string, error) {
	// Get endpoint and access key from environment variables
	agentEndpoint := os.Getenv("DO_AGENT_PIAUI_URL")
	if agentEndpoint == "" {
		return "", fmt.Errorf("DO_AGENT_PIAUI_URL environment variable not set")
	}

	agentAccessKey := os.Getenv("DO_AGENT_PIAUI_ACCESS_KEY")
	if agentAccessKey == "" {
		return "", fmt.Errorf("DO_AGENT_PIAUI_ACCESS_KEY environment variable not set")
	}

	// Construct the url for the request
	url := fmt.Sprintf("%s/api/v1/chat/completions", agentEndpoint)

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
		log.Printf("❌ Error marshaling AI agent payload: %v", err)
		return "", err
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("❌ Error creating AI agent request: %v", err)
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+agentAccessKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error sending message to AI agent: %v", err)
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
		log.Printf("❌ Error reading AI agent response: %v", err)
		return "", err
	}

	if err := json.Unmarshal(responseBody, &aiResponse); err != nil {
		log.Printf("❌ Error parsing AI agent response: %v", err)
		return "", err
	}

	// Extract the response message
	if len(aiResponse.Choices) == 0 {
		return "", fmt.Errorf("no response from AI agent")
	}

	log.Printf("✅ Got response from AI agent: %s", aiResponse.Choices[0].Message.Content)
	return aiResponse.Choices[0].Message.Content, nil
}

// This function has been moved to the WhatsAppService
