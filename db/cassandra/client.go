package cassandra

import (
	"errors"
	"strings"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"

	"github.com/ahab94/streaming-service/config"
	"github.com/ahab94/streaming-service/db"
)

func init() {
	db.Register("cassandra", NewStore)
}

// conn defines the memory store spec
type client struct {
	conn *gocql.Session
}

// NewStore initializes memory store
func NewStore() (db.DataStore, error) {
	cluster, err := newCluster()
	if err != nil {
		return nil, err
	}

	session, err := cluster.CreateSession()
	if err != nil {
		logger.Errorf("unable to connect to cassandra db: %s", err)
		return nil, err
	}

	return &client{conn: session}, nil
}

func newCluster() (*gocql.ClusterConfig, error) {
	cluster := gocql.NewCluster(strings.Split(viper.GetString(config.DbNodes), ",")...)
	cluster.Keyspace = viper.GetString(config.DbName)
	if viper.GetString(config.DbUser) == "" || viper.GetString(config.DbPass) == "" {
		return nil, errors.New("empty username or password")
	}
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: viper.GetString(config.DbUser),
		Password: viper.GetString(config.DbPass),
	}

	return cluster, nil
}
