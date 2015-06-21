#!/bin/bash
: ${INFECT_ID:=$1}
: ${DOCKER_HOST:=http://127.0.0.1:8080}

socat TCP-LISTEN:8080,fork UNIX:/var/run/docker.sock &

IMAGES=$(wget $DOCKER_HOST/images/json -O - -q |jq -r ".[].Id" | tr "\n" " ")
for image in $IMAGES ; do
  CANDIDATE=$(wget $DOCKER_HOST/images/${image}/history -O - -q |jq -r ".[].Id" | grep $INFECT_ID )
  if [ "${CANDIDATE}z" != "z" ]; then
    wget $DOCKER_HOST/images/json -O - -q | jq '.[]' | jq '{image_id: .Id, tags: .RepoTags}' |jq "if .image_id == \"${image}\" then . else empty end"
  fi
done

CONTAINERS=$(wget $DOCKER_HOST/containers/json -O - -q | jq -r ".[].Id" | tr "\n" " " )

for container in $CONTAINERS ; do
  IMAGE=$(wget $DOCKER_HOST/containers/${container}/json -O - -q | jq -r ".Image")
  CANDIDATE=$(wget $DOCKER_HOST/images/${IMAGE}/history -O - -q |jq -r ".[].Id" | grep $INFECT_ID )
  if [ "${CANDIDATE}z" != "z" ]; then
     echo "{ \"container_id\": \"${container}\", \"image_id\": \"${IMAGE}\" }" | jq "."
  fi
done
