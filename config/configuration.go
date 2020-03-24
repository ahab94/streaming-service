package config

import (
	"github.com/spf13/viper"
)

// keys for database configuration
const (
	LogLevel = "log.level"

	DbName  = "db.name"
	DbNodes = "db.nodes"
	DbPort  = "db.port"
	DbUser  = "db.user"
	DbPass  = "db.pass"

	ServerHost = "server.host"
	ServerPort = "server.port"

	StreamSrvHost = "stream.host"
	StreamSrvPort = "stream.port"

	ChannelSourceHost = "cs.host"
	ChannelSourcePort = "cs.port"

	PacketBatchSize = "packet.batch.size"
)

func init() {
	_ = viper.BindEnv(LogLevel, "LOG_LEVEL")

	// env var for db
	_ = viper.BindEnv(DbName, "DB_NAME")
	_ = viper.BindEnv(DbNodes, "DB_NODES")
	_ = viper.BindEnv(DbUser, "DB_USER")
	_ = viper.BindEnv(DbPass, "DB_PASS")

	// env var for server
	_ = viper.BindEnv(ServerHost, "SERVER_HOST")
	_ = viper.BindEnv(ServerPort, "SERVER_PORT")

	// env var for server
	_ = viper.BindEnv(StreamSrvHost, "STREAM_HOST")
	_ = viper.BindEnv(StreamSrvPort, "STREAM_PORT")

	// env var for channel source
	_ = viper.BindEnv(ChannelSourceHost, "CHANNEL_SOURCE_HOST")
	_ = viper.BindEnv(ChannelSourcePort, "CHANNEL_SOURCE_PORT")

	_ = viper.BindEnv(PacketBatchSize, "PACKET_BATCH_SIZE")

	// defaults
	viper.SetDefault(LogLevel, "debug")

	viper.SetDefault(DbName, "streaming_service")
	viper.SetDefault(DbNodes, "localhost")
	viper.SetDefault(DbPort, "9042")
	viper.SetDefault(DbUser, "cassandra")
	viper.SetDefault(DbPass, "cassandra")

	viper.SetDefault(ServerHost, "127.0.0.1")
	viper.SetDefault(ServerPort, "8080")

	viper.SetDefault(StreamSrvHost, "127.0.0.1")
	viper.SetDefault(StreamSrvPort, "9111")

	viper.SetDefault(ChannelSourceHost, "localhost")
	viper.SetDefault(ChannelSourcePort, "9110")

	viper.SetDefault(PacketBatchSize, 10)
}
