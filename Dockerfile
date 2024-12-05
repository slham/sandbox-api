FROM golang:alpine

ARG ex_path

# move files to container
ADD ./build/$ex_path /go/src/app/bin
WORKDIR /go/src/app

# give permission to run executable
RUN chmod +x ./bin/sandbox-api

CMD ./bin/sandbox-api
