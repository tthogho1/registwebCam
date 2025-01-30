package webcam

import (
	"time"
)

type Webcam struct {
	Title         string    `json:"title"`
	ViewCount     int       `json:"viewCount"`
	WebcamID      int       `json:"webcamId"`
	Status        string    `json:"status"`
	LastUpdatedOn time.Time `json:"lastUpdatedOn"`
	Categories    []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
	Images struct {
		Current struct {
			Icon      string `json:"icon"`
			Thumbnail string `json:"thumbnail"`
			Preview   string `json:"preview"`
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
		} `json:"sizes"`
		Daylight struct {
			Icon      string `json:"icon"`
			Thumbnail string `json:"thumbnail"`
			Preview   string `json:"preview"`
		} `json:"daylight"`
	} `json:"images"`
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
	} `json:"location"`
	Player struct {
		Day      string `json:"day"`
		Month    string `json:"month"`
		Year     string `json:"year"`
		Lifetime string `json:"lifetime"`
	} `json:"player"`
}

type WebCameraInfo struct {
	Total   int      `json:"total"`
	Webcams []Webcam `json:"webcams"`
}

type WebcamWithEmbedding struct {
	Webcam
	Embedding []float32 `bson:"embedding"`
}
