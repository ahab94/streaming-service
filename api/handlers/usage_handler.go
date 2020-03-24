package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ahab94/streaming-service/service"
)

type usageHandler struct {
	ctx  context.Context
	svc  service.Svc
	path string
}

func NewGetUsageHandler(ctx context.Context, svc service.Svc, path string) *usageHandler {
	return &usageHandler{ctx: ctx, svc: svc, path: path}
}

func (u usageHandler) Path() string {
	return u.path
}

func (u *usageHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		usage, err := u.svc.Usage(u.ctx, strings.Replace(r.URL.EscapedPath(), u.path, "", -1))
		if err != nil {
			logger.Errorf("unable to get usage req:%+v, err:%+v", r, err)
			w.WriteHeader(500)
			return
		}

		if err := json.NewEncoder(w).Encode(usage); err != nil {
			logger.Errorf("unable to encode json req:%+v, err:%+v", r, err)
			w.WriteHeader(500)
			return
		}
	default:
		w.WriteHeader(405)
		return
	}
}
