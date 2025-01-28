package services

import (
	"log"
	"os"
	"time"

	"swap/apperrors"
	"swap/models"
	"swap/utils"
	"strconv"

	"github.com/sethvargo/go-password/password"
)

type userService struct {
	UserRepository models.IUserRepository
}

func NewUserService(UserRepository models.IUserRepository) models.IUserService {
	return &userService{
		UserRepository: UserRepository,
	}
}

func (s *userService) Register(user *models.User) (*models.User, error) {
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		log.Printf("Error hashing password. Unable to register user with email: %s\n", user.Email)
		return nil, apperrors.NewInternal()
	}

	user.Password = hashedPassword
	
	return s.UserRepository.CreateUser(user)
}


func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.UserRepository.FindUserByEmailOrUsername(email)

	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	match, err := comparePassword(user.Password, password)

	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return user, nil
}


func (s *userService) GetUserById(id int) (*models.User, error) {
	user, err := s.UserRepository.GetUserById(id)
	if err != nil {
		return nil, apperrors.NewNotFound("user", strconv.Itoa(id))
	}
	return user, nil
}


func (s *userService) GetUserByUUID(id string) (*models.User, error) {
	user, err := s.UserRepository.GetUserByUUID(id)
	if err != nil {
		return nil, apperrors.NewNotFound("user", id)
	}
	return user, nil
}


func (s *userService) InitLoginWithOneTimePassword(email string) (error) {
	user, err := s.UserRepository.FindUserByEmail(email)
	if err != nil {
		log.Print(err)
		return err
	}

	otp, err := s.GenerateOneTimePasswordForUser(user, models.OneTimeLoginOTPType, time.Hour)
	if err != nil {
		log.Print(err)
		return err
	}

	err = utils.SendEmail(os.Getenv("EMAIL_SENDER_EMAIL"), user.Email, "One-time-password", otp)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}


func (s *userService) GenerateOneTimePasswordForUser(user *models.User, otpType models.OTPType, duration time.Duration) (string, error) {

	otp, err := s.generateOTPPasskey(otpType)
	if err != nil {
		log.Print(err)
		return "", err
	}

	expiry := time.Now().Add(duration)

	hashedPassword, err := hashPassword(otp)
	if err != nil {
		log.Print(err)
		return "", err
	}

	err = s.UserRepository.CreateOneTimePassword(user, hashedPassword, expiry)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return otp, nil
}


func (s *userService) generateOTPPasskey(otpType models.OTPType) (string, error) {
	var passKey string
	var err error

	if otpType == models.OneTimeLoginOTPType {
		passKey, err = password.Generate(6, 2, 0, false, false)
	}

	if err != nil {
		log.Print(err)
		return "", err
	}
	return passKey, nil
}

func (s *userService) LoginWithOneTimePassword(email, password string) (*models.User, error) {
	user, err := s.UserRepository.FindUserByEmail(email)

	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	if user.OneTimePasswordExpiry.Before(time.Now()) {
		return nil, apperrors.NewBadRequest("One time password already expired")
	}

	if !user.OneTimePasswordValid {
		return nil, apperrors.NewBadRequest(apperrors.InvalidCredentials)
	}

	match, err := comparePassword(user.OneTimePassword, password)
	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	err = s.UserRepository.InvalidateOneTimePassword(user)
	if err != nil {
		return nil, apperrors.NewInternal()
	}

	return user, nil
}


func (s *userService) InvalidateOneTimePassword(user *models.User) error {
	err := s.UserRepository.InvalidateOneTimePassword(user)
	if err != nil {
		log.Printf("error invalidating one time password for user %s: %+v\n", user.UUID, err)
		return apperrors.NewBadRequest("Unable to invalidate OTP")
	}

	return nil
}


func (s *userService) UpdatePassword(userId uint, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Print(err)
		return err
	}

	err = s.UserRepository.UpdatePassword(userId, hashedPassword)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}


