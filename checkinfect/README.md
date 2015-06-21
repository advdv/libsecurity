# Check infect docker images

Author: Peter Rossbach <peter.rossbach@bee42.com> @PRossbach

* identify infect images id or layer!
* List all repo/images from current host
* List all layer from image
* find if an image layer is equal infected layer
* result list all image IDs and repo tags

## build and usage

### build your own machine

```
$ docker-machine create -d virtualbox --engine-insecure-registry 127.0.0.1:5000 dockercheck
$ eval $(docker-machine env dockercheck)
```

### build and test checkinfect

```
$ ./build.sh
$ mkdir test
$ cd test
$ cat >Dockerfile <<EOF
FROM busybox
CMD ["echo", "hello"]
EOF
$ docker build -t hello .
docker images --no-trunc |grep hello
hello                     latest              fbeddc9c42d2340ee66049d246415b2c4df77ea9fdd153c2499638669fa2df00   6 hours ago         2.433 MB

$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock infrabricks/checkinfect fbeddc9c42d2340ee66049d246415b2c4df77ea9fdd153c2499638669fa2df00
{
  "id": "fbeddc9c42d2340ee66049d246415b2c4df77ea9fdd153c2499638669fa2df00",
  "tags": [
    "hello:latest",
  ]
}
$ alias checkinfect="docker run --rm -v /var/run/docker.sock:/var/run/docker.sock infrabricks/checkinfect"
# docker images --no-trunc
$ docker tag hello 127.0.0.1:5000/hello
$ IMAGE_LAYER=fbeddc9c42d2340ee66049d246415b2c4df77ea9fdd153c2499638669fa2df00
$ checkinfect $IMAGE_LAYER
{
  "id": "fbeddc9c42d2340ee66049d246415b2c4df77ea9fdd153c2499638669fa2df00",
  "tags": [
    "hello:latest",
    "127.0.0.1:5000/hello:latest"
  ]
}
```

## More examples

```
$ docker-machine create -d virtualbox --engine-insecure-registry 127.0.0.1:5000 dockercheck
$ eval $(docker-machine env dockercheck)
$ docker run -d -p 5000:5000 registry:2.0
$ docker images --digests=true
REPOSITORY          TAG                 DIGEST              IMAGE ID            CREATED             VIRTUAL SIZE
registry            2.0                 <none>              ec94325cd2c4        3 days ago          548.6 MB
$ wget --no-check-certificate --certificate=$DOCKER_CERT_PATH/cert.pem --private-key=$DOCKER_CERT_PATH/key.pem https://$(docker-machine ip dockercheck):2376/images/json?digest=1 -O - -q | jq "."
[
  {
    "Id": "ec94325cd2c49e50c3ce74d16169c536f5d423158b10e4481edc909c460301d3",
    "ParentId": "dea71ce2cbe00e98ca4413c7e93c9542513f011c588f6bcc9ed78b5431f30ad1",
    "RepoTags": [
      "registry:2.0"
    ],
    "RepoDigests": [],
    "Created": 1434531372,
    "Size": 0,
    "VirtualSize": 548616301,
    "Labels": {}
  }
]
$ cat >Dockerfile <<EOF
FROM busybox
CMD ["echo", "hello"]
EOF
$ docker build -t hello .
$ docker tag hello 127.0.0.1:5000/hello
$ docker push 127.0.0.1:5000/hello
$ docker pull 127.0.0.1:5000/hello@sha256:<digest>
$ docker pull registry@sha256:e62d5cdc270975f39eb469faf5cd7ba8c63f71ecfadb1f3e47557301c489516a

```

* Only see digests, if you pull with digest
* Image Layer Id are the same in all registries
* ID's are auto generated at every docker host

## Todo

* check also container that are available

## Reference

* https://docs.docker.com/registry/spec/api/
* https://github.com/advanderveer/docksec
* https://docs.docker.com/reference/api/docker_remote_api_v1.18/#inspect-a-container
* https://docs.docker.com/reference/api/docker_remote_api_v1.18/#list-images
* https://docs.docker.com/reference/api/docker_remote_api_v1.18/#get-the-history-of-an-image
