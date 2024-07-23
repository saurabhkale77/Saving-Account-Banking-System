package repository

import (
	"Saving-Account-Banking-System/app/specs"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type AdminStore struct {
	BaseRepository
}

type AdminStorer interface {
	RepositoryTrasanctions

	ListUsers(ctx context.Context) ([]specs.Response, error)
	ListBranches(ctx context.Context) ([]specs.BranchDetails, error)
	UpdateUserInfo(req specs.UpdateUserInfo) (specs.UpdateUserInfo, error)
	CreateAccount(req specs.CreateAccountReq) (specs.CreateAccountReq, error)
}

func NewAdminRepo(db *sql.DB) AdminStorer {
	return &AdminStore{
		BaseRepository: BaseRepository{db},
	}
}

const (
	getUserDetailsQuery   string = `SELECT user_id, name, address, email, password, mobile, role FROM user ORDER BY user_id DESC`
	userUpdateQuery       string = `UPDATE user SET name=?, address=?,email=?, password=?, mobile=?,role=?, updated_at=? WHERE user_id=?`
	getBranchDetailsQuery string = `SELECT id, name, location FROM branch ORDER BY id`
)

func (db *AdminStore) CreateAccount(req specs.CreateAccountReq) (specs.CreateAccountReq, error) {

	//To get Existing value
	var count int64
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	row := QueryExecuter.QueryRow(getTopAccNo)
	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			count = 100
		}
	}
	//For Inserting
	stmt, err := QueryExecuter.Prepare(insertAccountDetailsQuery)
	if err != nil {
		return specs.CreateAccountReq{}, fmt.Errorf("errror While inserting CreateAccount data in db")
	}
	acc_no := (count + 1)
	stmt.Exec(acc_no, req.User_id, req.Branch_id, req.Account_type, req.Balance, time.Now().Unix(), time.Now().Unix())

	res := specs.CreateAccountReq{
		Account_no:   int(acc_no),
		Account_type: req.Account_type,
		Balance:      req.Balance,
		Branch_id:    req.Branch_id,
		User_id:      req.User_id,
	}
	return res, nil
}

func (db *AdminStore) ListUsers(ctx context.Context) ([]specs.Response, error) {
	var result []specs.Response

	//To get All user values
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	rows, err := QueryExecuter.Query(getUserDetailsQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var res specs.Response
		if err := rows.Scan(&res.User_id, &res.Name, &res.Address, &res.Email, &res.Password, &res.Mobile, &res.Role); err != nil {
			log.Print("error while scanning row: ", err)
			continue
		}
		result = append(result, res)
	}
	return result, nil
}

func (db *AdminStore) UpdateUserInfo(req specs.UpdateUserInfo) (specs.UpdateUserInfo, error) {
	// For Updating User Info.
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	stmt, err := QueryExecuter.Prepare(userUpdateQuery)
	if err != nil {
		return specs.UpdateUserInfo{}, fmt.Errorf("error while updating user data in db: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Name, req.Address, req.Email, req.Password, req.Mobile, req.Role, time.Now().Unix(), req.User_id)
	if err != nil {
		return specs.UpdateUserInfo{}, fmt.Errorf("error executing updateUserInfo at Admin side: %v", err)
	}

	res := specs.UpdateUserInfo{
		User_id:  req.User_id,
		Name:     req.Name,
		Address:  req.Address,
		Email:    req.Email,
		Password: req.Password,
		Mobile:   req.Mobile,
		Role:     req.Role,
	}
	return res, nil
}

func (db *AdminStore) ListBranches(ctx context.Context) ([]specs.BranchDetails, error) {
	var result []specs.BranchDetails

	QueryExecuter := db.initiateQueryExecutor(db.DB)
	rows, err := QueryExecuter.Query(getBranchDetailsQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var res specs.BranchDetails
		if err := rows.Scan(&res.Id, &res.Name, &res.Location); err != nil {
			log.Print("error while getting branches: ", err)
			continue
		}
		result = append(result, res)
	}
	return result, nil
}
