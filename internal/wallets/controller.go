package wallets

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/dzyanis/go-service-example/pkg/controllers"
	"github.com/dzyanis/go-service-example/pkg/logger"
)

type Controller struct {
	helper  controllers.JSONHelper
	log     *logger.Logger
	service *Service
}

func NewController(log *logger.Logger,
	service *Service, helper controllers.JSONHelper) *Controller {
	return &Controller{
		log:     log,
		service: service,
		helper:  helper,
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := controllers.FormInt64(r, "user_id")
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	amount, err := controllers.FormValueAmount(r)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	resp, err := c.service.Create(r.Context(), userID, amount)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}
	c.helper.Response(w, resp)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	c.helper.Error(w, http.StatusNotImplemented, controllers.ErrNotImplemented)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		c.helper.Error(w, http.StatusBadRequest, fmt.Errorf("id params: %w", controllers.ErrRequredValue))
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}

	if err := c.service.Delete(r.Context(), userID); err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}
	c.helper.Response(w, struct {
		Status string
	}{
		Status: "OK",
	})
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		c.helper.Error(w, http.StatusBadRequest, fmt.Errorf("id params: %w", controllers.ErrRequredValue))
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}

	resp, err := c.service.Get(r.Context(), userID)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}
	c.helper.Response(w, resp)
}
