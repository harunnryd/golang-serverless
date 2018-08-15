package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

/**
 * ------------------------
 * NOTE : request body #1
 * ------------------------
 */
type RequestLPR struct {
	LPR             LPR             `json:"lpr"`
	IncomingWebHook IncomingWebHook `json:"incoming_webhook"`
}

/**
 * ------------------------
 * NOTE : request body #2
 * ------------------------
 */
type IncomingWebHook struct {
	EndPoint      string `json:"end_point"`
	Authorization string `json:"authorization"`
	Method        string `json:"method"`
}

/**
 * ------------------------
 * NOTE : request body #3
 * ------------------------
 */
type LPR struct {
	ImageURL  string `json:"image_url"`
	SecretKey string `json:"secret_key"`
	Type      string `json:"type"`
}

/**
 * ----------------------------
 * NOTE : handler lambda
 *
 * @param request RequestLPR
 * @return error
 * ----------------------------
 */
func handler(request RequestLPR) error {

	/**
	 * --------------------------------------------
	 * NOTE : binding request body kedalam values
	 * --------------------------------------------
	 */
	values := map[string]interface{}{
		"image_url":         request.LPR.ImageURL,
		"secret_key":        request.LPR.SecretKey,
		"recognize_vehicle": 1,
		"country":           "id",
		"state":             "id",
		"return_image":      0,
		"topn":              10,
	}

	/**
	 * -------------------------------------------
	 * NOTE : url dan parameter yang dibutuhkan
	 * -------------------------------------------
	 */
	urlAndParams := fmt.Sprintf("https://api.openalpr.com/v2/recognize_url?image_url=%s&secret_key=%s&recognize_vehicle=%d&country=%s&state=%s&return_image=%d&topn=%d",
		values["image_url"],
		values["secret_key"],
		values["recognize_vehicle"],
		values["country"],
		values["state"],
		values["return_image"],
		values["topn"],
	)

	/**
	 * --------------------------------------------------
	 * NOTE : set method, header, url dan parameternya
	 * --------------------------------------------------
	 */
	req, _ := http.NewRequest("POST", urlAndParams, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {

		/**
		 * -------------------
		 * NOTE : cetak error
		 * -------------------
		 */
		log := fmt.Sprintf("[%v ERROR] %v", time.Now().Format("2006-01-02 15:04:05"), err)
		fmt.Println(log)

		return err
	}

	defer resp.Body.Close()

	return nil
}

func main() {
	lambda.Start(handler)
}
