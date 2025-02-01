package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

var WINDY_API_KEY = os.Getenv("WINDY_API_KEY")

func GetImage(downloadUrl string) (data []byte) {

	// リクエストを作成
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ヘッダーを追加
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-windy-api-key", WINDY_API_KEY)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	imageData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("Status Code is " + strconv.Itoa(res.StatusCode))
		return
	}

	return imageData

}
