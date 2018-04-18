package main

import (
    "fmt"
    "encoding/json"
    "errors"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

    "myservice/lib/stringutil"
    "myservice/lib/dbconfig"
    d "myservice/lib/devicecrud"
)

func AddDevice(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    device, err := ValidateRequest(request)
    if err != nil {
        return events.APIGatewayProxyResponse{ 
            Body: err.Error(),
            StatusCode: 400,
        }, nil
    }

    svc, err := dbconfig.ConfigureDynamoDB(); 
    if err != nil {
        return events.APIGatewayProxyResponse{
            Body: "Internal Server Error",
            StatusCode: 500,
        }, nil
    }

    return AddDeviceToDBAndReturnResponse(svc, device)
}

func ValidateRequest(request events.APIGatewayProxyRequest)(*d.Device, error){

    var device d.Device
    json.Unmarshal([]byte(request.Body),&device)

    if stringutil.IsEmpty(device.ID) {
        return &device, errors.New("Bad Request: Id Can Not Be Empty")
    }

    if stringutil.IsEmpty(device.DeviceModel) {
        return &device, errors.New("Bad Request: DeviceModel Can Not Be Empty")
    }

    if stringutil.IsEmpty(device.Name) {
        return &device, errors.New("Bad Request: Name Can Not Be Empty")
    }

    if stringutil.IsEmpty(device.Note) {
        return &device, errors.New("Bad Request: Note Can Not Be Empty")
    }

    if stringutil.IsEmpty(device.Serial) {
        return &device, errors.New("Bad Request: Serial Can Not Be Empty")
    }
    return &device, nil 
}

func AddDeviceToDBAndReturnResponse(svc dynamodbiface.DynamoDBAPI, device *d.Device) (events.APIGatewayProxyResponse, error) {

    result, err := d.GetDeviceFromDB(svc,device.ID)

    if err != nil {
        fmt.Println(fmt.Sprintf("Failed to get device: %s", err.Error()))
        return events.APIGatewayProxyResponse{
            Body: "Internal Server Error",
            StatusCode: 500,
        }, nil
    }

    if len(result.Item) > 0  {
        return events.APIGatewayProxyResponse{
            Body: "Bad Request: Device with this id already exist",
            StatusCode: 400,
        }, nil
    }

    if _, err := d.AddDeviceToDB(svc, device); err != nil {
        fmt.Println(fmt.Sprintf("Failed to put item: %s", err.Error()))
        return events.APIGatewayProxyResponse{
            Body: "Internal Server Error",
            StatusCode: 500,
        }, nil
    }

    body, _ := json.Marshal(device)
    return events.APIGatewayProxyResponse{ // Success HTTP response
        Body: string(body),
        StatusCode: 201,
    }, nil
}

func main() {
    lambda.Start(AddDevice)
}
