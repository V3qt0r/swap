package repository

import (
	"swap/models"
	"swap/apperrors"

	"strconv"
	"time"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)


type swapRepository struct {
	DB *gorm.DB
}


func NewSwapRepository(db *gorm.DB) models.ISwapRepository{
	return &swapRepository{
		DB: db,
	}
}


func (r *swapRepository) InitiateSwapRequest(id1, id2, initiatorId uint) (*models.SwapRequest, error) {
	item1 := &models.Item{}
	item1Id := strconv.Itoa(int(id1))
	initiator := &models.User{}


	if err := r.DB.Where("id = ?", initiatorId).First(&initiator).Error; err != nil {
		return nil, apperrors.NewBadRequest("Error retrieving initiator's details")
	}

	if err := r.DB.Where("id = ?", item1Id).First(&item1).Error; err != nil {
		return nil, apperrors.NewNotFound("Item", item1Id)
	}

	if item1.Sold == true {
		fmt.Errorf("Your item has already been sold")
		return nil, apperrors.NewBadRequest("Your item has already been sold")
	}

	item2 := &models.Item{}
	item2Id := strconv.Itoa(int(id2))

	if err := r.DB.Where("id = ?", item2Id).First(&item2).Error; err != nil {
		return nil, apperrors.NewNotFound("Item", item2Id)
	}

	if item2.Sold == true {
		fmt.Errorf("Item has already been sold")
		return nil, apperrors.NewBadRequest("Item has already been sold")
	} 

	if item1.OwnerId != initiatorId {
		return nil, apperrors.NewBadRequest("You must be the owner of first item")
	}

	if item2.OwnerId == initiatorId {
		return nil, apperrors.NewBadRequest("You cannot swap with your own item")
	}

	swapRequest := &models.SwapRequest{}

	if err := r.DB.Where("initiator_id = ? AND item2_id = ?", initiatorId, item2.ID).First(&swapRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			swapRequest.Item1Id = item1.ID
			swapRequest.Item2Id = item2.ID
			swapRequest.OwnerId = item2.OwnerId
			swapRequest.InitiatorId = initiatorId
			swapRequest.Status = "PENDING"


			if err := r.DB.Create(&swapRequest).Error; err != nil{
				fmt.Errorf("Failed initialize swap process")
				return nil, apperrors.NewBadRequest("Failed to initialize swap process")
			}
			return swapRequest, nil
		} 
			fmt.Errorf("Failed to create swap process")
			return swapRequest, apperrors.NewBadRequest("Failed to create swap process")
		}

		if swapRequest.Item1Id == item1.ID{
			log.Print("You have already initiated a swap request with this item")
			return swapRequest, apperrors.NewBadRequest("You have already initiated a swap request with this item")
		}

		if err := r.DB.Model(&swapRequest).Updates(models.SwapRequest{Item1Id: item1.ID, Status: "PENDING"}).Error; err != nil {
			return nil, apperrors.NewBadRequest("Failed to update swap")
		}
		return swapRequest, nil
	
}



