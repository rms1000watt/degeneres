## {{.ProjectName}}

Description for {{.ProjectName}} goes here.

### Installation

```sh
go get -u -v {{.ImportPath}}
```

### Usage

```sh
go run main.go
```

### Deploy

```sh
go build
docker build --rm -t --no-cache {{.DockerPath}}:{{.Version}} .
docker push {{.DockerPath}}:{{.Version}}
```


*Built using [Degeneres](https://www.github.com/rms1000watt/degeneres)*
