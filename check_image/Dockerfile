FROM google/golang:latest
ADD main.go /
RUN mkdir /cve
RUN  go get -v github.com/fsouza/go-dockerclient
ADD cve/. /cve/
CMD go run main.go CVE-2014-6271
