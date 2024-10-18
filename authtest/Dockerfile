# syntax=docker/dockerfile:1

# Declare the Docker image that we want to base our container on
#
# We will build the app in the container to be the correct format
# so choose a Golang base.
FROM golang:1.23-bookworm

# Don't put all the code in the root directory, put it here instead
WORKDIR /source

# Pull down the Go modules we will need before bringing our Go source over
COPY go.* ./
RUN go mod download

# Bring the Go source into the container
COPY *.go ./

# Build the service
RUN go build -o /service

# Make port 8080 visible outside the container
EXPOSE 8080

# The command to run when the container starts
CMD [ "/service" ]