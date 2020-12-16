package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/logger"
)

func RegisterHandlers(log *logger.Logger, srv *Server) (*Server, error) {
	r := mux.NewRouter()

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
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Wallet
	{
		route := v1.PathPrefix("/wallet").Subrouter()
		controller := wallets.NewController(log)
		route.HandleFunc("/{id}", controller.Create).Methods(http.MethodGet)
		route.HandleFunc("/{id}", controller.Update).Methods(http.MethodPost)
		route.HandleFunc("/", controller.Create).Methods(http.MethodPost)
		route.HandleFunc("/{id}", controller.Delete).Methods(http.MethodDelete)
	}

	srv.Handler = r

	return srv, nil
}
