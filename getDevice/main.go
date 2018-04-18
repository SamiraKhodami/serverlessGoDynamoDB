package main

import (
    "fmt"
    "encoding/json"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

    "myservice/lib/dbconfig"
    d "myservice/lib/devicecrud"
)

func GetDevice(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    var id = request.PathParameters["id"]

    svc, err := dbconfig.ConfigureDynamoDB(); 
    if err != nil {
        return events.APIGatewayProxyResponse{
            Body: "Internal Server Error",
            StatusCode: 500,
        }, nil
    }

    return GetDeviceFromDBAndReturnResponse(svc,id)
}

func GetDeviceFromDBAndReturnResponse(svc dynamodbiface.DynamoDBAPI, id string) (events.APIGatewayProxyResponse, error) {

    result, err := d.GetDeviceFromDB(svc,id) 

    if err != nil {
        fmt.Println(fmt.Sprintf("Failed to get item: %s", err.Error()))
        return events.APIGatewayProxyResponse{
            Body: "Internal Server Error",
            StatusCode: 500,
        }, nil
    }

    if len(result.Item) <= 0  {
        return events.APIGatewayProxyResponse{
            Body: "Not Found",
            StatusCode: 404,
        }, nil
    }

    device := d.Device{}
    if err := dynamodbattribute.UnmarshalMap(result.Item, &device); err != nil {
        fmt.Println(fmt.Sprintf("Failed to unmarshal item: %s", err.Error()))
        return events.APIGatewayProxyResponse{
            Body: "Internal Server Error",
            StatusCode: 500,
        }, nil
    }

    body, _ := json.Marshal(device);
    return events.APIGatewayProxyResponse{ // Success HTTP response
        Body: string(body),
        StatusCode: 200,
    }, nil
}

func main() {
	lambda.Start(GetDevice)
}
