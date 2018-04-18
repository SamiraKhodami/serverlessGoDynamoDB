build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/lib/dbconfig lib/dbconfig/dbconfig.go  
	env GOOS=linux go build -ldflags="-s -w" -o bin/lib/devicecrud lib/devicecrud/devicecrud.go  
	env GOOS=linux go build -ldflags="-s -w" -o bin/lib/stringutil lib/stringutil/stringutil.go  
	env GOOS=linux go build -ldflags="-s -w" -o bin/addDevice addDevice/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getDevice getDevice/main.go
