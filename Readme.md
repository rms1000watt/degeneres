## Degeneres

Degeneres, the boilerplate generator for REST-like servers in Go!

### Motivation

Golang is a fantastic language! However, you often find yourself writing a ton of boilerplate when writing the same functionality across your multiple struct types. Degeneres was built to generate the boilerplate whenever your structs change or you require different functionality on your structs. 

### Usage

#### First Time Usage

In one terminal:

```bash
# Get and build the project
go get -u -v github.com/rms1000watt/degeneres
cd $(go env GOPATH)/src/github.com/rms1000watt/degeneres
go build

# Create a protobuf file
cat << EOF > pb/test.proto
syntax = "proto3";

package pb;

option (dg.version) = "v0.1.0";
option (dg.author) = "Ryan Smith";
option (dg.project_name) = "Test Server";
option (dg.docker_path) = "docker.io/rms1000watt/test-server";
option (dg.import_path) = "github.com/rms1000watt/test-server";

service Echo {
    option (dg.middleware.logger) = true;

    rpc Echo(EchoIn) returns (EchoOut) {
        option (dg.method) = "POST";
    }
}

message EchoIn {
    string in = 1 [(dg.validate) = "maxLength=100", (dg.transform) = "hash"];
}

message EchoOut {
    string out = 1;
}
EOF

# Generate the server code and go to it
./degeneres generate -f pb/test.proto -o ../test-server
cd ../test-server

# Install govendor and get vendored libraries
go get -u -v github.com/kardianos/govendor
govendor sync

# Build and start your server with debug level logging
go build
./test-server echo --log-level debug
```

Open a new terminal:

```bash
# Send a cURL to the server
# You should get a 200 with an empty JSON response: {}
curl -d `{"in":"Hello World"}` http://localhost:8080/echo
```

You get an empty JSON response because the logic to go from Input -> Output is up to you. Edit the handler to fill in the Output logic

```bash
open $(go env GOPATH)/src/github.com/rms1000watt/test-server/echo/echoHandler.go
```

And add 1 line to echo the response in `EchoHandlerPOST`:

```go
echoOut.Out = echoIn.In
```

Now rebuild and run the server again:

```bash
# Rebuild and run
cd $(go env GOPATH)/src/github.com/rms1000watt/test-server
go build
./test-server echo --log-level debug
```

Send the cURL request again:

```bash
curl -d `{"in":"Hello World"}` http://localhost:8080/echo
```

You should get a JSON with a hashed value back!

#### Repeated Usage

Naturally, you want to update your protobuf file and regenerate.

```bash
# Go to your project
cd $(go env GOPATH)/src/github.com/rms1000watt/test-server

# Update your protobuf file
cat << EOF > pb/test.proto
syntax = "proto3";

package pb;

option (dg.version) = "v0.1.0";
option (dg.author) = "Ryan Smith";
option (dg.project_name) = "Test Server";
option (dg.docker_path) = "docker.io/rms1000watt/test-server";
option (dg.import_path) = "github.com/rms1000watt/test-server";

service Echo {
    option (dg.middleware.logger) = true;
    option (dg.middleware.no_cache) = true;

    rpc Echo(EchoIn) returns (EchoOut) {
        option (dg.method) = "GET";
        option (dg.method) = "POST";
    }
}

message EchoIn {
    string in   = 1 [(dg.validate) = "maxLength=100", (dg.transform) = "hash"];
    int age     = 2;
    string name = 3;
}

message EchoOut {
    string out  = 1;
    int age     = 2;
    string name = 3;
}
EOF

# Regenerate
go generate

# Rebuild and run
go build
./test-server echo --log-level debug
```

### Features

Degeneres generates the boilerplate for you! From the `test.proto` file defined above, you get a complete go server:

```
.
├── Dockerfile
├── License
├── Readme.md
├── cmd
│   ├── echo.go
│   ├── root.go
│   └── version.go
├── data
│   ├── data.go
│   └── input.go
├── echo
│   ├── config.go
│   ├── echo.go
│   └── echoHandler.go
├── helpers
│   ├── handler.go
│   ├── helpers.go
│   ├── middlewares.go
│   ├── transform.go
│   ├── unmarshal.go
│   └── validate.go
├── main.go
└── vendor
    └── vendor.json
```

#### Validations

#### Transformations

#### Middleware

- `dg.middleware.cors`

#### Self-Signed Keys

```bash
./degeneres generate certs
```

### Limitations

- Server generation for Golang only
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
- [x] `go generate` to self regen generated code
- [x] Copy proto file into generated project
- [ ] More docs
- [ ] More examples
- [ ] Generate unit tests
- [ ] Workout kinks in workflow
- [ ] Better stubbing of handlers

### Dev Commands...

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

# Run the project with the default protobuf as `pb/main.proto`
clear; rm -rf out; go run main.go generate -f pb/main.proto

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
