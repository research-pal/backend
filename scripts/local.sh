#!/bin/sh


go mod vendor

go fmt ./...

GCP_PROJECT=research-pal-2
echo $GCP_PROJECT 

rm -f research-pal-backend 

go build -o research-pal-backend github.com/research-pal/backend/cmd

GOOGLE_APPLICATION_CREDENTIALS=/Users/p/go/src/github.com/research-pal/backend/tmp/research-pal-2-3ffc6017bbbb.json bash -c './research-pal-backend'

