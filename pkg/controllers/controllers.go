package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dzyanis/go-service-example/pkg/currencies"
	"github.com/dzyanis/go-service-example/pkg/money"
)

var (
	ErrRequredValue   = errors.New("value is required")
	ErrNotImplemented = errors.New("not implemented")
)

func FormValueAmount(r *http.Request) (money.Money, error) {
	currency, err := currencies.FromString(r.FormValue("currency"))
	if err != nil {
		return money.Money{}, fmt.Errorf("form currency: %w", err)
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		return money.Money{}, fmt.Errorf("form amount: %w", err)
	}

	return money.FromFloat64(amount, currency)
}

func FormInt64(r *http.Request, key string) (int64, error) {
	value, err := strconv.ParseInt(r.FormValue(key), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("form value %s: %w", key, err)
	}

	return value, nil
}
