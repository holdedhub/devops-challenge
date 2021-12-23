FROM golang:1.17-alpine
# Copy source code
WORKDIR /go/src/app
#COPY app .
ADD https://raw.githubusercontent.com/mikakatua/devops-challenge/master/app/server.go .
# Compile and install the app
RUN go mod init && go install -v ./...

EXPOSE 8080

USER nobody

CMD ["app"]
