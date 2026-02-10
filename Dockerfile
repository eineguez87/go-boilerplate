FROM golang:1.25-alpine

WORKDIR /app

# Required for go install + HTTPS
RUN apk add --no-cache \
    git \
    ca-certificates \
    build-base

# Install air
RUN go install github.com/air-verse/air@latest

# Ensure GOPATH/bin is on PATH
ENV PATH="/go/bin:${PATH}"

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN ls -la

EXPOSE 8080

CMD ["air"]



