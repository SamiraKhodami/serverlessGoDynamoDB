package devicecrud

import (
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Device struct {
    ID          string    `json:"id"`
    DeviceModel string    `json:"devicemodel"`
    Name        string    `json:"name"`
    Note        string    `json:"note"`
    Serial      string    `json:"serial"`
}

var (
    tableName = aws.String(os.Getenv("DEVICE_TABLE_NAME"))
)

func AddDeviceToDB(svc dynamodbiface.DynamoDBAPI, device *Device) (*dynamodb.PutItemOutput, error) {

    fmt.Println("AddDeviceToDB")
    item, _ := dynamodbattribute.MarshalMap(device)
    input := &dynamodb.PutItemInput{
        Item: item,
        TableName: tableName,
    }

    return svc.PutItem(input)
}

func GetDeviceFromDB(svc dynamodbiface.DynamoDBAPI, id string) (*dynamodb.GetItemOutput, error){

    fmt.Println("GetDeviceFromDB")
    input := &dynamodb.GetItemInput{
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                S: aws.String(id),
            },
        },
        TableName: tableName,
    }

    return svc.GetItem(input)
}

