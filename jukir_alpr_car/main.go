package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	sls "github.com/harunnryd/serverless"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// Request : is
	Request sls.Request
	// Response : is
	Response sls.Response
	// ALPR : is
	ALPR *sls.ALPR
	// IncomingWebhook : is
	IncomingWebhook sls.IncomingWebhook
)

/**
 * ----------------------------------------
 * NOTE : cetak log pada CloudWatch AWS
 *
 * @param status string
 * @param message interface{}
 * ----------------------------------------
 */
func printLog(status string, message interface{}) {
	fmt.Printf("[%v %s] %v\n", time.Now().Format("2006-01-02 15:04:05"), status, message)
}

/**
 * -------------------------------------
 * NOTE : handler serverless
 *
 * @param ctx context.context
 * @param snsEvent events.SNSEvent
 * ------------------------------------
 */
func handler(ctx context.Context, snsEvent events.SNSEvent) {

	/**
	 * ------------------------------------------
	 * NOTE : cetak log pada saat memulai proses
	 * ------------------------------------------
	 */
	printLog("INFO", fmt.Sprintf("PROCESS MESSAGE ID = %s", snsEvent.Records[0].SNS.MessageID))
	printLog("INFO", fmt.Sprintf("BODY = %d", len(snsEvent.Records[0].SNS.Message)))

	snsMessage := snsEvent.Records[0].SNS.Message

	/**
	 * -----------------------------------------------
	 * NOTE : decode string json, menjadi json object
	 * -----------------------------------------------
	 */
	err := json.Unmarshal([]byte(snsMessage), &Request)

	/**
	 * -------------------------------------
	 * NOTE : cetak log saat terjadi error
	 * -------------------------------------
	 */
	if err != nil {
		printLog("ERROR", fmt.Sprintf("JSON UNMARSHAL = %v", err))
		panic(err)
	}

	// ----------------------------------------------------------------------------------------------------------

	/**
	 * ------------------------------------------------------
	 * NOTE : extract plat nomor kendaraan menjadi informasi
	 * ------------------------------------------------------
	 */

	imageURL := Request.ALPR.ImageURL
	secretKey := Request.ALPR.SecretKey

	printLog("INFO", fmt.Sprintf("IMAGE URL = %s, SECRET KEY = %s", imageURL, secretKey))

	payload := generatePayload(imageURL, secretKey)

	if err := extractPlate(payload); err != nil {
		printLog("ERROR", fmt.Sprintf("EXTRACT PLATE %v", err))
		panic(err)
	}

	// ----------------------------------------------------------------------------------------------------------

	/**
	 * -------------------------------
	 * NOTE : invoke sns lain
	 * -------------------------------
	 */
}

/**
 * -----------------------------------------
 * NOTE : generate payload openalpr
 *
 * @param imageURL string
 * @param secretKey string
 *
 * @return string
 * -----------------------------------------
 */
func generatePayload(imageURL, secretKey string) string {

	param := map[string]interface{}{
		"image_url":         imageURL,
		"secret_key":        secretKey,
		"country":           "id",
		"recognize_vehicle": 1,
		"state":             "id",
		"return_image":      0,
		"topn":              10,
	}

	payload := fmt.Sprintf("https://api.openalpr.com/v2/recognize_url?image_url=%s&secret_key=%s&country=%s&recognize_vehicle=%dstate=%s&return_image=%d&topn=%d",
		param["image_url"],
		param["secret_key"],
		param["country"],
		param["recognize_vehicle"],
		param["state"],
		param["return_image"],
		param["topn"],
	)

	return payload
}

/**
 * ----------------------------------------------------------
 * NOTE : extract plat nomor kendaraan menggunakan 3rd party
 * openalpr
 *
 * @param payload string
 *
 * @return error
 * ----------------------------------------------------------
 */
func extractPlate(payload string) error {

	/**
	 * --------------------------------------------
	 * NOTE : inisialisasi client untuk open alpr
	 * --------------------------------------------
	 */
	req, _ := http.NewRequest("POST", payload, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	/**
	 * -------------------------------------
	 * NOTE : cetak log saat terjadi error
	 * -------------------------------------
	 */
	if err != nil {
		printLog("ERROR", fmt.Sprintf("PROCESS EXTRACT PLATE = %v", err))
		return err
	}

	defer resp.Body.Close()

	printLog("INFO", fmt.Sprintf("BODY = %s", resp.Body))

	/**
	 * -----------------------------------------
	 * NOTE : ubah response body menjadi []byte
	 * -----------------------------------------
	 */
	respData, err := ioutil.ReadAll(resp.Body)

	/**
	 * ------------------------------------
	 * NOTE : cetak log saat terjadi error
	 * ------------------------------------
	 */
	if err != nil {
		printLog("ERROR", fmt.Sprintf("PROCESS IOUTIL READ ALL RESPONSE BODY = %v", err))
		return err
	}

	/**
	 * -----------------------------------------
	 * NOTE : decode json kedalam struct ALPR
	 * -----------------------------------------
	 */
	if err := json.Unmarshal(respData, &ALPR); err != nil {
		printLog("ERROR", fmt.Sprintf("UNMARSHAL ALPR = %v", err))
		return err
	}

	/**
	 * ------------------------------
	 * NOTE : pretty print JSON ALPR
	 * ------------------------------
	 */
	fmt.Println(string(respData))
	fmt.Println(ALPR)

	return nil

}

/**
 * -----------------------
 * NOTE : jalankan lambda
 * -----------------------
 */
func main() {
	lambda.Start(handler)
}
