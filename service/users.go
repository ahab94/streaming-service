package service

import (
	"context"
	"errors"

	domain "github.com/ahab94/streaming-service/models"
)

// GetUser gets user from the store
func (s *service) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, errors.New("must have defined user_id to get user")
	}

	return s.store.GetUser(ctx, userID)
}

// SaveUser save user to the store
func (s *service) SaveUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("undefined user")
	}

	if err := s.store.SaveUser(ctx, user); err != nil {
		return errors.New("failed to save user")
	}

	return nil
}
