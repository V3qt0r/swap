package handler

import (
	"net/http"
	"log"
	"strconv"

	"swap/models"
	"swap/middleware"
	"swap/apperrors"
	"swap/api"

	"github.com/gin-gonic/gin"
)


type SwapHandler struct {
	swapService models.ISwapService
}


func NewSwapHandler(SwapService models.ISwapService) *SwapHandler {
	h := &SwapHandler{swapService: SwapService}
	return h
}


func (h *SwapHandler) InitiateSwapRequest(c *gin.Context) {
	var request api.InitiateSwapRequestPayload
	if ok := api.BindData(c, &request); !ok {
		log.Printf("Error binding data: %v\n", request)
		e := apperrors.NewBadRequest("Invalid swap request payload")
		c.JSON(e.Status(), e)
		return
	}

	userDetails, _  := c.Get("id")

	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not found", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID
	initiatorId := userId

	swapRequest, err := h.swapService.InitiateSwapRequest(request.Item1Id, request.Item2Id, initiatorId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Could not initialize swap request", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successfully initialized swap request", swapRequest))
}



func (h *SwapHandler) GetPendingSwapRequests(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not found", nil))
		return
	}
	ownerId := userDetails.(*middleware.User).ID

	swapRequests, err := h.swapService.GetPendingSwapRequests(int(ownerId), limitValue, pageValue)

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Could not get pending swap requests", nil))
		return
	}
	
	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", swapRequests))
}


func (h *SwapHandler) RejectSwapRequest(c *gin.Context) {
	routeId := c.Param("id")
	swapId, _ := strconv.Atoi(routeId)

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "User not found", nil))
		return
	}

	ownerId := int(userDetails.(*middleware.User).ID)

	err := h.swapService.RejectSwapRequest(ownerId, swapId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Could not reject swap request", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *SwapHandler) AcceptSwapRequest(c *gin.Context) {
	routeId := c.Param("id")
	swapId, _ := strconv.Atoi(routeId)

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "User not found", nil))
		return
	}

	ownerId := int(userDetails.(*middleware.User).ID)

	result, err := h.swapService.AcceptSwapRequest(ownerId, swapId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Could not accept swap request", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", result))
}



func (h *SwapHandler) CompleteSwapRequest(c *gin.Context) {
	routeId := c.Param("id")
	swapId, _ := strconv.Atoi(routeId)

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "User not found", nil))
		return
	}

	ownerId := userDetails.(*middleware.User).ID


	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Could not bind JSON data", nil))
		return
	}

	amount, ok := request["amount"].(float64)

	if !ok {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid payload", nil))
		return
	}

	result, err := h.swapService.CompleteSwapRequest(ownerId, amount, uint(swapId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Could not complete swap request", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successfull", result))
}


func (h *SwapHandler) GetIncompleteSwapByInitiatorId(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Could not bind JSON data", nil))
		return
	}

	itemId, ok := request["itemId"].(float64)

	if !ok || itemId < 0{
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest,"Invalid payload", nil))
		return
	}

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "User not found", nil))
		return
	}

	initiatorId := int(userDetails.(*middleware.User).ID)

	swapRequest, err := h.swapService.GetIncompleteSwapByInitiatorId(initiatorId, int(itemId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Could not get incomplete swap", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", swapRequest))
}


func (h *SwapHandler) GetAllIncompleteSwapByOwnerId(c *gin.Context) {
	limit := c.Param("limit")
	page := c.Param("page")

	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "User not found", nil))
		return
	}

	ownerId := int(userDetails.(*middleware.User).ID)

	swapRequests, err := h.swapService.GetAllIncompleteSwapByOwnerId(ownerId, limitValue, pageValue)

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusInternalServerError, "Could not get incomplete swaps", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", swapRequests))
}