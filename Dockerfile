#############################################
# Builder go
#############################################
FROM golang:1.22-alpine as builder
ARG APP_VERSION=v0.0.0

WORKDIR /app/build

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${APP_VERSION}" -o /app/build/akatran .

#############################################
# Builder go
#############################################
FROM alpine:3.20 as release

RUN adduser -D gorunner

USER gorunner

WORKDIR /app

COPY --chown=gorunner:gorunnner --from=builder /app/build/akatran /app/akatran

ENTRYPOINT [ "/app/akatran" ]
