package repository

import (
	"Saving-Account-Banking-System/app/specs"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	BaseRepository
}

const (
	insertUserQuery      string = "INSERT INTO user VALUES(?,?,?,?,?,?,?,?,?)"
	updateUserQuery      string = `UPDATE user SET name=?, address=?, password=?, mobile=?, updated_at=? WHERE user_id=?`
	getTokenDetailsQuery string = "SELECT user_id,role FROM user where email=?"
	getLoginDetails      string = "SELECT email, password FROM user"
	getCountofUser       string = "SELECT COUNT(user_id) FROM user"
	getUser              string = "SELECT name, address, email, password, mobile, role FROM user where user_id=?"
	GetMyAccounts        string = "SELECT acc_no, branch_id, acc_type, balance FROM account WHERE user_id=? ORDER BY created_at "
)

type UserStorer interface {
	RepositoryTrasanctions

	GetLoginDetails() (response map[string]string, err error)
	AddUser(req specs.CreateUser) (specs.Response, error)
	UpdateUser(req specs.UpdateUser, user_id int) (specs.UpdateUser, error)
	GetUser(user_id int) (specs.CreateUser, error)
	GetMyAccounts(user_id int) ([]specs.GetMyAccounts, error)

	TokenDetails(email string) (user_id int, role string, err error)
}

func NewUserRepo(db *sql.DB) UserStorer {

	return &UserStore{
		BaseRepository: BaseRepository{db},
	}
}

func (db *UserStore) GetLoginDetails() (response map[string]string, err error) {
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	rows, err := QueryExecuter.Query(getLoginDetails)
	if err != nil {
		log.Print("error while fetching login details from database: ", err)
		return nil, fmt.Errorf("error while fetching login details from database")
	}

	LoginMap := make(map[string]string)
	for rows.Next() {
		var email, pwd string
		if err := rows.Scan(&email, &pwd); err != nil {
			log.Print("error while scanning row: ", err)
			continue
		}
		LoginMap[email] = pwd
	}

	return LoginMap, nil

}

func (db *UserStore) AddUser(req specs.CreateUser) (specs.Response, error) {

	//To get Existing value
	var count int64
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	row := QueryExecuter.QueryRow(getCountofUser)
	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			count = 0
		}
		return specs.Response{}, fmt.Errorf("something went wrong")
	}
	//Hashing of Password
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Println("error !! while hashing password, Error : ", err)
		return specs.Response{}, fmt.Errorf("error !! while hashing password, Error : %v", err)
	}

	//For Inserting
	stmt, err := QueryExecuter.Prepare(insertUserQuery)
	if err != nil {
		return specs.Response{}, fmt.Errorf("errror While inserting sign-up data in db")
	}
	user_id := count + 1
	stmt.Exec(user_id, req.Name, req.Address, req.Email, string(hashPwd), req.Mobile, strings.ToLower(req.Role), time.Now().Unix(), time.Now().Unix())

	res := specs.Response{
		User_id:  int(user_id),
		Name:     req.Name,
		Address:  req.Address,
		Email:    req.Email,
		Password: req.Password,
		Mobile:   req.Mobile,
		Role:     req.Role,
	}
	return res, nil
}

func (db *UserStore) UpdateUser(req specs.UpdateUser, user_id int) (specs.UpdateUser, error) {
	// For Updating User Info.
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	stmt, err := QueryExecuter.Prepare(updateUserQuery)
	if err != nil {
		return specs.UpdateUser{}, fmt.Errorf("error while updating user data in db: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Name, req.Address, req.Password, req.Mobile, time.Now().Unix(), user_id)
	if err != nil {
		return specs.UpdateUser{}, fmt.Errorf("error executing update statement at user level, Error: %v", err)
	}

	res := specs.UpdateUser{
		User_id:  user_id,
		Name:     req.Name,
		Address:  req.Address,
		Mobile:   req.Mobile,
		Password: req.Password,
	}
	return res, nil
}

func (db *UserStore) GetUser(user_id int) (specs.CreateUser, error) {

	QueryExecuter := db.initiateQueryExecutor(db.DB)
	rows, err := QueryExecuter.Query(getUser, user_id)
	if err != nil {
		log.Println(err)
		return specs.CreateUser{}, err
	}

	var res specs.CreateUser
	for rows.Next() {
		if err := rows.Scan(&res.Name, &res.Address, &res.Email, &res.Password, &res.Mobile, &res.Role); err != nil {
			log.Print("error while scanning row: ", err)
			continue
		}
	}
	return res, nil
}

func (db *UserStore) GetMyAccounts(user_id int) ([]specs.GetMyAccounts, error) {
	var result []specs.GetMyAccounts
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	rows, err := QueryExecuter.Query(GetMyAccounts, user_id)
	if err != nil {
		log.Println(err)
		return []specs.GetMyAccounts{}, err
	}

	for rows.Next() {
		var res specs.GetMyAccounts
		if err := rows.Scan(&res.Acc_no, &res.Branch_id, &res.Acc_Type, &res.Balance); err != nil {
			log.Print("error while scanning row: ", err)
			continue
		}
		result = append(result, res)
	}
	return result, nil
}

func (db *UserStore) TokenDetails(email string) (user_id int, role string, err error) {
	QueryExecuter := db.initiateQueryExecutor(db.DB)
	row := QueryExecuter.QueryRow(getTokenDetailsQuery, email)
	err = row.Scan(&user_id, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("record not found")
		}
		return 0, "", fmt.Errorf("something went wrong")
	}
	return user_id, role, nil
}
