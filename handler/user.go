package handler

import (
	"time"
	"log"
	"net/http"
	"strings"
	"swap/utils"
	"strconv"
	"swap/api"
	"swap/apperrors"
	"swap/middleware"
	"swap/models"

	"github.com/gin-gonic/gin"
)


type UserHandler struct {
	userService models.IUserService
}


func NewUserHandler(UserService models.IUserService) *UserHandler {
	// Create a handler (which will later have injected services)
	h := &UserHandler{userService : UserService}
	return h
}


func (h *UserHandler) Register(c *gin.Context) {
	var request api.RegisterPayload

	if ok := api.BindData(c, &request); !ok {
		return
	}

	dob, err := time.Parse("2006-01-02", request.DOB)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Error parsing endDate", "Please pass date in the format: YYYY-MM-DD HH:MM:SS"))
		return
	}

	c.Get("id")
	request.Sanitize()
	registerUserPayload := &models.User{
		Name: request.Name,
		UserName: request.UserName,
		PhoneNumber: request.PhoneNumber,
		Email: request.Email,
		DOB: dob,
		Gender: request.Gender,
		Password: request.Password,
		Location: request.Location,
		BVN: request.BVN,
	}
	user, err := h.userService.Register(registerUserPayload)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.DuplicateEmail).Error() {
			ToFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}
		c.JSON(apperrors.Status(err), api.NewResponse(apperrors.Status(err), "Unable to register user", gin.H{
			"error": err,
		}))
		return
	}
	c.JSON(http.StatusCreated, api.NewResponse(http.StatusCreated, "Successful", user))
}


func (h *UserHandler) GetUserById(c *gin.Context) {
	routeId := c.Param("id")

	userId, err := strconv.Atoi(routeId)
	if err != nil {
		log.Printf("error processing user id: %v\n", err)
		e := apperrors.NewBadRequest("invalid user id provided")
		c.JSON(e.Status(), e)
		return
	}

	h.getUserById(c, userId)
}


func (h *UserHandler) GetLoggedInUser(c *gin.Context) {
	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}
	userId := userDetails.(*middleware.User).ID

	h.getUserById(c, int(userId))
}


func (h *UserHandler) getUserById(c *gin.Context, userId int) {
	user, err := h.userService.GetUserById(userId)

	if err != nil {
		log.Printf("Unable to find user with ID: %v\n%v\n", userId, err)
		e := apperrors.NewNotFound("User", strconv.Itoa(userId))

		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Message, gin.H{
			"error" : e,
		}))
		return
	}
	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", user))
}


func (h *UserHandler) UpdateUser(c *gin.Context) {
	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Error retrieving user details", 
			nil))
			return
	}

	userId := userDetails.(*middleware.User).ID

	var request api.UserUpdatePayload
	if ok := api.BindData(c, &request); !ok {
		e := apperrors.NewBadRequest("Invalid user payload")
		c.JSON(e.Status(), e)
		return
	}
	request.Sanitize()

	user := request.ToEntity()
	user.ID = userId

	err := h.userService.UpdateUser(user)
	if err != nil {
		e := apperrors.GetAppError(err, "Update User failed!")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) Search(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	var request api.UserSearchRequest
	if ok := api.BindData(c, &request); !ok {
		return
	}

	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	users, err := h.userService.Search(request.SearchTerm, limitValue, pageValue)
	if err != nil {
		e := apperrors.NewInternal()

		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Message, gin.H{ "error" : e, }))
		return
	}

	var searchResponses []api.UserSearchResponse
	
	for _, user := range users {
		searchResponse := api.UserSearchResponse{
			UserName: 		user.UserName,
			Email: 			user.Email,
			PhoneNumber:	user.PhoneNumber,
			UUID: 			user.UUID.String(),
		}
		searchResponses = append(searchResponses, searchResponse)
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", searchResponses))
}


