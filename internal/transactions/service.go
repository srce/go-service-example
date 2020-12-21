package transactions

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzyanis/go-service-example/internal/users"
	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/currencies"
	"github.com/dzyanis/go-service-example/pkg/logger"
)

const (
	CompanyBeneficiaryEmail = "company@beneficiary"
	CompanyFeePercent       = 1.5
)

var (
	ErrDifferentCurrency = errors.New("different currency")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSameWallet        = errors.New("same wallet")
)

func fee(amount int64, percent float32) int64 {
	return int64(float64(amount) / 100 * float64(percent))
}

func units(a float64, currency currencies.Currency) int64 {
	u := a * float64(currency.Units())
	return int64(u)
}

func calculate(sender, beneficiary *wallets.Wallet, amount int64, feePercent float32) (*Transaction, int64, error) {
	if sender.Currency != beneficiary.Currency {
		return nil, 0, ErrDifferentCurrency
	}

	if sender.ID == beneficiary.ID {
		return nil, 0, ErrSameWallet
	}

	feeAmount := fee(amount, feePercent)
	amountWithFee := amount + feeAmount

	if sender.Amount < amountWithFee {
		return nil, 0, ErrInsufficientFunds
	}
	sender.Amount -= amountWithFee
	beneficiary.Amount += amount
	return &Transaction{
		SenderID:      sender.ID,
		BeneficiaryID: beneficiary.ID,
		Amount:        amount,
		Currency:      sender.Currency,
	}, feeAmount, nil
}

type Service struct {
	log         *logger.Logger
	repoTrans   *Repository
	repoUsers   *users.Repository
	repoWallets *wallets.Repository
	uow         UOWStartFunc
}

func NewService(log *logger.Logger, repoTrans *Repository,
	repoUsers *users.Repository, repoWallets *wallets.Repository, uow UOWStartFunc) *Service {
	return &Service{
		log:         log,
		repoTrans:   repoTrans,
		repoUsers:   repoUsers,
		repoWallets: repoWallets,
	}
}

func (s *Service) Transfer(ctx context.Context, senderID, beneficiaryID int64,
	amount float64, currency currencies.Currency) error {
	senderWallet, err := s.getWalletByUserID(ctx, senderID, currency)
	if err != nil {
		return fmt.Errorf("getting sender wallet: %w", err)
	}

	beneficiaryWallet, err := s.getWalletByUserID(ctx, beneficiaryID, currency)
	if err != nil {
		return fmt.Errorf("getting beneficiary wallet: %w", err)
	}

	amountUnits := units(amount, currency)

	trans, fees, err := calculate(senderWallet, beneficiaryWallet, amountUnits, CompanyFeePercent)
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

	if fees > 0 {
		companyBeneficiary, err := uow.Users().GetByEmail(ctx, CompanyBeneficiaryEmail)
		if err != nil {
			return fmt.Errorf("getting company beneficiary: %w", err)
		}

		companyWallet, err := uow.Wallets().GetByUserIDAndCurrency(ctx, companyBeneficiary.ID, currency.String())
		if err != nil {
			return fmt.Errorf("getting wallet: %w", err)
		}

		_, err = uow.Trans().Create(ctx, &Transaction{
			SenderID:      senderWallet.ID,
			BeneficiaryID: companyWallet.ID,
			Amount:        fees,
			Currency:      senderWallet.Currency,
		})
		if err != nil {
			return fmt.Errorf("saving fee transaction: %w", err)
		}

		companyWallet.Amount += fees
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
