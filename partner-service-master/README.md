# Partner Settings
Backend microservice written in Go to create, edit, and return specific preferences wholesale partners have when it comes to EDI, currency, etc.
### Setup still in progress

_________

## Getting Started

. [Install gRPC protoc](https://github.com/google/protobuf/releases)
    1. Scroll to `Downloads` and get the binary appropriate for your system (probably [64 bit for OS X](https://github.com/google/protobuf/releases/download/v3.3.0/protoc-3.3.0-osx-x86_64.zip)).
    1. Unzip the file (may have been done automatically)
    1. place the executable in the protoc-3/bin directory somewhere in your $PATH (probably /usr/local/bin)
1. Install grpc and grpc-gateway packages.
    * `$ go get -u google.golang.org/grpc` is the base grpc package
    * `$ go get -u github.com/golang/protobuf/protoc-gen-go` allows the protoc binary from above to generate Go code
    * `$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway` creates the gateway during compilation
    * `$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger` generates swagger file
1. From the tutorial root directory run the command `$ bash ./pb/compile.sh`. This generates 3 files in pb:
    * partner_service.pb.go, which contains the Go interfaces for the main gRPC service
    * partner_service.pb.gw.go, which contains the Go interfaces for the REST proxy interfaces
    * partner_service.swagger.json, which is swagger documentation for our REST endpoints

`go get github.com/golang/dep/cmd/dep`
`dep ensure`
`go run`

### `Go build` is mostly for distribution. 

Go can cross compile, so if you were going to put your program up for download by different systems you could "go build" for Mac, Linux and windows all from your machine, or a circle instance, etc.


### `Go run` for development
Go run just gets your program running so you can work with it.



