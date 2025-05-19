FROM golang:1.23.5 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify && go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /server cmd/main.go

FROM alpine:latest AS release-stage

WORKDIR /

COPY --from=build-stage /server /server
COPY --from=build-stage /app/configs/config.yml /configs/config.yml

EXPOSE 8080

ENTRYPOINT ["/server"]
