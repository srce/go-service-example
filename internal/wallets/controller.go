package wallets

import (
	"errors"
	"net/http"

	"github.com/dzyanis/go-service-example/pkg/controllers"
	"github.com/dzyanis/go-service-example/pkg/logger"
)

type Controller struct {
	log    *logger.Logger
	helper controllers.JSONHelper
}

func NewController(log *logger.Logger, helper controllers.JSONHelper) *Controller {
	return &Controller{
		log:    log,
		helper: helper,
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.helper.Error(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.helper.Error(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.helper.Error(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.helper.Error(w, http.StatusNotImplemented, errors.New("not implemented"))
}
