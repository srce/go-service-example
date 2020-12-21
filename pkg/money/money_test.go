package money

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dzyanis/go-service-example/pkg/currencies"
)

func Test(t *testing.T) {
	negZero := float64(-0)
	tests := []struct {
		arg float64
		exp float64
	}{
		{-0.49999999999999994, negZero}, // -0.5+epsilon
		{-0.5, negZero},
		{-0.5000000000000001, -1}, // -0.5-epsilon
		{0, 0},
		{0.49999999999999994, 0}, // 0.5-epsilon
		{0.5, 0},
		{0.5000000000000001, 1},  // 0.5+epsilon
		{1.390671161567e-309, 0}, // denormal
		{2.2517998136852485e+15, 2.251799813685248e+15}, // 1 bit fraction
		{4.503599627370497e+15, 4.503599627370497e+15},  // large integer
		{math.Inf(-1), math.Inf(-1)},
		{math.Inf(1), math.Inf(1)},
		{negZero, negZero},
	}

	for i, tc := range tests {
		r := Round(tc.arg)
		assert.Equalf(t, tc.exp, r, "case %d: %f", i+1, tc.arg)
	}
}

func TestPercent(t *testing.T) {
	cents := func(i int64) Money {
		return Money{units: i, currency: currencies.USD}
	}
	cases := []struct {
		amount Money
		fee    float64
		expect Money
	}{
		{amount: cents(1), fee: 0, expect: cents(0)},
		{amount: cents(0), fee: 1, expect: cents(0)},
		{amount: cents(1), fee: 1, expect: cents(0)},
		{amount: cents(1), fee: 5, expect: cents(0)},
		{amount: cents(1), fee: 100, expect: cents(1)},
		{amount: cents(-1), fee: 100, expect: cents(-1)},
		{amount: cents(1), fee: 1000, expect: cents(10)},
		{amount: cents(10), fee: 5, expect: cents(0)},
		{amount: cents(100), fee: 5, expect: cents(5)},
		{amount: cents(1000), fee: 5, expect: cents(50)},
		{amount: cents(10000), fee: 1.25, expect: cents(125)},
		{amount: cents(1000000), fee: 5, expect: cents(50000)},
		{amount: cents(1000000), fee: 500, expect: cents(5000000)},
		{amount: cents(1000000), fee: 33.33, expect: cents(333300)},
		{amount: cents(-1000000), fee: 33.33, expect: cents(-333300)},
	}

	for _, tc := range cases {
		got := tc.amount.Percent(tc.fee)
		assert.Equal(t, tc.expect, got)
	}
}

func TestFromFloat64(t *testing.T) {
	cases := []struct {
		amount   float64
		currency currencies.Currency
		expect   Money
	}{
		{
			amount: 0, currency: currencies.USD,
			expect: Money{units: 0, currency: currencies.USD},
		},
		{
			amount: 1, currency: currencies.USD,
			expect: Money{units: 100, currency: currencies.USD},
		},
		{
			amount: 10, currency: currencies.USD,
			expect: Money{units: 1000, currency: currencies.USD},
		},
		{
			amount: -10, currency: currencies.USD,
			expect: Money{units: -1000, currency: currencies.USD},
		},
	}

	for _, tc := range cases {
		got, err := FromFloat64(tc.amount, tc.currency)
		assert.NoError(t, err)
		assert.Equal(t, tc.expect, got)
	}

	_, err := FromFloat64(float64(3.333), currencies.USD)
	assert.Error(t, err)
}
