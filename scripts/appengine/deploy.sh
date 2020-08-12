#!/bin/sh

go mod vendor

cd cmd

gcloud config set project research-pal-2

gcloud app deploy 


