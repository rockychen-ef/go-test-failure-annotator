# Download dependencies
FROM golang:1.12 as dependencies
ENV GO111MODULE=on
WORKDIR $GOPATH/src/el-sample-go-actions
# COPY go.mod .
# COPY go.sum .
# RUN go mod download
COPY . .

# Compile
FROM dependencies as compile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $GOPATH/bin/tfa -mod vendor

# Bin
FROM alpine:3.9 as binary
RUN apk update && apk add ca-certificates

WORKDIR /opt
COPY --from=compile /go/bin/tfa ./tfa
CMD /opt/tfa