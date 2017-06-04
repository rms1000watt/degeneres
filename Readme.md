## Degeneres

Degeneres, the Golang Code generator. Use Protobuf definitions to generate complete microservices.

```proto
syntax = "proto3";

package pb;

import "github.com/rms1000watt/degeneres/dg.proto";

option (dg.rpc_to_endpoints) = true;
option (dg.out_dir) = "./out";
option (dg.author) = "Ryan Smith";

service BallparkAPI {
    option (dg.middleware.cors) = "localhost,127.0.0.1,www.example.com";

    rpc Person(PersonIn) returns (PersonOut) {
        option (dg.middleware.no_cache) = true;
        option (dg.method) = "GET";
        option (dg.method) = "POST";
    }

    rpc Ticket(TicketIn) returns (TicketOut) {
        option (dg.middleware.no_cache) = true;
        option (dg.method) = "GET";
        option (dg.method) = "POST";
    }

}

message PersonIn {
    int64 id          = 1;
    string first_name = 2 [(dg.validate.max_length) = 100, (dg.transform.truncate) = 50];
    string last_name  = 3 [(dg.validate.max_length) = 100, (dg.transform.truncate) = 50];
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
```

### Development Commands

```sh
clear; rm -rf out; go run main.go generate
PROJECT_PATH=$(go env GOPATH)/src/github.com/rms1000watt/degeneres-test bash -c 'rm -rf $PROJECT_PATH && mkdir $PROJECT_PATH  && mkdir $PROJECT_PATH/certs && cp -r out/* $PROJECT_PATH && cp -r certs/* $PROJECT_PATH/certs && cp out/.gitignore $PROJECT_PATH/'

cd ../degeneres-test; clear; go run main.go ballpark --certs-path ./certs

# Fail
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' localhost:8080/person

# Success
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' localhost:8080/person
curl -X POST -d '{"first_name":"ChetChetChetChet","middle_name":"Darf","last_name":"Star","age":33,"account":123.123,"password":"pASSword"}' --insecure https://localhost:8080/person

```

### TODO

- [x] Fix lexer to include `repeated`
- [x] Move `data` to different dir
- [] Identify if message is input & create inputP
- [] Continue refactoring templates
