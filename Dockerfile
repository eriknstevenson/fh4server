# Stage 1
FROM golang as build
WORKDIR /fh4server/src
COPY go.mod .
COPY go.sum .
# Download dependencies and cache results.
RUN go mod download
COPY . .
# This setting is needed for the go binary to work properly on alpine.
ENV CGO_ENABLED=0
RUN ls
RUN go build -o ../bin/fh4server cmd/fh4server.go


# Stage 2
FROM ubuntu

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /fh4server/bin/fh4server /usr/bin/fh4server

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD ["fh4server", "-alsologtostderr", "-log_dir=logs/"]

EXPOSE 10001/udp