func (r *swapRepository) GetPendingSwapRequests(ownerId int, limit, page int) ([]models.EnrichedSwapRequest, error) {
	var requests []models.SwapRequest
	var enrichedRequests []models.EnrichedSwapRequest

	if err := r.DB.Where("owner_id = ? AND status = ?", ownerId, "PENDING").Find(&requests).Error; err != nil {
		log.Print("Unable to retrieve pending swap requests")
		return enrichedRequests, apperrors.NewBadRequest("Unable to retrieve pending swap requests")
	}

	for _, swap := range requests {
		item1 := &models.Item{}
		item2 := &models.Item{}
		initiator := &models.User{}

		if err := r.DB.Where("id = ?", swap.Item1Id).First(&item1).Error; err != nil {
			log.Print("Unable to get destination item")
			return nil, apperrors.NewBadRequest("Unable to get destination item")
		}

		if err := r.DB.Where("id = ?", swap.Item2Id).First(&item2).Error; err != nil {
			log.Print("Unable to get source item")
			return nil, apperrors.NewBadRequest("Unable to get source item")
		}

		if err := r.DB.Where("id = ?", swap.InitiatorId).First(&initiator).Error; err != nil {
			log.Print("Unable to get swap initiator")
			return nil, apperrors.NewBadRequest("Unable to get swap initiator")
		}

		enrichedRequests = append(enrichedRequests, models.EnrichedSwapRequest{
			ID				: 	swap.ID,
			Item1Id 		:	item1.ID,
			Item2Id			:	item2.ID,
			Item1Details:   models.ItemDetails{
				Name 			:	item1.Name,
				Description 	:	item1.Description,
				Category		:   item1.CategoryName,
				Prize 			:	item1.Prize,
			},
			Item2Details:   models.ItemDetails{
				Name 			:	item2.Name,
				Description 	:	item2.Description,
				Category		:   item2.CategoryName,
				Prize 			:	item2.Prize,
			},
			InitiatorId 	:	initiator.ID,
			InitiatorDetails:   models.UserDetails{
				Name 			:	initiator.Name,
				UserName 		:	initiator.UserName,
				PhoneNumber 	:	initiator.PhoneNumber,
				Email 			:	initiator.Email,
				Gender 			:	initiator.Gender,
				Location 		:	initiator.Location,
				ProfileUrl 		:	initiator.ProfileUrl,
				ProfileIcon 	:	initiator.ProfileIcon,
			},
			Status 			:	swap.Status,
			CreatedAt 		:	swap.CreatedAt,
		})
	}
	return enrichedRequests, nil
}



func (r *swapRepository) RejectSwapRequest(ownerId, swapId int) error {
	request := &models.SwapRequest{}

	if err := r.DB.Where("id = ? AND owner_id = ?", swapId, ownerId).First(&request).Error; err != nil {
		log.Print("Swap request not found")
		return apperrors.NewNotFound("user", strconv.Itoa(ownerId))
	}
	if request.Status == "REJECTED" || request.CompletionStatus == "INCOMPLETE" || request.CompletionStatus == "COMPLETED"{
		log.Print("You can reject only pending and accepted swap requests.")
		return apperrors.NewBadRequest("You can reject only pending and accepted swap requests.")
	}

	if err := r.DB.Model(&request).Updates(models.SwapRequest{Status: "REJECTED"}).Error; err != nil {
		log.Print("Could not reject swap request. Please try again.")
		return apperrors.NewBadRequest("Could not reject swap request. Please try again.")
	}
	
	if err := r.DB.Where("id = ?", request.ID).Delete(&request).Error; err != nil {
		log.Print("Could not delete swap request. Please try again")
		return apperrors.NewBadRequest("Could not delete swap request. Please try again")
	}
	
	return nil
}



