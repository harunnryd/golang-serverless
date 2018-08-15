package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Request struct {
	LPR ReqLPR `json:"lpr"`
}

type ReqLPR struct {
	ImageURL         string `json:"image_url"`
	SecretKey        string `json:"secret_key"`
	Country          string `json:"country"`
	RecognizeVehicle int    `json:"recognize_vehicle"`
	State            string `json:"state"`
	ReturnImage      int    `json:"return_image"`
	Topn             int    `json:"topn"`
}

// type ResponseLPR struct {
// 	UUID           string `json:"uuid"`
// 	DataType       string `json:"data_type"`
// 	EpochTime      int64  `json:"epoch_time"`
// 	ProcessingTime struct {
// 		Plates float64 `json:"plates"`
// 		Total  float64 `json:"total"`
// 	} `json:"processing_time"`
// 	ImgHeight int `json:"img_height"`
// 	ImgWidth  int `json:"img_width"`
// 	Results   []struct {
// 		Plate            string  `json:"plate"`
// 		Confidence       float64 `json:"confidence"`
// 		RegionConfidence int     `json:"region_confidence"`
// 		VehicleRegion    struct {
// 			Y      int `json:"y"`
// 			X      int `json:"x"`
// 			Height int `json:"height"`
// 			Width  int `json:"width"`
// 		} `json:"vehicle_region"`
// 		Region           string  `json:"region"`
// 		PlateIndex       int     `json:"plate_index"`
// 		ProcessingTimeMs float64 `json:"processing_time_ms"`
// 		Candidates       []struct {
// 			MatchesTemplate int     `json:"matches_template"`
// 			Plate           string  `json:"plate"`
// 			Confidence      float64 `json:"confidence"`
// 		} `json:"candidates"`
// 		Coordinates []struct {
// 			Y int `json:"y"`
// 			X int `json:"x"`
// 		} `json:"coordinates"`
// 		MatchesTemplate int `json:"matches_template"`
// 		RequestedTopn   int `json:"requested_topn"`
// 	} `json:"results"`
// 	CreditsMonthlyUsed  int  `json:"credits_monthly_used"`
// 	Version             int  `json:"version"`
// 	CreditsMonthlyTotal int  `json:"credits_monthly_total"`
// 	Error               bool `json:"error"`
// 	RegionsOfInterest   []struct {
// 		Y      int `json:"y"`
// 		X      int `json:"x"`
// 		Height int `json:"height"`
// 		Width  int `json:"width"`
// 	} `json:"regions_of_interest"`
// 	CreditCost int `json:"credit_cost"`
// }

type ResponseLPR struct {
	UUID           string `json:"uuid"`
	DataType       string `json:"data_type"`
	EpochTime      int64  `json:"epoch_time"`
	ProcessingTime struct {
		Total    float64 `json:"total"`
		Plates   float64 `json:"plates"`
		Vehicles float64 `json:"vehicles"`
	} `json:"processing_time"`
	ImgHeight int `json:"img_height"`
	ImgWidth  int `json:"img_width"`
	Results   []struct {
		Plate            string  `json:"plate"`
		Confidence       float64 `json:"confidence"`
		RegionConfidence int     `json:"region_confidence"`
		VehicleRegion    struct {
			Y      int `json:"y"`
			X      int `json:"x"`
			Height int `json:"height"`
			Width  int `json:"width"`
		} `json:"vehicle_region"`
		Region           string  `json:"region"`
		PlateIndex       int     `json:"plate_index"`
		ProcessingTimeMs float64 `json:"processing_time_ms"`
		Candidates       []struct {
			MatchesTemplate int     `json:"matches_template"`
			Plate           string  `json:"plate"`
			Confidence      float64 `json:"confidence"`
		} `json:"candidates"`
		Coordinates []struct {
			Y int `json:"y"`
			X int `json:"x"`
		} `json:"coordinates"`
		Vehicle struct {
			Orientation []struct {
				Confidence float64 `json:"confidence"`
				Name       string  `json:"name"`
			} `json:"orientation"`
			Color []struct {
				Confidence float64 `json:"confidence"`
				Name       string  `json:"name"`
			} `json:"color"`
			Make []struct {
				Confidence float64 `json:"confidence"`
				Name       string  `json:"name"`
			} `json:"make"`
			BodyType []struct {
				Confidence float64 `json:"confidence"`
				Name       string  `json:"name"`
			} `json:"body_type"`
			Year []struct {
				Confidence float64 `json:"confidence"`
				Name       string  `json:"name"`
			} `json:"year"`
			MakeModel []struct {
				Confidence float64 `json:"confidence"`
				Name       string  `json:"name"`
			} `json:"make_model"`
		} `json:"vehicle"`
		MatchesTemplate int `json:"matches_template"`
		RequestedTopn   int `json:"requested_topn"`
	} `json:"results"`
	CreditsMonthlyUsed  int  `json:"credits_monthly_used"`
	Version             int  `json:"version"`
	CreditsMonthlyTotal int  `json:"credits_monthly_total"`
	Error               bool `json:"error"`
	RegionsOfInterest   []struct {
		Y      int `json:"y"`
		X      int `json:"x"`
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"regions_of_interest"`
	CreditCost int `json:"credit_cost"`
}

func TestOpenJSON(t *testing.T) {
	_, err := os.Open("alpr_request.json")

	assert.NoError(t, err)
}

func TestBindJSON(t *testing.T) {
	jsonFile, _ := os.Open("alpr_request.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var req Request

	json.Unmarshal(byteValue, &req)

	assert.Equal(t, "https://camargus.com/magazine/data/2/174/3650.jpeg", req.LPR.ImageURL)

}

func TestPostToALPR(t *testing.T) {
	values := map[string]interface{}{
		"image_url":         "https://camargus.com/magazine/data/2/174/3650.jpeg",
		"secret_key":        "sk_61ef6b64a64e70f6225075cb",
		"recognize_vehicle": 1,
		"country":           "id",
		"state":             "id",
		"return_image":      0,
		"topn":              10,
	}
	// https://api.openalpr.com/v2/recognize_url?image_url=https%3A%2F%2Fcamargus.com%2Fmagazine%2Fdata%2F2%2F174%2F3650.jpeg&secret_key=sk_61ef6b64a64e70f6225075cb&recognize_vehicle=1&country=id&state=id&return_image=0&topn=10
	urlAndParams := fmt.Sprintf("https://api.openalpr.com/v2/recognize_url?image_url=%s&secret_key=%s&recognize_vehicle=%d&country=%s&state=%s&return_image=%d&topn=%d",
		values["image_url"],
		values["secret_key"],
		values["recognize_vehicle"],
		values["country"],
		values["state"],
		values["return_image"],
		values["topn"],
	)

	req, err := http.NewRequest("POST", urlAndParams, nil)

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()
	responseData, _ := ioutil.ReadAll(resp.Body)
	// assert.Equal(t, "", string(responseData))

	var responseLPR *ResponseLPR

	json.Unmarshal(responseData, &responseLPR)
	j, _ := json.Marshal(responseLPR)

	assert.NotNil(t, responseLPR)
	assert.Equal(t, "", string(j))

}
