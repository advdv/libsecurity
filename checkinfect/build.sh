#!/bin/bash
IMAGE=checkinfect
ACCOUNT=infrabricks
TAG_LONG=0.0.1
docker build -t="${ACCOUNT}/$IMAGE" .
DATE=`date +'%Y%m%d%H%M'`
IID=$(docker inspect -f "{{.Id}}" ${ACCOUNT}/$IMAGE)
docker tag -f $IID ${ACCOUNT}/$IMAGE:$DATE
docker tag -f $IID ${ACCOUNT}/$IMAGE:$TAG_LONG
