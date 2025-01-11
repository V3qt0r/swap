package api

import (
	validation "github.com/go-ozzo/ozzo-validation"
)


type InitiateSwapRequestPayload struct {
	Item1Id     uint 		`json:"item1Id"`
	Item2Id     uint 		`json:"item2Id"`
}


func (r InitiateSwapRequestPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Item1Id, validation.Required),
		validation.Field(&r.Item2Id, validation.Required),
	)
}
