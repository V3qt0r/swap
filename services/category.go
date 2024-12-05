package services

import (
	"swap/models"
)


type categoryService struct {
	CategoryRepository models.ICategoryRepository
}


func NewCategoryService(CategoryRepository models.ICategoryRepository) models.ICategoryService {
	return &categoryService{
		CategoryRepository: CategoryRepository,
	}
}



func (s *categoryService) CreateCategory(name string) (*models.Category, error) {
	return s.CategoryRepository.CreateCategory(name)
}


func (s *categoryService) DeleteCategory(name string) error {
	return s.CategoryRepository.DeleteCategory(name)
}


func (s *categoryService) GetAllValidCategories()([]models.Category, error) {
	return s.CategoryRepository.GetAllValidCategories()
}


func (s *categoryService) BanCategory(name string) error {
	return s.CategoryRepository.BanCategory(name)
}


func (s *categoryService) UnBanCategory(name string) error {
	return s.CategoryRepository.UnBanCategory(name)
}


func (s *categoryService) CheckStatus(name string) (bool, error) {
	return s.CategoryRepository.CheckStatus(name)
}