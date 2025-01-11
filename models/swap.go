package models

import "time"


type SwapRequest struct {
	Base
	Item1Id  			uint 			`json:"item1Id" gorm:"not null"`
	Item2Id				uint 			`json:"item2Id" gorm:"not null"`
	OwnerId 			uint 			`json:"ownerId" gorm:"not null"` //Owner of target item item2
	InitiatorId         uint 			`json:"initiatorId" gorm:not null"` //Initiator of swap request and owner of item1
	Status				string 			`json:"status" gorm:"default: PENDING"` //PENDING, APPROVED, REJECTED
	CompletionStatus    string			`json:"completionStatus"`
}


type EnrichedSwapRequest struct {
	ID 					uint 			`json:"id"`
	Item1Id  			uint 			`json:"item1Id" gorm:"not null"`
	Item2Id				uint 			`json:"item2Id" gorm:"not null"`
	Item1Details        ItemDetails		`json:"item1Details"`
	Item2Details 		ItemDetails		`json:"item2Details"`
	InitiatorId         uint 			`json:"initiatorId" gorm:not null"` //Owner of item1
	InitiatorDetails    UserDetails 	`json:"initiatorDetails"`
	Status				string 			`json:"status" gorm:"default: PENDING"` //PENDING, APPROVED, REJECTED
	CreatedAt           time.Time 		`json:"createdAt"`
}


type IncompleteSwaps struct {
	ID 					uint 			`json:"id"`
	BalanceOwed 		float64			`json:"balanceOwed"`
	ItemDetails        ItemDetails		`json:"itemDetails"`
}


type ItemDetails struct {
	Name				string		`json:"name"`
	Description			string 		`json: "description"`
	Category            string   	`json:"category"`
	Prize				float64 	`json: "prize"`
}


type UserDetails struct {
	Name 				  string  	`json:"name"`
	UserName			  string 	`json:"userName" gorm:"unique"`
	PhoneNumber			  string 	`json:"phoneNumber" gorm:"unique"`
	Email				  string 	`json:"email" gorm:"unique"`
	Gender                string 	`json: "gender"`
	Location              string 	`json:"location"`
	ProfileUrl            string    `json:"profileUrl"`
	ProfileIcon           string    `json:"profileIcon"`
}



type ISwapRepository interface {
	InitiateSwapRequest(item1Id, item2Id, initiatorId uint) (*SwapRequest, error)
	GetPendingSwapRequests(ownerId int, limit, page int) ([]EnrichedSwapRequest, error)
	RejectSwapRequest(ownerId, swapId int) error
	AcceptSwapRequest(ownerId, swapId int) (string, error)
	CompleteSwapRequest(ownerId uint, amount float64, swapId uint) (string, error)
	GetIncompleteSwapByInitiatorId(initiatorId, itemId int) (IncompleteSwaps, error)
	GetAllIncompleteSwapByOwnerId(ownerId, limit, page int) ([]IncompleteSwaps, error)
}


type ISwapService interface {
	InitiateSwapRequest(item1Id, item2Id, initiatorId uint) (*SwapRequest, error)
	GetPendingSwapRequests(ownerId int, limit, page int) ([]EnrichedSwapRequest, error)
	RejectSwapRequest(ownerId, swapId int) error
	AcceptSwapRequest(ownerId, swapId int) (string, error)
	CompleteSwapRequest(ownerId uint, amount float64, swapId uint) (string, error)
	GetIncompleteSwapByInitiatorId(initiatorId, itemId int) (IncompleteSwaps, error)
	GetAllIncompleteSwapByOwnerId(ownerId, limit, page int) ([]IncompleteSwaps, error)
}