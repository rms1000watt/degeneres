## Degeneres

Degeneres, the Golang Code generator. Use Protobuf definitions to generate complete microservices (not using RPC).

```proto
syntax = "proto3";

package pb;

option (dg.version) = "v0.1.0";
option (dg.author) = "Ryan Smith";
option (dg.project_name) = "Degeneres Test";
option (dg.docker_path) = "docker.io/rms1000watt/degeneres-test";
option (dg.import_path) = "github.com/rms1000watt/degeneres-test";
option (dg.origins) = "http://localhost,https://localhost,http://127.0.0.1,https://127.0.0.1";

service BallparkAPI {
    option (dg.short_description) = "Ballpark Service API for stadium information";
    option (dg.middleware.cors) = "true";
    option (dg.middleware.no_cache) = true;

    rpc Person(PersonIn) returns (PersonOut) {
        option (dg.method) = "GET";
        option (dg.method) = "POST";
    }

    rpc Ticket(TicketIn) returns (TicketOut) {
        option (dg.method) = "GET";
        option (dg.method) = "POST";
        option (dg.method) = "PUT";
    }

    rpc Management(ManagementIn) returns (ManagementOut) {}
}

message PersonIn {
    int64 id          = 1;
    string first_name = 2 [(dg.validate) = "maxLength=100", (dg.transform) = "truncate=50"];
    string last_name  = 3 [(dg.validate) = "maxLength=1000,minLength=1,required", (dg.transform) = "truncate=50,hash"];
}

message PersonOut {
    string first_name = 1;
    string last_name  = 2;
}

message TicketIn {
    string id = 1;
}

message TicketOut {
    string row  = 1;
    string seat = 2;
}

message ManagementIn {
    repeated bool power = 1;
}

message ManagementOut {
    repeated bool success = 1;
}
```

### Usage

In one terminal:

```bash
# Get the project
go get github.com/rms1000watt/degeneres
cd $(go env GOPATH)/src/github.com/rms1000watt/degeneres

# Get vendored projects
go get -u -v github.com/kardianos/govendor
govendor sync

# Generate self signed certs
go run main.go generate certs

# Run the project with the default protobuf as `pb/test.proto`
clear; rm -rf out; go run main.go generate

# Copy the output to a test directory
PROJECT_PATH=$(go env GOPATH)/src/github.com/rms1000watt/degeneres-test bash -c 'rm -rf $PROJECT_PATH && mkdir $PROJECT_PATH  && mkdir $PROJECT_PATH/certs && cp -r out/* $PROJECT_PATH && cp -r certs/* $PROJECT_PATH/certs && cp out/.gitignore $PROJECT_PATH/'

# Go to the test directory
cd $(go env GOPATH)/src/github.com/rms1000watt/degeneres-test

# Run the project 
cd ../degeneres-test; clear; go run main.go ballpark
cd ../degeneres-test; clear; go run main.go ballpark --certs-path ./certs --cert-name server.cer --key-name server.key
```

In another terminal:

```bash
# Run a Successful command
curl -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' http://localhost:8080/person
curl -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' -H "Origin: http://www.example.com" --verbose http://localhost:8080/person
curl -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' --insecure https://localhost:8080/person

# Run a Failing command
curl -d '{"first_name":"Chet"}' http://localhost:8080/person
curl -d '{"first_name":"Chet"}' --insecure https://localhost:8080/person
```



### TODO

- [x] Fix lexer to include `repeated`
- [x] Move `data` to different dir
- [x] Identify if message is input & create inputP
- [x] Continue refactoring templates
- [x] Check for `required` tag first then continue in order
- [] Use a logging package
- [] Use default Options method
- [x] CORS middleware
- [] Check true/false on middleware
- [] Vendoring in generated code
- [] More middleware: hsts, ssl redirect, xss protection, method logging
- [] More docs
- [] Create test repo
