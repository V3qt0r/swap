package api

import (
	"strings"
	"swap/models"

	validation "github.com/go-ozzo/ozzo-validation"
)


type RegisterItemPayload struct {
	Name		string	`json:"name"`
	Category	string	`json:"category"`
	Description	string	`json:"description"`
	Prize		float64	`json:"prize"`
	OwnerId     uint     `json:"ownerId"`
}


func (r RegisterItemPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Description, validation.Length(10, 300)),
		validation.Field(&r.Category, validation.Required),
		validation.Field(&r.Prize, validation.Required, validation.Min(0.00)),
	)
}

func (r RegisterItemPayload) Sanitize() {
	r.Category = strings.TrimSpace(r.Category)
	r.Category = strings.ToUpper(r.Category)
}


type ItemSearchRequest struct {
	SearchTerm string `json:"term"`
}


type ItemSearchResponse struct {
	Name		string `json:"name"`
	Category	string `json:"category"`
	Description	string `json:"description"`
	Prize		float64`json:"prize"`
	UUID        string `json:"uuid"`
	ID          int    `json:"id"`
}


func (r ItemSearchRequest) Sanitize() {
	r.SearchTerm = strings.TrimSpace(r.SearchTerm)
}


func (r ItemSearchRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SearchTerm, validation.Required),
	)
}


type ItemUpdatePayload struct {
	Name		string 	`json:"name"`
	Category	string 	`json:"category"`
	Description	string 	`json:"description"`
	Prize		float64 `json:"prize"`
}


func (r ItemUpdatePayload) Validate() error {
	if r.Category != "" {
		return validation.Validate(r.Category, validation.Required)
	}
	return nil
}


func (r ItemUpdatePayload) Sanitize() {
	r.Category = strings.TrimSpace(r.Category)
	r.Category = strings.ToUpper(r.Category)
}


func (r ItemUpdatePayload) ToEntity() models.Item {
	var item models.Item

	if r.Name != "" {
		item.Name = r.Name
	}
	if r.Category != "" {
		item.Category = r.Category
	}
	if r.Description != "" {
		item.Description = r.Description
	}
	if r.Prize >= 0.00 {
		item.Prize = r.Prize
	}
	return item
}


type BuyItemPayload struct {
	ID 		int			`json:"id"`
	Amount  float64		`json:"amount"`
}


func (r BuyItemPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(r.ID, validation.Required),
		validation.Field(r.Amount, validation.Required, validation.Min(0.00)),
	)
}


type SwapItemPayload struct {
	Item1Id		int		`json:"item1Id"`
	Item2Id		int		`json:"item2Id"`
	Amount		float64	`json:"amount"`
}


func (r SwapItemPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(r.Item1Id, validation.Required),
		validation.Field(r.Item2Id, validation.Required),
		validation.Field(r.Amount, validation.Min(0.00)),
	)
}
