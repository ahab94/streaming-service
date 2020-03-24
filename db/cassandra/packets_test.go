package cassandra

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	domain "github.com/ahab94/streaming-service/models"
)

func TestCassandra_ListPackets(t *testing.T) {
	os.Setenv("DB_NODES", "cassandra-db")
	m, err := NewStore()

	if err != nil {
		t.Fatal(err)
	}
	packet1 := &domain.Packet{
		SequenceNum: 11111111,
		ChannelID:   "1",
		UserID:      "1",
		Timestamp:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Millisecond),
	}
	_ = m.SavePacket(context.TODO(), packet1)

	packet2 := &domain.Packet{
		SequenceNum: 22222222,
		ChannelID:   "1",
		UserID:      "1",
		Timestamp:   time.Date(2021, 2, 2, 2, 2, 2, 2, time.UTC).Round(time.Millisecond),
	}
	_ = m.SavePacket(context.TODO(), packet2)

	packet3 := &domain.Packet{
		SequenceNum: 33333333,
		ChannelID:   "3",
		UserID:      "3",
		Timestamp:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Millisecond),
	}
	_ = m.SavePacket(context.TODO(), packet3)

	packet4 := &domain.Packet{
		SequenceNum: 22222222,
		ChannelID:   "1",
		UserID:      "4",
		Timestamp:   time.Date(2021, 2, 2, 2, 2, 2, 2, time.UTC).Round(time.Millisecond),
	}
	_ = m.SavePacket(context.TODO(), packet4)

	type args struct {
		filter map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[int64]domain.Packet
		wantErr bool
	}{
		{
			name: "list by user 3",
			args: args{filter: map[string]interface{}{
				"user_id": "3",
			}},
			want: map[int64]domain.Packet{
				packet3.SequenceNum: *packet3,
			},
		},
		{
			name: "list by start date and end date",
			args: args{filter: map[string]interface{}{
				"user_id":    "1",
				"start_date": time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
				"end_date":   time.Date(2021, 2, 2, 2, 2, 2, 2, time.UTC),
			}},
			want: map[int64]domain.Packet{
				packet1.SequenceNum: *packet1,
				packet2.SequenceNum: *packet2,
			},
		},
		{
			name: "list by end date",
			args: args{filter: map[string]interface{}{
				"user_id":  "1",
				"end_date": time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			}},
			want: map[int64]domain.Packet{
				packet1.SequenceNum: *packet1,
			},
		},
		{
			name: "list by start date",
			args: args{filter: map[string]interface{}{
				"user_id":    "1",
				"start_date": time.Date(2021, 2, 2, 2, 2, 2, 2, time.UTC),
			}},
			want: map[int64]domain.Packet{
				packet2.SequenceNum: *packet2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.ListPackets(context.TODO(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.ListPackets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(packetsToMap(got), tt.want) {
				t.Errorf("client.ListPackets() = %v, want %v", packetsToMap(got), tt.want)
			}
		})
	}
}

func TestCassandra_SavePackets(t *testing.T) {
	os.Setenv("DB_NODES", "cassandra-db")
	m, _ := NewStore()

	packet := &domain.Packet{
		SequenceNum: 111111111,
		ChannelID:   "1",
		UserID:      "1",
		Timestamp:   time.Now(),
	}

	type args struct {
		packet *domain.Packet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "save Packet",
			args: args{
				packet: packet,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.SavePacket(context.TODO(), tt.args.packet)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.SavePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func packetsToMap(packets []domain.Packet) map[int64]domain.Packet {
	m := make(map[int64]domain.Packet)
	for _, val := range packets {
		m[val.SequenceNum] = val
	}
	return m
}
