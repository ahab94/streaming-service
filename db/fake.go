package db

import (
	"context"
	"errors"
	"time"

	domain "github.com/ahab94/streaming-service/models"
)

func NewFakeStore() (DataStore, error) {
	return &FakeStore{}, nil
}

type FakeStore struct {
	Fail bool
}

// ListPackets gets Packet by id
func (f *FakeStore) ListPackets(ctx context.Context, filter map[string]interface{}) ([]domain.Packet, error) {
	if f.Fail {
		return nil, errors.New("Fail")
	}
	return nil, nil
}

// SavePacket saves Packet in store
func (f *FakeStore) SavePacket(ctx context.Context, stats *domain.Packet) error {
	if f.Fail {
		return errors.New("Fail")
	}
	return nil
}

func (f *FakeStore) IncrementUsage(ctx context.Context, packet *domain.Packet) error {
	if f.Fail {
		return errors.New("Fail")
	}
	return nil
}

// GetUsage gets hourly usages
func (f *FakeStore) GetUsage(ctx context.Context, userID string, start, end time.Time) ([]domain.HourlyUsage, error) {
	if f.Fail {
		return nil, errors.New("Fail")
	}
	return nil, nil
}

func (f *FakeStore) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	if f.Fail {
		return nil, errors.New("Fail")
	}
	return nil, nil
}

func (f *FakeStore) SaveUser(ctx context.Context, user *domain.User) error {
	if f.Fail {
		return errors.New("Fail")
	}
	return nil
}
