package specs

import (
	"fmt"
	"unicode"
)

type UpdateUserInfo struct {
	User_id  int    `json:"user_id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
	Role     string `json:"role"`
}

func (req *UpdateUserInfo) ValidateUpdate() error {
	if len(req.Name) == 0 {
		return fmt.Errorf("name field cannot be empty")
	}
	if len(req.Name) > 50 {
		return fmt.Errorf("name field length cannot exceed 50 characters")
	}
	if len(req.Address) == 0 {
		return fmt.Errorf("address field cannot be empty")
	}
	if len(req.Address) > 100 {
		return fmt.Errorf("address field length cannot exceed 100 characters")
	}
	if len(req.Password) == 0 {
		return fmt.Errorf("password field cannot be empty")
	}
	if len(req.Password) < 3 || len(req.Password) > 16 {
		return fmt.Errorf("password field length must be between 3 and 16 characters")
	}
	if len(req.Mobile) == 0 {
		return fmt.Errorf("mobile field cannot be empty")
	}
	if len(req.Mobile) != 10 {
		return fmt.Errorf("mobile field must be 10 digits long")
	}
	for _, char := range req.Mobile {
		if !unicode.IsDigit(char) {
			return fmt.Errorf("mobile field must contain only digits")
		}
	}
	if len(req.Role) == 0 {
		return fmt.Errorf("role field cannot be empty")
	}
	if req.Role != "Admin" && req.Role != "Customer" {
		return fmt.Errorf("role field must be either 'Admin' or 'Customer'")
	}
	return nil
}
