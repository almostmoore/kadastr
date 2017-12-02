#!/bin/bash

VERSION=$1

# Kadastr
docker build -t iamsalnikov/kadastr:$VERSION -f Dockerfile .
docker build -t iamsalnikov/kadastr:latest -f Dockerfile .
docker push iamsalnikov/kadastr:$VERSION
docker push iamsalnikov/kadastr:latest

# API
docker build -t iamsalnikov/kadastr_api:$VERSION -f Dockerfile_api .
docker build -t iamsalnikov/kadastr_api:latest -f Dockerfile_api .
docker push iamsalnikov/kadastr_api:$VERSION
docker push iamsalnikov/kadastr_api:latest

# FP
docker build -t iamsalnikov/kadastr_fp:$VERSION -f Dockerfile_fp .
docker build -t iamsalnikov/kadastr_fp:latest -f Dockerfile_fp .
docker push iamsalnikov/kadastr_fp:$VERSION
docker push iamsalnikov/kadastr_fp:latest

# QUARTER
docker build -t iamsalnikov/kadastr_quarter:$VERSION -f Dockerfile_quarter .
docker build -t iamsalnikov/kadastr_quarter:latest -f Dockerfile_quarter .
docker push iamsalnikov/kadastr_quarter:$VERSION
docker push iamsalnikov/kadastr_quarter:latest

# TG
docker build -t iamsalnikov/kadastr_tg:$VERSION -f Dockerfile_tg .
docker build -t iamsalnikov/kadastr_tg:latest -f Dockerfile_tg .
docker push iamsalnikov/kadastr_tg:$VERSION
docker push iamsalnikov/kadastr_tg:latest