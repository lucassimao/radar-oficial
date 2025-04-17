package whatsapp

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/model"
)

// UserSessionService handles user session operations
type UserSessionService struct {
	DB *pgxpool.Pool
}

// NewUserSessionService creates a new UserSessionService
func NewUserSessionService(db *pgxpool.Pool) *UserSessionService {
	return &UserSessionService{DB: db}
}

// GetOrCreateUserSession retrieves an existing user session or creates a new one
func (s *UserSessionService) GetOrCreateUserSession(ctx context.Context, phoneNumber string) (*model.UserSession, error) {
	// First, try to find an existing session
	userSession, err := s.GetUserSession(ctx, phoneNumber)
	if err == nil {
		// Found an existing session
		return userSession, nil
	}
	
	// If no session exists, create a new one
	userSession = &model.UserSession{
		PhoneNumber:   phoneNumber,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}
	
	// Insert the new session
	query := `
		INSERT INTO user_sessions (
			phone_number, 
			created_at, 
			last_updated_at
		) VALUES ($1, $2, $3)
		RETURNING id, created_at, last_updated_at;
	`
	
	err = s.DB.QueryRow(ctx, query,
		userSession.PhoneNumber,
		userSession.CreatedAt,
		userSession.LastUpdatedAt,
	).Scan(
		&userSession.ID,
		&userSession.CreatedAt,
		&userSession.LastUpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user session: %w", err)
	}
	
	log.Printf("✅ Created new user session for phone number: %s", phoneNumber)
	return userSession, nil
}

// GetUserSession retrieves a user session by phone number
func (s *UserSessionService) GetUserSession(ctx context.Context, phoneNumber string) (*model.UserSession, error) {
	query := `
		SELECT 
			id, 
			phone_number, 
			created_at, 
			last_updated_at, 
			state
		FROM user_sessions
		WHERE phone_number = $1
	`
	
	userSession := &model.UserSession{}
	err := s.DB.QueryRow(ctx, query, phoneNumber).Scan(
		&userSession.ID,
		&userSession.PhoneNumber,
		&userSession.CreatedAt,
		&userSession.LastUpdatedAt,
		&userSession.State,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}
	
	return userSession, nil
}

// UpdateUserState updates the state of a user session
func (s *UserSessionService) UpdateUserState(ctx context.Context, phoneNumber string, state string) error {
	query := `
		UPDATE user_sessions
		SET 
			state = $2,
			last_updated_at = NOW()
		WHERE phone_number = $1
	`
	
	result, err := s.DB.Exec(ctx, query, phoneNumber, state)
	if err != nil {
		return fmt.Errorf("failed to update user state: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		// No rows affected means no matching user session was found
		// Create a new session with the state
		_, err = s.GetOrCreateUserSession(ctx, phoneNumber)
		if err != nil {
			return err
		}
		
		// Now try updating again
		_, err = s.DB.Exec(ctx, query, phoneNumber, state)
		if err != nil {
			return fmt.Errorf("failed to update user state after creation: %w", err)
		}
	}
	
	log.Printf("✅ Updated state to '%s' for phone number: %s", state, phoneNumber)
	return nil
}