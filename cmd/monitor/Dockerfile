# Start from the latest golang base image
FROM golang:latest as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the application
WORKDIR /app/cmd/monitor
RUN go build

###########################################

# Build image from alpine
FROM alpine:latest  

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app /app

# Build Args
#ARG REPORT_DIR=/reports

# Create Report Directory
#WORKDIR ${REPORT_DIR}
#RUN mkdir -p ${REPORT_DIR}

# Declare volumes to mount
#VOLUME [${REPORT_DIR}]

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/app/cmd/monitor/monitor"] 
