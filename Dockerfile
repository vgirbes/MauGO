# Base off official golang image
FROM golang:latest

# Create app directory
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# Bundle app source
COPY src /go/src/app

# Download and install any required third party dependencies into the container.
RUN go get github.com/pilu/fresh
RUN go get -u github.com/go-sql-driver/mysql
RUN go get github.com/julienschmidt/httprouter
RUN go-wrapper download
RUN go-wrapper install

# Set the PORT environment variable inside the container
ENV PORT 3000

# Expose port 8080 to the host so we can access our application
EXPOSE 3000

# Now tell Docker what command to run when the container starts
CMD ["fresh"]
