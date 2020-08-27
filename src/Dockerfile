FROM xo-api-base as builder
ADD ./api /go/src/api
WORKDIR /go/src/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /go/src/api
COPY --from=builder /go/src/api/api .
EXPOSE 1323

ENTRYPOINT ["./api"]
