package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"log"
	"swap/models"


	"io/ioutil"
)

type Utils struct {
	imageRepository models.IImageRepository
}


func NewUtils(imageRepo models.IImageRepository) *Utils {
	return &Utils{
		imageRepository : imageRepo,
	}
}


func UploadFile(file *multipart.FileHeader, uploadDir string) (string, error) {
	err := os.MkdirAll(uploadDir, os.ModePerm)

	if err != nil {
		return "", fmt.Errorf("Failed to create upload directory: %v", err)
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)

	filePath := filepath.Join(uploadDir, filename)

	if err := saveFile(file, filePath); err != nil {
		return "", err
	}
	return filePath, nil
}


func (u *Utils)UploadItemFile(itemId int, file *multipart.FileHeader, uploadDir string) (string, error) {
	log.Printf("Uploading image in item ID %d", itemId)
	err := os.MkdirAll(uploadDir, os.ModePerm)


	if err != nil {
		return "", fmt.Errorf("Failed to create upload directory: %v", err)
	}
	

	if u.imageRepository == nil {
		log.Print("Item repository not initialized")
		return "", fmt.Errorf("Item repository not initialized")
	}


	itemDir := filepath.Join(uploadDir, fmt.Sprintf("item_%d", itemId))
	err = os.MkdirAll(itemDir, os.ModePerm)

	
	if err != nil {
		return "", fmt.Errorf("Failed to create item directory: %v", err)
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)

	filePath := filepath.Join(itemDir, filename)

	err = u.imageRepository.UploadImage(itemId, itemDir, filename)
	if err != nil {
		log.Print("Failed to update images in item")
		return "", fmt.Errorf("Failed to upload image")
	}

	if err := saveFile(file, filePath); err != nil {
		return "", fmt.Errorf("Failed to save file: %v", err)
	}

	return filePath, nil
}


func saveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("Failed to open file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v", err)
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return fmt.Errorf("Failed to save file: %v", err)
	}
	return nil
}


func ReadImage(folderName, fileName string) ([]byte, error) {
	filePath := fmt.Sprintf("uploads/%s/%s", folderName, fileName)


	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return nil, fmt.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}

	return data, nil
}


func ReadImageByPath(path string) ([]byte, error) {
	filePath := fmt.Sprintf("uploads/%s", path)


	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return nil, fmt.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}

	return data, nil
}