func (s *userService) ConfirmPassword(userId uint, password string) (*models.User, error) {
	user, err := s.UserRepository.GetUserById(int(userId))
	if err != nil {
		return nil, apperrors.NewNotFound("user", strconv.Itoa(int(userId)))
	}

	match, err := comparePassword(password, user.Password)

	if err != nil {
		return user, apperrors.NewInternal()
	}

	if !match {
		return user, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return user, nil
}


func (s *userService) Search(term string, limit, page int) ([]models.User, error) {
	return s.UserRepository.Search(term, limit, page)
}


func (s *userService) FindUserByEmailOrUsername(email string) (*models.User, error) {
	user, err := s.UserRepository.FindUserByEmailOrUsername(email)
	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return user, nil
}


func (s *userService) FindUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	return s.UserRepository.FindUserByPhoneNumber(phoneNumber)
}


func (s *userService) UpdateUser(user models.User) error {
	return s.UserRepository.UpdateUser(user)
}


func (s *userService) EnrollTOTP(userId int) ([]byte, error) {
	user, err := s.UserRepository.GetUserById(userId)
	if err != nil {
		return nil, apperrors.NewBadRequest("Cannot enroll totp. User does exist!")
	}

	totpSecret := utils.GenerateTOTPSecret()
	hashedTotpSecret, err := utils.Encrypt(totpSecret)
	if err != nil {
		log.Printf("Error on user totp enrollment: %v\n", err)
		return nil, apperrors.NewInternalWithMessage("Unable to enroll user totp. Pleae try again later")
	}

	_ = s.UserRepository.UpdateUserTOTP(*user, hashedTotpSecret, false)
	if err != nil {
		return nil, apperrors.NewInternalWithMessage("Unable to enroll user totp. Pleae try again later")
	}

	return utils.GenerateTOTPQRCode(totpSecret, user.Email)
}


func (s *userService) VerifyTOTP(userId int, totp models.VerifyTOTPRequest) error {
	user, err := s.UserRepository.GetUserById(userId)
	if err != nil {
		return apperrors.NewBadRequest("Unable to verify totp. User does not exist")
	}

	match, err := s.verifyUserTotp(user, totp.Totp)
	if err != nil {
		log.Printf("Error verifying totp: %v\n", err)
		return err
	}

	if !match {
		return apperrors.NewBadRequest("One time password is not valid")
	}

	var updatedUser models.User
	updatedUser.ID = user.ID
	updatedUser.TotpEnabled = true

	_ = s.UpdateUser(updatedUser)
	if err != nil {
		return apperrors.NewInternalWithMessage("Unable to enroll user totp. Try again later.")
	}

	return nil
}


func (s *userService) verifyUserTotp(user *models.User, totp string) (bool, error) {
	if totp == "" {
		log.Printf("User totp is empty. Totp should be be enrolled!")
		return false, apperrors.NewBadRequest("User totp is invalid. Please enroll your totp.")
	}

	totpSecret, err := utils.Decrypt(totp)
	if err != nil {
		return false, apperrors.NewInternalWithMessage("Unable to verify totp. Please try again.")
	}

	return utils.VerifyTOTP(totpSecret, totp), nil
}


func (s *userService) DisableTOTP(userId int) error {
	user, err := s.UserRepository.GetUserById(userId)
	if err != nil {
		return apperrors.NewBadRequest("Cannot disable totp. User does not exits!")
	}

	_ = s.UserRepository.UpdateUserTOTP(*user, "", false)
	if err != nil {
		return apperrors.NewInternalWithMessage("Unable to disable user totp. Pls try again")
	}
	return nil
}


func (s *userService) EnableTOTP(userId int) error {
	user, err := s.UserRepository.GetUserById(userId)
	if err != nil {
		return apperrors.NewBadRequest("Cannot enable totp. User does not exist!")
	}

	if user.TotpSecret == "" {
		return apperrors.NewBadRequest("Totp is not enrolled for this user!")
	}

	_ = s.UserRepository.UpdateUserTOTP(*user, user.TotpSecret, true)
	if err != nil {
		return apperrors.NewInternalWithMessage("Unable to enable user totp. Please trt again.")
	}

	return nil
}


func (s *userService) GetUserTransactions(userId, limit, page int) ([]models.Transactions, error){
	return s.UserRepository.GetUserTransactions(userId, limit, page)
}


func (s *userService) GetUserByItemId(id int) (*models.User, error) {
	return s.UserRepository.GetUserByItemId(id)
}