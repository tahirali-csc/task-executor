#!/bin/sh

env GOOS=linux GOARCH=amd64 go build -o repo-cloner main.go

docker build -t repo-cloner:latest .