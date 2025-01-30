package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"registWebCam/util"
	"registWebCam/webcam"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	offset_max int = 50
	increment  int = 0
)

var (
	baseurl    string = "https://api.windy.com/webcams/api/v3/webcams?lang=en&limit=" + strconv.Itoa(offset_max) + "&offset=%s&regions=%s"
	parameters string = "&sortDirection=asc&include=categories,images,location,player"
	logfile    string = "c:\\temp\\logrus.log"

	requestUrl string = baseurl + parameters
)

var logger = logrus.New()

func main() {

	err := godotenv.Load(".env") // .envファイルを読み込む
	if err != nil {
		panic(err)
	}

	// ファイル出力
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to Open log to file, using default stderr")
	}

	regionCodes := extractRegionCode()
	logger.Println(regionCodes)
	//fmt.Println(regionCodes)

	for _, regionCode := range regionCodes {
		increment = 0

		registerWebCameraToMongoDB(0, regionCode)
		logger.Println(regionCode)
	}

}

var (
	mongouri string = "mongodb+srv://webcam:webcam@cluster0.pizmgb2.mongodb.net/?retryWrites=true&w=majority"
)

func extractRegionCode() []string {
	url := "https://api.windy.com/webcams/api/v3/regions?lang=en"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("accept", "application/json")

	WINDY_API_KEY := os.Getenv("WINDY_API_KEY")
	req.Header.Set("x-windy-api-key", WINDY_API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	jsonArray := []map[string]interface{}{}
	err = json.Unmarshal(body, &jsonArray)
	if err != nil {
		panic(err)
	}

	codes := []string{}
	for _, item := range jsonArray {
		if item["code"] != nil {
			codes = append(codes, item["code"].(string))
		}
	}

	//fmt.Println(codes)
	return codes
}

func registerWebCameraToMongoDB(id int, regionCode string) {

	// defer wg.Done()
	//  - バックグラウンドで接続する。タイムアウトは10秒
	ctx := context.TODO()

	// Create a new client and connect to the server
	con, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
	defer con.Disconnect(ctx)
	if err != nil {
		panic(err)
	}

	WEBCAM_DB := os.Getenv("WEBCAM_DB")
	WEBCAM_COLLECTION := os.Getenv("WEBCAM_COLLECTION")
	coll := con.Database(WEBCAM_DB).Collection(WEBCAM_COLLECTION)
	//rpgClient, _ := util.CreateClient()

	for {
		// offset 取得
		offset := increment * offset_max
		increment++

		// url作成
		url := fmt.Sprintf(requestUrl, strconv.Itoa(offset), regionCode)
		fmt.Println(url)

		// リクエストを作成
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// ヘッダーを追加
		req.Header.Add("Accept", "application/json")
		WINDY_API_KEY := os.Getenv("WINDY_API_KEY")
		req.Header.Add("x-windy-api-key", WINDY_API_KEY)

		// リクエストを送信
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		// レスポンスを取得
		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		if res.StatusCode != 200 {
			fmt.Println("Status Code is " + strconv.Itoa(res.StatusCode) + " Body is " + string(data))
			log.Info("Status Code is " + strconv.Itoa(res.StatusCode) + " Body is " + string(data))
			break
		}

		var webCameraInfo webcam.WebCameraInfo
		if err := json.Unmarshal(data, &webCameraInfo); err != nil {
			panic(err)
		}
		res.Body.Close()

		result_len := len(webCameraInfo.Webcams)
		fmt.Println("Id : " + strconv.Itoa(id) + " data length : " + strconv.Itoa(result_len))

		for _, webCam := range webCameraInfo.Webcams {

			var webCamWithEmd webcam.WebcamWithEmbedding
			webCamWithEmd.Webcam = webCam

			imgUrl := webCam.Images.Daylight.Thumbnail
			//var imageData []byte
			if imgUrl != "" {
				// 画像をダウンロード
				//imageData = util.GetImage(imgUrl)
				util.GetImage(imgUrl)
			}
			filename := "image_" + strconv.Itoa(webCam.WebcamID) + ".jpg"
			print(filename)
			//webCamWithEmd.Embedding = util.GetEmbedding(rpgClient, imageData, filename)

			_, err := coll.InsertOne(ctx, webCamWithEmd)
			if err != nil {
				panic(err)
			}

			//fmt.Println(webCam.WebcamID)
			log.Info(webCam.WebcamID)
			fmt.Println(webCam.WebcamID)
		}

		if result_len < offset_max {
			fmt.Println("exit Id : " + strconv.Itoa(id) + " data length : " + strconv.Itoa(result_len))
			fmt.Println(webCameraInfo.Webcams)
			break
		}

	}
}
