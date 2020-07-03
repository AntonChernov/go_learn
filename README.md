# Learning of Go programming language 
API service Just for fun


## go.mod
go mod init  --> create a "go.mod", "go.sum" files this files something 
like requirements in python but one thing about this files it is used as a replacement for GOPATH

detailed about packages https://blog.golang.org/using-go-modules and https://golangbot.com/go-packages/

## Go, vendor 
go mod vendor --> isolate env it means all packages will be installed in project dir not in main go dir

NOTE add vendor dir to .gitignore before use command!

## Run server (Server no blocking operations! Use channels and goroutines.)
go run main.go --> server start at 127.0.0.1:8800 

#### CLI:

go run main.go -hostport=<some_another_host:port>


## TODO:

#### 0) Fix Routers (Collect routers and put in to main.go use function for it)
#### 1) Tests
#### 2) Settings
#### 3) Validation
#### 4) Check exist ORM (Load test for understand spped of GO RAW query and ORM representation)
#### 5) API doc (Find some "Swagger" for docwlent API`s)
#### 6) add Pagination for API`s (normal pagination on the interface or something like this need investigate)
#### 7) send an email for "verify" user email(Two places when user-created and when user change the email) 
#### 8) Add Permissions(maybe this part already implement in some lib need to investigate)
