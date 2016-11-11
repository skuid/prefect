# To build:
# $ docker run --rm -v $(pwd):/go/src/github.com/skuid/prefect -w /go/src/github.com/skuid/prefect golang:1.7 go build -v -a -tags netgo -installsuffix netgo -ldflags '-w'
# $ docker build -t skuid/prefect .
#
# To run:
# $ docker run skuid/prefect

FROM busybox

MAINTAINER Micah Hausler, <micah@skuid.com>

COPY prefect /bin/prefect
RUN chmod 755 /bin/prefect

ENTRYPOINT ["/bin/prefect"]
