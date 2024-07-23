package account

import (
	"Saving-Account-Banking-System/app/specs"
	"Saving-Account-Banking-System/repository"
	"context"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type service struct {
	AccountRepo repository.AccountStorer
}

type Service interface {
	Authenticate(tknStr string) (user_id int, response string, err error)

	DeleteAccount(ctx context.Context, req specs.DeleteAccountReq, user_id int) (specs.DeleteAccount, error)
	DepositMoney(ctx context.Context, req specs.Transaction, user_id int) (specs.TransactionResponse, error)
	WithdrawalMoney(ctx context.Context, req specs.Transaction, user_id int) (specs.TransactionResponse, error)
	ViewBalance(ctx context.Context, req specs.TransactionResponse, user_id int) (specs.TransactionResponse, error)
}

func NewService(AccountRepo repository.AccountStorer) Service {
	return &service{
		AccountRepo: AccountRepo,
	}
}

// All Account related bussiness logic here onwards=>
func (us *service) Authenticate(tknStr string) (user_id int, response string, err error) {

	jwtkey := []byte(os.Getenv("jwtkey"))
	claims := &specs.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("error in parsing claims")
	}

	if !tkn.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}
	return claims.User_id, claims.Username, nil
}

func (as *service) DeleteAccount(ctx context.Context, req specs.DeleteAccountReq, user_id int) (specs.DeleteAccount, error) {

	tx, _ := as.AccountRepo.BeginTx(ctx)

	response, err := as.AccountRepo.DeleteAccount(req, user_id)
	if err != nil {
		return specs.DeleteAccount{}, err
	}

	defer func() {
		txErr := as.AccountRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}

func (as *service) DepositMoney(ctx context.Context, req specs.Transaction, user_id int) (response specs.TransactionResponse, err error) {

	tx, _ := as.AccountRepo.BeginTx(ctx)
	defer func() {
		txErr := as.AccountRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	response, err = as.AccountRepo.DepositMoney(req, user_id)
	if err != nil {
		return specs.TransactionResponse{}, err
	}

	return response, nil
}

func (as *service) WithdrawalMoney(ctx context.Context, req specs.Transaction, user_id int) (specs.TransactionResponse, error) {

	tx, _ := as.AccountRepo.BeginTx(ctx)

	response, err := as.AccountRepo.WithdrawalMoney(req, user_id)
	if err != nil {
		return specs.TransactionResponse{}, err
	}

	defer func() {
		txErr := as.AccountRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}

func (as *service) ViewBalance(ctx context.Context, req specs.TransactionResponse, user_id int) (specs.TransactionResponse, error) {

	tx, _ := as.AccountRepo.BeginTx(ctx)

	response, err := as.AccountRepo.ViewBalance(req, user_id)
	if err != nil {
		return specs.TransactionResponse{}, err
	}

	defer func() {
		txErr := as.AccountRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}
