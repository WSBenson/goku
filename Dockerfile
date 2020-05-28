FROM golang:1.14.3-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    PORT=3000


# Move to working directory /build
WORKDIR /build

# Copy the code into the container
ADD . /build

# Build the application
RUN go mod download
RUN go build -o main .

# Export necessary port
EXPOSE 3000

# Command to run when starting the container
CMD ["/build/main"]