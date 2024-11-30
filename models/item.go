package models

type Item struct {
	Base
	Name				string		`json:"name"`
	Description			string 		`json: "description"`
	Category			string 		`json: "category" gorm:"default:'GENERAL'"` // MOBILE, AUTO, ELECTRONICS, GENERAL
	Prize				float64 	`json: "prize" gorm:"type:numeric(19,2);default:0"`
	Sold				bool		`json: "sold" gorm:"type:boolean; default:false"`
	User 				User        `gorm:"foreignKey:OwnerId" json:"-"`
	OwnerId    			uint			`json: "-"`
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
	BuyItem(itemId int, amount float64) (string, error)
	SwapItem(item1Id, item2Id int) (string, error)
}


type IItemService interface {
	RegisterItem(item *Item) (*Item, error)
	GetItemById(id int) (*Item, error)
	GetItemByUUID(uuid string) (*Item, error)
	GetItemsByCategory(category string, limit, page int) ([]Item, error)
	GetUnsoldItemsByCategory(category string, limit, page int) ([]Item, error)
	UpdateItem(item Item) error
	DeleteItem(itemId int) error
	BuyItem(itemId int, amount float64) (string, error)
	SwapItem(item1Id, item2Id int, amount float64) (string, error)
	GetItemsByOwnerId(ownerId uint, limit, page int) ([]Item, error)
	GetCategories() []string
	AddCategory(category string) error
	RemoveCategory(category string) error
}

