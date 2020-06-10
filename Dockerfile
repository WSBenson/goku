# First stage of multi-stage build
FROM golang:1.14.3 as stageOne

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


# Move to working directory /build
WORKDIR /src

# Copy the code into the container
ADD . /src

# install go
RUN go install


# Second stage of multi-stage build
FROM alpine:latest

# Copy files from the first stage into the current
COPY --from=stageOne /go/bin/goku /bin/goku
COPY --from=stageOne /src/mapping/mapping.json /etc/mapping.json

# Export default port
EXPOSE 3000

# Run to start the server
CMD ["goku", "app"]