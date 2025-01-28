package models

import (

	"time"
	"golang.org/x/crypto/bcrypt"
)

const (
	OneTimeLoginOTPType     OTPType = "OneTimeLoginOTPType"
)

type OTPType string


type Transactions struct {
	Base
	Name 			string 		`json:"name"`
	Email 			string 		`json:"email"`
	PhoneNumber 	string 		`json:"phoneNumber"`
	OwnerId         uint        `gorm:"not null;index" json: "-"`
	ItemId 			uint 		`json: "_"`
	ItemName        string      `json:"itemName"`
	Bought          bool        `json:"bought"`
	Swapped         bool        `json:"swapped"`
	AmountPaid      float64     `json: "amount:"amountPaid"`
	BalanceAvailabe float64     `json:"balanceAvailable"`
	BalanceOwed     float64		`json:"balanceOwed"`
}


type User struct {
	Base
	Name 				  string  	`json:"name"`
	UserName			  string 	`json:"userName" gorm:"unique"`
	PhoneNumber			  string 	`json:"phoneNumber" gorm:"unique"`
	Email				  string 	`json:"email" gorm:"unique"`
	DOB 				  time.Time `json: "dob"`
	ProfileUrl            string    `json:"profileUrl"`
	ProfileIcon           string    `json:"profileIcon"`
	Gender                string 	`json: "gender"`
	Password			  string 	`json: "-"`
	Location              string 	`json:"location"`
	TotpEnabled           bool      `json:"totpEnabled" gorm:"type:bool;default:false"`
	TotpSecret            string    `json:"_"`
	OneTimePassword       string    `json:"-"`
	OneTimePasswordExpiry time.Time `json:"oneTimePasswordExpiry"`
	OneTimePasswordValid  bool      `json:"oneTimePasswordValid" gorm:"type:bool;default:false"`
	// TransactionsId        uint      `json:"-"`
	// Transactions		  Transactions `gorm:"foreignKey:TransactionsId" json:"-"`
}



type IUserRepository interface {
	GetUserById(id int) (*User, error)
	GetUserByUUID(uuid string) (*User, error)
	CreateUser(user *User) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindUserByPhoneNumber(phoneNumber string) (*User, error)
	FindUserByEmailOrUsername(email string) (*User, error)
	InvalidateOneTimePassword(user *User) error
	CreateOneTimePassword(user *User, password string, expiry time.Time) error
	Search(term string, limit, page int) ([]User, error)
	UpdateUser(user User) error
	UpdateUserTOTP(user User, totpSecret string, totpEnabled bool) error
	UpdatePassword(userId uint, password string) error
	GetUserTransactions(userId, limit, page int) ([]Transactions, error)
	GetUserByItemId(id int) (*User, error)
}

type IUserService interface {
	Register(user *User) (*User, error)
	Login(email, passsword string) (*User, error)
	GetUserById(id int) (*User, error)
	GetUserByUUID(uuid string) (*User, error)
	InitLoginWithOneTimePassword(email string) (error)
	LoginWithOneTimePassword(email, code string) (*User, error)
	UpdatePassword(userId uint, password string) error
	ConfirmPassword(userId uint, password string) (*User, error)
	Search(term string, limit, page int) ([]User, error)
	FindUserByEmailOrUsername(email string) (*User, error)
	FindUserByPhoneNumber(phoneNumber string) (*User, error)
	UpdateUser(user User) error
	GenerateOneTimePasswordForUser(user *User, otpType OTPType, duration time.Duration) (string, error)
	InvalidateOneTimePassword(user *User) error
	EnrollTOTP(userId int) ([]byte, error)
	VerifyTOTP(userId int, verifyTOTP VerifyTOTPRequest) error
	DisableTOTP(userId int) error
	EnableTOTP(userId int) error
	GetUserTransactions(userId, limit, page int) ([]Transactions, error)
	GetUserByItemId(id int) (*User, error)
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

type VerifyTOTPRequest struct {
	Totp string `json: "totp"`
}

func (r VerifyTOTPRequest) Validate() error {
	return nil
}

// func (r VerifyTOTPRequest) Sanitize() {

// }