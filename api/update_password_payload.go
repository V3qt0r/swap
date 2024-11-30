package api 

import "github.com/go-ozzo/ozzo-validation"

type UpdatePasswordPayload struct {
	Password		string	`json:"password"`
	ConfirmPassword	string	`json:"confirmPassword"`
}


func (r UpdatePasswordPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.ConfirmPassword, validation.Required),	
	)
}


type ConfirmPasswordPayload struct {
	Password	string	`json:"password"`
}


func (c ConfirmPasswordPayload) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Password, validation.Required),
	)
}