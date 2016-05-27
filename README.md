# henge
Transform multi container spec across providers: Docker Compose, Kubernetes, Openshift, etc.

## Usage

```
# This takes a Docker compose spec file and prints generated kubernetes artifacts
henge compose.yml
```


## Developing and building from source

### Setting up GOPATH

Follow instructions [here](https://golang.org/doc/code.html#GOPATH) to setup GO developer environment.


### Getting sources

If you are building upstream code
```bash
go get github.com/rtnpro/henge
cd $GOPATH/src/github.com/rtnpro/henge/
```

If you developing and using your own fork
```bash
mkdir -p $GOPATH/src/github.com/rtnpro
cd $GOPATH/src/github.com/rtnpro
git clone https://github.com/<forkid>/henge
cd henge/
git remote add upstream https://github.com/rtnpro/henge
```

### Build
Check your Go version `go version`

#### using Go v1.6
```
go build henge.go
```

#### using Go v1.5
```
GO15VENDOREXPERIMENT=1 go build henge.go
```

### Debug
You can run henge with verbose logging by adding `-v 5` option
```
./henge -v 5 docker-compose.yml
```

