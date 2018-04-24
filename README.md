# serverlessGoDynamoDB
# 1.Install go
$ tar -C /usr/local -xzf go1.10.1.linux-amd64.tar.gz

$ export PATH=$PATH:/usr/local/go/bin

$ export GOPATH=$HOME/go

$ export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

$ export GOBIN=$HOME/go/bin

# 2. Install nodejs
$ curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -

$ sudo apt-get install -y nodejs

# 3.Install serverless
$ npm install -g serverless

# 4.Create AWS account
https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/

# 5.Create an I AM user and configure credentials
https://serverless.com/framework/docs/providers/aws/guide/credentials/

$ serverless config credentials (-o) --provider aws --key XXX --secret YYY

# 6.Test
Change directory to addDevice (or getDevice) for running unit test of it :

$ go test

# 7.Compile
Change into service directory and compile :

$ make

# 8.Deploy
$ serverless deploy (Or $ sls deploy)

# 9.Add Device
Enter data that you want to insert in /lib/data.json.
$ curl -X POST -H "Content-Type: application/json" -d @lib/data.json  https://XXX.execute-api.us-east-1.amazonaws.com/dev/devices

Note : put the url which is created when you deploy service.

Result :

Success: Inserted record in json format 

data.json contains duplicate id or empty fields: Bad Request message 

Internal error: Internal Server Error message 

# 10.Get Device By ID
$ curl https://XXX.execute-api.us-east-1.amazonaws.com/dev/devices/"id7"

Result :

Success: Record in json format 

Id does not exist: Not Found message 

Internal error: Internal Server Error message
