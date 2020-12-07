# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.13 as builder

# Create and change to the app directory.
WORKDIR /go-samples/kms/

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o kmsservice

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

RUN mkdir -p /tmp/kms   # for logging log
RUN mkdir -p /tmp/config/

# env variable
ENV KMS_PORT=8280
ENV KMS_PROJECT_ID=kms-idm-232323
ENV KMS_LOCATION=global
ENV KMS_LOGGING_LEVEL=debug
ENV KMS_LOG_PATH_FILENAME=/tmp/kms/kms.log
ENV KMS_CREDENTIAL_PATH_FILENAME=/tmp/kms/config/Credential.json

# Copy Config file
COPY --from=builder /kms/misc/config  /tmp/kms/config/

# Copy the binary to the production image from the builder stage.
COPY --from=builder /kms/kmsservice /kmsservice

# Expose port 8080 to the outside world
EXPOSE 8280

# Run the web service on container startup.
CMD ["./kmsservice"]
