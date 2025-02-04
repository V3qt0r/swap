package repository

import (
	"swap/models"
	"swap/apperrors"
	
	"errors"
	"strconv"
	"fmt"
	"log"
	"time"

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


func (r *itemRepository) UpdateCategory(itemId int, categoryName string) error {
	category := &models.Category{}
	item,err := r.GetItemById(itemId)

	if err != nil {
		return apperrors.NewNotFound("ID", strconv.Itoa(itemId))
	}

	if err := r.DB.Where("name = ?", categoryName).First(&category).Error; err != nil {
		return apperrors.NewBadRequest("Category does not exist")
	}

	item.CategoryId = &category.ID
	item.CategoryName = category.Name
	if err := r.DB.Save(&item).Error; err != nil {
		return apperrors.NewBadRequest("Failed to update category")
	}
	return nil
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

	err := r.DB.Joins("JOIN categories ON categories.id = items.category_id").
				Where("categories.name = ?", category).
				Find(&items)

	if err.Error != nil {
		return items, apperrors.NewInternal()
	}
	return items, nil
}


func (r *itemRepository) RegisterItem(item *models.Item) (*models.Item, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var category models.Category

		if err := tx.Where("name = ? AND ban = ?", item.CategoryName, false).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound){
				log.Print("Invalid category")
				return fmt.Errorf("Invalid category: %v", item.CategoryName)
			}
			return err
		}

		item.CategoryId = &category.ID
		item.CategoryName = category.Name
		if err := tx.Create(&item).Error; err != nil {
			log.Print("Error")
			return err
		}
		return nil
	})

	if err != nil {
		log.Print("Error")
		return nil, apperrors.NewInternal()
	}

	return item, nil
}


func (r *itemRepository) GetUnsoldItemsByCategory(category string, limit, page int) ([]models.Item, error) {
	var items []models.Item

	err := r.DB.Joins("JOIN categories ON categories.id = items.category_id").
				Where("categories.name = ? AND categories.ban = ? AND items.sold = ?", category, false, false).
				Find(&items)

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
	
	err := r.DB.Select("name", "description", "category_name", "prize", "sold","ID").Where("owner_id = ?", ownerId).Find(&items)
	if err.Error != nil {
		return items, apperrors.NewInternal()
	}
	return items, nil
}


func (r *itemRepository) BuyItem(userID, id int, amount float64) (string, error){
	item := &models.Item{}
	itemId := strconv.Itoa(id)
	userId := strconv.Itoa(userID)
	user := &models.User{}
	owner := &models.User{}

	if err := r.DB.Where("id = ?", itemId).Find(&item).Error; err != nil {
		return "", apperrors.NewBadRequest("ID not found")
	}

	if err := r.DB.Where("id = ?", item.OwnerId).First(&owner).Error; err != nil {
		return "", apperrors.NewBadRequest("Unable to find owner")
	}

	if err := r.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return "", apperrors.NewBadRequest("Unable to find buyer")
	}

	if item.OwnerId == uint(userID) {
		log.Print("Cannot purchase own item!")
		return "", apperrors.NewBadRequest("Cannot purchase own item!")
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

	transactions := models.Transactions{
		Name:   owner.Name,
		Email: owner.Email,
		PhoneNumber: owner.PhoneNumber,
		OwnerId: user.ID,
		ItemId: item.ID,
		ItemName: item.Name,
		Bought:  true,
		Swapped: false,
		AmountPaid: amount,
		BalanceAvailabe: balance,
		BalanceOwed: 0.00,
	}

	if err := r.DB.Create(&transactions).Error; err != nil {
		return "", apperrors.NewBadRequest("Unable to create transaction")
	}

	if err := r.DB.Model(&item).Updates(models.Item{Sold: true, SoldAt : time.Now().Truncate(time.Second)}).Error; err != nil {
		return "", apperrors.NewInternal()
	}
	

	result := fmt.Sprintf("Item ID: %v\nItem Name: %s\nSwapped: %v\nPrize: $%.2f\nAmount Paid: $%.2f\nBalance To Retreive: $%.2f\n",
	item.ID, item.Name, true, item.Prize, amount, balance)

	return result, nil
}

