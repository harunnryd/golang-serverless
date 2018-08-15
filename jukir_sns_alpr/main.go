package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
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
 * dan response body #3
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
 * ------------------------
 * NOTE : response body #1
 * ------------------------
 */
type Response struct {
	LPR             *ResponseFromLPR `json:"lpr"`
	IncomingWebHook IncomingWebHook  `json:"incoming_webhook"`
}

/**
 * ------------------------
 * NOTE : response body #2
 * ------------------------
 */
type ResponseFromLPR struct {
	UUID           string `json:"uuid"`
	DataType       string `json:"data_type"`
	EpochTime      int64  `json:"epoch_time"`
	ProcessingTime struct {
		Plates float64 `json:"plates"`
		Total  float64 `json:"total"`
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

/**
 * --------------------------------
 * NOTE : handler
 * --------------------------------
 */
func handler(ctx context.Context, snsEvent events.SNSEvent) {

	/**
	 * -----------------------------
	 * NOTE : cetak log proses SNS
	 * -----------------------------
	 */
	fmt.Printf("[%v INFO] Process MessageID: %s\n", time.Now().Format("2006-01-02 15:04:05"), snsEvent.Records[0].SNS.MessageID)
	fmt.Printf("[%v INFO] Body: %d\n", time.Now().Format("2006-01-02 15:04:05"), len(snsEvent.Records[0].SNS.Message))

	requestLPR := snsEvent.Records[0].SNS.Message

	var formRequestLPR RequestLPR

	/**
	 * --------------------------------------------------------
	 * NOTE : binding requestLPR kedalam form request (struct)
	 * --------------------------------------------------------
	 */
	err := json.Unmarshal([]byte(requestLPR), &formRequestLPR)

	if err != nil {

		/**
		 * ------------------------
		 * NOTE : cetak log gagal
		 * ------------------------
		 */
		log := fmt.Sprintf("[%v ERROR] %v", time.Now().Format("2006-01-02 15:04:05"), err)
		fmt.Println(log)

		panic(err)
	}

	err = postToOpenLPR(getURLAndParams(formRequestLPR), formRequestLPR.IncomingWebHook)

	if err != nil {

		/**
		 * ------------------------
		 * NOTE : cetak log gagal
		 * ------------------------
		 */
		log := fmt.Sprintf("[%v handler ERROR] %v", time.Now().Format("2006-01-02 15:04:05"), err)
		fmt.Println(log)

		panic(err)
	}

}

/**
 * ----------------------------------------------
 * NOTE : ekstrak value pada form request lpr
 * ke dalam sebuah url parameter openlpr
 *
 * @param req RequestLPR
 *
 * @return string
 * ----------------------------------------------
 */
func getURLAndParams(req RequestLPR) string {

	fmt.Println("IMAGEURL = ", req.LPR.ImageURL)
	fmt.Println("SECRET_KEY = ", req.LPR.SecretKey)

	/**
	 * --------------------------------------------
	 * NOTE : binding request body kedalam values
	 * --------------------------------------------
	 */
	values := map[string]interface{}{
		"image_url":         req.LPR.ImageURL,
		"secret_key":        req.LPR.SecretKey,
		"country":           "id",
		"recognize_vehicle": 1,
		"state":             "id",
		"return_image":      0,
		"topn":              10,
	}

	/**
	 * -------------------------------------------
	 * NOTE : url dan parameter yang dibutuhkan
	 * -------------------------------------------
	 */
	urlAndParams := fmt.Sprintf("https://api.openalpr.com/v2/recognize_url?image_url=%s&secret_key=%s&country=%s&recognize_vehicle=%dstate=%s&return_image=%d&topn=%d",
		values["image_url"],
		values["secret_key"],
		values["country"],
		values["recognize_vehicle"],
		values["state"],
		values["return_image"],
		values["topn"],
	)

	return urlAndParams
}

func init() {

	/**
	 * --------------------------------
	 * NOTE : inisialisasi credential
	 * milik account aws
	 * --------------------------------
	 */
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAJP6D7T2LLRWIJ6NQ")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_REGION", "ap-southeast-1")
}

/**
 * --------------------------------------------------------
 * NOTE : hit end point open lpr untuk memproses gambar
 * plat nomor kendaraan menjadi berbagai macam informasi
 *
 * @param urlAndparams string
 *
 * @return ResponseLPR
 * @return error
 * --------------------------------------------------------
 */
func postToOpenLPR(urlAndParams string, incomingWebHook IncomingWebHook) error {
	/**
	 * --------------------------------------------------
	 * NOTE : set method, header, url dan parameternya
	 * --------------------------------------------------
	 */
	fmt.Println("URL & PARAMS OPEN LRP = ", urlAndParams)

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
		log := fmt.Sprintf("[%v postToOpenLPR ERROR] %v", time.Now().Format("2006-01-02 15:04:05"), err)
		fmt.Println(log)

		return err
	}

	/**
	 * ----------------------------
	 * NOTE : tutup response body
	 * ----------------------------
	 */
	defer resp.Body.Close()

	responseData, _ := ioutil.ReadAll(resp.Body)
	var responseLPR *ResponseFromLPR

	json.Unmarshal(responseData, &responseLPR)

	err = invokeSNS(responseLPR, incomingWebHook)

	fmt.Println(responseLPR)

	if err != nil {

		/**
		 * -------------------
		 * NOTE : cetak error
		 * -------------------
		 */
		log := fmt.Sprintf("[%v InvokeSNS ERROR] %v", time.Now().Format("2006-01-02 15:04:05"), err)
		fmt.Println(log)

		return err
	}

	return nil
}

