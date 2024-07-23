package specs

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

type CreateAccountReq struct {
	Account_no   int     `json:"acc_no,omitempty"`
	User_id      int     `json:"user_id"`
	Branch_id    int     `json:"branch_id"`
	Account_type string  `json:"acc_type"`
	Balance      float64 `json:"balance"`
}

type DeleteAccountReq struct {
	Account_no int `json:"acc_no"`
	User_id    int `json:"user_id"`
}

type Transaction struct {
	Account_no int     `json:"acc_no"`
	Amount     float64 `json:"amount"`
}

type TransactionResponse struct {
	Account_no int     `json:"acc_no"`
	Balance    float64 `json:"balance,omitempty"`
}

type DeleteAccount struct {
	Msg string `json:"msg"`
}

type Acc_Type string

const (
	Savings Acc_Type = "savings"
	Current Acc_Type = "current"
	Loan    Acc_Type = "loan"
	Salary  Acc_Type = "salary"
)

func isValidateAcc_Type(t Acc_Type) bool {
	switch t {
	case Savings, Current, Loan, Salary:
		return true
	default:
		return false
	}
}

func (req *CreateAccountReq) Validate() error {
	if req.Branch_id <= 0 {
		return errors.New("branch_id must be greater than 0")
	}
	if req.Account_no < 0 {
		return errors.New("account_no cannot be negative")
	}
	if len(req.Account_type) <= 0 {
		return fmt.Errorf("please provide Valid Account type")
	}
	if !isValidateAcc_Type(Acc_Type(strings.ToLower(req.Account_type))) {
		return fmt.Errorf(" invalid role, accepted account types are - savings, current,loan,salary")
	}
	if req.Balance < 0 {
		return errors.New("balance cannot be negative")
	}
	return nil
}

func (req *DeleteAccountReq) ValidateDeleteReq() error {
	if req.Account_no <= 0 {
		return fmt.Errorf("please provide Valid Account No")
	}
	if req.User_id < 0 {
		return fmt.Errorf("please provide Valid User_ID")
	}
	return nil
}

func (req *Transaction) ValidateTransaction() error {

	if req.Account_no < 0 {
		return fmt.Errorf("account number never be negative")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount never be negative or zero")
	}
	return nil
}

func SendMail(data CreateAccountReq) {
	from := "saurabhskale242001@gmail.com"
	password := "xnpwowohwjqwbvwn"
	to := "saurabh.kale@joshsoftware.com"

	subject := "Welcome to SBM Bank | test@mail"
	body := "Dear Customer,\n\n"
	body += "Thank you for choosing SBM Bank. We are thrilled to serve you the best banking facility!\n\n"
	body += "At SBM Bank, we are dedicated to providing you with the best banking experience possible.\n Please find below the details of your new account:\n\n"
	body += "Account Number: " + strconv.Itoa(data.Account_no) + "\n"
	body += "Account Type: " + data.Account_type + "\n"
	body += "Opening Balance: " + strconv.FormatFloat(data.Balance, 'f', 2, 64) + "\n"
	body += "Branch ID: " + strconv.Itoa(data.Branch_id) + "\n\n"
	body += "If you have any questions or need assistance, feel free to reach out to us.\n\n"
	body += "Welcome once again to SBM Bank!\n\n"
	body += "Best regards,\n"
	body += "The SBM Bank Team"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: " + subject + "\r\n" +
		"\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Fatal("Error sending mail:", err)
	}
	log.Println("Mail sent successfully!")
}
