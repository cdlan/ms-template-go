FROM golang:1.21

RUN mkdir repo-dir
RUN cd repo-dir
WORKDIR /go/repo-dir

COPY . .
RUN go install cmd/server/server.go

RUN cd ..
RUN rm -rf repo-dir
WORKDIR /go

# update packages to latest version
RUN apt update && apt -y upgrade

#ENTRYPOINT ["sleep", "infinity"] # x debug
ENTRYPOINT ["server"]