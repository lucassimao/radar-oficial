package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

type WhatsAppService struct {
	phoneNumberID    string
	token            string
	userSessionSvc   *UserSessionService
}

type InteractiveListRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type InteractiveListSection struct {
	Title string              `json:"title"`
	Rows  []InteractiveListRow `json:"rows"`
}

func NewWhatsAppService(db *pgxpool.Pool) (*WhatsAppService, error) {
	token := os.Getenv("WHATSAPP_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("WHATSAPP_TOKEN environment variable not set")
	}

	phoneNumberID := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	if phoneNumberID == "" {
		return nil, fmt.Errorf("WHATSAPP_PHONE_NUMBER_ID environment variable not set")
	}

	userSessionSvc := NewUserSessionService(db)

	return &WhatsAppService{
		phoneNumberID:  phoneNumberID,
		token:          token,
		userSessionSvc: userSessionSvc,
	}, nil
}

func (s *WhatsAppService) SendWelcomeMessage(ctx context.Context, recipientID string, userName string) error {
	// Create or update user session when sending welcome message
	_, err := s.userSessionSvc.GetOrCreateUserSession(ctx, recipientID)
	if err != nil {
		log.Printf("⚠️ Failed to create/update user session: %v", err)
		// Continue anyway, don't fail the welcome message
	}
	welcomeMsg := fmt.Sprintf("Olá %s! 👋\n\nBem-vindo ao *Radar Oficial*. "+
		"Estou aqui para ajudar você a encontrar informações nos Diários Oficiais do estado do Piauí.\n\n"+
		"Você pode me perguntar sobre:\n"+
		"✅ Licitações e contratos\n"+
		"✅ Nomeações e exonerações\n"+
		"✅ Legislação estadual\n"+
		"✅ Outras publicações oficiais\n\n"+
		"Como posso ajudar você hoje?", userName)

	return s.SendTextMessage(recipientID, welcomeMsg)
}

func (s *WhatsAppService) SendStateSelectionList(ctx context.Context, recipientID string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", s.phoneNumberID)

	// Define the interactive list payload
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                recipientID,
		"type":              "interactive",
		"interactive": map[string]interface{}{
			"type": "list",
			"header": map[string]interface{}{
				"type": "text",
				"text": "Estados Disponíveis",
			},
			"body": map[string]interface{}{
				"text": "Selecione um estado para acessar os diários oficiais:",
			},
			"footer": map[string]interface{}{
				"text": "Radar Oficial - Consulte diários oficiais facilmente",
			},
			"action": map[string]interface{}{
				"button": "Ver Estados",
				"sections": []InteractiveListSection{
					{
						Title: "Estados",
						Rows: []InteractiveListRow{
							{
								ID:          "piaui",
								Title:       "Piauí",
								Description: "Diários Oficiais do Estado do Piauí",
							},
							// Add more states as they become available
							{
								ID:          "coming_soon",
								Title:       "Mais estados em breve",
								Description: "Estamos trabalhando para adicionar mais estados",
							},
						},
					},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Error marshaling list message payload: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("❌ Error creating request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error sending list message: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("❌ WhatsApp API error (status %d): %s", resp.StatusCode, string(responseBody))
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	log.Printf("✅ Interactive list message sent successfully to %s", recipientID)
	return nil
}

// UpdateUserState updates the user's selected state
func (s *WhatsAppService) UpdateUserState(ctx context.Context, phoneNumber string, state string) error {
	return s.userSessionSvc.UpdateUserState(ctx, phoneNumber, state)
}

// GetUserState gets the user's current state
func (s *WhatsAppService) GetUserState(ctx context.Context, phoneNumber string) (string, error) {
	session, err := s.userSessionSvc.GetUserSession(ctx, phoneNumber)
	if err != nil {
		return "", err
	}
	
	if session.State == nil {
		return "", nil
	}
	
	return *session.State, nil
}

func (s *WhatsAppService) SendTextMessage(recipientID, message string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", s.phoneNumberID)

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
	req.Header.Set("Authorization", "Bearer "+s.token)

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