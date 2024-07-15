FROM golang:1.22-alpine

COPY go.mod ./
COPY go.sum ./

## Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY *.go ./
COPY config ./
COPY controllers ./
COPY middlewares ./
COPY models ./
COPY routers ./
COPY services ./
COPY static ./
COPY templates ./

RUN go build -o /portal_app

EXPOSE 8080

CMD [ "/portal_app" ]