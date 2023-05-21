package apiserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/eqimd/transbyte-site/internal/equivcheck"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CheckEquivalenceRequest struct {
	FirstClassCode  string
	SecondClassCode string
}

type CheckEquivalenceResponse struct {
	Output string
}

type APIServer struct {
	config *Config
	router *chi.Mux
}

func NewServer(cfg *Config) *APIServer {
	server := &APIServer{
		config: cfg,
		router: chi.NewRouter(),
	}

	server.setupRouter()

	return server
}

func (s *APIServer) Start() error {
	return http.ListenAndServe(s.config.BindAddr, s.router)
	// return http.Serve(autocert.NewListener(s.config.BindAddr), s.router)
}

func (s *APIServer) setupRouter() {
	s.router.Use(middleware.Logger)

	s.setupRouterHandlers()
}

func (s *APIServer) setupRouterHandlers() {
	s.router.Post(RouteRoot, s.handleRoot())
}

func (s *APIServer) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		equivReq := new(CheckEquivalenceRequest)
		if err := json.Unmarshal(b, equivReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		output, _ := equivcheck.CheckEquivalence(equivReq.FirstClassCode, equivReq.SecondClassCode)
		resp := CheckEquivalenceResponse{
			Output: output,
		}

		b, _ = json.Marshal(resp)
		_, _ = w.Write(b)
	}
}