func (r *swapRepository) AcceptSwapRequest(ownerId, swapId int) (string, error){
	request := &models.SwapRequest{}

	if err := r.DB.Where("id = ? AND owner_id = ?", swapId, ownerId).First(&request).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not find swap request. Please try again")
	}

	if request.Status != "PENDING" {
		return "", apperrors.NewBadRequest("You can only accept pending requests")
	}

	item1 := &models.Item{}
	if err := r.DB.Where("id = ?", request.Item1Id).First(&item1).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not find target item")
	}

	item2 := &models.Item{}
	if err := r.DB.Where("id = ?", request.Item2Id).First(&item2).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not find source item")
	}

	if item1.Prize != item2.Prize {
		if err := r.DB.Model(&request).Updates(models.SwapRequest{Status: "ACCEPTED", CompletionStatus: "INCOMPLETE"}).Error; err != nil {
			return "Swap request accepted waiting for balance payment to complete", nil
		}
		return "Swap request accepted. Incomplete till payment of balance is confirmed", nil
	}

	if err := r.DB.Model(&item1).Updates(models.Item{Sold: true, SoldAt : time.Now().Truncate(time.Second)}).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not update source item sold status")
	}

	if err := r.DB.Model(&item2).Updates(models.Item{Sold: true, SoldAt : time.Now().Truncate(time.Second)}).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not update target item sold status")
	}

	transaction1 := &models.Transactions{}
	transaction2 := &models.Transactions{}
	owner1 := &models.User{}
	owner2 := &models.User{}

	if err := r.DB.Where("id = ?", request.InitiatorId).First(&owner1).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not find initiator details")
	}

	if err := r.DB.Where("id = ?", request.OwnerId).First(&owner2).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not find owner details")
	}

	transaction1.Name = owner1.Name
	transaction1.Email = owner1.Email
	transaction1.PhoneNumber = owner1.PhoneNumber
	transaction1.OwnerId = owner1.ID
	transaction1.ItemId = item2.ID
	transaction1.ItemName = item2.Name
	transaction1.Bought = false
	transaction1.Swapped = true
	transaction1.AmountPaid = 0.00
	transaction1.BalanceAvailabe = 0.00
	transaction1.BalanceOwed = 0.00

	transaction2.Name = owner2.Name
	transaction2.Email = owner2.Email
	transaction2.PhoneNumber = owner2.PhoneNumber
	transaction2.OwnerId = owner2.ID
	transaction2.ItemId = item1.ID
	transaction2.ItemName = item1.Name
	transaction2.Bought = false
	transaction2.Swapped = true
	transaction2.AmountPaid = 0.00
	transaction2.BalanceAvailabe = 0.00
	transaction2.BalanceOwed = 0.00

	if err := r.DB.Create(&transaction1).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not create a transaction receipt for user")
	}

	if err := r.DB.Create(&transaction2).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not create a transaction receipt for user")
	}

	if err := r.DB.Model(&request).Updates(models.SwapRequest{Status: "ACCEPTED", CompletionStatus: "COMPLETED"}).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not update swap status")
	}

	completedRequest := &models.SwapRequest{}

	if err := r.DB.Where("item2_id", request.Item2Id).Delete(&completedRequest).Error; err != nil {
		return "", apperrors.NewBadRequest("Could not delete swap requests")
	}

	result := fmt.Sprintf("Item ID: %v\nItem Name: %s\nSwapped: %v\nPrize: $%.2f\nAmount Paid: $%.2f\nBalance To Retreive: $%.2f\n",
	item1.ID, item1.Name, true, item1.Prize, 0.00, 0.00)

	return result, nil
}



