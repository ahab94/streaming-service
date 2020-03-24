package models

import (
	"encoding/json"
	"time"

	"github.com/fatih/structs"
)

// Packet this struct holds information about a packet
type Packet struct {
	UserID      string    `db:"user_id" structs:"user_id" json:"user_id"`
	SequenceNum int64     `db:"sequence_number" structs:"sequence_number" json:"sequence_number"`
	ChannelID   string    `db:"channel_id" structs:"channel_id" json:"channel_id"`
	Timestamp   time.Time `db:"timestamp" structs:"timestamp" json:"timestamp"`
	Raw         []byte    `db:"-" structs:"-" json:"-"`
}

// MarshalBinary instantiates json bytes of Packet struct
func (p *Packet) MarshalBinary() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		logger.Errorf("unable to marshal error because = %v ", err)
		return nil
	}
	return data
}

// Map instantiates a map of struct
func (p *Packet) Map() map[string]interface{} {
	return structs.Map(p)
}

// Names instantiates a list of fields
func (p *Packet) Names() []string {
	fields := structs.Fields(p)
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
