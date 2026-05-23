#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o app_linux main.go
