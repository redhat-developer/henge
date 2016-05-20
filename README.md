# henge
Transform multi container spec across providers: Docker Compose, Kubernetes, Openshift, etc.

## Usage

```
# This takes a Docker compose spec file and prints generated kubernetes artifacts
henge compose.yml
```


## Developing and building from source
### Getting sources
```
mkdir -p $GOPATH/src/github.com/rtnpro
cd $GOPATH/src/github.com/rtnpro


# if you are building upstream code
git clone https://github.com/rtnpro/henge


# if you developing and using your own fork
git clone git@github.com:<forkid>/henge.git # Replace <forkid> with the your github id
git remote add upstream https://github.com/rtnpro/henge
```

### Build
Check your Go version `go version`

#### using Go v1.6
```
cd $GOPATH/src/github.com/rtnpro
go build henge.go
```

#### using Go v1.5
```
cd $GOPATH/src/github.com/rtnpro
GO15VENDOREXPERIMENT=1 go build henge.go
```

### Debug
You can run henge with verbose logging by adding `-v 5` option
```
./henge -v 5 docker-compose.yml
```

