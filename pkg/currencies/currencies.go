package currencies

type Currency string

func (c Currency) String() string {
	return string(c)
}

func (c Currency) Units() int {
	return currencies[c].Units
}

func (c Currency) Name() string {
	return currencies[c].Name
}

var currencies = map[Currency]struct {
	Name  string
	Units int
}{
	USD: {Name: "US Dollar", Units: 100},
	EUR: {Name: "Euro", Units: 100},
}

const (
	USD Currency = "USD"
	EUR          = "EUR"
)
