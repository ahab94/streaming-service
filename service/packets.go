package service

import (
	"context"
	"errors"
	"time"

	domain "github.com/ahab94/streaming-service/models"
)

// ListPackets list packets by constraint
func (s *service) ListPackets(ctx context.Context, filter map[string]interface{}) ([]domain.Packet, error) {
	if _, ok := filter["user_id"]; !ok {
		return nil, errors.New("must have user_id to list packets")
	}

	return s.store.ListPackets(ctx, filter)
}

// SavePacket saves Packet in store
func (s *service) SavePacket(ctx context.Context, packet *domain.Packet) error {
	if packet == nil {
		return errors.New("undefined packet")
	}
	if err := s.store.SavePacket(ctx, packet); err != nil {
		return errors.New("failed to save packet")
	}
	if err := s.store.IncrementUsage(ctx, packet); err != nil {
		return errors.New("failed to increment usage")
	}

	return nil
}

// Usage processes user's streaming usage for the last 24h
func (s *service) Usage(ctx context.Context, userID string) (map[string]map[string]time.Duration, error) {
	if userID == "" {
		return nil, errors.New("id must be defined")
	}

	now := time.Now()
	hourly, err := s.store.GetUsage(ctx, userID, now.Add(-24*time.Hour), now.Add(1*time.Hour))
	if err != nil {
		return nil, err
	}

	usage := make(map[string]time.Duration)
	for _, use := range hourly {
		// assuming each packet had 20ms of payload
		usage[use.ChannelID] += time.Duration(use.PacketCount / 50)

	}

	return map[string]map[string]time.Duration{"sessions": usage}, nil
}
