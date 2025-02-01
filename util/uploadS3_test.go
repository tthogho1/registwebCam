package util

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestUploadS3(t *testing.T) {

	err := godotenv.Load("..//.env") // .envファイルを読み込む
	if err != nil {
		panic(err)
	}
	println("test TestUploadS3 start")
	imgUrl := "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_120x44dp.png"

	println(imgUrl)
	imageData := GetImage(imgUrl)

	println("test TestUploadS3 start")
	UploadS3(imageData, "test.png")
	println("test TestUploadS3 end")
}
