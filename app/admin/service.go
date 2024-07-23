package admin

import (
	"Saving-Account-Banking-System/app/specs"
	"Saving-Account-Banking-System/repository"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	AdminRepo repository.AdminStorer
}
type Service interface {
	Authenticate(tknStr string) (response string, err error)
	ListUsers(ctx context.Context) ([]specs.Response, error)
	ListBranches(ctx context.Context) ([]specs.BranchDetails, error)
	UpdateUser(ctx context.Context, req specs.UpdateUserInfo) (specs.UpdateUserInfo, error)
	CreateAccount(ctx context.Context, req specs.CreateAccountReq) (specs.CreateAccountReq, error)
}

func NewService(AdminRepo repository.AdminStorer) Service {
	return &service{
		AdminRepo: AdminRepo,
	}
}

func (adm *service) Authenticate(tknStr string) (response string, err error) {

	jwtkey := []byte(os.Getenv("jwtkey"))
	claims := &specs.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return "", fmt.Errorf("error in parsing claims")
	}
	if !tkn.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if !strings.EqualFold(claims.Role, "admin") {
		return "", fmt.Errorf("access denied")
	}
	return claims.Username, nil
}

func (adm *service) ListUsers(ctx context.Context) ([]specs.Response, error) {
	tx, err := adm.AdminRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	response, err := adm.AdminRepo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		txErr := adm.AdminRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}

func (us *service) UpdateUser(ctx context.Context, req specs.UpdateUserInfo) (specs.UpdateUserInfo, error) {
	tx, err := us.AdminRepo.BeginTx(ctx)
	if err != nil {
		return specs.UpdateUserInfo{}, fmt.Errorf(err.Error())
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return specs.UpdateUserInfo{}, fmt.Errorf(err.Error())
	}
	req.Password = string(hashPwd)

	response, err := us.AdminRepo.UpdateUserInfo(req)
	if err != nil {
		return specs.UpdateUserInfo{}, err
	}

	defer func() {
		txErr := us.AdminRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}

func (adm *service) CreateAccount(ctx context.Context, req specs.CreateAccountReq) (specs.CreateAccountReq, error) {

	response, err := adm.AdminRepo.CreateAccount(req)

	if err != nil {
		return specs.CreateAccountReq{}, err
	}
	// specs.SendMail(response)

	return response, nil
}

func (adm *service) ListBranches(ctx context.Context) ([]specs.BranchDetails, error) {
	tx, err := adm.AdminRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	response, err := adm.AdminRepo.ListBranches(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		txErr := adm.AdminRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}
