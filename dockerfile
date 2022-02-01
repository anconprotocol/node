
# We build from the ubuntu long term support image
FROM ubuntu:20.04
LABEL description="This is a custom Docker Image for Ancon Protocol Node Server."

# Package updates and variable configurations
ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.17.6

RUN apt-get update
RUN apt install -y curl
RUN apt install -y build-essential
RUN curl -OL https://golang.org/dl/go1.17.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.17.6.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# We create a directory called app where our application will live.
RUN mkdir /app

# Make /app the work directory
WORKDIR /app

# Copy the go.mod go.sum and init.sh to the /app folder which is our application.
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./init.sh /app/init.sh
COPY . .
# Instaling go packages
RUN go mod tidy
RUN go build main.go

# Giving executable permissions to the script
RUN chmod +x main

# Exposes port 7788 where Ancon Protocol Node will be run (this is the port of the server inside the docker container)
EXPOSE 7788

# Run init script

CMD ./main --ipfshost=$IPFSHOST  --rootkey=$ROOTKEY --peeraddr $PEERADDR  --cors $CORS  --origins=$ORIGINS --quic $QUIC --tlscert=$TLSCERT --tlskey=$TLSKEY

