package controllers

import (
	"encoding/json"
	"net/http"
)

type JSONHelper struct{}

func (j JSONHelper) Response(w http.ResponseWriter, obj interface{}) error {
	w.WriteHeader(http.StatusOK)
	return toJSON(w, obj)
}

func (j JSONHelper) Error(w http.ResponseWriter, status int, err error) error {
	w.WriteHeader(status)
	resp := struct {
		Error string
	}{
		Error: err.Error(),
	}

	return toJSON(w, resp)
}

func toJSON(w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("content-type", "application/json")
	return json.NewEncoder(w).Encode(obj)
}