func (h *UserHandler) SendOneTimePassword(c *gin.Context) {
	var request api.InitOnetimePasswordPayload
	if ok := api.BindData(c, &request); !ok {
		return
	}

	err := h.userService.InitLoginWithOneTimePassword(request.Email)

	if err != nil {
		log.Print(err)
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var request api.UpdatePasswordPayload

	if ok := api.BindData(c, &request); !ok {
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Failed", "Passwords do not match"))
		return
	}

	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Error retieving user details",
			nil))
		return
	}

	userId := userDetails.(*middleware.User).ID

	err := h.userService.UpdatePassword(userId, request.Password)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Password could not be updated",
			nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) ConfirmUserPassword(c *gin.Context) {
	var request api.ConfirmPasswordPayload

	if ok := api.BindData(c, &request); !ok {
		return
	}

	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Error retieving user details",
			nil))
		return
	}

	userId := userDetails.(*middleware.User).ID

	_, err := h.userService.ConfirmPassword(userId, request.Password)
	if err != nil {
		log.Print(err)
		e := apperrors.GetAppError(err, "Password could not be verified")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), gin.H{ "isValid" : false}))
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", gin.H{ "isValid" : true}))
}


// func (h *UserHandler) UploadFile(c *gin.Context) {
// 	UploadFile(c)
// 	return
// }

func (h *UserHandler) UploadFile(c *gin.Context){
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "No file uploaded", nil))
		return
	}

	uploadDir := "./uploads"

	filePath, err := utils.UploadFile(file, uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "File not upoaded", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "File uploaded successfully", filePath))
}

func (h *UserHandler) EnrollTOTP(c *gin.Context) {
	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, 
				"Error retrieving user details", nil))
		return
	}
	userId := userDetails.(*middleware.User).ID

	totpQRCode, err := h.userService.EnrollTOTP(int(userId))

	if err != nil {
		e := apperrors.GetAppError(err, "User totp enrollment failed")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}
	c.Data(http.StatusOK, "image/png", totpQRCode)
}


func (h *UserHandler) VerifyTOTP(c *gin.Context) {
	var request models.VerifyTOTPRequest
	if ok := api.BindData(c, &request); !ok {
		return
	}

	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, 
			"Error getting user details", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID

	err := h.userService.VerifyTOTP(int(userId), request)
	if err != nil {
		e := apperrors.GetAppError(err, "Error verifying user totp")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) DisableTOTP(c *gin.Context) {
	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, 
			"Error getting user details", nil))
		return
	}
	userId := userDetails.(*middleware.User).ID

	err := h.userService.DisableTOTP(int(userId))
	if err != nil {
		e := apperrors.GetAppError(err, "Cannot disable totp, please try again")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) EnableTOTP(c *gin.Context) {
	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Error getting user details", nil))
			return
	}
	userId := userDetails.(*middleware.User).ID

	err := h.userService.EnableTOTP(int(userId))
	if err != nil {
		e := apperrors.GetAppError(err, "Cannot enable totp, please try again")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) FindUserByEmailOrUsername(c *gin.Context) {
	var request api.FindUserByEmailOrUsernamePayload
	if ok := api.BindData(c, &request); !ok {
		return
	}

	request.Sanitize()
	request.Validate()
	user, err := h.userService.FindUserByEmailOrUsername(request.EmailOrUsername)

	if err != nil {
		e := apperrors.GetAppError(err, "Cannot find user by details provided")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	var response api.UserSearchResponse
	response.UserName = user.UserName
	response.Email = user.Email
	response.PhoneNumber = user.PhoneNumber
	response.UUID = user.UUID.String()
	
	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", response))
}


func (h *UserHandler) FindUserByPhoneNumber(c *gin.Context) {
	phoneNumber := c.Query("phoneNumber")
	if phoneNumber == "" {
		return
	}

	phoneNumber = strings.TrimSpace(phoneNumber)

	user, err := h.userService.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		e := apperrors.GetAppError(err, "Could not find user by phone number")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}
	var response api.UserSearchResponse
	response.UserName = user.UserName
	response.Email = user.Email
	response.PhoneNumber = phoneNumber
	response.UUID = user.UUID.String()

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", response))
}


func (h *UserHandler) GetUserTransactions(c *gin.Context) {
	userDetails, _ := c.Get("id")
	limit := c.Query("limit")
	page := c.Query("page")

	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Error getting user details", nil))
		return
	}
	userId := int(userDetails.(*middleware.User).ID)
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	transactions, err := h.userService.GetUserTransactions(userId, limitValue, pageValue)

	if transactions == nil {
		log.Print("No transaction available")
		c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "No transaction available", nil))
		return
	}
	if err != nil {
		log.Print("Error getting user transactions")
		e := apperrors.NewInternal()

		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Message, gin.H{ "error" : e, }))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", transactions))
}