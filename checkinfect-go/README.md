# check infect as go function

Author: Peter Rossbach <peter.rossbach@bee42.com> @PRossbach

list all infected images like checkinfect

## build
```
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=DARWIN -e GOARCH=amd64 golang:1.4-cross go get -d -v ; go build -v
```

build for my mac

```
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=DARWIN -e GOARCH=amd64 golang:1.4-cross go get -d -v ; go build -v
```

## run

```
INFECT_ID=8c2e06607696bd4afb3d03b687e361cc43cf8ec1a4a725bc96e39f05ba97dd55 ./images
&{[{fbeddc9c42d2340ee66049d246415b2c4df77ea9fdd153c2499638669fa2df00 [hello:latest 127.0.0.1:5000/hello:latest] 1434829473 0 2433303 8c2e06607696bd4afb3d03b687e361cc43cf8ec1a4a725bc96e39f05ba97dd55} {8c2e06607696bd4afb3d03b687e361cc43cf8ec1a4a725bc96e39f05ba97dd55 [busybox:latest] 1429308073 0 2433303 6ce2e90b0bc7224de3db1f0d646fe8e2c4dd37f1793928287f6074bc451a57ea}] [{773e5a0d97aa8a8bde0dd879de63a63b34a646e9da2a08dc24146d108792c7ee busybox /bin/sh -c 'while true ; do echo hello; sleep 10 ; done' 1434855854 Up About an hour [] 0 0 [/high_goldstine]}]}
```
