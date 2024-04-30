package server

import (
	"log"
	"net/http"
	"os"
)

const stdPort = "7540"

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr: getPort(),
	}

	log.Printf("Server started on %s", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}

func getPort() string {
	port, exists := os.LookupEnv("TODO_PORT")
	if !exists || port == "" {
		port = stdPort
	}
	return ":" + port
}
