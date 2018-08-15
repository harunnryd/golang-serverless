package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

/**
 * -----------------------------------------------
 * NOTE : handler APIGatewayProxy
 *
 * @param ctx context.Context
 * @param request events.APIGatewayProxyRequest
 *
 * @return events.APIGatewayProxyResponse
 * @return error
 * -----------------------------------------------
 */
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	/**
	 * ------------------------------------------------
	 * NOTE : setiap hal yang dicetak ke console maka
	 * akan otomatis tercetak di CloudWatch milik AWS
	 * ------------------------------------------------
	 */
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	return events.APIGatewayProxyResponse{
		Body:       request.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
