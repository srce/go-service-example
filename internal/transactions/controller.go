package transactions

import (
	"net/http"
	"strconv"

	"github.com/dzyanis/go-service-example/pkg/controllers"
	"github.com/dzyanis/go-service-example/pkg/currencies"
	"github.com/dzyanis/go-service-example/pkg/logger"
)

type Controller struct {
	helper  controllers.JSONHelper
	service *Service
	log     *logger.Logger
}

func NewController(log *logger.Logger,
	service *Service, helper controllers.JSONHelper) *Controller {
	return &Controller{
		log:     log,
		service: service,
		helper:  helper,
	}
}

func (c *Controller) Transfer(w http.ResponseWriter, r *http.Request) {
	senderID, err := strconv.ParseInt(r.FormValue("sender_id"), 10, 64)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	beneficiaryID, err := strconv.ParseInt(r.FormValue("beneficiary_id"), 10, 64)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	currencyName := r.FormValue("currency")
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}
	currency, err := currencies.FromString(currencyName)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	err = c.service.Transfer(r.Context(), senderID, beneficiaryID, amount, currency)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}
	c.helper.Response(w, struct {
		Status string
	}{
		Status: "OK",
	})
}
