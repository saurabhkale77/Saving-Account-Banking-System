package specs

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/dgrijalva/jwt-go"
)

type CreateLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	User_id  int
	Role     string
	jwt.StandardClaims
}

type CreateUser struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Mobile   string `json:"mobile"`
	Role     string `json:"role"`
}

type UpdateUser struct {
	User_id  int    `json:"user_id,omitempty"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Password string `json:"password,omitempty"`
	Mobile   string `json:"mobile"`
}

type Response struct {
	User_id  int             `json:"user_id"`
	Name     string          `json:"name"`
	Address  string          `json:"address"`
	Email    string          `json:"email"`
	Password string          `json:"password,omitempty"`
	Mobile   string          `json:"mobile"`
	Role     string          `json:"role"`
	Accounts []GetMyAccounts `json:"accounts"`
}

type BranchDetails struct {
	Id       int    `json:"branch_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type GetMyAccounts struct {
	Acc_no          int     `json:"acc_no"`
	Branch_id       int     `json:"branch_id"`
	Branch_name     string  `json:"branch_name"`
	Branch_location string  `json:"branch_location"`
	Acc_Type        string  `json:"acc_type"`
	Balance         float64 `json:"balance"`
}

type LoginToken struct {
	Issuespecsken string `json:"token"`
}

type Role string

const (
	Customer Role = "customer"
	Admin    Role = "admin"
)

func (req *CreateLoginRequest) Validate() error {
	if len(req.Username) <= 0 || len(req.Password) <= 0 {
		return fmt.Errorf("please provide anything as input")
	}
	if !isValidEmail(req.Username) {
		return fmt.Errorf("please provide a proper Email credentials")
	}
	if !isValidPassword(req.Password) {
		return fmt.Errorf("password should contains atleast one UpperCase,Lowercase letter, one special char and one digit")
	}
	return nil
}

func (req *CreateUser) ValidateUser() error {
	if len(req.Name) <= 0 {
		return fmt.Errorf("name field cannot be empty")
	}
	if len(req.Address) <= 0 {
		return fmt.Errorf("address field cannot be empty")
	}
	if len(req.Email) <= 0 {
		return fmt.Errorf("email field cannot be empty")
	}
	if !isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if len(req.Password) <= 0 {
		return fmt.Errorf("password field cannot be empty")
	}
	if len(req.Password) < 3 || len(req.Password) > 16 {
		return fmt.Errorf("length of the password field must be between 3 and 16 characters")
	}
	if strings.Contains(req.Password, " ") {
		return fmt.Errorf("password cannot contain spaces")
	}
	if len(req.Mobile) <= 0 {
		return fmt.Errorf("mobile field cannot be empty")
	}
	if len(req.Mobile) != 10 {
		return fmt.Errorf("mobile field must be of 10 digits only")
	}
	// Check each character is a digit
	for _, char := range req.Mobile {
		if !unicode.IsDigit(char) {
			return fmt.Errorf("mobile field must be contains digits only")
		}
	}
	if len(req.Role) <= 0 || (strings.EqualFold(req.Role, "Customer") && strings.EqualFold(req.Role, "Admin")) {
		return fmt.Errorf("role field can't be empty")
	}
	if !isValidateRole(Role(strings.ToLower(req.Role))) {
		return fmt.Errorf(" invalid role, accepted roles are customer and admin only")
	}
	if !isValidPassword(req.Password) {
		return fmt.Errorf("password should contains atleast one UpperCase,Lowercase letter, one special char and one digit")
	}
	return nil
}

func isValidateRole(role Role) bool {
	switch role {
	case Customer, Admin:
		return true
	default:
		return false
	}
}

func isValidEmail(email string) bool {
	// Basic Checks
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(regex).MatchString(email) {
		return false
	}
	if strings.Contains(email, " ") {
		return false
	}
	if strings.Count(email, "@") != 1 {
		return false
	}
	return true
}

func (req *UpdateUser) ValidateUpdate() error {
	if len(req.Name) <= 0 {
		return fmt.Errorf("name field cannot be empty")
	}
	if len(req.Address) <= 0 {
		return fmt.Errorf("address field cannot be empty")
	}
	if len(req.Password) <= 0 {
		return fmt.Errorf("password field cannot be empty")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("length of the password field must be greater than 6 chars")
	}
	if !isValidPassword(req.Password) {
		return fmt.Errorf("password should contains atleast one UpperCase,Lowercase letter, one special char and one digit")
	}
	if len(req.Mobile) <= 0 {
		return fmt.Errorf("mobile field cannot be empty")
	}
	for _, char := range req.Mobile {
		if !unicode.IsDigit(char) {
			return fmt.Errorf("mobile field must be contains digits only")
		}
	}
	return nil
}

func isValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	containsUpperCase := false
	containsLowerCase := false
	containsDigit := false
	containsSpecialChar := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			containsUpperCase = true
		} else if unicode.IsLower(char) {
			containsLowerCase = true
		} else if unicode.IsDigit(char) {
			containsDigit = true
		} else if !unicode.IsSpace(char) {
			containsSpecialChar = true
		}
	}
	return containsUpperCase && containsLowerCase && containsDigit && containsSpecialChar
}
