### multistage dockerfile ###

### Build Stage ###

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang as builder

# Set GOPATH for container
ENV $GOPATH=/home/rross/go

# Set working directory for the container
WORKDIR $GOPATH/src/github.com/nazufel/raepublishing-website-api

# Copy the local package files to the container's workspace.
COPY . $GOPATH/src/github.com/nazufel/raepublishing-website-api

# Build the the app
RUN go install github.com/nazufel/raepublishing-website-api

WORKDIR /bin/

COPY --from=builder /bin/raepublishing-website-api .

### Final running stage ###
FROM alpine
# Run the outyet command by default when the container starts.
CMD ["./go/bin/raepublishing-website-api"]

# Document that the service listens on port 8080.
EXPOSE 3000
