package handler

import (
	"log"
	"net/http"
	"strconv"

	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"

	"swap/api"
	"swap/models"

	"github.com/gin-gonic/gin"
)


type ImageHandler struct {
	imageService models.IImageService
}

func NewImageHandler(ImageService models.IImageService) *ImageHandler {
	h := &ImageHandler{ imageService: ImageService }
	return h
}


func (h *ImageHandler) ReadFirstImageById(c *gin.Context) {
	id := c.Param("id")
	itemId, _ := strconv.Atoi(id)

	image, err := h.imageService.ReadFirstImageById(itemId)

	if err != nil {
		log.Print("Could not get any image")
		placeHolderImage := generatePlaceHolderImage()
		c.Data(http.StatusOK, "image/jpeg", placeHolderImage)
		return
	}

	c.Data(http.StatusOK, "image/jpeg", image)
}

func generatePlaceHolderImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)

	return buf.Bytes()
}

func (h *ImageHandler) ReadAllImagesByItemId(c *gin.Context) {
	id := c.Param("id")
	limit := c.Query("limit")
	page := c.Query("page")

	itemId, _  := strconv.Atoi(id) 
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	images, err := h.imageService.ReadAllImagesByItemId(itemId, limitValue, pageValue)

	if err != nil {
		log.Print("Cannot get images for this item")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Cannot get images for this item", nil))
		return
	}

	var imagePaths []string

	for _, image := range images {
		imagePaths = append(imagePaths, image.FilePath+"/"+image.FileName)
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", imagePaths))
}

func (h *ImageHandler) ReadImage(c *gin.Context) {
	data, err := h.imageService.ReadImage("item_17", "1735237103_FB_IMG_16159860997307807.jpg")
	
	if err != nil {
		log.Print("Failed to read image")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Failed to read image", nil))
	}

	c.Data(http.StatusOK, "image/jpeg", data)
}


func (h *ImageHandler) ReadImageByPath(c *gin.Context) {
	imagePath := c.Query("imagePath")

	if imagePath == "" {
		log.Print("Image path query parameter is required")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Image path query parameter is required", nil))
	}

	data, err := h.imageService.ReadImageByPath(imagePath)
	
	if err != nil {
		log.Print("Failed to read image")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Failed to read image", nil))
	}

	c.Data(http.StatusOK, "image/jpeg", data)
}
