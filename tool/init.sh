#!/bin/sh

go get -u github.com/golang/dep/cmd/dep
go get -u github.com/golang/lint/golint
dep ensure
dep status
