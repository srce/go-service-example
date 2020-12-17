package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dzyanis/go-service-example/pkg/logger"
)

func RegisterHandlers(log *logger.Logger, s *Server) error {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("not found: %s", r.URL.String())
	})

	// CORS
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status": "OK"}`))
		if err != nil {
			log.Error(err)
		}
	}).Methods(http.MethodGet)

	// protected routes
	v1 := r.PathPrefix("/v1").Subrouter()

	// Users
	sr := v1.PathPrefix("/users").Subrouter()
	sr.HandleFunc("/{id}", s.users.Get).Methods(http.MethodPost)
	sr.HandleFunc("/{id}", s.users.Update).Methods(http.MethodPost)
	sr.HandleFunc("/", s.users.Create).Methods(http.MethodPost)
	sr.HandleFunc("/{id}", s.users.Delete).Methods(http.MethodDelete)

	// Wallets
	{
		sr := v1.PathPrefix("/wallets").Subrouter()
		sr.HandleFunc("/{id}", s.wallets.Create).Methods(http.MethodGet)
		sr.HandleFunc("/{id}", s.wallets.Update).Methods(http.MethodPost)
		sr.HandleFunc("/", s.wallets.Create).Methods(http.MethodPost)
		sr.HandleFunc("/{id}", s.wallets.Delete).Methods(http.MethodDelete)
	}

	s.Handler = r

	return nil
}
