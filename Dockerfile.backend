# Create a first temporary image named "builder"
FROM golang:1.19.1 AS builder
# Copy Go code (in "builder")
WORKDIR /src/
COPY . .
# Build executables (in "builder")
RUN go build -o /app/webapi ./cmd/webapi
# Create final container
FROM debian:bookworm
# Inform Docker about which port is used
EXPOSE 3000
# Copy the executable from the "builder" image
WORKDIR /app/
COPY --from=builder /app/webapi ./
# Define volumes for database and image files
VOLUME /app/media/img
VOLUME /tmp
# Set the default program to our Go backend
# CMD ["/bin/ls", "-la", "/media/img"]
CMD ["/app/webapi"]

# To start docker: docker run --rm -p 3000:3000 -v /tmp/:/tmp -v ./media/img/:/app/media/img --user $(id -u):$(id -g) -it backend:latest