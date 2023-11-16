package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/arimatakao/encrypted-messages/cmd/config"
	"github.com/gorilla/mux"
)

type Server struct {
	l   *log.Logger
	srv *http.Server
}

func NewServer(logger *log.Logger) *Server {
	server := &Server{
		l: logger,
	}

	return server
}

func (s *Server) Run() error {
	baseRoute := mux.NewRouter().PathPrefix(config.App.BaseUrl).Subrouter()

	userRoutes := baseRoute.PathPrefix(config.USER_ROUTE).Subrouter()

	userRoutes.HandleFunc(config.REGISTRATION_ROUTE, s.RegisterHandler).
		Methods(http.MethodPost)

	userRoutes.HandleFunc(config.LOGIN_ROUTE, s.EmptyHandler).
		Methods(http.MethodPost)

	userRoutes.HandleFunc(config.USER_MESSAGES_ROUTE, s.EmptyHandler).
		Methods(http.MethodGet)

	messageRoutes := baseRoute.PathPrefix(config.MESSAGE_ROUTE).Subrouter()

	messageRoutes.HandleFunc(config.GET_MESSAGE_ROUTE, s.EmptyHandler).
		Methods(http.MethodPost)

	_ = baseRoute.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			methods, err := route.GetMethods()
			if err == nil {
				s.l.Printf("registered route | %s %s", methods, path)
			}
		}
		return nil
	})

	s.srv = &http.Server{
		Addr:         ":" + config.App.Port,
		Handler:      baseRoute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	s.l.Printf("start listening port: %s", config.App.Port)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s Server) EmptyHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("empty handler triggered by route %s", r.URL.Path)
	WriteStatus(w, http.StatusLocked)
}

func (s Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("registration handler triggering")
	resp := map[string]string{
		"bearer_token": "qwer1234asdf",
	}

	if err := WriteJSON(w, http.StatusOK, resp); err != nil {
		WriteStatus(w, http.StatusInternalServerError)
		s.l.Printf("get internal error %s", err.Error())
		return
	}
	s.l.Printf("success triggering")
}
