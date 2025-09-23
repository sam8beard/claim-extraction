#!/bin/bash
set -x 

go mod tidy
go build -o upload
./upload -file ../test-files/raw/sample-local-pdf.pdf -source 1104