package util

import (
	"app/internal/models"
	"bufio"
	"encoding/base64"
	"fmt"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func CropAndSaveImage(wg *sync.WaitGroup, sizeChan <-chan int, img image.Image, userID uint, imageFormat string, imageName string, store *gorm.DB) {
	defer wg.Done()

	size := <-sizeChan
	//Cropping the image
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	cropSize := image.Rect((width/2)-(size/2), (height/2)-(size/2), (width/2)+size-(size/2), (height/2)+size-(size/2))
	croppedImage := img.(SubImager).SubImage(cropSize)

	//Create a folder for the user if it does not already exist
	dirPath := fmt.Sprintf("./images/%v", userID)
	err := os.Mkdir(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Println(err)
		return
	}

	filePath := fmt.Sprintf("./images/%v/%v_%v.%v", userID, size, imageName, imageFormat)
	out, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()

	store.Create(&models.ImageURLs{
		UserID:   userID,
		ImageURL: filePath,
	})

	//Saving the image
	if imageFormat == "jpeg" {
		err = jpeg.Encode(out, croppedImage, nil)
	} else if imageFormat == "png" {
		err = png.Encode(out, croppedImage)
	}

}

func ConvertToBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return "", err
	}

	imageBytes := make([]byte, fileStat.Size())
	_, err = bufio.NewReader(file).Read(imageBytes)
	if err != nil && err != io.EOF {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageBytes), nil

}

func GetImageName(path string) string {
	splittedPath := strings.FieldsFunc(path, split)
	return splittedPath[len(splittedPath)-2]
}

func split(r rune) bool {
	return r == '.' || r == '/'
}
