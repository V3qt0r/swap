package repository

import (
	"swap/models"
	"swap/apperrors"
	"swap/utils"

	"log"
	"errors"
	"strconv"

	"gorm.io/gorm"
)


type imageRepository struct {
	DB *gorm.DB
}


func NewImageRepository(db *gorm.DB) models.IImageRepository {
	return &imageRepository{
		DB: db,
	}
}

func (r *imageRepository) UploadImage(itemId int, folderName, fileName string) error{
	item := &models.Item{}
	image := &models.Image{}

	if err := r.DB.Where("id = ?", itemId).First(&item).Error; err != nil {
		log.Printf("Could not find item with ID: %d\n", itemId)
		return apperrors.NewBadRequest("Could not find item with provided ID")
	}

	


	image.FilePath = folderName
	image.FileName = fileName
	image.ItemId = uint(itemId)
	image.OwnerId = item.OwnerId

	if err := r.DB.Create(&image).Error; err != nil {
		log.Print("Could not save image to database")
		return apperrors.NewBadRequest("Could not save image to database")
	}
	return nil
}


func (r *imageRepository) ReadFirstImageById(id int) ([]byte, error) {
	item := &models.Item{}
	image := &models.Image{}

	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		log.Print("Item with ID %v does not exist\n", id)
		return nil, apperrors.NewBadRequest("Item with provided ID does not exist")
	}

	if err := r.DB.Where("item_id = ?", id).First(&image).Error; err != nil {
		log.Print("Could not find image with item ID %v\n", id)
		return nil, apperrors.NewBadRequest("Could not find image with provided item ID")
	}

	folderName := "item_" + strconv.Itoa(id)
	fileName :=  image.FileName

	return utils.ReadImage(folderName, fileName)
}


func (r *imageRepository) ReadAllImagesByItemId(id, limit, page int) ([]models.Image, error) { 
	item := &models.Item{}
	var images []models.Image


	if err := r.DB.Where("id = ?", id).Find(&item).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound){
			log.Print("Error getting item with ID %d\n", id)
			return images, apperrors.NewBadRequest("Could not find Item with provided ID")
		}
		log.Printf("Error getting item")
		return images, apperrors.NewInternal()
	}

	if err := r.DB.Where("item_id = ?", id).Find(&images).Error; err != nil {
		log.Print("Could not find images for this item")
		return images, apperrors.NewBadRequest("Could not find images for this item")
	}

	return images, nil
}
