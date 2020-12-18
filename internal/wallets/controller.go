package wallets

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dzyanis/go-service-example/pkg/controllers"
	"github.com/dzyanis/go-service-example/pkg/currencies"
	"github.com/dzyanis/go-service-example/pkg/logger"
	"github.com/gorilla/mux"
)

type Controller struct {
	log     *logger.Logger
	service *Service
	helper  controllers.JSONHelper
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
	userID, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}

	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}

	currencyName := r.FormValue("currency")
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}
	currency, err := currencies.FromString(currencyName)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
	}

	resp, err := c.service.Create(r.Context(), userID, amount, currency)
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