func (r *swapRepository) CompleteSwapRequest(ownerId uint, amount float64, swapId uint) (string, error) {
	request := &models.SwapRequest{}

	if err := r.DB.Where("id = ?",swapId).First(&request).Error; err != nil {
		log.Print("Could not find swap request")
		return "", apperrors.NewBadRequest("Could not find swap request")
	}

	item1 := &models.Item{}
	item2 := &models.Item{}

	if err := r.DB.Where("id = ?", request.Item1Id).First(&item1).Error; err != nil {
		log.Print("Could not find item")
		return "", apperrors.NewBadRequest("Could not find item")
	}

	if err := r.DB.Where("id = ?", request.Item2Id).First(&item2).Error; err != nil {
		log.Print("Could not find item")
		return "", apperrors.NewBadRequest("Could not find item")
	}

	transaction1 := &models.Transactions{}
	transaction2 := &models.Transactions{}
	owner1 := &models.User{}
	owner2 := &models.User{}

	if err := r.DB.Where("id = ?", item1.OwnerId).First(&owner1).Error; err != nil {
		log.Print("Could not find owner")
		return "", apperrors.NewBadRequest("Could not find owner")
	}

	if err := r.DB.Where("id = ?", item2.OwnerId).First(&owner2).Error; err != nil {
		log.Print("Could not find owner")
		return "", apperrors.NewBadRequest("Could not find owner")
	}

	result := ""

	if item1.Prize > item2.Prize && item2.OwnerId == ownerId{
		balance, err := r.PayOff(item1.Prize, item2.Prize, amount)

		if err != nil{
			log.Print("Could not get balance")
			return "", apperrors.NewBadRequest("Could not get balance")
		}

		err = r.AssignTransaction(transaction1, transaction2, owner2, owner1, request, amount, balance)
		if err != nil {
			log.Print("Could not assign transactions history 1")
			return "", apperrors.NewBadRequest("Could not assign transactions history 1")
		}

		result = fmt.Sprintf("Item ID: %v\nItem Name: %s\nSwapped: %v\nPrize: $%.2f\nAmount Paid: $%.2f\nBalance To Retreive: $%.2f\n",
		item1.ID, item1.Name, true, item1.Prize, amount, balance)

	} else if item2.Prize > item1.Prize && item1.OwnerId == ownerId{
		balance, err := r.PayOff(item2.Prize, item1.Prize, amount)

		if err != nil {
			log.Print("Could not not get balance")
			return "", apperrors.NewBadRequest("Could not not get balance")
		}

		err = r.AssignTransaction(transaction1, transaction2, owner1, owner2, request, amount, balance)
		if err != nil {
			log.Print("Could not assign transactions history 2")
			return "", apperrors.NewBadRequest("Could not assign transactions history 2")
		}

		result = fmt.Sprintf("Item ID: %v\nItem Name: %s\nSwapped: %v\nPrize: $%.2f\nAmount Paid: $%.2f\nBalance To Retreive: $%.2f\n",
	item2.ID, item2.Name, true, item2.Prize, amount, balance)

	} else {
		log.Print("You do not have any balance owed")
		return "", apperrors.NewBadRequest("You do not have any balance owed")
	}

	if err := r.DB.Model(&item1).Updates(models.Item{Sold: true, SoldAt : time.Now().Truncate(time.Second)}).Error; err != nil {
		log.Print("Unable to update item status")
		return "", apperrors.NewBadRequest("Unable to update item status")
	}

	if err := r.DB.Model(&item2).Updates(models.Item{Sold: true, SoldAt : time.Now().Truncate(time.Second)}).Error; err != nil {
		log.Print("Unable to update item status")
		return "", apperrors.NewBadRequest("Unable to update item status")
	}

	if err := r.DB.Model(&request).Updates(models.SwapRequest{CompletionStatus: "COMPLETED"}).Error; err != nil {
		log.Print("Could not update completion status")
		return "", apperrors.NewBadRequest("Could not update completion status")
	}

	if err := r.DB.Where("item2_id = ?", request.Item2Id).Delete(&request).Error; err != nil {
		log.Print("Unable to delete swap request")
		return "", apperrors.NewBadRequest("Unable to delete swap request")
	}
	return result, nil
}


func (r *swapRepository) PayOff(prize1, prize2, amount float64) (float64, error){
	prize := prize1 - prize2
	balance := amount - prize

	if prize > 0 && balance < 0 {
		log.Print("Amount too low for prize")
		return 0.00, apperrors.NewBadRequest("Amount too low for prize")
	}
	return balance, nil
}


