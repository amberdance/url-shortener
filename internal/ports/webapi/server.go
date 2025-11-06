package webapi

import (
	"log"
	"net/http"

	infr "github.com/amberdance/url-shortener/internal/infrastructure/storage"
)

type Server struct {
	addr    string
	httpSrv *http.Server
}

func NewServer(addr string) *Server {
	storage := infr.NewInMemoryStorage()
	handler := NewHandler(storage, addr)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return &Server{
		addr: addr,
		httpSrv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("Server is running on %s\n", s.addr)
	return s.httpSrv.ListenAndServe()
}
