package repository

import (
	"swap/models"
	"swap/apperrors"
	"log"

	"gorm.io/gorm"
)


type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) models.ICategoryRepository {
	return &categoryRepository{
		DB: db,
	}
}


func (r *categoryRepository) CreateCategory(name string)(*models.Category, error){
	category := &models.Category{}
	
	result := r.DB.Where("name = ?", name).Find(&category)
	if result.Error != nil {
		log.Print("Category already exists")
		return nil, apperrors.NewBadRequest("Category already exists")
	}

	category.Name = name
	if err := r.DB.Create(&category).Error; err != nil {
		return nil, apperrors.NewBadRequest("Unable to create category")
	}
	return category, nil
}


func (r *categoryRepository) DeleteCategory(category string) error {
	result := r.DB.Model(&models.Item{}).
			  Where("category_name = ?", category).
			  Updates(map[string]interface{}{
				"category_name": nil,
				"category_id": nil,
			  })
	
	if result.Error != nil {
		return apperrors.NewBadRequest("Unable to delete category")
	}

	if err := r.DB.Where("name = ?", category).Delete(&models.Category{}).Error; err != nil {
		return apperrors.NewBadRequest("Unable to delete category")
	}
	return nil
}


func (r *categoryRepository) GetAllValidCategories() ([]models.Category, error) {
	var categories []models.Category
	
	if err := r.DB.Select("name", "ban", "id").Where("ban = ?", false).Find(&categories).Error; err != nil {
		return categories, apperrors.NewInternal()
	}
	return categories, nil
}


func (r *categoryRepository) BanCategory(name string) error {
	category := models.Category{}

	if err := r.DB.Where("name = ?", name).Find(&category).Error; err != nil {
		return apperrors.NewBadRequest("Invalid category")
	}

	category.Ban = true
	err := r.DB.Save(&category)
	if err.Error != nil {
		return apperrors.NewInternal()
	}
	return nil
}


func (r *categoryRepository) UnBanCategory(name string) error {
	category := models.Category{}

	if err := r.DB.Where("name = ?", name).Find(&category).Error; err != nil {
		return apperrors.NewBadRequest("Invalid category")
	}

	category.Ban = false
	err := r.DB.Save(&category)
	if err.Error != nil {
		return apperrors.NewInternal()
	}
	return nil
}


func (r *categoryRepository) CheckStatus(name string) (bool, error) {
	category := &models.Category{}

	if err := r.DB.Where("name = ?", name).First(category).Error; err != nil {
		if err == gorm.ErrRecordNotFound{
			return false, nil
		}
		return false, apperrors.NewBadRequest("Invalid category")
	}

	return category.Ban, nil
}


func (r *categoryRepository) GetAllItemsInCategory(id, limit, page int) ([]models.Item, error) {
	var items []models.Item
	category := &models.Category{}

	if err := r.DB.Where("id = ?", id).Find(&category).Error; err != nil {
		log.Print("Could not find category")
		return items, apperrors.NewBadRequest("Could not find category")
	}

	if err := r.DB.Where("category_id = ? AND sold = ?", category.ID, false).Find(&items).Error; err != nil {
		log.Print("Could not find items")
		return items, apperrors.NewInternal()
	}

	return items, nil
}