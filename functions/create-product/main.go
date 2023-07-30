package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var tableName string

func init() {
	_ = godotenv.Load()
	name, ok := os.LookupEnv("TABLE_NAME")
	if !ok {
		panic("Need TABLE environment variable")
	}
	tableName = name
}

func CreateProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var product Product

	err := json.Unmarshal([]byte(request.Body), &product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	if strings.TrimSpace(request.Body) == "" {
		return events.APIGatewayProxyResponse{Body: "empty request body", StatusCode: 400}, nil
	}

	product.ID = uuid.New().String()

	item, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := json.Marshal(product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(CreateProduct)
}
