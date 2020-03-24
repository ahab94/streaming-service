package handlers

import (
	"bufio"
	"context"
	"net/http"
	"time"

	"github.com/ahab94/streaming-service/models"
	"github.com/ahab94/streaming-service/service"
)

type registerHandler struct {
	ctx  context.Context
	svc  service.Svc
	path string
}

func NewRegisterHandler(ctx context.Context, svc service.Svc, path string) *registerHandler {
	return &registerHandler{ctx: ctx, svc: svc, path: path}
}

func (r registerHandler) Path() string {
	return r.path
}

func (r *registerHandler) Handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		uid, _, err := bufio.NewReader(req.Body).ReadLine()
		if err != nil || len(uid) < 1 {
			w.WriteHeader(400)
			return
		}

		if err := r.svc.SaveUser(context.TODO(), &models.User{
			UserID:    string(uid),
			Timestamp: time.Now(),
		}); err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	default:
		w.WriteHeader(405)
	}
}
