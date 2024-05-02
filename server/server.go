package server

import (
	"log"
	"net/http"
	"os"
)

const (
	stdPort    = "7540"
	varEnvPort = "TODO_PORT"
)

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

// Function to get the port from the environment variable TODO_PORT
func getPort() string {
	port, exists := os.LookupEnv(varEnvPort)
	if !exists || port == "" {
		port = stdPort
	}

	log.Printf(`Retrieved port %s from env variable "%s"`, port, varEnvPort)

	return ":" + port
}
