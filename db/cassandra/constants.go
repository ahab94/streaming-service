package cassandra

const (
	packetCF   = "packets"
	timeFormat = "2006-01-02T15:04:05.000Z"
	usersCF    = "users"
	usageCF    = "hourly_usage"

	// model key constants
	userIDKey    = "user_id"
	startDateKey = "start_date"
	endDateKey   = "end_date"
	hourKey      = "hour"
	channelIDKey = "channel_id"
	timestampKey = "timestamp"
)
