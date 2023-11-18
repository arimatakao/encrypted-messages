package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arimatakao/encrypted-messages/cmd/config"
	"github.com/arimatakao/encrypted-messages/internal/server/storage"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_SUBROUTE      = "/user"
	REGISTRATION_ROUTE = "/registration"
	LOGIN_ROUTE        = "/login"

	MESSAGE_ROUTE    = "/message"
	MESSAGE_ID_ROUTE = "/message/{id}"
)

type Server struct {
	l   *log.Logger
	db  storage.Storager
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
	r := mux.NewRouter().PathPrefix(config.App.BaseUrl).Subrouter()

	// DELETE EMPTY HANDLER
	// AND
	// MAKE NORMAL STRUCT WITCH USE AS FILTER TO DB

	r.HandleFunc(REGISTRATION_ROUTE, s.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc(LOGIN_ROUTE, s.LoginHandler).Methods(http.MethodPost)

	r.HandleFunc(MESSAGE_ROUTE, s.EmptyHandler).Methods(http.MethodGet)
	r.HandleFunc(MESSAGE_ID_ROUTE, s.EmptyHandler).Methods(http.MethodGet)

	userSubroute := r.PathPrefix(USER_SUBROUTE).Subrouter()
	userSubroute.Use(s.AuthMiddleware)
	userSubroute.HandleFunc(MESSAGE_ROUTE, s.AddMessageHandler).Methods(http.MethodPost)
	userSubroute.HandleFunc(MESSAGE_ROUTE, s.EmptyHandler).Methods(http.MethodGet)

	_ = r.Walk(func(route *mux.Route,
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
		Handler:      r,
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
	uId, err := GetUserId(r)
	if err != nil {
		s.l.Printf("can't get user_id from context: %v", err)
		WriteStatus(w, http.StatusUnauthorized)
		return
	}
	s.l.Printf("empty handler triggered by route %s by user", uId)
	WriteStatus(w, http.StatusOK)
}

func (s Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("registration handler triggering")

	user := &storage.UserReq{}
	if err := ReadJSON(r, user); err != nil {
		s.l.Printf("failed to unmarshal request body: %s", err.Error())
		WriteStatus(w, http.StatusBadRequest)
		return
	}
	_, err := s.db.ReadUserByUsername(user.Username)
	if err != mongo.ErrNoDocuments {
		s.l.Printf("failed to register user due already exist: %s", err.Error())
		WriteStatus(w, http.StatusConflict)
		return
	}

	if err = s.db.AddUser(user); err != nil {
		s.l.Printf("failed to add user: %s", err.Error())
		WriteStatus(w, http.StatusInternalServerError)
		return
	}

	WriteStatus(w, http.StatusCreated)

	s.l.Printf("success triggering: %s", r.URL.Path)
}

func (s Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("triggering %s", r.URL.Path)

	user := &storage.UserReq{}
	if err := ReadJSON(r, user); err != nil {
		s.l.Printf("failed to unmarshal request body: %s", err.Error())
		WriteStatus(w, http.StatusBadRequest)
		return
	}

	u, err := s.db.ReadUser(user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		s.l.Printf("failed to read user while loggining: %s", err.Error())
		WriteStatus(w, http.StatusNotFound)
		return
	}
	s.l.Printf("found user %s with id %s", u.Username, u.Id)

	resp := map[string]string{
		"auth_token": u.Id,
	}

	if err := WriteJSON(w, http.StatusOK, resp); err != nil {
		WriteStatus(w, http.StatusInternalServerError)
		s.l.Printf("error due send response: %s", err.Error())
		return
	}

	s.l.Printf("end triggering: %s", r.URL.Path)
}

func (s Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.l.Printf("auth middleware triggering")
		userId := r.Header.Get("Authorization")
		if userId == "" {
			WriteStatus(w, http.StatusUnauthorized)
			return
		}

		ctx, err := NewRequestContext(r.Context(), userId)
		if err != nil {
			s.l.Printf("failed to create request context: %s", err.Error())
			WriteStatus(w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
		s.l.Printf("success authorization triggering")
	})
}

func (s Server) ReadUserMessagesHandler(w http.ResponseWriter, r *http.Request) {
}

func (s Server) ReadPublicMessagesHandler(w http.ResponseWriter, r *http.Request) {
}

func (s Server) DeleteAllMessagesHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("triggering %s", r.URL.Path)

	id := "abc"

	err := s.db.DeleteAllMessages(id)
	if err != nil {
		s.l.Printf("failed to add message to db: %s", err.Error())
		WriteStatus(w, http.StatusInternalServerError)
		return
	}

	WriteStatus(w, http.StatusNoContent)

	s.l.Printf("end triggering: %s", r.URL.Path)
}

func (s Server) AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("triggering %s", r.URL.Path)

	uId, err := GetUserId(r)
	if err != nil {
		s.l.Printf("can't get user_id from context: %v", err)
		WriteStatus(w, http.StatusUnauthorized)
		return
	}

	msg := &storage.Message{}
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

	msg.Id = ""
	msg.OwnerId = uId

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
		s.l.Printf("failed to read message from db: %s", err.Error())
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