/**
 * -----------------------------------------------------------------------------------------------
 * NOTE : invoke sns topic arn:aws:sns:ap-southeast-1:929307424935:topic-record-response-alpr-car
 *
 * @param responseLPR ResponseFormLPR
 * @param incomingWebHook IncomingWebHook
 *
 * @return error
 * -----------------------------------------------------------------------------------------------
 */
func invokeSNS(responseLPR *ResponseFromLPR, incomingWebHook IncomingWebHook) error {

	/**
	 * ----------------------------
	 * NOTE : set credentials aws
	 * ----------------------------
	 */
	cred := credentials.NewStaticCredentials("AKIAJP6D7T2LLRWIJ6NQ", "ktxchSBUjsofw/S4JjXxCD0DAkP1dQVf7numuT3l", "")
	config := aws.NewConfig().WithRegion("ap-southeast-1").WithCredentials(cred)
	sess := session.New(config)

	svc := sns.New(sess)

	response := Response{
		LPR:             responseLPR,
		IncomingWebHook: incomingWebHook,
	}

	message, err := json.Marshal(response)

	if err != nil {

		/**
		 * -------------------
		 * NOTE : cetak error
		 * -------------------
		 */
		log := fmt.Sprintf("[%v MARSHALL INVOKE SNS ERROR] %v", time.Now().Format("2006-01-02 15:04:05"), err)
		fmt.Println(log)

		return err
	}

	/**
	 * ----------------------------------------
	 * NOTE : bind params yang akan digunakan
	 * untuk di include pada sns topic
	 * ----------------------------------------
	 */
	params := &sns.PublishInput{

		/**
		 * --------------------------------------------------------------------
		 * NOTE : message bisa diisi dengan (XML / JSON / Text - atau apapun)
		 * --------------------------------------------------------------------
		 */
		Message: aws.String(string(message)),

		/**
		 * -----------------
		 * NOTE : topic arn
		 * -----------------
		 */
		TopicArn: aws.String("arn:aws:sns:ap-southeast-1:929307424935:topic-record-response-alpr-car"),
	}

	/**
	 * -----------------------
	 * NOTE : publish message
	 * -----------------------
	 */
	resp, err := svc.Publish(params)

	if err != nil {

		/**
		 * -------------------
		 * NOTE : cetak error
		 * -------------------
		 */
		log := fmt.Sprintf("[%v PUBLISH TOPIC ERROR] %v", time.Now().Format("2006-01-02 15:04:05"), err)
		fmt.Println(log)

		return err
	}

	/**
	 * ---------------------------------
	 * NOTE : cetak log pretty response
	 * ---------------------------------
	 */
	log := fmt.Sprintf("[%v INFO] %s", time.Now().Format("2006-01-02 15:04:05"), resp)
	fmt.Println(log)

	return nil
}

func main() {
	lambda.Start(handler)
}
