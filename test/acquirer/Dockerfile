FROM golang:1.21.4-alpine as build
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o test/acquirer/build/acquirer test/acquirer/cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/test/acquirer/build/acquirer .
CMD [ "./acquirer" ]
