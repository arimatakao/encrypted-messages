package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/arimatakao/encrypted-messages/cmd/config"
	"github.com/arimatakao/encrypted-messages/internal/server/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	l   *log.Logger
	db  storage.MessageStorager
	srv *http.Server
}

func NewServer(logger *log.Logger) (*Server, error) {
	logger.Print("creating connection to DB")
	database, err := storage.NewMongoDB(config.App.DbUrl)
	if err != nil {
		logger.Print("get error while connecting to DB")
		return nil, err
	}
	logger.Print("connected to DB")

	server := &Server{
		l:  logger,
		db: database,
	}

	return server, nil
}

func (s *Server) Run() error {
	baseRoute := mux.NewRouter().PathPrefix(config.App.BaseUrl).Subrouter()

	// Create messasge
	baseRoute.HandleFunc(config.MESSAGE_ROUTE, s.AddMessageHandler).
		Methods(http.MethodPost)

	// Get user info
	baseRoute.HandleFunc(config.USER_ROUTE, s.EmptyHandler)

	userRoutes := baseRoute.PathPrefix(config.USER_ROUTE).Subrouter()

	// Register user
	userRoutes.HandleFunc(config.REGISTRATION_ROUTE, s.RegisterHandler).
		Methods(http.MethodPost)

	// Login user. Return access token
	userRoutes.HandleFunc(config.LOGIN_ROUTE, s.EmptyHandler).
		Methods(http.MethodPost)

	// Get all user message
	userRoutes.HandleFunc(config.USER_MESSAGES_ROUTE, s.EmptyHandler).
		Methods(http.MethodGet)

	messageRoutes := baseRoute.PathPrefix(config.MESSAGE_ROUTE).Subrouter()

	// Get message by id in url
	messageRoutes.HandleFunc(config.MESSAGE_ID_ROUTE, s.ReadMessageHandler).
		Methods(http.MethodGet)

	messageRoutes.HandleFunc(config.MESSAGE_ID_ROUTE, s.DeleteMessageHandler).
		Methods(http.MethodDelete)

	_ = baseRoute.Walk(func(route *mux.Route,
		router *mux.Router,
		ancestors []*mux.Route) error {

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
	if err := s.db.Disconnect(ctx); err != nil {
		return err
	}
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
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
	s.l.Printf("success triggering: %s", r.URL.Path)
}

func (s Server) AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("triggering %s", r.URL.Path)

	msg := &storage.MessageReq{}
	if err := ReadJSON(r, msg); err != nil {
		s.l.Printf("failed to unmarshal request body: %s", err.Error())
		WriteStatus(w, http.StatusBadRequest)
		return
	}

	if msg.IsEmpty() {
		s.l.Print("request body without correct fields")
		WriteStatus(w, http.StatusBadRequest)
		return
	}

	res, err := s.db.AddMessage(msg)
	if err != nil {
		s.l.Printf("failed to add message to db: %s", err.Error())
		WriteStatus(w, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"id": res,
	}

	if err := WriteJSON(w, http.StatusOK, resp); err != nil {
		WriteStatus(w, http.StatusInternalServerError)
		s.l.Printf("error due send response: %s", err.Error())
		return
	}

	s.l.Printf("end triggering: %s", r.URL.Path)
}

func (s Server) ReadMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("triggering %s", r.URL.Path)

	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		WriteStatus(w, http.StatusInternalServerError)
		s.l.Printf("failed to get var id from path: %s", r.URL.Path)
		return
	}

	msg, err := s.db.ReadMessage(id)
	if err != nil {
		s.l.Printf("failed to read message to db: %s", err.Error())
		WriteStatus(w, http.StatusNotFound)
		return
	}

	if err := WriteJSON(w, http.StatusOK, msg); err != nil {
		WriteStatus(w, http.StatusInternalServerError)
		s.l.Printf("error due send response: %s", err.Error())
		return
	}

	s.l.Printf("end triggering: %s", r.URL.Path)
}

func (s Server) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("triggering %s", r.URL.Path)

	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		WriteStatus(w, http.StatusInternalServerError)
		s.l.Printf("failed to get var id from path: %s", r.URL.Path)
		return
	}

	err := s.db.DeleteMessage(id)
	if err != nil {
		s.l.Printf("failed to delete message from db: %s", err.Error())
		WriteStatus(w, http.StatusInternalServerError)
		return
	}

	WriteStatus(w, http.StatusNoContent)

	s.l.Printf("end triggering: %s", r.URL.Path)
}
