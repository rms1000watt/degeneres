## Degeneres

Degeneres, the microservice generator. Use Protobuf definitions to generate complete REST-like microservices in Golang.

### Features

- Protobuf data/service configuration
- Input validation on `json` tags
- Input transformation on `json` tags
- CLI Commander (spf13/cobra)
- Middleware (security, logging, cors, no-cache)
- Self signed key generation
- Vendored libraries

### Example

Take an input protobuf file and generate your full API

```proto
syntax = "proto3";

package pb;

option (dg.version) = "v0.1.0";
option (dg.author) = "Ryan Smith";
option (dg.project_name) = "Degeneres Test";
option (dg.docker_path) = "docker.io/rms1000watt/degeneres-test";
option (dg.import_path) = "github.com/rms1000watt/degeneres-test";
option (dg.origins) = "http://localhost,https://localhost,http://127.0.0.1,https://127.0.0.1";

service Ballpark {
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

After you `go build` the generated code, you can run...

```
degeneres-test ballpark
```

...and have endpoints accesible at `/person`, `/ticket`, and `/mangement`.

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

# Run the project with or without TLS
cd ../degeneres-test; govendor sync; clear; go run main.go ballpark --log-level debug
cd ../degeneres-test; govendor sync; clear; go run main.go ballpark --log-level debug --certs-path ./certs --cert-name server.cer --key-name server.key
```

In another terminal:

```bash
# Run a Successful command
curl -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' -H "Origin: http://www.example.com" -D - http://localhost:8080/person
curl -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' -H "Origin: https://www.example.com" -D - --insecure https://localhost:8080/person

# Run a Failing command
curl -d '{"first_name":"Chet"}' http://localhost:8080/person
curl -d '{"first_name":"Chet"}' -H "Origin: http://www.example.com" http://localhost:8080/person
curl -d '{"first_name":"Chet"}' --insecure https://localhost:8080/person
```

### Motivation

gRPC is a brilliant system: describe data and services in Protobuf definitions that are used to generate servers in the language of your choosing. Need to make a change to the data or service? Update the Protobufs, regenerate the servers, rinse, repeat. This has been a fantastic workflow in production environments.

The only downside at the moment is the gRPC ecosystem isn't readily available to many web toolsets & systems people use regularly like Postman, cURL, Angular/React, legacy/production JSON servers, etc. (although this will change in short order as the development of gRPC-web has been active for some time). So, Degeneres solves this by changing the data serialization from Protobuf to JSON and exposing REST-like endpoints.

Also, for convenience and the reduction of hand-typed boilerplate, input data validation & transformation have been added along with configurable middleware at the service or endpoint levels.

Finally, there is no way for the templates and configurable parameters to fit everyones' needs. The templates and configurable parameters have been written to hit a sweet spot between ease-of-use and excessiveness. So, to fit your needs--please fork the project and adjust accordingly. Degeneres should be incorporated into your development workflow; **generate the boilerplate**, don't write it.

### Limitations

- Golang server generation only
- Less performant than gRPC (JSON vs Protobuf)
- Not production ready

### TODO

- [x] Fix lexer to include `repeated`
- [x] Move `data` to different dir
- [x] Identify if message is input & create inputP
- [x] Continue refactoring templates
- [x] Check for `required` tag first then continue in order
- [x] Use a logging package
- [x] Use default Options method
- [x] Convert generator warnings to errors
- [x] CORS middleware
- [x] Check true/false on middleware
- [x] Vendoring in generated code
- [x] More middleware: hsts, ssl redirect, xss protection, method logging
- [x] Add Catch-all, root handler with debug "path not found"
- [ ] More docs
- [ ] Test protobuf with gRPC
- [ ] Generate unit tests
- [ ] Create test repo
- [ ] Static file handling + gzip
