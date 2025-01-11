package handler

import (
	"log"
	"net/http"
	"strings"
	"strconv"
	"swap/utils"
	"swap/api"
	"swap/apperrors"
	"swap/middleware"
	"swap/models"

	"github.com/gin-gonic/gin"
)


type ItemHandler struct {
	itemService models.IItemService
	Utils *utils.Utils
}


func NewItemHandler(ItemService models.IItemService, util *utils.Utils) *ItemHandler{
	h := &ItemHandler{ itemService: ItemService, Utils: util }
	return h
}


func (h *ItemHandler) RegisterItem(c *gin.Context) {
	var request api.RegisterItemPayload
	if ok := api.BindData(c, &request); !ok {
		log.Printf("Error binding data: %v\n", request)
		e := apperrors.NewBadRequest("Invalid item payload")
		c.JSON(e.Status(), e)
		return
	}

	request.Sanitize()

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID
	request.OwnerId = userId

	registerItemPayload := &models.Item{
		Name:			request.Name,
		CategoryName:	request.CategoryName,
		Description:	request.Description,
		Prize:			request.Prize,
		OwnerId:		request.OwnerId,
	}

	item, err := h.itemService.RegisterItem(registerItemPayload)
	if err != nil {
		log.Print("Unable to register item")
		c.JSON(apperrors.Status(err), api.NewResponse(apperrors.Status(err), "Unable to register item", gin.H{
			"error": err,
		}))
		return
	}
	
	c.JSON(http.StatusCreated, api.NewResponse(http.StatusCreated, "Successful", item))
}


func (h *ItemHandler) GetItemById(c *gin.Context) {
	routeId := c.Param("id")
	itemId, _ := strconv.Atoi(routeId)

	item, err := h.itemService.GetItemById(itemId)

	if err != nil {
		log.Printf("Unable to find item with ID: %v\n", itemId)
		id := strconv.Itoa(itemId)
		e := apperrors.NewNotFound("item", id)
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Message, gin.H{ "error" : e, }))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", item))
}


func (h *ItemHandler) GetItemsByCategory(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	var request api.ItemSearchRequest
	if ok := api.BindData(c, &request); !ok {
		return
	}

	request.Sanitize()
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	items, err := h.itemService.GetItemsByCategory(strings.ToUpper(request.SearchTerm), limitValue, pageValue)
	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), api.NewResponse(e.Status(), "Couldnt get items", gin.H{ "error" : e, }))
		return
	}

	var responses []api.ItemSearchResponse

	for _, item := range items {
		response := api.ItemSearchResponse{
			Name:	 		item.Name,
			CategoryName:	item.CategoryName,
			Description:	item.Description,
			Prize:			item.Prize,
			UUID:			item.UUID.String(),
			ID:				int(item.ID),
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", responses))
}


func (h *ItemHandler) GetUnsoldItemsByCategory(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	var request api.ItemSearchRequest
	if ok := api.BindData(c, &request); !ok {
		return
	}

	request.Sanitize()
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)
    
	items, err := h.itemService.GetUnsoldItemsByCategory(strings.ToUpper(request.SearchTerm), limitValue, pageValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Couldnt get items", nil))
		return
	}

	var responses []api.ItemSearchResponse

	for _, item := range items {
		response := api.ItemSearchResponse{
			Name:	 		item.Name,
			CategoryName:	item.CategoryName,
			Description:	item.Description,
			Prize:			item.Prize,
			UUID:			item.UUID.String(),
			ID:				int(item.ID),
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", responses))
}


func (h *ItemHandler) UpdateItem(c *gin.Context) {
	routeId := c.Param("id")
	itemId := routeId

	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID

	var request api.ItemUpdatePayload
	if ok := api.BindData(c, &request); !ok {
		log.Printf("Error binding data: %v\n", request)
		e := apperrors.NewBadRequest("Invalid item payload")
		c.JSON(e.Status(), e)
		return
	}

	item := request.ToEntity()
	id, _ := strconv.Atoi(itemId)
	item.ID = uint(id)

	foundItem, _ := h.itemService.GetItemById(id)
	
	if foundItem.OwnerId != userId {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "You do not own this item", nil))
		return
	}

	err := h.itemService.UpdateItem(item)
	if err != nil {
		e := apperrors.GetAppError(err, "Update item failed!")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Update item failed!", e.Error()))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *ItemHandler) DeleteItem(c *gin.Context) {
	routeId := c.Param("id")
	if routeId == "" {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "ID is required", nil))
		return
	}

	itemId, err := strconv.Atoi(routeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid ID format", nil))
		return
	}


	err = h.itemService.DeleteItem(itemId)
	if err != nil {
		e := apperrors.GetAppError(err, "Item delete failed!")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Item delete failed!", e.Error()))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *ItemHandler) GetItemsByOwnerId(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	items, err := h.itemService.GetItemsByOwnerId(userId, limitValue, pageValue)

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), api.NewResponse(e.Status(), "Couldnt get items", gin.H{ "error": e }))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", items))
}


func (h *ItemHandler) BuyItem(c *gin.Context) {
	var request map[string] interface{}
	id := c.Param("id")

	itemId, _ := strconv.Atoi(id)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid item payload", nil))
		return
	}

	amount, ok := request["amount"].(float64)
	if !ok || amount < 0.00 {
		log.Print("Invalid or missing amount")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid or missing Amount", nil))
		return
	}

	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}
	// userEmail := userDetails.(*middleware.User).Email
	userId := userDetails.(*middleware.User).ID

	result, err := h.itemService.BuyItem(int(userId), int(itemId), amount)

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), api.NewResponse(e.Status(), "Failed to buy item", gin.H{ "error": e }))
		return
	}
	
	// err = utils.SendEmailWithDefaultSender(userEmail, "Purchase of Item", result)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Unable to send mail", nil))
	// 	return
	// }

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", result))
}



func (h *ItemHandler) UploadFile(c *gin.Context){
	itemIdParam := c.Param("id")
	itemId, _ := strconv.Atoi(itemIdParam)
	

	_, err := h.itemService.GetItemById(itemId)
	if err != nil {
		c.JSON(http.StatusNotFound, api.NewResponse(http.StatusNotFound, "Item not found", nil))
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "No file uploaded", nil))
		return
	}

	uploadDir := "./uploads"

	filePath, err := h.Utils.UploadItemFile(itemId, file, uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "File not upoaded", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "File uploaded successfully", filePath))
}


func (h *ItemHandler) UpdateCategory(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid item payload", nil))
		return
	}

	itemId, ok := request["id"].(float64)
	if !ok || itemId < 0 {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid or missing ID", nil))
		return
	}
	
	categoryName, ok := request["name"].(string)
	if !ok || itemId < 0 {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid or category name", nil))
		return
	}
	

	err := h.itemService.UpdateCategory(int(itemId), categoryName)

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to update category", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}