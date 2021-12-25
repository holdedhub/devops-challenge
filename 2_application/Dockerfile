# build stage
FROM golang:1.17-alpine AS build
ARG SOURCE
WORKDIR /go/src/app
ADD $SOURCE .
RUN go mod init && go build -o /app

# final stage
FROM alpine:3.15
WORKDIR /
COPY --from=build /app .
COPY static ./static
EXPOSE 8080
USER nobody
CMD ["/app"]
