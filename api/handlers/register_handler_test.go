package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahab94/streaming-service/db"
	"github.com/ahab94/streaming-service/service"
)

func Test_registerHandler_Handle(t *testing.T) {
	type fields struct {
		ctx  context.Context
		fail bool
	}
	type args struct {
		body       string
		wantStatus int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success - register new user",
			fields: fields{
				ctx:  context.TODO(),
				fail: false,
			},
			args: args{
				body:       "new-user",
				wantStatus: 200,
			},
		},
		{
			name: "fail - invalid body",
			fields: fields{
				ctx:  context.TODO(),
				fail: false,
			},
			args: args{
				body:       "",
				wantStatus: 400,
			},
		},
		{
			name: "fail - db failure",
			fields: fields{
				ctx:  context.TODO(),
				fail: true,
			},
			args: args{
				body:       "user-1",
				wantStatus: 500,
			},
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(tt.args.body)))
		t.Run(tt.name, func(t *testing.T) {
			r := &registerHandler{
				ctx: tt.fields.ctx,
				svc: service.NewService(&db.FakeStore{Fail: tt.fields.fail}),
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(r.Handle)
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.args.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v with response: %v",
					status, tt.args.wantStatus, rr.Body.String())
			}
		})
	}
}
