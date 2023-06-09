FROM golang:1.20.4-alpine as build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w"  -o /URand

FROM scratch as app
COPY --from=build /URand /URand
ENTRYPOINT [ "/URand" ]