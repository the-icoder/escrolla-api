package models

import (
	"errors"
	goval "github.com/go-passwd/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Username       string `json:"username,omitempty" gorm:"unique;default:null"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	Email          string `json:"email,omitempty" gorm:"unique"`
	Password       string `json:"password,omitempty" gorm:"-" binding:"required,min=8,max=15"`
	PasswordHash   string `json:"-" gorm:"password_hash"`
	Phone          string `json:"phone,omitempty"`
	Status         string `json:"status,omitempty"`
	EmailConfirmed bool   `json:"-" gorm:"default:false"`
	BusinessName   string `json:"business_name"`
	CustomerID     int    `json:"-"`
	ReferralCode   string `json:"referral_code"`
}

// VerifyPassword verifies the collected password with the user's hashed password
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

type SignupRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required,e164"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=15"`
	//ReferralCode    string `json:"referral_code"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserResponse
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type ForgotPassword struct {
	Email string `json:"email" binding:"required,email"`
}
type ResetPassword struct {
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func (u *User) LoginUserToDto(token string) *LoginResponse {
	return &LoginResponse{
		UserResponse: UserResponse{
			ID:          u.ID,
			Name:        u.FirstName + " " + u.LastName,
			PhoneNumber: u.Phone,
			Email:       u.Email,
		},
		AccessToken: token,
	}
}

func ValidatePassword(password string) error {
	passwordValidator := goval.New(goval.MinLength(6, errors.New("password cant be less than 6 characters")),
		goval.MaxLength(15, errors.New("password cant be more than 15 characters")))
	err := passwordValidator.Validate(password)
	return err
}
