FROM golang:1.19.3-alpine
WORKDIR /user-balance

COPY ./ ./
RUN go mod download

RUN go build -o bin/user-balance github.com/bifrurcated/user-balance/cmd/main
EXPOSE 1234
CMD ["bin/user-balance"]