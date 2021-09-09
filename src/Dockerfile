FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /go/src/app

# Copy the code into the container
COPY . .

RUN go get github.com/cespare/reflex

RUN go mod download

# Run test
RUN go test ./...

# # Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /go/src/app/main .

WORKDIR /go/src/app

############################
# STEP 2 build a small image
############################
FROM scratch AS production

COPY --from=builder /dist/main /

# Command to run the executable
ENTRYPOINT [ "/main" ]