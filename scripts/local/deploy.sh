#!/bin/sh


go mod vendor

go fmt ./...

GCP_PROJECT=research-pal-2
echo $GCP_PROJECT 

rm -f research-pal-backend 

go build -o research-pal-backend github.com/research-pal/backend/cmd

# TODO: need to remove the hardcoding, and save this in local env. and in this script just verify if a value is set for this env or not
GOOGLE_APPLICATION_CREDENTIALS=/Users/muly/go/src/github.com/research-pal/backend/tmp/research-pal-2-a30e5d878434.json bash -c './research-pal-backend'

