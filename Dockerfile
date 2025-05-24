FROM golang:1.24 AS build

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod tidy && go mod download && go mod verify

COPY . .
RUN go test -race -v ./...
RUN CGO_ENABLED=0 go build -o parking-system .

FROM alpine:3.21.0
RUN apk --no-cache update && apk --no-cache  add bash
WORKDIR /app
COPY --from=build /app/parking-system .
RUN chmod +x parking-system

ENTRYPOINT ["/app/parking-system"]