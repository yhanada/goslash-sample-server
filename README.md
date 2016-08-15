# goslash-sample-server

## install dependencies

```shell
$ make install
```

## run dev mode

```shell
$ make run
```

## build executable

```shell
$ make build
```

## build on docker

```shell
$ docker run -it --rm -v $PWD:/go/src/app golang:onbuild /bin/bash -c "make clean && make build"
```
