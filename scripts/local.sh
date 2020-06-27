#!/bin/sh



go get ./...

go fmt ./...


GCP_PROJECT=research-pal-2

echo $GCP_PROJECT 

rm -f research-pal-backend 

go build -o research-pal-backend github.com/research-pal/backend/cmd

GOOGLE_APPLICATION_CREDENTIALS=/Users/muly/go/src/github.com/research-pal/backend/tmp/youtube-6d4fe2872c49.json bash -c './research-pal-backend'

