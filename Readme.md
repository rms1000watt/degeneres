## Degeneres

Degeneres, the boilerplate generator for REST-like servers in Go!

### Motivation

Golang is a fantastic language! However, you often find yourself writing a ton of boilerplate when writing the same functionality across your multiple struct types. Degeneres was built to generate the boilerplate whenever your structs change or you require different functionality on your structs. 

While gRPC (leveraging Protobuf) would handle a lot of this functionality, many developer tools and in-production systems aren't able to communicate to gRPC servers (cURL, Postman, Javascript fetch, paying business-customers, ...). So, Degeneres leverages Protobuf definitions to generate REST-like servers with JSON serialization that majority of systems can communicate with.

### Usage

#### First Time Usage

In one terminal, get and build the project:

```bash
go get -u -v github.com/rms1000watt/degeneres
cd $(go env GOPATH)/src/github.com/rms1000watt/degeneres
go build
```

Create a protobuf file at `pb/test.proto`:

```protobuf
syntax = "proto3";

package pb;

option (dg.version) = "v0.1.0";
option (dg.author) = "Ryan Smith";
option (dg.project_name) = "Test Server";
option (dg.docker_path) = "docker.io/rms1000watt/test-server";
option (dg.import_path) = "github.com/rms1000watt/test-server";

service Echo {
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
```

Generate the server code and cd to it:

```bash
./degeneres generate -f pb/test.proto -o ../test-server
cd ../test-server
```

You should now have complete server code:

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
├── doc.go
├── echo
│   ├── config.go
│   ├── echoHandler.go
│   └── preServe.go
├── helpers
│   ├── handler.go
│   ├── helpers.go
│   ├── middlewares.go
│   ├── transform.go
│   ├── unmarshal.go
│   └── validate.go
├── main.go
├── pb
│   └── test.proto
├── server
│   ├── config.go
│   └── echo.go
└── vendor
    └── vendor.json
```

Install govendor and get vendored libraries:

```bash
go get -u -v github.com/kardianos/govendor
govendor sync
```

Build and start your server with debug level logging:

```bash
go build
./test-server echo --log-level debug
```

Open a new terminal. Send a cURL to the server. (You should get a 200 with an empty JSON response: `{}`)

```bash
curl -d '{"in":"Hello World"}' http://localhost:8080/echo
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
cd $(go env GOPATH)/src/github.com/rms1000watt/test-server
go build
./test-server echo --log-level debug
```

Send the cURL request again:

```bash
curl -d '{"in":"Hello World"}' http://localhost:8080/echo
```

You should get a JSON with a hashed value back!

#### Repeated Usage

Naturally, you'll want to update your protobuf file and regenerate.

Go to your project:

```bash
cd $(go env GOPATH)/src/github.com/rms1000watt/test-server
```

Update your protobuf file `pb/test.proto`:

```protobuf
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
        option (dg.middleware.no_cache) = true;
        option (dg.method) = "POST";
    }
}

message EchoIn {
    string in   = 1 [(dg.validate) = "maxLength=100", (dg.transform) = "hash"];
    int64 age   = 2;
    string name = 3;
}

message EchoOut {
    string out  = 1;
    int64 age   = 2;
    string name = 3;
}
```

Regenerate the code:

```bash
go generate
```

Update the handler to go from Input -> Ouput

```bash
open $(go env GOPATH)/src/github.com/rms1000watt/test-server/echo/echoHandler.go
```

Rebuild and run

```bash
go build
./test-server echo --log-level debug
```

Send a cURL 

```bash
curl -d '{"in":"Hello World","age":88,"name":"Darf"}' http://localhost:8080/echo
```

Should be all good to go!

### Features

Degeneres generates boilerplate so you don't have to. It handles some field level validations & transformations and has some useful HTTP middleware. You can also generate self-signed certs.

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
- [x] Pull handlers into separate package for easier regen
- [ ] Generator validation on types to handle duplication infile and across imports
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
rm -rf out; go run main.go generate -f pb/main.proto

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
