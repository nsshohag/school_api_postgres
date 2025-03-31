FROM ubuntu:latest

LABEL "Author"="Sadat"
LABEL "Project" = "Web-Server"

RUN apt update && apt-get install -y wget

RUN wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz && rm go1.24.1.linux-amd64.tar.gz


# Set Go environment variables
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin


WORKDIR /school_server
COPY go.mod go.sum /school_server/

RUN go mod download

COPY . ./

RUN go build -o server main.go

# Expose the port your webserver runs on
EXPOSE 8080

# CMD ["go","run","main.go"]
CMD ["./server"]