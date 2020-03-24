package cassandra

import (
	"context"
	"os"
	"testing"
	"time"

	domain "github.com/ahab94/streaming-service/models"

	"github.com/google/go-cmp/cmp"
)

func TestClient_GetUser(t *testing.T) {
	os.Setenv("DB_NODES", "cassandra-db")

	newConn, _ := NewStore()

	user := &domain.User{
		UserID:    "1dc2b7e9b-47be-11e9-b0be-0800270099",
		Timestamp: time.Now().Round(time.Millisecond),
	}
	_ = newConn.SaveUser(context.TODO(), user)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name:    "fetch user by valid id",
			args:    args{id: "1dc2b7e9b-47be-11e9-b0be-0800270099"},
			want:    user,
			wantErr: false,
		},
		{
			name:    "fetch user by invalid id",
			args:    args{id: "some-invalid-id"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConn.GetUser(context.TODO(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("conn.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("conn.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SaveUser(t *testing.T) {
	os.Setenv("DB_NODES", "cassandra-db")

	newConn, _ := NewStore()
	user := &domain.User{
		UserID:    "1d2b7e9b-47be-11e9-b0be-080027009b40",
		Timestamp: time.Now(),
	}

	type args struct {
		user *domain.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "save user",
			args: args{
				user: user,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := newConn.SaveUser(context.TODO(), tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("conn.SaveEmployee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
