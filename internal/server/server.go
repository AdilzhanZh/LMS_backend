package server

import (
	"log/slog"
	"net/http"
)

type Server struct {
	httpServ *http.Server
}

func New(router http.Handler, port string) *Server {
	return &Server{
		httpServ: &http.Server{
			Addr:    ":" + port,
			Handler: router,
			//TODO configs
		},
	}
}

func (s *Server) Run() error {
	slog.Info("server is started", "ADDR", s.httpServ.Addr)
	return s.httpServ.ListenAndServe()
}
