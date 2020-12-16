package wallets

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dzyanis/go-service-example/pkg/logger"
)

func toJSON(w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("content-type", "application/json")
	return json.NewEncoder(w).Encode(obj)
}

type Controller struct {
	log *logger.Logger
}

func NewController(log *logger.Logger) *Controller {
	return &Controller{
		log: log,
	}
}

func (c *Controller) response(w http.ResponseWriter, obj interface{}) {
	w.WriteHeader(http.StatusOK)
	if err := toJSON(w, obj); err != nil {
		c.log.Error(err)
	}
}

func (c *Controller) error(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	resp := struct {
		Error string
	}{
		Error: err.Error(),
	}

	if err := toJSON(w, resp); err != nil {
		c.log.Error(err)
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.error(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.error(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.error(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.error(w, http.StatusNotImplemented, errors.New("not implemented"))
}
