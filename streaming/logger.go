package streaming

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ahab94/streaming-service/config"
)

var logger = log.New()

func init() {
	level, err := log.ParseLevel(viper.GetString(config.LogLevel))
	if err != nil {
		log.Warnf("failed to parse loglevel... switching to debug level")
		level = log.DebugLevel
	}

	logger.SetLevel(level)
	logger.SetFormatter(&log.JSONFormatter{})
	logger.WithFields(log.Fields{
		"pkg": "streaming",
	})
}
