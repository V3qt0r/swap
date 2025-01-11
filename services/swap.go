package services

import (
	"swap/models"
)

type swapService struct {
	SwapRepository models.ISwapRepository
}


func NewSwapService(SwapRepository models.ISwapRepository) models.ISwapService {
	return &swapService{
		SwapRepository: 	SwapRepository,
	}
}


func (s *swapService) InitiateSwapRequest(item1Id, item2Id, initiatorId uint) (*models.SwapRequest, error) {
	return s.SwapRepository.InitiateSwapRequest(item1Id, item2Id, initiatorId)
}


func (s *swapService) GetPendingSwapRequests(ownerId int, limit, page int) ([]models.EnrichedSwapRequest, error) {
	return s.SwapRepository.GetPendingSwapRequests(ownerId, limit, page)
}


func (s *swapService) RejectSwapRequest(ownerId, swapId int) error {
	return s.SwapRepository.RejectSwapRequest(ownerId, swapId)
}


func (s *swapService) AcceptSwapRequest(ownerId, swapId int) (string, error) {
	return s.SwapRepository.AcceptSwapRequest(ownerId, swapId)
}


func (s *swapService) CompleteSwapRequest(ownerId uint, amount float64, swapId uint) (string, error) {
	return s.SwapRepository.CompleteSwapRequest(ownerId, amount, swapId)
}


func (s *swapService) GetIncompleteSwapByInitiatorId(initiatorId, itemId int) (models.IncompleteSwaps, error) {
	return s.SwapRepository.GetIncompleteSwapByInitiatorId(initiatorId, itemId)
}


func (s *swapService) GetAllIncompleteSwapByOwnerId(ownerId int, limit, page int) ([]models.IncompleteSwaps, error) {
	return s.SwapRepository.GetAllIncompleteSwapByOwnerId(ownerId, limit, page)
}