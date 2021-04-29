# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build cmd/main.go

# Copy config files to Working Directory
RUN cp cmd/*.json ./ 

# Expose port 8080 to the outside world
EXPOSE 8080

RUN chmod +x ./main

# Command to run the executable
CMD ["./main","-env","docker"]

# docker build -t event-auth-docker .
# docker run -d -p 5555:8080 event-auth-docker