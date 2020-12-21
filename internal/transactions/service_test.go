package transactions

import (
	"testing"

	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/currencies"
	"github.com/dzyanis/go-service-example/pkg/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func usd(tb testing.TB, value float64) money.Money {
	m, err := money.FromFloat64(value, currencies.USD)
	require.NoErrorf(tb, err, "usd %f", value)
	return m
}

func Test_calculate(t *testing.T) {
	type args struct {
		sender      *wallets.Wallet
		beneficiary *wallets.Wallet
		amount      money.Money
		feePercent  float64
	}

	type exp struct {
		trans  *Transaction
		fee    money.Money
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
				amount:     usd(t, 42),
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
					Amount:   10000,
					Currency: currencies.USD.String(),
				},
				beneficiary: &wallets.Wallet{
					ID:       2,
					Amount:   10000,
					Currency: currencies.USD.String(),
				},
				amount:     usd(t, 10),
				feePercent: 10,
			},
			exp: exp{
				fee: usd(t, 1),
				trans: &Transaction{
					SenderID:      1,
					BeneficiaryID: 2,
					Amount:        1000,
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
