# syntax=docker/dockerfile:1

#
# Build
#
FROM golang:1.20.4-buster AS build
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat -gcflags "all=-N -l" main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
# FROM gcr.io/distroless/base:debug

WORKDIR /app

COPY --from=build /app/chat /app/chat
COPY --from=build ./app/app.env /app/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/chat", "server", "run"]