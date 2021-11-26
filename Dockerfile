FROM golang:buster as builder
WORKDIR /builder
RUN apt update
RUN apt install -y make protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
ENV CGO_ENABLED=0
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
RUN make mod-download
COPY . .
RUN make compile

FROM alpine
WORKDIR /k8trics
COPY --from=builder /builder/bin/k8trics ./k8trics-bin
CMD [ "/k8trics/k8trics-bin" ]