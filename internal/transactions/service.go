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

type Service struct {
	log         *logger.Logger
	repo        *Repository
	repoUsers   *users.Repository
	repoWallets *wallets.Repository
}

func NewService(log *logger.Logger,
	repo *Repository, repoUsers *users.Repository, repoWallets *wallets.Repository) *Service {
	return &Service{
		log:         log,
		repo:        repo,
		repoUsers:   repoUsers,
		repoWallets: repoWallets,
	}
}

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

func (s *Service) Transfer(ctx context.Context,
	senderID, beneficiaryID int64, amount float64, currency currencies.Currency) error {

	senderWallet, err := s.getWalletByUserID(ctx, senderID, currency)
	if err != nil {
		return fmt.Errorf("getting sender wallet: %w", err)
	}

	beneficiaryWallet, err := s.getWalletByUserID(ctx, beneficiaryID, currency)
	if err != nil {
		return fmt.Errorf("getting beneficiary wallet: %w", err)
	}

	amountUnits := units(amount, currency)

	// TODO: needs db transaction
	trans, fees, err := calculate(senderWallet, beneficiaryWallet, amountUnits, CompanyFeePercent)
	if err != nil {
		return fmt.Errorf("transfering: %w", err)
	}

	transID, err := s.repo.Create(ctx, trans)
	if err != nil {
		return fmt.Errorf("saving transaction: %w", err)
	}
	s.log.Printf("transaction %d created", transID)

	if fees > 0 {
		companyWallet, err := s.getWalletByEmail(ctx, CompanyBeneficiaryEmail, currency)
		if err != nil {
			return fmt.Errorf("getting company wallet: %w", err)
		}

		transID, err = s.repo.Create(ctx, &Transaction{
			SenderID:      senderWallet.ID,
			BeneficiaryID: companyWallet.ID,
			Amount:        fees,
			Currency:      senderWallet.Currency,
		})
		if err != nil {
			return fmt.Errorf("saving fee transaction: %w", err)
		}
		s.log.Printf("transaction %d created", transID)
	}

	// Updating wallets
	err = s.repoWallets.Update(ctx, senderWallet)
	if err != nil {
		return fmt.Errorf("update sender wallet: %w", err)
	}

	err = s.repoWallets.Update(ctx, beneficiaryWallet)
	if err != nil {
		return fmt.Errorf("update sender wallet: %w", err)
	}

	return nil
}

func (s *Service) getWalletByUserID(ctx context.Context, id int64, currency currencies.Currency) (*wallets.Wallet, error) {
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

func (s *Service) getWalletByEmail(ctx context.Context, email string, currency currencies.Currency) (*wallets.Wallet, error) {
	companyBeneficiary, err := s.repoUsers.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("getting company beneficiary: %w", err)
	}

	wallet, err := s.repoWallets.GetByUserIDAndCurrency(ctx, companyBeneficiary.ID, currency.String())
	if err != nil {
		return nil, fmt.Errorf("getting wallet: %w", err)
	}

	return wallet, nil
}
