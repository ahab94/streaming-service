package cassandra

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	domain "github.com/ahab94/streaming-service/models"
)

func TestCassandra_GetUsages(t *testing.T) {
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
	_ = m.IncrementUsage(context.TODO(), packet1)

	packet2 := &domain.Packet{
		SequenceNum: 222222222,
		ChannelID:   "1",
		UserID:      "1",
		Timestamp:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Millisecond),
	}
	_ = m.IncrementUsage(context.TODO(), packet2)

	packet3 := &domain.Packet{
		SequenceNum: 333333333,
		ChannelID:   "3",
		UserID:      "3",
		Timestamp:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Millisecond),
	}
	_ = m.IncrementUsage(context.TODO(), packet3)

	packet4 := &domain.Packet{
		SequenceNum: 222222222,
		ChannelID:   "1",
		UserID:      "4",
		Timestamp:   time.Date(2020, 1, 1, 2, 1, 1, 1, time.UTC).Round(time.Millisecond),
	}
	_ = m.IncrementUsage(context.TODO(), packet4)

	type args struct {
		userID    string
		channelID string
		start     time.Time
		end       time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "get user 3 usages for channel 3",
			args: args{userID: "3", channelID: "3",
				start: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Hour),
				end:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Hour)},
			want: 1,
		},
		{
			name: "get user 1 usages for channel 1",
			args: args{userID: "1", channelID: "1",
				start: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Hour),
				end:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Hour)},
			want: 2,
		},
		{
			name: "get user 4 usages for channel 1",
			args: args{userID: "4", channelID: "1",
				start: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Hour),
				end:   time.Date(2020, 1, 1, 2, 1, 1, 1, time.UTC).Round(time.Hour)},
			want: 1,
		},
		{
			name: "get user 1 usages bad channel",
			args: args{userID: "1", channelID: "2",
				start: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Hour),
				end:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Round(time.Hour)},
			want: 0,
		},
		{
			name: "get user 1 usages for bad time",
			args: args{userID: "1", channelID: "2",
				start: time.Date(2020, 1, 1, 2, 1, 1, 1, time.UTC).Round(time.Hour),
				end:   time.Date(2020, 1, 1, 3, 1, 1, 1, time.UTC).Round(time.Hour)},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.GetUsage(context.TODO(), tt.args.userID, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.ListPackets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(usagesToCount(got, tt.args.channelID), tt.want) {
				t.Errorf("client.GetUsages() = %v, want %v", usagesToCount(got, tt.args.channelID), tt.want)
			}
		})
	}
}

func TestCassandra_IncrementUsage(t *testing.T) {
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
			err := m.IncrementUsage(context.TODO(), tt.args.packet)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.IncrementUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func usagesToCount(packets []domain.HourlyUsage, channelID string) int {
	m := 0
	for _, val := range packets {
		if val.ChannelID == channelID {
			m += val.PacketCount
		}
	}

	return m
}
