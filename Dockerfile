# First stage of multi-stage build
FROM golang:1.14.3

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


# Move to working directory /build
WORKDIR /build

# Copy the code into the container
ADD . /build

# Build the application, one line to minimize the # of layers
RUN go mod download && go build -o main .




# Second stage of multi-stage build
FROM alpine:latest

# Get needed certificates
RUN apk --no-cache add ca-certificates

# Move to working directory /root
WORKDIR /root/

# Copy the built main from the first stage into the WORKDIR
COPY --from=0 /build/main .

# Export necessary port
EXPOSE 3000

# Run main and the server will be up
CMD ["./main"]