package repository

import (
	"swap/models"
	"errors"
	"strconv"
	"fmt"
	"log"
	"strings"
	"swap/apperrors"

	"gorm.io/gorm"
)

type itemRepository struct {
	DB *gorm.DB
}

func NewItemRepository(db *gorm.DB) models.IItemRepository {
	return &itemRepository{
		DB: db,
	}
}


func (r *itemRepository) GetItemById(id int) (*models.Item, error) {
	item := &models.Item{}
	itemId := strconv.Itoa(id)

	if err := r.DB.Where("id = ?", itemId).Find(&item).Error; err != nil {
		log.Printf("Error getting item with ID %d\n", itemId)

		if errors.Is(err, gorm.ErrRecordNotFound){
			return item, apperrors.NewNotFound("ID", itemId)
		}
		return item, apperrors.NewInternal()
	}
	return item, nil
}


func (r * itemRepository) GetByUUID(uuid string) (*models.Item, error) {
	item := &models.Item{}

	if err := r.DB.Where("uuid = ?", uuid).Find(&item).Error; err != nil {
		log.Printf("Error getting item with UUID %s\n", uuid)
		
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, apperrors.NewNotFound("UUID", uuid)
		}
		return item, apperrors.NewInternal()
	}
	return item, nil
}


func (r *itemRepository) GetItemsByCategory(category string, limit, page int) ([]models.Item, error) {
	var items []models.Item

	err := r.DB.Select("name", "description", "category", "prize", "ID", "UUID","owner_id").Where("category = ?", category).Find(&items) 
	if err.Error != nil {
		return items, apperrors.NewInternal()
	}
	return items, nil
}


func (r *itemRepository) RegisterItem(item *models.Item) (*models.Item, error) {
	item.Category = strings.ToUpper(item.Category)
	if result := r.DB.Create(&item); result.Error != nil {
			log.Printf("Could not register item %v. Reason: %v\n", item.Name, result.Error)
			return nil, apperrors.NewInternal()
		}
	
	return item, nil
}


func (r *itemRepository) GetUnsoldItemsByCategory(category string, limit, page int) ([]models.Item, error) {
	var items []models.Item

	err := r.DB.Select("name", "description", "category", "prize", "ID", "UUID", "owner_id").Where("category = ? AND sold = false", category).Find(&items)
	if err.Error != nil {
		return items, apperrors.NewInternal()
	}
	return items, nil
}


func (r *itemRepository) UpdateItem(item models.Item) error {
	itemId := int(item.ID)
	foundItem, err := r.GetItemById(itemId)

	if err != nil {
		return apperrors.NewInternal()
	}

	if foundItem == nil {
		return apperrors.NewNotFound("item ID", strconv.Itoa(itemId))
	}

	updatedDetails := map[string] interface {}{}
		if item.Name != "" {
			updatedDetails["Name"] = item.Name
		}
		if  item.Description != "" {
			updatedDetails["Description"] = item.Description
		}
		if item.Category != "" {
			updatedDetails["Category"] = strings.ToUpper(item.Category)
		}
		if item.Prize >= 0.00 {
			updatedDetails["Prize"] = item.Prize
		}
		if  !item.Sold {
			updatedDetails["Sold"] = item.Sold
		}

	if err := r.DB.Model(&foundItem).Updates(updatedDetails).Error; err != nil {
		return apperrors.NewInternal()
	}
	return nil
}


func (r *itemRepository) DeleteItem(id int) error {
	if err := r.DB.Where("id = ?", id).Delete(&models.Item{}).Error; err != nil {
		return err
	}
	return nil
}


func (r *itemRepository) GetItemsByOwnerId(ownerId uint, limit, page int) ([]models.Item, error) {
	var items []models.Item
	
	err := r.DB.Select("name", "description", "category", "prize", "sold","ID").Where("owner_id = ?", ownerId).Find(&items)
	if err.Error != nil {
		return items, apperrors.NewInternal()
	}
	return items, nil
}


func (r *itemRepository) BuyItem(id int, amount float64) (string, error){
	item := &models.Item{}
	itemId := strconv.Itoa(id)

	if err := r.DB.Where("id = ?", itemId).Find(&item).Error; err != nil {
		return "", apperrors.NewBadRequest("ID not found")
	}


	if item.Sold == true {
		log.Print("Item have already been sold\n")
		return "", apperrors.NewBadRequest("Item have aleady been sold")
	}

	balance := amount - item.Prize
	if balance < 0 {
		log.Printf("Insufficient amount: %.2f required\n", item.Prize)
		return "", apperrors.NewBadRequest("Insufficient amount: " + strconv.Itoa(int(item.Prize)) + " required")
	}

	if err := r.DB.Model(&item).Updates(models.Item{Sold: true}).Error; err != nil {
		return "", apperrors.NewInternal()
	}
	

	result := fmt.Sprintf("Item ID: %v\nItem Name: %s\nSold: %v\nPrize: $%.2f\nAmount Paid: $%.2f\nBalance To Retreive: $%.2f\n",
	item.ID, item.Name, item.Sold, item.Prize, amount, balance)


	return result ,nil
}


func (r *itemRepository) SwapItem(item1Id, item2Id int) (string, error) {
	item1, err := r.GetItemById(item1Id)

	if err != nil {
		return "", apperrors.NewBadRequest("Item with ID "+ strconv.Itoa(item1Id) + " not found")
	}

	item2, err := r.GetItemById(item2Id)

	if err != nil {
		return "", apperrors.NewBadRequest("Item with ID "+ strconv.Itoa(item2Id) + " not found")
	}

	if err := r.DB.Model(&item1).Updates(models.Item{Sold: true}).Error; err != nil {
		return "", apperrors.NewInternal()
	}

	if err := r.DB.Model(&item2).Updates(models.Item{Sold: true}).Error; err != nil {
		return "", apperrors.NewInternal()
	}

	return "", nil
}