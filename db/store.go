package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	domain "github.com/ahab94/streaming-service/models"
)

// DataStore is an interface for store implementation logic
type DataStore interface {
	ListPackets(ctx context.Context, filter map[string]interface{}) ([]domain.Packet, error)

	SavePacket(ctx context.Context, stats *domain.Packet) error

	IncrementUsage(ctx context.Context, packet *domain.Packet) error

	GetUsage(ctx context.Context, userID string, start, end time.Time) ([]domain.HourlyUsage, error)

	GetUser(ctx context.Context, userID string) (*domain.User, error)

	SaveUser(ctx context.Context, user *domain.User) error
}

type dataStoreFactory func() (DataStore, error)

// datastoreFactories is a key/value pair for the registered datastores
var datastoreFactories = make(map[string]dataStoreFactory)

// Register registers a store
func Register(name string, factory dataStoreFactory) {
	if factory == nil {
		log.Panicf("DataStore factory %s does not exist.", name)
	}
	_, registered := datastoreFactories[name]
	if registered {
		log.Errorf("DataStore factory %s already registered. Ignoring.", name)
	}
	datastoreFactories[name] = factory
}

// CreateDatastore instantiates datastore of particular type
func CreateDatastore(datastoreType string) (DataStore, error) {

	engineFactory, ok := datastoreFactories[datastoreType]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableDatastores := make([]string, len(datastoreFactories))
		for k := range datastoreFactories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, fmt.Errorf("invalid DataStore name. Must be one of: %s", strings.Join(availableDatastores, ", "))
	}

	// Run the factory with the configuration.
	return engineFactory()
}
