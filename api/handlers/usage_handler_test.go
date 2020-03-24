package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahab94/streaming-service/db"
	"github.com/ahab94/streaming-service/service"
)

func Test_getStatsHandler_Handle(t *testing.T) {
	type fields struct {
		ctx  context.Context
		fail bool
	}
	type args struct {
		wantStatus int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success - get usage",
			fields: fields{
				ctx:  context.TODO(),
				fail: false,
			},
			args: args{
				wantStatus: 200,
			},
		},
		{
			name: "fail - db failure",
			fields: fields{
				ctx:  context.TODO(),
				fail: true,
			},
			args: args{
				wantStatus: 500,
			},
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", "/usage", nil)
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r := &usageHandler{
				ctx: tt.fields.ctx,
				svc: service.NewService(&db.FakeStore{Fail: tt.fields.fail}),
			}
			handler := http.HandlerFunc(r.Handle)
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.args.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v with response %v",
					status, tt.args.wantStatus, rr.Body.String())
			}
		})
	}
}
