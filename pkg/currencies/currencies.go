package currencies

import "errors"

type Currency string

func (c Currency) String() string {
	return string(c)
}

func (c Currency) Units() int64 {
	return currencies[c].Units
}

func (c Currency) Name() string {
	return currencies[c].Name
}

var currencies = map[Currency]struct {
	Name  string
	Units int64
}{
	USD: {Name: "US Dollar", Units: 100},
	EUR: {Name: "Euro", Units: 100},
}

const (
	USD Currency = "USD"
	EUR          = "EUR"
)

var ErrUnknownCurrency = errors.New("unknown currency")

func FromString(s string) (Currency, error) {
	cur := Currency(s)
	if _, ok := currencies[cur]; !ok {
		return cur, ErrUnknownCurrency
	}
	return cur, nil
}
