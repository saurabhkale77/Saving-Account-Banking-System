package enduser

import (
	"Saving-Account-Banking-System/app/specs"
	"Saving-Account-Banking-System/repository"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	UserRepo repository.UserStorer
}

// All User related funcs that processing From DB in Bussiness Logic
type Service interface {
	Authenticate(tknStr string) (user_id int, response string, err error)
	CreateLogin(ctx context.Context, req specs.CreateLoginRequest) (res string, err error)
	CreateSignup(ctx context.Context, req specs.CreateUser) (specs.Response, error)
	UpdateUser(ctx context.Context, req specs.UpdateUser, user_id int) (specs.UpdateUser, error)
	GetUser(ctx context.Context, user_id int) (specs.CreateUser, error)
	GetMyAccounts(ctx context.Context, user_id int) ([]specs.GetMyAccounts, error)
}

func NewService(UserRepo repository.UserStorer) Service {
	return &service{
		UserRepo: UserRepo,
	}
}

// All Bussiness Logic Related funcs of USER with logic here onwards=>
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

func (us *service) CreateLogin(ctx context.Context, req specs.CreateLoginRequest) (res string, err error) {

	tx, err := us.UserRepo.BeginTx(ctx)
	LoginMap, err := us.UserRepo.GetLoginDetails()
	if err != nil {
		return "", err
	}
	var jwtkey = os.Getenv("jwtkey")

	expectedPwd, ok := LoginMap[req.Username]
	bcrErr := bcrypt.CompareHashAndPassword([]byte(expectedPwd), []byte(req.Password))
	if !ok || bcrErr != nil {
		fmt.Println(ok)
		return "", fmt.Errorf("invalid credentials")
	}

	expirationTime := time.Now().Add(time.Minute * 100)
	//Getting Additional data from DB like user_id, role
	uid, role, err := us.UserRepo.TokenDetails(req.Username)

	claims := &specs.Claims{
		Username: req.Username,
		User_id:  uid,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtkey))

	if err != nil {
		return "", fmt.Errorf("error in parse token, %s", err)
	}

	defer func() {
		txErr := us.UserRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return tokenStr, nil
}

func (us *service) CreateSignup(ctx context.Context, req specs.CreateUser) (specs.Response, error) {
	tx, err := us.UserRepo.BeginTx(ctx)
	if err != nil {
		return specs.Response{}, err
	}

	response, err := us.UserRepo.AddUser(req)
	if err != nil {
		return specs.Response{}, err
	}

	defer func() {
		txErr := us.UserRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}

func (us *service) UpdateUser(ctx context.Context, req specs.UpdateUser, user_id int) (specs.UpdateUser, error) {
	tx, err := us.UserRepo.BeginTx(ctx)
	if err != nil {
		return specs.UpdateUser{}, err
	}

	//Hashing of password
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return specs.UpdateUser{}, fmt.Errorf(err.Error())
	}
	req.Password = string(hashPwd)

	response, err := us.UserRepo.UpdateUser(req, user_id)
	if err != nil {
		return specs.UpdateUser{}, err
	}

	defer func() {
		txErr := us.UserRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}

func (us *service) GetUser(ctx context.Context, user_id int) (specs.CreateUser, error) {
	tx, err := us.UserRepo.BeginTx(ctx)
	if err != nil {
		return specs.CreateUser{}, fmt.Errorf(err.Error())
	}

	response, err := us.UserRepo.GetUser(user_id)
	if err != nil {
		return specs.CreateUser{}, err
	}

	defer func() {
		txErr := us.UserRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}

func (us *service) GetMyAccounts(ctx context.Context, user_id int) ([]specs.GetMyAccounts, error) {
	tx, err := us.UserRepo.BeginTx(ctx)
	if err != nil {
		return []specs.GetMyAccounts{}, fmt.Errorf(err.Error())
	}

	response, err := us.UserRepo.GetMyAccounts(user_id)
	if err != nil {
		return []specs.GetMyAccounts{}, err
	}

	defer func() {
		txErr := us.UserRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()
	return response, nil
}
