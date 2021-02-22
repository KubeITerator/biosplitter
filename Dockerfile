FROM golang:alpine AS builder
# Install git.
RUN apk update && apk add --no-cache git
# Add certs
RUN apk --no-cache add ca-certificates

# Create appuser.
ENV USER=biokubeuser
ENV UID=10001
# Statically link libc
ENV CGO_ENABLED=0
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/go/biosplitter/
COPY . .
# Get dependencies
RUN go get -d -v
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bio-splitter
# 2. Build small image
FROM scratch
# Copy certificates for S3
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy our static executable.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /bio-splitter /bio-splitter
# Run the hello binary.
USER biokubeuser:biokubeuser
ENTRYPOINT ["/bio-splitter"]