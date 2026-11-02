SHELL := /bin/bash
PREFIX ?= $(HOME)/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3
.PHONY: help install test example memory-test

help:
	@bash installer.sh --help

install:
	bash installer.sh

test:
	. "$(PREFIX)/env.sh"; go test -buildvcs=false ./...

example:
	. "$(PREFIX)/env.sh"; go run -buildvcs=false ./examples/kmeans

memory-test:
	. "$(PREFIX)/env.sh"; GOMAXPROCS=1 OMP_NUM_THREADS=1 OPENBLAS_NUM_THREADS=1 go test -buildvcs=false -tags mlpack_memtest -run 'TestNative' -v .
