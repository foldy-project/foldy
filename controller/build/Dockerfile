FROM golang:1.13 AS builder
ENV GO111MODULE=on
WORKDIR /go/src/github.com/foldy-project/foldy
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /go/src/github.com/foldy-project/foldy/controller/cmd/controller
#RUN go build \
#    -o foldy-controller \
#    -gcflags all=-trimpath=/go/src/github.com/foldy-project/foldy \
#    -asmflags all=-trimpath=/go/src/github.com/foldy-project/foldy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o foldy-controller
RUN ls -al
RUN pwd

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
ENV OPERATOR=/usr/local/bin/foldy-controller \
    USER_UID=1001 \
    USER_NAME=foldy-controller
COPY --from=builder /go/src/github.com/foldy-project/foldy/controller/cmd/controller/foldy-controller ${OPERATOR}
COPY build/bin /usr/local/bin
RUN chmod +x /usr/local/bin/entrypoint \
    && chmod +x /usr/local/bin/user_setup \
    && /usr/local/bin/user_setup
ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER ${USER_UID}
