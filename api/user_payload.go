package api

import (
	"strings"
	"swap/models"
	"swap/apperrors"
	"time"

	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)


type RegisterPayload struct {
	Name		string  	`json:"name"`
	UserName	string		`json:"userName"`
	PhoneNumber string		`json:"phoneNumber"`
	Email 		string		`json:"email"`
	DOB 	    string		`json:"dob"`
	Gender 		string      `json:"gender"`
	Password 	string		`json:"password"`
	Location 	string 		`json:"location"`
	IsAbove18 	bool		`json:"isAbove18"`
}

func (r RegisterPayload) Sanitize() {
	r.Name			=   strings.TrimSpace(r.Name)
	r.UserName 		= 	strings.TrimSpace(r.UserName)
	r.Email			= 	strings.TrimSpace(r.Email)
	r.Email			=	strings.ToLower(r.Email)
	r.PhoneNumber	=	strings.TrimSpace(r.PhoneNumber)
	r.Password		=	strings.TrimSpace(r.Password)
	r.Gender		=	strings.TrimSpace(r.Gender)
	r.Gender		=   strings.ToLower(r.Gender)
}

func (r RegisterPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.PhoneNumber /*validation.By(validatePhoneNumber)*/),
		validation.Field(&r.UserName, validation.Length(3, 30)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
		validation.Field(&r.IsAbove18, validation.Required, validation.In(true).Error("You must be above 18 to register")),
		validation.Field(&r.Name,  validation.Length(3, 30)),
		validation.Field(&r.Gender, validation.In("male", "female").Error("You must either be male or female")),
		validation.Field(&r.Location, ),
	)
}


func isDigits(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return apperrors.NewBadRequest("BVN must be a valid string")
	}
	match, _ := regexp.MatchString("^[0-9]+$", s)
	if !match {
		return apperrors.NewBadRequest("BVN must contain only numeric characters")
	}
	return nil
}


type UserSearchResponse struct {
	UserName    string  `json:"userName"`
	Email		string	`json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	UUID        string 	`json:"password"`
}


type UserSearchRequest struct {
	SearchTerm	string  `json:"searchTerm"`
}


func (r UserSearchRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SearchTerm, validation.Required),
	)
}


type UserUpdatePayload struct {
	Name		string	   	 `json:"name"`
	UserName	string		 `json:"userName"`
	PhoneNumber	string		 `json:"phoneNumber"`
	Email 		string		 `json:"email"`
	DOB 		string	 `json:"dob"`
	Gender		string		 `json:"gender"`
	Location    string		 `json:"location"`
	ProfileUrl 	string		 `json:"profileUrl"`
	ProfileIcon	string		 `json:"profileIcon"`		
}


func (r UserUpdatePayload) Validate() error {
	if r.Email != "" {
		return validation.Validate(r.Email, validation.Required, is.EmailFormat)
	}

	if r.UserName != "" {
		return validation.Validate(r.UserName, validation.Required, validation.Length(3, 30))
	}

	if r.PhoneNumber != "" {
		return validation.Validate(r.PhoneNumber, validation.Required /*validation.By(validatePhoneNumber)*/)
	}

	return nil
}


func (r UserUpdatePayload) Sanitize() {
	r.Name			=   strings.TrimSpace(r.Name)
	r.UserName 		= 	strings.TrimSpace(r.UserName)
	r.Email			= 	strings.TrimSpace(r.Email)
	r.Email			=	strings.ToLower(r.Email)
	r.PhoneNumber	=	strings.TrimSpace(r.PhoneNumber)
	r.Gender		=	strings.TrimSpace(r.Gender)
	r.Gender		=   strings.ToLower(r.Gender)
}



func (r UserUpdatePayload) ToEntity() models.User{
	var user models.User

	if r.Name != "" {
		user.Name = r.Name
	}
	if r.UserName != "" {
		user.UserName = r.UserName
	}
	if r.PhoneNumber != "" {
		user.PhoneNumber = r.PhoneNumber
	}
	if r.Email != "" {
		user.Email = r.Email
	}
	
	dob, _ := time.Parse("2006-01-02", r.DOB)
	if !dob.IsZero() {
		user.DOB = dob
	}
	if r.Gender != "" {
		user.Gender = r.Gender
	}
	if r.Location != "" {
		user.Location = r.Location
	}
	if r.ProfileUrl != "" {
		user.ProfileUrl = r.ProfileUrl
	}
	if r.ProfileIcon != "" {
		user.ProfileIcon = r.ProfileIcon
	}
	return user
}

type FindUserByEmailOrUsernamePayload struct {
	EmailOrUsername		string		`json:"email`
}

func (r FindUserByEmailOrUsernamePayload) Sanitize() {
	r.EmailOrUsername = strings.TrimSpace(r.EmailOrUsername)
}

func (r FindUserByEmailOrUsernamePayload) Validate() error {	
	return validation.ValidateStruct(&r,
		validation.Field(r.EmailOrUsername, validation.Required, validation.Length(3, 30)),
	)
}