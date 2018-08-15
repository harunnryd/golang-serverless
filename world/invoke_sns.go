package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type InvokeSNS struct {
	ID       string
	Secret   string
	Token    string
	Region   string
	TopicArn string
}

var instanceInvoiceSNS *InvokeSNS

/**
 * -------------------------------
 * NOTE : initialize invoice sns
 *
 * @return *InvokeSNS
 * -------------------------------
 */
func NewInvoiceSNS() *InvokeSNS {
	if instanceInvoiceSNS == nil {
		instanceInvoiceSNS = new(InvokeSNS)
	}
	return instanceInvoiceSNS
}

/**
 * -----------------------------
 * NOTE : initialize credential
 *
 * @param id string
 * @param secret string
 * @param token string
 * @param region string
 * @param tipicArn string
 *
 * @return *InvokeSNS
 * -----------------------------
 */
func (i *InvokeSNS) initialize(id, secret, token, region, topicArn string) *InvokeSNS {
	i.ID = id
	i.Secret = secret
	i.Token = token
	i.Region = region
	i.TopicArn = topicArn
	return i
}

/**
 * -------------------------------
 * NOTE : publish message ke sns
 * -------------------------------
 */
func (i *InvokeSNS) Publish(message string) {
	cred := credentials.NewStaticCredentials(
		i.ID,
		i.Secret,
		i.Token,
	)

	config := aws.NewConfig().WithRegion(i.Region).WithCredentials(cred)

	sess := session.New(config)

	svc := sns.New(sess)

	params := &sns.PublishInput{
		/**
		 * --------------------------------------------------------------------
		 * NOTE : message bisa diisi dengan (XML / JSON / Text - atau apapun)
		 * --------------------------------------------------------------------
		 */
		Message: aws.String(message),

		/**
		 * -----------------
		 * NOTE : topic arn
		 * -----------------
		 */
		TopicArn: aws.String(i.TopicArn),
	}

	resp, err := svc.Publish(params)

	if err != nil {
		printLog("ERROR", fmt.Sprintf("PUBLISH SNS TOPIC %v", err))
		panic(err)
	}

	/**
	 * ------------------------------------------------
	 * NOTE : cetak log ketika sukses mem publish sns
	 * ------------------------------------------------
	 */
	printLog("INFO", resp)

}

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
