package handler

import (
	"strings"
	"net/http"
	"swap/api"
	"swap/models"

	"strconv"

	"github.com/gin-gonic/gin"
)


type CategoryHandler struct {
	categoryService models.ICategoryService
}


func NewCategoryHandler(CategoryService models.ICategoryService) *CategoryHandler {
	h := &CategoryHandler { categoryService : CategoryService }
	return h
}



func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid category payload", nil))
		return
	}

	name, ok := request["name"].(string)
	if !ok || name == "" {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Cannot add empty or invalid category", nil))
		return
	}

	category, err := h.categoryService.CreateCategory(strings.ToUpper(name))

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to add category", nil))
		return
	}

	c.JSON(http.StatusCreated, api.NewResponse(http.StatusCreated, "Successful", category))
}


func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid category payload", nil))
		return
	}

	name, ok := request["name"].(string)
	if !ok || name == "" {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Cannot add empty or invalid category", nil))
		return
	}

	name = strings.TrimSpace(name)
	err := h.categoryService.DeleteCategory(strings.ToUpper(name))

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to delete category", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *CategoryHandler) GetAllValidCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllValidCategories()

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "No valid category found", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", categories))
}


func (h *CategoryHandler) BanCategory(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid category payload", nil))
		return
	}

	name, ok := request["name"].(string)
	if !ok || name == "" {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Cannot ban empty or invalid category", nil))
		return
	}

	name = strings.TrimSpace(name)
	err := h.categoryService.BanCategory(strings.ToUpper(name))

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to ban category", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *CategoryHandler) UnBanCategory(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid category payload", nil))
		return
	}

	name, ok := request["name"].(string)
	if !ok || name == "" {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Cannot unban empty or invalid category", nil))
		return
	}

	name = strings.TrimSpace(name)
	err := h.categoryService.UnBanCategory(strings.ToUpper(name))

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to unban category", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}



func (h *CategoryHandler) CheckStatus(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid category payload", nil))
		return
	}

	name, ok := request["name"].(string)
	if !ok || name == "" {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Cannot check empty or invalid category", nil))
		return
	}

	name = strings.TrimSpace(name)
	result, err := h.categoryService.CheckStatus(strings.ToUpper(name))

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to check category status", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", result))
}


func (h *CategoryHandler) GetAllItemsInCategory(c *gin.Context) {
	routeId := c.Param("id")
	categoryId, _ := strconv.Atoi(routeId)

	limit := c.Query("limit")
	page := c.Query("page")

	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	items, err := h.categoryService.GetAllItemsInCategory(categoryId, limitValue, pageValue)

	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Couldnt get items", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", items))
}