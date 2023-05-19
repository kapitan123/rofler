# syntax=docker/dockerfile:1

# ===================================================================================
# === Stage 1:Builder container =====================================================
# ===================================================================================
FROM golang:1.19-alpine AS builder

WORKDIR /build

# Fetch modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy Go source files
COPY convertor/config/ ./convertor/config/
COPY convertor/api/ ./convertor/api/ 
COPY convertor/internal/ ./convertor/internal/

# All GCE instances are x86_64 linux based machines. Dynamic linking is disabled because it will be copied to scratch.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o server github.com/kapitan123/telegrofler/convertor/api

# ===================================================================================
# === Stage 2: Create a run container =======================================
# ===================================================================================
FROM golang:1.19-alpine 

RUN apk add --no-cache python3 py3-pip
&& apk add --no-cache ffmpeg
&& pip3 install youtube-dl

# Add certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Add binary
COPY --from=builder /build/server /server 

EXPOSE 9001

ENTRYPOINT [ "/server"]