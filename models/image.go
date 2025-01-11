package models

type Image struct {
	Base
	FilePath 		string		`json:"filePath"`
	FileName        string      `json:"fileName"`
	ItemId 			int			`json:"itemId"`
	OwnerId 		int 		`json:"ownerId"`
}


type IImageRepository interface {
	UploadImage(itemId int, folderName, fileName string) error
	ReadFirstImageById(id int) ([]byte, error)
	ReadAllImagesByItemId(id, limit, page int) ([]Image, error)
}


type IImageService interface {
	// UploadImage(itemId int, folderName, fileName string) error
	ReadImage(folderName, fileName string) ([]byte, error)
	ReadFirstImageById(id int) ([]byte, error)
	ReadAllImagesByItemId(id, limit, page int) ([]Image, error)
	ReadImageByPath(path string) ([]byte, error)
}