package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)


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


func UploadItemFile(itemId int, file *multipart.FileHeader, uploadDir string) (string, error) {
	err := os.MkdirAll(uploadDir, os.ModePerm)

	if err != nil {
		return "", fmt.Errorf("Failed to create upload directory: %v", err)
	}


	itemDir := filepath.Join(uploadDir, fmt.Sprintf("item_%d", itemId))
	err = os.MkdirAll(itemDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Failed to create item directory: %v", err)
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)

	filePath := filepath.Join(itemDir, filename)

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