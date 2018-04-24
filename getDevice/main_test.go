package main_test

import (
    "testing"
    "errors"
    "encoding/json"

    "github.com/stretchr/testify/assert"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

    d "github.com/SamiraKhodami/serverlessGoDynamoDB/lib/devicecrud"
    main "github.com/SamiraKhodami/serverlessGoDynamoDB/getDevice"
)

type mockDynamoDBClient struct {
    dynamodbiface.DynamoDBAPI
}

type InputData struct {
    ID   string    `json:"id"`
}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
    inputData := InputData{} 
    _ = dynamodbattribute.UnmarshalMap(input.Key, &inputData);

    if inputData.ID == "id5" {
        result := dynamodb.GetItemOutput{
            Item: map[string]*dynamodb.AttributeValue{
                "ID" : {
                    S: aws.String("id5"),
                },
                "DeviceModel" : {
                    S: aws.String("/devicemodels/id5"),
                },
                "Name" : {
                    S: aws.String("Sensor"),
                },
                "Note" : {
                    S: aws.String("Testing a sensor."),
                },
                "Serial" : {
                    S: aws.String("A020000105"),
                },
            },
        }
        return &result,nil
    } 

    result := dynamodb.GetItemOutput{
        Item: map[string]*dynamodb.AttributeValue{},
    }

    if inputData.ID == "corrupt" {
        return &result,errors.New("corrupt device")
    }

    return &result,nil
}


func TestGetDeviceFromDBAndReturnResponse(t *testing.T) {

    mockSvc := &mockDynamoDBClient{}

    device := d.Device{
        ID:"id5",
        DeviceModel: "/devicemodels/id5",
        Name: "Sensor",
        Note: "Testing a sensor.",
        Serial: "A020000105",
    }
    body, _ := json.Marshal(device)

    // empty request id will never happen
    tests := []struct {
        id      string 
        expect  string
        code    int
    }{
        {
            id: "10",
            expect: "Not Found", 
            code: 404,
        },
        {
            id: "id5",
            expect: string(body),
            code: 200,
        },
        {
            id: "corrupt",
            expect: "Internal Server Error",
            code: 500,
        },
    }

    for _, test := range tests {
        response, _ := main.GetDeviceFromDBAndReturnResponse(mockSvc,test.id)
        assert.Equal(t, test.expect, response.Body)
        assert.Equal(t, test.code, response.StatusCode)
    }
}
