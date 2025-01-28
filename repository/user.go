package repository

import (
	"swap/models"
	"errors"
	"log"
	"strconv"
	"time"
	"swap/apperrors"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) models.IUserRepository {
	return &userRepository{
		DB: db,
	}
}

var userId uint

func (r *userRepository) GetUserById(id int) (*models.User, error) {
	user := &models.User{}
	userId := strconv.Itoa(id)
	
	if err := r.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		log.Printf("Could not find user with id: %d\n", id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("user ID" ,userId)
		}
		return user, apperrors.NewInternal()
	}
	return user, nil 
}


func (r *userRepository) GetUserByUUID(id string) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("uuid = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("id", id)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}


func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {

		log.Printf("Could not create user with email %v. Reason: %v\n", user.Email, result.Error)
		return nil, apperrors.NewInternal()
	}
	
	return user, nil
}


func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("Could not find user with email %s\n", email)
		if errors.Is(err, gorm.ErrRecordNotFound){
			return user, apperrors.NewNotFound("email", email)
		}
		return user, apperrors.NewInternal()
	}
	return user, nil
}


func (r *userRepository) FindUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		log.Printf("Could not find user with phone number %s\n", phoneNumber)
		if errors.Is(err, gorm.ErrRecordNotFound){
			return user, apperrors.NewNotFound("phonenumber", phoneNumber)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}


func (r *userRepository) FindUserByEmailOrUsername(email string) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("email = ?", email).Or("user_name = ?", email).Find(&user).Error; err != nil {
		log.Printf("Could not find user with provided email or username \n")
		if errors.Is(err, gorm.ErrRecordNotFound){
			value := email
			return user, apperrors.NewNotFound("email or username", value)
		}
		return user, apperrors.NewInternal()
	}
	return user, nil
}


func (r *userRepository) InvalidateOneTimePassword(user *models.User) error {

	result := r.DB.Model(&user).Updates(models.User{OneTimePasswordValid: false})
	if result.Error != nil {
		return result.Error
	}
	return nil
}


func (r *userRepository) CreateOneTimePassword(user *models.User, password string, expiryTime time.Time) error {

	result := r.DB.Model(&user).Updates(models.User{OneTimePassword: password, OneTimePasswordValid: true,
		 OneTimePasswordExpiry: expiryTime})

		 if result.Error != nil {
			return result.Error
		 }
		 return nil
}


func (r *userRepository) Search(term string, limit, page int) ([]models.User, error) {
	var users []models.User

	err := r.DB.Select("email", "userName", "phoneNumber", "bvn", "ID", "UUID").Where("email like ?", "%"+term+"%").Or("user_name = ?", "%"+term+"%").Find(&users)
	if err.Error != nil {
		return users, apperrors.NewInternal()
	}
	return users, nil
}


func (r *userRepository) UpdateUser(user models.User) error {
	userId := int(user.ID)
	foundUser, err := r.GetUserById(userId)

	if err != nil {
		return apperrors.NewInternal()
	}

	if foundUser == nil {
		return apperrors.NewNotFound("user ID", strconv.Itoa(userId))
	}
	
	updatedDetails := map[string] interface{}{}
	if user.Name != "" {
		updatedDetails["Name"] = user.Name
	}
	if user.UserName != "" {
		updatedDetails["UserName"] = user.UserName
	}
	if user.PhoneNumber != "" {
		updatedDetails["PhoneNumber"] = user.PhoneNumber
	}
	if user.Email != "" {
		updatedDetails["Email"] = user.Email
	}
	if !user.DOB.IsZero() {
		updatedDetails["DOB"] = user.DOB
	}
	if user.ProfileUrl != "" {
		updatedDetails["ProfileUrl"] = user.ProfileUrl
	}
	if user.ProfileIcon != "" {
		updatedDetails["ProfileIcon"] = user.ProfileIcon
	}
	if user.Gender != "" {
		updatedDetails["Gender"] = user.Gender
	}
	if user.Location != "" {
		updatedDetails["Location"] = user.Location 
	}

	if err := r.DB.Model(&foundUser).Updates(updatedDetails).Error; err != nil {
		return apperrors.NewInternal()
	}
	return nil
}


func (r *userRepository) UpdateUserTOTP(user models.User, totpSecret string, totpEnabled bool) error {
	if err := r.DB.Model(&user).Updates(map[string] interface{}{"totp_secret":totpSecret, "totp_enabled":totpEnabled}).Error;
		err != nil {
			log.Printf("error while updating user totp %v\n", err)
			if errors.Is(err, gorm.ErrRecordNotFound){
				return apperrors.NewNotFound("ID", strconv.Itoa(int(user.ID)))
			}
			return apperrors.NewInternal()
		}
		return nil 
}


func (r *userRepository) UpdatePassword(id uint, password string) error {
	foundUser, err := r.GetUserById(int(id))

	if err != nil {
		return apperrors.NewInternal()
	}

	if foundUser == nil {
		return apperrors.NewNotFound("user ID", strconv.Itoa(int(id)))
	}

	erro := r.DB.Model(&foundUser).Updates(models.User{Password: password})

	if erro.Error != nil {
		return erro.Error
	}
	return nil
}


func (r *userRepository) GetUserTransactions(userId, limit, page int) ([]models.Transactions, error) {
	var transactions []models.Transactions
	foundUser, _ := r.GetUserById(userId)

	if foundUser == nil {
		log.Print("Could not find user with provided id")
		return transactions, apperrors.NewBadRequest("Could not find user with provided id")
	}

	if err := r.DB.Where("owner_id = ?", userId).Find(&transactions).Error; err != nil {
		log.Print("Could not find user transactions\n")
		return transactions, apperrors.NewBadRequest("Could not find user transaction")
	}

	return transactions, nil
}


func (r userRepository) GetUserByItemId(id int) (*models.User, error) {
	item := &models.Item{}
	user := &models.User{}

	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		log.Print("Could not find item with provided id")
		return nil, apperrors.NewBadRequest("Couldnt find item with provided id")
	}

	if err := r.DB.Where("id = ?", item.OwnerId).First(&user).Error; err != nil {
		log.Print("Could not find user with provided item id")
		return nil, apperrors.NewBadRequest("Could not find user with provided item id")
	}

	return user, nil
}