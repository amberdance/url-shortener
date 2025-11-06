package webapi

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type Server struct {
	addr    string
	storage map[string]string
	httpSrv *http.Server
}

func NewServer(addr string) *Server {
	s := &Server{
		addr:    addr,
		storage: make(map[string]string),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)

	s.httpSrv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s
}

func (s *Server) Run() error {
	log.Printf("Server is running on %s\n", s.addr)
	return s.httpSrv.ListenAndServe()
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodGet:
		s.handleGet(w, r)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := generateShortID()
	s.storage[id] = originalURL

	shortURL := fmt.Sprintf("http://localhost:%s/", s.addr) + id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, exists := s.storage[id]
	if !exists {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateShortID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
