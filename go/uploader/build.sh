#!/bin/bash
set -x 

go mod tidy
go build -o upload