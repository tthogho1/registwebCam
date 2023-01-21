package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	max        int   = 80000
	offset_max int   = 50
	offsets    []int = make([]int, max/offset_max)
	increment  int   = 0
)

var (
	baseurl    string = "https://api.windy.com/api/webcams/v2/list/limit=" + strconv.Itoa(offset_max) + ","
	parameters string = "?lang=en&key=4tpguJklGSjb3f0nVny1wwR9bqHquToz&show=webcams:image,player,location"

	client = new(http.Client)
	//var URL = baseurl + "offset=" + strconv.Itoa(t) + parameters
	wg sync.WaitGroup
)

var log = logrus.New()

func main() {

	// ファイル出力
	file, err := os.OpenFile("c:\\temp\\logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	for i := 0; i < max/offset_max; i++ {
		var t = i * offset_max
		offsets[i] = t
	}

	/*	for i := 0; i < max/offset_max; i++ {
		var t = getOffset()
		fmt.Println(t)
	} */

	increment = 0
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go registerWebCameraToMongoDB(&wg, i)
	}

	wg.Wait()
}

var (
	mu sync.Mutex
)

func getOffset(increment *int) int {

	if *increment >= max/offset_max {
		return -1
	}

	fmt.Println("increment: ", *increment)
	t := offsets[*increment]

	return t
}

type WebCameraInfo struct {
	Status string `json:"status"`
	Result struct {
		Offset  int `json:"offset"`
		Limit   int `json:"limit"`
		Total   int `json:"total"`
		Webcams []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Title  string `json:"title"`
			Image  struct {
				Current struct {
					Icon      string `json:"icon"`
					Thumbnail string `json:"thumbnail"`
					Preview   string `json:"preview"`
					Toenail   string `json:"toenail"`
				} `json:"current"`
				Sizes struct {
					Icon struct {
						Width  int `json:"width"`
						Height int `json:"height"`
					} `json:"icon"`
					Thumbnail struct {
						Width  int `json:"width"`
						Height int `json:"height"`
					} `json:"thumbnail"`
					Preview struct {
						Width  int `json:"width"`
						Height int `json:"height"`
					} `json:"preview"`
					Toenail struct {
						Width  int `json:"width"`
						Height int `json:"height"`
					} `json:"toenail"`
				} `json:"sizes"`
				Daylight struct {
					Icon      string `json:"icon"`
					Thumbnail string `json:"thumbnail"`
					Preview   string `json:"preview"`
					Toenail   string `json:"toenail"`
				} `json:"daylight"`
				Update int `json:"update"`
			} `json:"image"`
			Location struct {
				City          string  `json:"city"`
				Region        string  `json:"region"`
				RegionCode    string  `json:"region_code"`
				Country       string  `json:"country"`
				CountryCode   string  `json:"country_code"`
				Continent     string  `json:"continent"`
				ContinentCode string  `json:"continent_code"`
				Latitude      float64 `json:"latitude"`
				Longitude     float64 `json:"longitude"`
				Timezone      string  `json:"timezone"`
				Wikipedia     string  `json:"wikipedia"`
			} `json:"location"`
			Player struct {
				Live struct {
					Available bool   `json:"available"`
					Embed     string `json:"embed"`
				} `json:"live"`
				Day struct {
					Available bool   `json:"available"`
					Link      string `json:"link"`
					Embed     string `json:"embed"`
				} `json:"day"`
				Month struct {
					Available bool   `json:"available"`
					Link      string `json:"link"`
					Embed     string `json:"embed"`
				} `json:"month"`
				Year struct {
					Available bool   `json:"available"`
					Link      string `json:"link"`
					Embed     string `json:"embed"`
				} `json:"year"`
				Lifetime struct {
					Available bool   `json:"available"`
					Link      string `json:"link"`
					Embed     string `json:"embed"`
				} `json:"lifetime"`
			} `json:"player"`
		} `json:"webcams"`
	} `json:"result"`
}

func registerWebCameraToMongoDB(wg *sync.WaitGroup, id int) {

	defer wg.Done()
	//  - バックグラウンドで接続する。タイムアウトは10秒
	ctx := context.TODO()

	// Create a new client and connect to the server
	con, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
	defer con.Disconnect(ctx)
	if err != nil {
		panic(err)
	}

	coll := con.Database("webcam").Collection("webcam")

	for {
		// offset 取得
		mu.Lock()
		offset := getOffset(&increment)
		increment++
		mu.Unlock()

		//fmt.Println(strconv.Itoa(id) + ":" + strconv.Itoa(offset))
		if offset == -1 {
			break
		}

		// url作成
		url := baseurl + strconv.Itoa(offset) + parameters
		fmt.Println(url)

		// リクエスト取得
		resp, err := client.Get(url)
		if err != nil {
			panic(err)
		}
		//defer
		// レスポンス取得
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		var webCameraInfo WebCameraInfo
		if err := json.Unmarshal(body, &webCameraInfo); err != nil {
			panic(err)
		}

		result_len := len(webCameraInfo.Result.Webcams)
		//fmt.Println(t)

		for _, webCam := range webCameraInfo.Result.Webcams {

			_, err := coll.InsertOne(ctx, webCam)
			if err != nil {
				panic(err)
			}
			fmt.Println(webCam.ID)

			log.Info(webCam.ID)
		}

		if result_len < offset_max {
			fmt.Println(result_len)
			break
		}

	}
}

var (
	mongouri string = "mongodb://webcam:webcam@localhost:27017/webcam"
)
