package dbconfig

import (
    "fmt"
    "os"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func ConfigureDynamoDB() (dynamodbiface.DynamoDBAPI, error) {
    var svc dynamodbiface.DynamoDBAPI 
    region := os.Getenv("AWS_REGION")
    session, err := session.NewSession(&aws.Config{ 
        Region: &region,
    });
    if err != nil {
        fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
    } else {
        svc = dynamodbiface.DynamoDBAPI(dynamodb.New(session))
    }
    return svc, err
}
