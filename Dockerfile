FROM golang:1.21.9-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/assessment-tax .

### ----------------------------------

FROM alpine:3.16.2
COPY --from=build-base /app/out/assessment-tax /app/assessment-tax

CMD ["/app/assessment-tax"]