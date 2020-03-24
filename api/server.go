package api

import (
	"context"
	"net/http"

	"github.com/ahab94/streaming-service/api/handlers"
	"github.com/ahab94/streaming-service/service"
)

func NewStreamService(context context.Context, svc service.Svc) *http.Server {
	srv := &http.Server{}

	srv.Handler = streamSvcHandler(context, svc)

	return srv
}

func streamSvcHandler(context context.Context, svc service.Svc) http.Handler {
	sm := http.NewServeMux()

	registerHandler := handlers.NewRegisterHandler(context, svc, "/register")
	sm.HandleFunc(registerHandler.Path(), registerHandler.Handle)

	statsHandler := handlers.NewGetUsageHandler(context, svc, "/usage/")
	sm.HandleFunc(statsHandler.Path(), statsHandler.Handle)

	return sm
}
