FROM golang:1.21.3-alpine as build
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o build/acquirer cmd/acquirer/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/build/acquirer .
CMD [ "./acquirer" ]
