#!/bin/sh
GOOS=windows GOACH=amd64 go build -o qbtchangetracker_v${1}_amd64.exe -tags forceposix
GOOS=windows GOARCH=386 go build -o qbtchangetracker_v${1}_i386.exe -tags forceposix
GOOS=linux GOARCH=amd64 go build -o qbtchangetracker_v${1}_amd64_linux -tags forceposix
GOOS=linux GOARCH=386 go build -o qbtchangetracker_v${1}_i386_linux -tags forceposix
GOOS=darwin GOARCH=amd64 go build -o qbtchangetracker_v${1}_amd64_macos -tags forceposix
