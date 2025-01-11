package services

import (
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


// func (s *itemService) SwapItem(item1Id, item2Id int, amount float64) (string, error){
// 	return s.ItemRepository.SwapItem(item1Id, item2Id, amount)
// }


func (s *itemService) GetItemsByOwnerId(ownerId uint, limit, page int) ([]models.Item, error) {
	return s.ItemRepository.GetItemsByOwnerId(ownerId, limit, page)
}


func (s *itemService) GetItemByUUID(uuid string) (*models.Item, error) {
	return s.ItemRepository.GetByUUID(uuid)
}


func (s *itemService) UpdateCategory(itemId int, categoryName string) error {
	return s.ItemRepository.UpdateCategory(itemId, categoryName)
}


// func (s *itemService) ReadImageByPath(path string) ([]byte, error) {
// 	return utils.ReadImageByPath(path)
// }	
