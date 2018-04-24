package main_test

import (
    "testing"
    "errors"
    "encoding/json"

    "github.com/stretchr/testify/assert"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

    d "github.com/SamiraKhodami/serverlessGoDynamoDB/lib/devicecrud"
    main "github.com/SamiraKhodami/serverlessGoDynamoDB/addDevice"
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

func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
    device := d.Device{}
    _ = dynamodbattribute.UnmarshalMap(input.Item, &device);

    result := dynamodb.PutItemOutput{}

    if device.ID == "corruptPut" {
        return &result,errors.New("corrupt device")
    }

    return &result, nil 
}

func TestAddDeviceToDBAndReturnResponse(t *testing.T) {

    mockSvc := &mockDynamoDBClient{}

    newDevice := d.Device{
        ID:"id8",
        DeviceModel: "/devicemodels/id8",
        Name: "Sensor",
        Note: "Testing a sensor.",
        Serial: "A020000108", 
    }
    availableDevice := d.Device{
        ID:"id5",
        DeviceModel: "/devicemodels/id5",
        Name: "Sensor",
        Note: "Testing a sensor.",
        Serial: "A020000105",
    }
    corruptDevice := d.Device{
        ID:"corrupt",
        DeviceModel: "-",
        Name: "-",
        Note: "-",
        Serial: "-",
    }
    corruptPutDevice := d.Device{
        ID:"corruptPut",
        DeviceModel: "-",
        Name: "-",
        Note: "-",
        Serial: "-",
    }


    tests := []struct {
        device  d.Device 
        code    int
    }{
        {
            device: newDevice,
            code: 201,
        },
        {
            device: availableDevice,
            code: 400,
        },
        {
            device: corruptDevice,
            code: 500,
        },
        {
            device: corruptPutDevice,
            code: 500,
        },
    }

    for _, test := range tests {
        response, _ := main.AddDeviceToDBAndReturnResponse(mockSvc,&test.device)
        assert.Equal(t, test.code, response.StatusCode)
    }
}

func TestValidateRequest(t *testing.T) {
    emptyIdDevice := d.Device{
        ID:"",
        DeviceModel: "/devicemodels/id8",
        Name: "Sensor",
        Note: "Testing a sensor.",
        Serial: "A020000108",
    }
    emptyIdDeviceBody, _ := json.Marshal(emptyIdDevice)

    emptyModelDevice := d.Device{
        ID:"id8",
        DeviceModel: "",
        Name: "Sensor",
        Note: "Testing a sensor.",
        Serial: "A020000108",
    }
    emptyModelDeviceBody, _ := json.Marshal(emptyModelDevice)

    emptyNameDevice := d.Device{
        ID:"id8",
        DeviceModel: "/devicemodels/id8",
        Name: "     ",
        Note: "Testing a sensor.",
        Serial: "A020000108",
    }
    emptyNameDeviceBody, _ := json.Marshal(emptyNameDevice)

    emptyNoteDevice := d.Device{
        ID:"id8",
        DeviceModel: "/devicemodels/id8",
        Name: "Sensor",
        Note: "  ",
        Serial: "A020000108",
    }
    emptyNoteDeviceBody, _ := json.Marshal(emptyNoteDevice)

    emptySerialDevice := d.Device{
        ID:"id8",
        DeviceModel: "/devicemodels/id8",
        Name: "Sensor",
        Note: "Testing a sensor.",
        Serial: "",
    }
    emptySerialDeviceBody, _ := json.Marshal(emptySerialDevice)

    tests := []struct {
        request  events.APIGatewayProxyRequest
        err      error
    }{
        {
            request: events.APIGatewayProxyRequest{Body: string(emptyIdDeviceBody)}, 
            err: errors.New("Bad Request: Id Can Not Be Empty"),
        },
        {
            request: events.APIGatewayProxyRequest{Body: string(emptyModelDeviceBody)},
            err: errors.New("Bad Request: DeviceModel Can Not Be Empty"),
        },
        {
            request: events.APIGatewayProxyRequest{Body: string(emptyNameDeviceBody)},
            err: errors.New("Bad Request: Name Can Not Be Empty"),
        },
        {
            request: events.APIGatewayProxyRequest{Body: string(emptyNoteDeviceBody)},
            err: errors.New("Bad Request: Note Can Not Be Empty"),
        },
        {
            request: events.APIGatewayProxyRequest{Body: string(emptySerialDeviceBody)},
            err: errors.New("Bad Request: Serial Can Not Be Empty"),
        },

    }

    for _, test := range tests {
        _, err := main.ValidateRequest(test.request)
        assert.Equal(t, test.err.Error(), err.Error())
    }

}
