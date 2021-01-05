package transactions

import (
	"net/http"

	"github.com/dzyanis/go-service-example/internal/transactions"
	"github.com/dzyanis/go-service-example/pkg/controllers"
	"github.com/dzyanis/go-service-example/pkg/logger"
)

type Controller struct {
	helper  controllers.JSONHelper
	service transactions.Service
	log     *logger.Logger
}

func NewController(log *logger.Logger,
	service transactions.Service, helper controllers.JSONHelper) *Controller {
	return &Controller{
		log:     log,
		service: service,
		helper:  helper,
	}
}

func (c *Controller) Transfer(w http.ResponseWriter, r *http.Request) {
	senderID, err := controllers.FormInt64(r, "sender_id")
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	beneficiaryID, err := controllers.FormInt64(r, "beneficiary_id")
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	amount, err := controllers.FormValueAmount(r)
	if err != nil {
		c.helper.Error(w, http.StatusBadRequest, err)
		return
	}

	err = c.service.Transfer(r.Context(), senderID, beneficiaryID, amount)
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
