FROM golang:1.20

# Work directory
WORKDIR /docker-go

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "cmd/main.go"]
EXPOSE 8087

# # Run
# CMD ["/chatbot"]
# TODO: The CMD should be updated in order to use a build binary version instead the current running one. 