package models

import (
	"reflect"
	"testing"
	"time"
)

func TestUsage_Map(t *testing.T) {
	kv := map[string]interface{}{
		"hour":       time.Time{},
		"channel_id": "1234",
		"user_id":    "1234",
		"pkt_count":  10,
	}
	type fields struct {
		userID      string
		hour        time.Time
		channelID   string
		packetCount int
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "Usage to map expect value",
			fields: fields{
				hour:        time.Time{},
				channelID:   "1234",
				userID:      "1234",
				packetCount: 10,
			},
			want: kv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &HourlyUsage{
				Hour:        tt.fields.hour,
				PacketCount: tt.fields.packetCount,
				ChannelID:   tt.fields.channelID,
				UserID:      tt.fields.userID,
			}
			if got := f.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usage.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsage_Names(t *testing.T) {

	tests := []struct {
		name string
		want []string
	}{
		{
			name: "field names of struct Usage",
			want: []string{
				"user_id",
				"channel_id",
				"hour",
				"pkt_count",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &HourlyUsage{}
			if got := f.Names(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usage.Names() = %v, want %v", got, tt.want)
			}
		})
	}
}
