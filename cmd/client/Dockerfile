FROM golang:1.18-alpine3.15 as Builder
COPY ../.. /go/src
WORKDIR /go/src
RUN go mod tidy -compat=1.18
RUN go build -o /go/bin/client /go/src/cmd/client/*.go

FROM alpine:3.15
COPY --from=Builder /go/bin/* /