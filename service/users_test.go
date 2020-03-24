package service

import (
	"context"
	"testing"
	"time"

	"github.com/ahab94/streaming-service/db"
	domain "github.com/ahab94/streaming-service/models"
)

func Test_service_GetUser(t *testing.T) {
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
			name: "success - get usar",
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
			_, err := s.GetUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_SaveUser(t *testing.T) {
	type fields struct {
		fail bool
	}
	type args struct {
		ctx  context.Context
		user *domain.User
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
				ctx: context.TODO(),
				user: &domain.User{
					UserID:    "1",
					Timestamp: time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "fail - save nil user",
			fields: fields{
				fail: false,
			},
			args: args{
				ctx:  context.TODO(),
				user: nil,
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
				user: &domain.User{
					UserID:    "1",
					Timestamp: time.Time{},
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
			if err := s.SaveUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SaveUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
