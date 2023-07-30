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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/joho/godotenv"
	"os"
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

func GetProducts(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	Idexists := expression.Name("Id").AttributeNotExists()
	proj := expression.NamesList(expression.Name("Id"), expression.Name("Name"), expression.Name("Price"))

	expr, err := expression.NewBuilder().WithCondition(Idexists).WithProjection(proj).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		Limit:                     aws.Int64(20),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	var products []Product
	for _, item := range result.Items {
		var product Product
		err = dynamodbattribute.UnmarshalMap(item, &product)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
		products = append(products, product)
		/*price, err := strconv.Atoi(*i["price"].N)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
		}
		product := Product{
			ID:    *i["id"].S,
			Name:  *i["name"].S,
			Price: price,
		}
		products = append(products, product)*/
	}

	body, err := json.Marshal(products)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(GetProducts)
}
