package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

/**
 * ---------------------
 * NOTE : request body
 * ---------------------
 */
type Request struct {
	IncomingWebhook `json:"incoming_webhook"`
	ALPR            `json:"alpr"`
}

type IncomingWebhook struct {
	Authorization string `json:"authorization"`
	EndPoint      string `json:"end_point"`
	Method        string `json:"method"`
}

type ALPR struct {
	ImageURL  string `json:"image_url"`
	SecretKey string `json:"secret_key"`
	Type      string `json:"type"`
}

func main() {

	req := Request{
		IncomingWebhook: IncomingWebhook{
			Authorization: "https://jukir.co/testing-alpr/create",
			EndPoint:      "JUKIR AIUEO-!@#$%-SKSKS-FVCK-YVO!",
			Method:        "POST",
		},
		ALPR: ALPR{
			ImageURL:  "https://camargus.com/magazine/data/2/174/3650.jpeg",
			SecretKey: "sk_757cf6785a0472290102148f",
			Type:      "car",
		},
	}

	byteValue, _ := json.Marshal(req)

	//Create a session object to talk to SNS (also make sure you have your key and secret setup in your .aws/credentials file)
	svc := sns.New(session.New())
	// params will be sent to the publish call included here is the bare minimum params to send a message.
	params := &sns.PublishInput{
		Message:  aws.String(string(byteValue)),                                                   // This is the message itself (can be XML / JSON / Text - anything you want)
		TopicArn: aws.String("arn:aws:sns:ap-southeast-1:929307424935:topic-processing-alpr-car"), //Get this from the Topic in the AWS console.
	}

	resp, err := svc.Publish(params) //Call to puclish the message

	if err != nil { //Check for errors
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
