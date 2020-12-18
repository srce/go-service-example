package transactions

import (
	"testing"

	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/currencies"
	"github.com/stretchr/testify/assert"
)

func Test_fee(t *testing.T) {
	cases := []struct {
		amount int64
		fee    float32
		expect int64
	}{
		{amount: 1, fee: 0, expect: 0},
		{amount: 0, fee: 1, expect: 0},
		{amount: 1, fee: 5, expect: 0},
		{amount: 1, fee: 100, expect: 1},
		{amount: 1, fee: 1000, expect: 10},
		{amount: 10, fee: 5, expect: 0},
		{amount: 100, fee: 5, expect: 5},
		{amount: 1000, fee: 5, expect: 50},
		{amount: 10000, fee: 1.25, expect: 125},
		{amount: 1000000, fee: 5, expect: 50000},
		{amount: 1000000, fee: 500, expect: 5000000},
		{amount: 1000000, fee: 33.33, expect: 333300},
	}

	for _, tc := range cases {
		got := fee(tc.amount, tc.fee)
		assert.Equal(t, tc.expect, got)
	}
}

func Test_calculate(t *testing.T) {
	type args struct {
		sender      *wallets.Wallet
		beneficiary *wallets.Wallet
		amount      int64
		feePercent  float32
	}

	type exp struct {
		trans  *Transaction
		fee    int64
		failed bool
	}

	cases := []struct {
		title string
		args  args
		exp   exp
	}{
		{
			title: "different currencies",
			args: args{
				sender:      &wallets.Wallet{Currency: currencies.USD.String()},
				beneficiary: &wallets.Wallet{Currency: currencies.EUR.String()},
			},
			exp: exp{failed: true},
		},

		{
			title: "insufficient funds",
			args: args{
				sender: &wallets.Wallet{
					Amount:   5,
					Currency: currencies.USD.String(),
				},
				beneficiary: &wallets.Wallet{
					Currency: currencies.USD.String(),
				},
				amount:     42,
				feePercent: 0,
			},
			exp: exp{failed: true},
		},

		{
			title: "the same account",
			args: args{
				sender: &wallets.Wallet{
					ID:       1,
					Amount:   1000,
					Currency: currencies.USD.String(),
				},
				beneficiary: &wallets.Wallet{
					ID:       1,
					Amount:   1000,
					Currency: currencies.USD.String(),
				},
			},
			exp: exp{failed: true},
		},

		{
			title: "success",
			args: args{
				sender: &wallets.Wallet{
					ID:       1,
					Amount:   1000,
					Currency: currencies.USD.String(),
				},
				beneficiary: &wallets.Wallet{
					ID:       2,
					Amount:   1000,
					Currency: currencies.USD.String(),
				},
				amount:     10,
				feePercent: 10,
			},
			exp: exp{
				fee: 1,
				trans: &Transaction{
					SenderID:      1,
					BeneficiaryID: 2,
					Amount:        10,
					Currency:      currencies.USD.String(),
				},
			},
		},
	}

	for _, tc := range cases {
		trans, fee, err := calculate(tc.args.sender, tc.args.beneficiary, tc.args.amount, tc.args.feePercent)
		if tc.exp.failed {
			assert.Error(t, err, tc.title)
		} else {
			assert.NoError(t, err, tc.title)
			assert.Equal(t, tc.exp.fee, fee, tc.title)
			assert.Equal(t, tc.exp.trans, trans, tc.title)
		}
	}
}
