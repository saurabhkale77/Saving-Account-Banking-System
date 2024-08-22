package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitializeDB() (*sql.DB, error) {

	// Load .env vars
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Print("error !! while loading env vars !! Due to : ", err)
	// }

	database, err := sql.Open("sqlite3", "repository/bank.db")
	if err != nil {
		log.Print("error !! while creating the database !!")
		return nil, err
	}
	db = database

	statement, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS user(
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(30) NOT NULL,
		address VARCHAR(50) NOT NULL,
		email VARCHAR(30) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		mobile CHAR(10) UNIQUE NOT NULL,
		role VARCHAR(10) NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	)	
	`)
	if err != nil {
		log.Print("error !! while creating user table !! Due to : ", err)
		return nil, err
	}
	statement.Exec()

	statement, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS branch(
		id INTEGER PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		location VARCHAR(30) NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	)
	`)
	if err != nil {
		log.Print("error !! while creating branch table !! Due to : ", err)
		return nil, err
	}
	statement.Exec()

	statement, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS account(
		acc_no INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		branch_id INTEGER NOT NULL,
		acc_type VARCHAR(10) NOT NULL,
		balance FLOAT NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES user(user_id),
		FOREIGN KEY(branch_id) REFERENCES branch(id)
	)
	`)
	if err != nil {
		log.Print("error !! while creating account table !! Due to : ", err)
		return nil, err
	}
	statement.Exec()

	statement, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS tbltransaction(
		transaction_id INTEGER PRIMARY KEY,
		acc_no INTEGER NOT NULL,
		amount FLOAT NOT NULL,
		type VARCHAR(10) NOT NULL,
		updated_at INTEGER NOT NULL,
		FOREIGN KEY(acc_no) REFERENCES account(acc_no)
	)
	`)
	if err != nil {
		log.Print("error !! while creating tbltransaction table !! Due to : ", err)
		return nil, err
	}
	statement.Exec()

	fmt.Println("Successfully Initialized Database !!")

	statement.Close()

	return db, nil
}

func InsertSeedData() {
	database, err := sql.Open("sqlite3", "repository/bank.db")
	if err != nil {
		log.Print("error !! while Connecting with database !!")
		return
	}
	defer db.Close()

	//Seed Data for User
	stmt, err := database.Prepare(`INSERT INTO user VALUES
	(1, "Saurabh", "Pune", "abc@gmail.com", "Pass@123", "9595601925", "customer", 1707391842, 1707391925)`)
	if err != nil {
		fmt.Println("errror While inserting User'Seed Data 1' into db !! Due to : ", err)
		return
	}
	stmt.Exec()

	//Seed Data for Branch
	stmt, err = database.Prepare(`INSERT INTO branch VALUES
		(1001, "Central Bank Of India", "Shivajinagar Pune", 1707465912,1707465912)`)
	if err != nil {
		fmt.Println("errror While inserting Branch'Seed Data 1' into db !! Due to : ", err)
		return
	}
	stmt.Exec()

	stmt, err = database.Prepare(`INSERT INTO branch VALUES
		(1002, "Central Bank Of India", "Balewadi,Pune", 1707465912,1707465912)`)
	if err != nil {
		fmt.Println("errror While inserting Branch'Seed Data 2' into db !! Due to : ", err)
		return
	}
	stmt.Exec()

	stmt.Close()

	fmt.Println("Seed Data Inserted !!")
}
