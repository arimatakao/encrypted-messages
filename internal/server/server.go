package server

import (
	"log"
	"net/http"
	"os"

	"github.com/arimatakao/encrypted-messages/cmd/config"
	"github.com/gorilla/mux"
)

type Server struct {
	l       *log.Logger
	addr    string
	baseUrl string
	router  *mux.Router
}

func NewServer(c config.Config) Server {
	logger := log.New(os.Stdout, "INFO ", log.LstdFlags)
	logger.Print("server logger inited")

	r := mux.NewRouter()

	server := Server{
		l:       logger,
		addr:    c.App.Port,
		baseUrl: c.App.BaseUrl,
		router:  r,
	}

	server.router.PathPrefix(c.App.BaseUrl).
		Path("/register").
		HandlerFunc(server.RegisterHandler).
		Methods(http.MethodPost)

	server.l.Print("path inited")
	return server
}

func (s *Server) Run() error {
	s.l.Printf("start listening port: %s", s.addr)
	return http.ListenAndServe(":"+s.addr, s.router)
}

func (s Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Printf("get request: %s", r.URL.Path)
	w.WriteHeader(http.StatusOK)
}
