package services

import (
	"swap/models"
	"swap/utils"
)


type imageService struct {
	ImageRepository models.IImageRepository
}


func NewImageService(imageRepository models.IImageRepository) *imageService {
	return &imageService{
		ImageRepository: imageRepository,
	}
}


func (s *imageService) ReadFirstImageById(id int) ([]byte, error) {
	return s.ImageRepository.ReadFirstImageById(id)
}


func (s *imageService) ReadAllImagesByItemId(id, limit, page int) ([]models.Image, error) {
	return s.ImageRepository.ReadAllImagesByItemId(id, limit, page)
}


// func (s *imageService) UploadImage(itemId int, folderName, fileName string) error {
// 	return s.ImageRepository.UploadImage(itemId, folderName, fileName)
// }

func (s *imageService) ReadImage(folderName, fileName string) ([]byte, error) {
	return utils.ReadImage(folderName, fileName)
}	


func (s *imageService) ReadImageByPath(path string) ([]byte, error) {
	return utils.ReadImageByPath(path)
}	
