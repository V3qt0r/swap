package services

import (
	"log"
	"strconv"
	"fmt"

	"swap/apperrors"
	"swap/models"
)


type itemService struct {
	ItemRepository models.IItemRepository
}

func NewItemService(ItemRepository models.IItemRepository) models.IItemService {
	return &itemService{
		ItemRepository: 	ItemRepository,
	}
}


func (s *itemService) RegisterItem(item *models.Item) (*models.Item, error) {
	return s.ItemRepository.RegisterItem(item)
}


func (s *itemService) GetItemById(id int) (*models.Item, error) {
	return s.ItemRepository.GetItemById(id)
}


func (s *itemService) GetItemsByCategory(category string, limit, page int) ([]models.Item, error){
	return s.ItemRepository.GetItemsByCategory(category, limit, page)
}


func (s *itemService) GetUnsoldItemsByCategory(category string, limit, page int) ([]models.Item, error){
	return s.ItemRepository.GetUnsoldItemsByCategory(category, limit, page)
}


func (s *itemService) UpdateItem(item models.Item) error {
	return s.ItemRepository.UpdateItem(item)
}


func (s *itemService) DeleteItem(itemId int) error {
	return s.ItemRepository.DeleteItem(itemId)
}


func (s *itemService) BuyItem(userId, itemId int, amount float64) (string, error) {
	return s.ItemRepository.BuyItem(userId, itemId, amount)
}


func (s *itemService) SwapItem(item1Id, item2Id int, amount float64) (string, error) {
	item1, err := s.ItemRepository.GetItemById(item1Id)
	if err != nil {
		log.Printf("Item with ID: %v does not exist\n", item1Id)
		return "", apperrors.NewBadRequest("Item with ID: " + strconv.Itoa(item1Id) + " does not exist\n")
	}

	item2, err := s.ItemRepository.GetItemById(item2Id)
	if err != nil {
		log.Printf("Item with ID: %v does not exist\n", item2Id)
		return "", apperrors.NewBadRequest("Item with ID: " + strconv.Itoa(item2Id) + " does not exist\n")
	}

	if item1.Sold == true || item2.Sold == true {
		if item1.Sold == true && item2.Sold == true {
			log.Printf("Both item with ID %v and ID %v have already been sold\n", item1Id, item2Id)
			return "", apperrors.NewBadRequest("Both items have already been sold")
		}

		if item1.Sold == true {
			log.Printf("Item with ID %v already been sold\n", item1Id)
			return "", apperrors.NewBadRequest("Item with ID " + strconv.Itoa(int(item1Id)) + " have already been sold")
		}

		log.Printf("Item with ID %v already been sold\n", item2Id)
		return "", apperrors.NewBadRequest("Item with ID " + strconv.Itoa(int(item2Id)) + " have already been sold")
	}

	if item1.CategoryName != item2.CategoryName {
		log.Print("Cannot swap with item outside of category")
		return "", apperrors.NewBadRequest("Cannot swap with item outside of category")
	}

	prize := 0.00

	if item1.Prize > item2.Prize {
		prize = item1.Prize - item2.Prize
	} else if item2.Prize > item1.Prize {
		prize = item2.Prize - item1.Prize
	} else {
		prize = item1.Prize - item2.Prize
	}

	balance := amount - prize

	if prize > 0 && balance < 0 {
		log.Printf("Insufficient amount: %.2f required\n", prize)
		return "", apperrors.NewBadRequest("Insufficient amount: " + strconv.Itoa(int(prize)) + " required")
	} 

	result, err := s.ItemRepository.SwapItem(item1Id, item2Id)
	if err != nil {
		return "", apperrors.NewInternal()
	}

	result = fmt.Sprintf("Item ID: %v\nItem Name: %s\nSold: true\nPrize: $%.2f\n\nSuccessfully swapped with\n\nItem ID: %v\nItem Name: %s\nSold: true\nPrize: $%.2f\nAmount paid: $%.2f\nBalance To Retreive: $%v\n",
	item1Id, item1.Name, item1.Prize, item2Id, item2.Name, item2.Prize, amount,balance)
	
	return result, nil
}


func (s *itemService) GetItemsByOwnerId(ownerId uint, limit, page int) ([]models.Item, error) {
	return s.ItemRepository.GetItemsByOwnerId(ownerId, limit, page)
}


func (s *itemService) GetItemByUUID(uuid string) (*models.Item, error) {
	return s.ItemRepository.GetByUUID(uuid)
}


func (s *itemService) UpdateCategory(itemId int, categoryName string) error {
	return s.ItemRepository.UpdateCategory(itemId, categoryName)
}