package entity

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type User struct {
	ID           string `json:"id,omitempty" db:"id"`
	Email        string `json:"email,omitempty" db:"email"`
	Password     string `json:"password,omitempty" db:"password"`
	PasswordHash string `json:"password_hash,omitempty"`
}

func (e *User) ValidateEmail() error {
	if e.Email == "" {
		return fmt.Errorf("email is required")
	}
	// Регулярное выражение для проверки формата email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(e.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (e *User) ValidateLogin() error {
	err := e.ValidateEmail()
	if err != nil {
		return err
	}
	if e.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (e *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.PasswordHash = string(bytes)
	return nil
}

func (e *User) CheckPassword() bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.PasswordHash), []byte(e.Password))
	if err != nil {
		return false
	}
	return true
}

type UserRegister struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Code     string `json:"code,omitempty"`
}

func (e *UserRegister) Validate() error {
	if e.Email == "" {
		return fmt.Errorf("email is required")
	}
	// Регулярное выражение для проверки формата email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(e.Email) {
		return fmt.Errorf("invalid email format")
	}
	if e.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(e.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if e.Code == "" {
		return fmt.Errorf("code is required")
	}
	return nil
}
