## Degeneres

Degeneres, the boilerplate generator for REST-like servers in Go!

### Example

An example server generated with Degeneres can be seen at [http://github.com/rms1000watt/degeneres-test](http://github.com/rms1000watt/degeneres-test)

### Motivation

Golang is a fantastic language! However, you often find yourself writing a ton of boilerplate when writing the same functionality across your multiple struct types. Degeneres was built to generate the boilerplate whenever your structs change or you require different functionality on your structs. 

While gRPC (leveraging Protobuf) would handle a lot of this functionality, many developer tools and in-production systems aren't able to communicate to gRPC servers (cURL, Postman, Javascript fetch, paying business-customers, etc.) So, Degeneres leverages Protobuf definitions to generate REST-like servers with JSON serialization that majority of systems can communicate with.

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

Validations have the Protobuf field level option syntax:

```proto
string first_name = 1 [(dg.validate) = "minLength=2,maxLength=100"];
```

String Validations:

| Validation | Usage | Example | Description |
| --- | --- | --- | --- |
| Max Length | `maxLength=VALUE` | `maxLength=100` | Fails if len(input) > maxLength |
| Min Length | `minLength=VALUE` | `minLength=2` | Fails if len(input) > maxLength |
| Must Have Chars | `mustHaveChars=VALUE` | `mustHaveChars=aeiou` | Fails if chars in VALUE are not in input |
| Can't Have Chars | `cantHaveChars=VALUE` | `cantHaveChars=aeiou` | Fails if chars in VALUE are in input |
| Only Have Chars | `onlyHaveChars=VALUE` | `onlyHaveChars=aeiou` | Fails if input has chars not in VALUE |

Float and Int Validations:

| Validation | Usage | Example | Description |
| --- | --- | --- | --- |
| Greater Than | `greaterThan=VALUE` | `greaterThan=100` | Fails if input < VALUE |
| Less Than | `lessThan=VALUE` | `lessThan=100` | Fails if input > VALUE |


#### Transformations

Transformations have the Protobuf field level option syntax:

```proto
string first_name = 1 [(dg.transform) = "truncate=50,hash"];
```

And can be combined with other options:

```proto
string first_name = 1 [(dg.validate) = "maxLength=100", (dg.transform) = "truncate=50,hash"];
```

General Trasformation:

| Transformation | Usage | Example | Description |
| --- | --- | --- | --- |
| Default | `default=VALUE` | `default=Darf` | Sets input as VALUE if input is nil |

String Transformations:

| Transformation | Usage | Example | Description |
| --- | --- | --- | --- |
| Hash | `hash` | `hash` | Essentially `hexEncode(sha256(input))` |
| Encrypt | `encrypt` | `encrypt` | `aesgcm.Seal` (**DONT USE DEFAULT VALUES OR SCHEME IN PRODUCTION!**) |
| Decrypt | `decrypt` | `decrypt` | `aesgcm.Open` (**DONT USE DEFAULT VALUES OR SCHEME IN PRODUCTION!**) |
| Trim Chars | `trimChars=VALUE` | `trimChars=xx` | Uses `strings.Trim(input, VALUE)` |
| Trim Space | `trimSpace` | `trimSpace` | Uses `strings.TrimSpace(input)` |
| Truncate | `truncate=VALUE` | `truncate=50` | Essentially `input[:VALUE]` |
| Password Hash | `passwordHash` | `passwordHash` | argon2 password hashing (**PLEASE INSPECT CODE THOROUGHLY**) |

#### Middleware

Middleware can be added as `service` or `rpc` options.

```proto
option (dg.middleware.no_cache) = true;
```

| Middleware | Usage | Description |
| --- | --- | --- |
| Logger | `dg.middleware.logger` | Logs the time for each request | 
| CORS | `dg.middleware.cors` | Adds CORS headers if `option (dg.origins)` is added | 
| Secure | `dg.middleware.secure` | Does TLS Redirect & adds HSTS, XSS Protection, Nosniff, Frame deny headers | 
| No Cache | `dg.middleware.no_cache` | Adds no cache headers | 

#### Self-Signed Keys

```bash
./degeneres generate certs
```

### Limitations

- Server generation for Golang only
- Less performant than gRPC (JSON vs Protobuf)

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
- [ ] Docs to show all options
- [ ] More examples
- [ ] Generate unit tests
- [ ] Example repo in github
- [x] Expvar
- [ ] DB connection example (inversion of control)
- [ ] Docker compose
- [ ] Workout kinks in workflow
