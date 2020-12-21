package money

import (
	"errors"
	"math"

	"github.com/dzyanis/go-service-example/pkg/currencies"
)

var (
	ErrAmountHasRemainder = errors.New("the number has a remainder")
	ErrDifferentCurrency  = errors.New("different currency")
)

type Money struct {
	units    int64
	currency currencies.Currency
}

func (m Money) Units() int64 {
	return m.units
}

func (m Money) Amount() float64 {
	return float64(m.units) / float64(m.currency.Units())
}

func (m Money) Currency() currencies.Currency {
	return m.currency
}

func (m Money) Percent(p float64) Money {
	a := (float64(m.units) / 100) * p
	a = Round(a)
	return Money{
		units:    int64(a),
		currency: m.currency,
	}
}

// Add returns new Money entity with sum of two.
func (m Money) Add(a Money) (Money, error) {
	if m.currency != a.currency {
		return Money{}, ErrDifferentCurrency
	}
	return Money{
		units:    m.units + a.units,
		currency: m.currency,
	}, nil
}

func Zero(currency currencies.Currency) Money {
	return Money{currency: currency}
}

func FromFloat64(amount float64, currency currencies.Currency) (Money, error) {
	units := int64(amount * float64(currency.Units()))

	abs := math.Abs(float64(units) / float64(currency.Units()))
	if math.Abs(amount)-abs > 0 {
		return Money{}, ErrAmountHasRemainder
	}

	return Money{
		units:    units,
		currency: currency,
	}, nil
}

// Round implements survey of rounding.
// See https://www.cockroachlabs.com/blog/rounding-implementations-in-go/
func Round(x float64) float64 {
	if math.IsNaN(x) {
		return x
	}
	if x == 0.0 {
		return x
	}
	roundFn := math.Ceil
	if math.Signbit(x) {
		roundFn = math.Floor
	}
	xOrig := x
	x -= math.Copysign(0.5, x)
	if x == 0 || math.Signbit(x) != math.Signbit(xOrig) {
		return math.Copysign(0.0, xOrig)
	}
	if x == xOrig-math.Copysign(1.0, x) {
		return xOrig
	}
	r := roundFn(x)
	if r != x {
		return r
	}
	return roundFn(x*0.5) * 2.0
}
