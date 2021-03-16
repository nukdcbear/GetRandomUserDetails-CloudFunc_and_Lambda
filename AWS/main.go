package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	UTC time.Time `json:"utc"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	randomUserClient := http.Client{
		Timeout: time.Second * 3,
	}

	req, err := http.NewRequest(http.MethodGet, "https://randomuser.me/api/", nil)
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, err
	}

	res, err2 := randomUserClient.Do(req)
	if err2 != nil {
		log.Fatal(err2)
		return events.APIGatewayProxyResponse{}, err
	}

	resBody, err3 := ioutil.ReadAll(res.Body)
	if err3 != nil {
		log.Fatal(err3)
	}

	var o map[string]interface{}
	json.Unmarshal([]byte(resBody), &o)

	results := o["results"].([]interface{})
	result := results[0].(map[string]interface{})

	result["generator"] = "aws-lambda-function"

	body, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(body), StatusCode: 200}, nil
}
func main() {
	lambda.Start(handleRequest)
}
