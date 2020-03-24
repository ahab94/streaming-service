package service

import (
	"context"
	"testing"
	"time"

	"github.com/ahab94/streaming-service/db"
	domain "github.com/ahab94/streaming-service/models"
)

func Test_service_ListPackets(t *testing.T) {
	type fields struct {
		fail bool
	}
	type args struct {
		ctx    context.Context
		filter map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success - filter packets",
			fields: fields{
				fail: false,
			},
			args: args{
				ctx:    context.TODO(),
				filter: map[string]interface{}{"user_id": "1234"},
			},
			wantErr: false,
		},
		{
			name: "fail - without user id in the filter",
			fields: fields{
				fail: false,
			},
			args: args{
				ctx:    context.TODO(),
				filter: map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "fail - db failure",
			fields: fields{
				fail: true,
			},
			args: args{
				ctx:    context.TODO(),
				filter: map[string]interface{}{"user_id": "1234"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: &db.FakeStore{Fail: tt.fields.fail},
			}
			_, err := s.ListPackets(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPackets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_SavePacket(t *testing.T) {
	type fields struct {
		fail bool
	}
	type args struct {
		ctx    context.Context
		packet *domain.Packet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success - save packet",
			fields: fields{
				fail: false,
			},
			args: args{
				ctx: context.TODO(),
				packet: &domain.Packet{
					UserID:      "1",
					SequenceNum: 1,
					ChannelID:   "1",
					Timestamp:   time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "fail - save nil packet",
			fields: fields{
				fail: false,
			},
			args: args{
				ctx:    context.TODO(),
				packet: nil,
			},
			wantErr: true,
		},
		{
			name: "fail - db failure",
			fields: fields{
				fail: true,
			},
			args: args{
				ctx: context.TODO(),
				packet: &domain.Packet{
					UserID:      "1",
					SequenceNum: 1,
					ChannelID:   "1",
					Timestamp:   time.Time{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: &db.FakeStore{Fail: tt.fields.fail},
			}
			if err := s.SavePacket(tt.args.ctx, tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("SavePacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_Usage(t *testing.T) {
	type fields struct {
		fail bool
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success - get usage",
			fields: fields{
				fail: false,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "1234",
			},
			wantErr: false,
		},
		{
			name: "fail - invalid user id",
			fields: fields{
				fail: false,
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "fail - db failure",
			fields: fields{
				fail: true,
			},
			args: args{
				ctx:    context.TODO(),
				userID: "1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: &db.FakeStore{Fail: tt.fields.fail},
			}
			_, err := s.Usage(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
