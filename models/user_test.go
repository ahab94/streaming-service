package models

import (
	"reflect"
	"testing"
	"time"
)

func TestUser_Map(t *testing.T) {
	kv := map[string]interface{}{
		"user_id":   "1234",
		"timestamp": time.Time{},
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
			name: "User to map expect value",
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
			f := &User{
				Timestamp: tt.fields.timestamp,
				UserID:    tt.fields.userID,
			}
			if got := f.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Names(t *testing.T) {

	tests := []struct {
		name string
		want []string
	}{
		{
			name: "field names of struct User",
			want: []string{
				"user_id",
				"timestamp",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &User{}
			if got := f.Names(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Names() = %v, want %v", got, tt.want)
			}
		})
	}
}
