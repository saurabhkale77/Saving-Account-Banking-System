package app

import (
	"Saving-Account-Banking-System/app/account"
	"Saving-Account-Banking-System/app/admin"
	"Saving-Account-Banking-System/app/enduser"
	"Saving-Account-Banking-System/repository"
	"database/sql"
)

type Dependencies struct {
	UserService    enduser.Service
	AccountService account.Service
	AdminService   admin.Service
}

func NewServices(db *sql.DB) Dependencies {

	//Initialize repo dependencies
	UserRepo := repository.NewUserRepo(db)
	AccountRepo := repository.NewAccountRepo(db)
	AdminRepo := repository.NewAdminRepo(db)

	//Initialize Service Dependencies
	userService := enduser.NewService(UserRepo)
	accountService := account.NewService(AccountRepo)
	adminService := admin.NewService(AdminRepo)

	return Dependencies{
		UserService:    userService,
		AccountService: accountService,
		AdminService:   adminService,
	}
}
