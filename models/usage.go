package models

import (
	"time"

	"github.com/fatih/structs"
)

// HourlyUsage this struct holds information about an stream usage
type HourlyUsage struct {
	UserID      string    `db:"user_id" structs:"user_id" json:"user_id"`
	ChannelID   string    `db:"channel_id" structs:"channel_id" json:"channel_id"`
	Hour        time.Time `db:"hour" structs:"hour" json:"hour"`
	PacketCount int       `db:"pkt_count" structs:"pkt_count" json:"pkt_count"`
}

// Map instantiates a map of struct
func (h *HourlyUsage) Map() map[string]interface{} {
	return structs.Map(h)
}

// Names instantiates a list of fields
func (h *HourlyUsage) Names() []string {
	fields := structs.Fields(h)
	names := make([]string, len(fields))
	for i, field := range fields {
		name := field.Name()
		tagName := field.Tag(structs.DefaultTagName)
		if tagName != "" {
			name = tagName
		}
		names[i] = name
	}
	return names
}
