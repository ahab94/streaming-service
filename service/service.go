package service

import (
	"context"
	"time"

	"github.com/ahab94/streaming-service/db"
	domain "github.com/ahab94/streaming-service/models"
)

// Svc exposes services for streaming-service
type Svc interface {
	// ListPackets gets Packet by id
	ListPackets(ctx context.Context, filter map[string]interface{}) ([]domain.Packet, error)

	// SavePacket saves Packet in store
	SavePacket(ctx context.Context, stats *domain.Packet) error

	// Usage processes user's streaming usage for the last 24h
	Usage(ctx context.Context, userID string) (map[string]map[string]time.Duration, error)

	// GetUser gets user from the store
	GetUser(ctx context.Context, userID string) (*domain.User, error)

	// SaveUser saves user to the store
	SaveUser(ctx context.Context, user *domain.User) error
}

type service struct {
	store db.DataStore
}

// NewService exports service struct
func NewService(store db.DataStore) Svc {
	return &service{store: store}
}
