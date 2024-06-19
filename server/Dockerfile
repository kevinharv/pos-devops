ARG GO_VERSION=1.22.3
ARG ALPINE_VERSION=3.19
ARG APP_VERSION=0.1.0
ARG APP_NAME="pos-server"

# Builder image
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add openssl

# Copy source files and build
WORKDIR /src
COPY go.mod .
RUN go mod download

# Generate SSL Certificate
RUN mkdir -p /certs && openssl req -x509 -newkey rsa:4096 -keyout /certs/key.pem -out /certs/cert.pem -sha256 -days 365 -nodes -subj "/CN=localhost"

COPY . .
RUN go build -o /bin/server .

# Runner image - copy binary and start
FROM alpine:${ALPINE_VERSION} AS runner

LABEL app.name=${APP_NAME}
LABEL app.version=${APP_VERSION}

RUN apk add curl
COPY --from=builder /bin/server /bin/server
COPY --from=builder /certs/key.pem /certs/cert.pem /bin/

EXPOSE 8080
ENTRYPOINT [ "/bin/server" ]