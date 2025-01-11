package models

import (
	"time"
)

type Item struct {
	Base
	Name				string		`json:"name"`
	Description			string 		`json: "description"`
	CategoryName 		string 		`json: "catgeoryName"`
	CategoryId          *uint       `json: "-"`
	Prize				float64 	`json: "prize" gorm:"type:numeric(19,2);default:0"`
	Sold				bool		`json: "sold" gorm:"type:boolean; default:false"`
	User 				User        `gorm:"foreignKey:OwnerId" json:"-"`
	OwnerId    			uint		`json: "-"`
	SoldAt 				time.Time   `json:"soldAt"`
	Images  			[]Image		`json:"images" gorm:"costraint: OnDelete: CASCADE"`
}


type Category struct {
	Base
	Name 	string `json:"name" gorm:"unique"`
	Ban     bool   `json:"ban" gorm:"type:boolean; defaut:false"`
}
type ICategoryRepository interface {
	CreateCategory(name string) (*Category, error)
	DeleteCategory(name string) error
	GetAllValidCategories() ([]Category, error)
	GetAllItemsInCategory(id, limit, page int) ([]Item, error)
	BanCategory(name string) error
	CheckStatus(name string) (bool, error)
	UnBanCategory(name string) error
}
type ICategoryService interface {
	CreateCategory(name string) (*Category, error)
	DeleteCategory(name string) error
	GetAllValidCategories() ([]Category, error)
	GetAllItemsInCategory(id, limit, page int) ([]Item, error)
	BanCategory(name string) error
	CheckStatus(name string) (bool, error)
	UnBanCategory(name string) error
}


type IItemRepository interface {
	GetItemById(id int) (*Item, error)
	GetByUUID(uuid string) (*Item, error)
	GetItemsByCategory(category string, limit, page int) ([]Item, error)
	RegisterItem(item *Item) (*Item, error)
	GetUnsoldItemsByCategory(category string, limit, page int) ([]Item, error)
	UpdateItem(item Item) error
	DeleteItem(itemId int) error
	GetItemsByOwnerId(ownerId uint, limit, page int) ([]Item, error)
	BuyItem(userId, itemId int, amount float64) (string, error)
	UpdateCategory(itemId int, categoryName string) error
}


type IItemService interface {
	RegisterItem(item *Item) (*Item, error)
	GetItemById(id int) (*Item, error)
	GetItemByUUID(uuid string) (*Item, error)
	GetItemsByCategory(category string, limit, page int) ([]Item, error)
	GetUnsoldItemsByCategory(category string, limit, page int) ([]Item, error)
	UpdateItem(item Item) error
	DeleteItem(itemId int) error
	BuyItem(userId, itemId int, amount float64) (string, error)
	GetItemsByOwnerId(ownerId uint, limit, page int) ([]Item, error)
	UpdateCategory(itemId int, categoryName string) error
}
