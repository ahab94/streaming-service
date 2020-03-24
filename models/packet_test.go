package models

import (
	"reflect"
	"testing"
	"time"
)

func TestPacket_Map(t *testing.T) {
	kv := map[string]interface{}{
		"sequence_number": int64(1),
		"channel_id":      "1234",
		"user_id":         "1234",
		"timestamp":       time.Time{},
	}
	type fields struct {
		userID      string
		timestamp   time.Time
		channelID   string
		sequenceNum int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "Packet to map expect value",
			fields: fields{
				sequenceNum: 1,
				channelID:   "1234",
				userID:      "1234",
				timestamp:   time.Time{},
			},
			want: kv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Packet{
				Timestamp:   tt.fields.timestamp,
				SequenceNum: tt.fields.sequenceNum,
				ChannelID:   tt.fields.channelID,
				UserID:      tt.fields.userID,
			}
			if got := f.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Packet.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_Names(t *testing.T) {

	tests := []struct {
		name string
		want []string
	}{
		{
			name: "field names of struct Packet",
			want: []string{
				"user_id",
				"sequence_number",
				"channel_id",
				"timestamp",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Packet{}
			if got := f.Names(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Packet.Names() = %v, want %v", got, tt.want)
			}
		})
	}
}
