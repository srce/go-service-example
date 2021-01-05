package transactions

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzyanis/go-service-example/internal/users"
	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/currencies"
	"github.com/dzyanis/go-service-example/pkg/logger"
	"github.com/dzyanis/go-service-example/pkg/money"
)

const (
	CompanyBeneficiaryEmail = "company@beneficiary"
	CompanyFeePercent       = 1.5
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSameWallet        = errors.New("same wallet")
)

func calculate(sender, beneficiary *wallets.Wallet,
	amount money.Money, feePercent float64) (*Transaction, money.Money, error) {
	var (
		currency = amount.Currency()
		zero     = money.Zero(currency)
		fee      = amount.Percent(feePercent)
	)

	if sender.ID == beneficiary.ID {
		return nil, zero, ErrSameWallet
	}

	amountWithFee, err := amount.Add(fee)
	if err != nil {
		return nil, zero, fmt.Errorf("adding: %w", err)
	}

	if sender.Amount < amountWithFee.Units() {
		return nil, zero, ErrInsufficientFunds
	}
	sender.Amount -= amountWithFee.Units()
	beneficiary.Amount += amount.Units()
	return &Transaction{
		SenderID:      sender.ID,
		BeneficiaryID: beneficiary.ID,
		Amount:        amount.Units(),
		Currency:      amount.Currency().String(),
	}, fee, nil
}

type Service struct {
	log         *logger.Logger
	repoTrans   *Repository
	repoUsers   users.Repository
	repoWallets wallets.Repository
	uow         UOWStartFunc
}

func NewService(log *logger.Logger, repoTrans *Repository,
	repoUsers users.Repository, repoWallets wallets.Repository, uow UOWStartFunc) *Service {
	return &Service{
		log:         log,
		repoTrans:   repoTrans,
		repoUsers:   repoUsers,
		repoWallets: repoWallets,
		uow:         uow,
	}
}

func (s *Service) Transfer(ctx context.Context,
	senderID, beneficiaryID int64, amount money.Money) error {
	senderWallet, err := s.getWalletByUserID(ctx, senderID, amount.Currency())
	if err != nil {
		return fmt.Errorf("getting sender wallet: %w", err)
	}

	beneficiaryWallet, err := s.getWalletByUserID(ctx, beneficiaryID, amount.Currency())
	if err != nil {
		return fmt.Errorf("getting beneficiary wallet: %w", err)
	}

	trans, fees, err := calculate(senderWallet, beneficiaryWallet, amount, CompanyFeePercent)
	if err != nil {
		return fmt.Errorf("transferring: %w", err)
	}

	// Using Unit Of Work Patter for making one transaction with many repositories
	uow, err := s.uow()
	if err != nil {
		return fmt.Errorf("creating uow: %w", err)
	}
	defer uow.Rollback()

	_, err = uow.Trans().Create(ctx, trans)
	if err != nil {
		return fmt.Errorf("saving transaction: %w", err)
	}

	if fees.Units() > 0 {
		companyBeneficiary, err := uow.Users().GetByEmail(ctx, CompanyBeneficiaryEmail)
		if err != nil {
			return fmt.Errorf("getting company beneficiary: %w", err)
		}

		companyWallet, err := uow.Wallets().GetByUserIDAndCurrency(ctx,
			companyBeneficiary.ID, amount.Currency().String())
		if err != nil {
			return fmt.Errorf("getting wallet: %w", err)
		}

		_, err = uow.Trans().Create(ctx, &Transaction{
			SenderID:      senderWallet.ID,
			BeneficiaryID: companyWallet.ID,
			Amount:        fees.Units(),
			Currency:      fees.Currency().String(),
		})
		if err != nil {
			return fmt.Errorf("saving fee transaction: %w", err)
		}

		companyWallet.Amount += fees.Units()
		if err = uow.Wallets().Update(ctx, companyWallet); err != nil {
			return fmt.Errorf("update sender wallet: %w", err)
		}
	}
	// Updating wallets
	if err = uow.Wallets().Update(ctx, senderWallet); err != nil {
		return fmt.Errorf("update sender wallet: %w", err)
	}

	if err = uow.Wallets().Update(ctx, beneficiaryWallet); err != nil {
		return fmt.Errorf("update sender wallet: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return fmt.Errorf("uof commit: %w", err)
	}

	return nil
}

func (s *Service) getWalletByUserID(ctx context.Context,
	id int64, currency currencies.Currency) (*wallets.Wallet, error) {
	user, err := s.repoUsers.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	wallet, err := s.repoWallets.GetByUserIDAndCurrency(ctx, user.ID, currency.String())
	if err != nil {
		return nil, fmt.Errorf("getting wallet: %w", err)
	}

	return wallet, nil
}
