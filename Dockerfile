FROM golang:latest AS build-env
WORKDIR $GOPATH/src/golang-block-chain
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -installsuffix cgo .

FROM scratch
ENV IS_CONTAINER=TRUE
COPY --from=build-env go/src/golang-block-chain/golang-block-chain .
EXPOSE 8081
CMD ["./golang-block-chain"]