func (r *swapRepository) AssignTransaction(transaction1, transaction2 *models.Transactions, owner1, owner2 *models.User, request *models.SwapRequest, amount, balance float64) error {
	item1 := &models.Item{}
	item2 := &models.Item{}

	log.Print(request.Item1Id)
	if err := r.DB.Where("id = ?", request.Item1Id).First(&item1).Error; err != nil {
		log.Print("Item not found 1")
		return apperrors.NewBadRequest("Item not found")
	}

	if err := r.DB.Where("id = ?", request.Item2Id).First(&item2).Error; err != nil {
		log.Print("Item not found 2")
		return apperrors.NewBadRequest("Item not found")
	}

	transaction1.Name = owner1.Name
	transaction1.Email = owner1.Email
	transaction1.PhoneNumber = owner1.PhoneNumber
	transaction1.OwnerId = owner1.ID
	transaction1.ItemId = request.Item2Id
	transaction1.ItemName = item2.Name
	transaction1.Bought = false
	transaction1.Swapped = true
	transaction1.AmountPaid = amount
	transaction1.BalanceAvailabe = balance
	transaction1.BalanceOwed = 0.00

	transaction2.Name = owner2.Name
	transaction2.Email = owner2.Email
	transaction2.PhoneNumber = owner2.PhoneNumber
	transaction2.OwnerId = owner2.ID
	transaction2.ItemId = request.Item1Id
	transaction2.ItemName = item1.Name
	transaction2.Bought = false
	transaction2.Swapped = true
	transaction2.AmountPaid = 0.00
	transaction2.BalanceAvailabe = 0.00
	transaction2.BalanceOwed = 0.00

	tx := r.DB.Begin()

	if err := tx.Create(transaction1).Error; err != nil {
		tx.Rollback()
		log.Print("Unable to assign transaction history 1")
		return apperrors.NewBadRequest("Unable to assign transaction history 1")
	}
	
	if err := tx.Create(transaction2).Error; err != nil {
		tx.Rollback()
		log.Print("Unable to assign transaction history 2")
		return apperrors.NewBadRequest("Unable to assign transaction history 2")
	}

	if err := tx.Commit().Error; err != nil {
		log.Print("Unable to commit transactions")
		return apperrors.NewBadRequest("Unable to commit transactions")
	}
	return nil
}




func (r *swapRepository) GetAllIncompleteSwapByOwnerId(ownerId int, limit, page int) ([]models.IncompleteSwaps, error) {
	var incompleteSwaps []models.IncompleteSwaps
	var requests []models.SwapRequest

	if err := r.DB.Where("owner_id = ?", ownerId).Find(&requests).Error; err != nil {
		return incompleteSwaps, apperrors.NewBadRequest("You have no incomplete swaps")
	}

	for _, swap := range requests {
		item1 := &models.Item{}

		if err := r.DB.Where("id = ?", swap.Item1Id).First(&item1).Error; err != nil {
			return nil, apperrors.NewBadRequest("Could not retrieve item")
		}

		item2 := &models.Item{}

		if err := r.DB.Where("id = ?", swap.Item2Id).First(&item2).Error; err != nil {
			return nil, apperrors.NewBadRequest("Could not retrieve item")
		}

		balanceOwed := item1.Prize - item2.Prize

		incompleteSwaps = append(incompleteSwaps, models.IncompleteSwaps{
			ID : swap.ID,
			BalanceOwed: balanceOwed,
			ItemDetails: models.ItemDetails{
				Name:  item1.Name,
				Description: item1.Description,
				Category: item1.CategoryName,
				Prize: item1.Prize,
			},
		})
	}

	return incompleteSwaps, nil
}



func (r *swapRepository) GetIncompleteSwapByInitiatorId(initiatorId, itemId int) (models.IncompleteSwaps, error) {
	request := &models.SwapRequest{}
	item1 := &models.Item{}
	item2 := &models.Item{}
	incompleteSwap := models.IncompleteSwaps{}

	if err := r.DB.Where("initiator_id = ? AND item2_id = ? AND completion_status = ?", initiatorId, itemId, "INCOMPLETE").First(&request).Error; err != nil {
		log.Print("Unable to find incomplete swap")
		return incompleteSwap, apperrors.NewBadRequest("Unable to find incomplete swap")
	}

	if err := r.DB.Where("id = ?", request.Item1Id).First(&item1).Error; err != nil{
		log.Print("Unable to find item 1")
		return incompleteSwap, apperrors.NewBadRequest("Unable to find item")
	}

	if err := r.DB.Where("id = ?", request.Item2Id).First(&item2).Error; err != nil{
		log.Print("Unable to find item 2")
		return incompleteSwap, apperrors.NewBadRequest("Unable to find item")
	}

	incompleteSwap.ID = request.ID
	incompleteSwap.BalanceOwed = item1.Prize - item2.Prize
	incompleteSwap.ItemDetails = models.ItemDetails{
		Name: item2.Name,
		Description: item2.Description,
		Category: item2.CategoryName,
		Prize: item2.Prize,
	}

	return incompleteSwap, nil
}
