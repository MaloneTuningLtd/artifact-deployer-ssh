FROM golang:alpine AS build-env
RUN apk --no-cache add git
COPY . /go/src/github.com/MaloneTuningLtd/artifact-deployer-ssh
RUN cd /go/src/github.com/MaloneTuningLtd/artifact-deployer-ssh \
  && go get \
  && go build -o artifact-deployer

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/MaloneTuningLtd/artifact-deployer-ssh/artifact-deployer /usr/local/bin/artifact-deployer
ENTRYPOINT /usr/local/bin/artifact-deployer
