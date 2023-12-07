FROM golang:1.21.4-alpine as build
WORKDIR /app
COPY . .
RUN GGOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o test/authentication/build/authentication test/authentication/cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/test/authentication/build/authentication .
CMD [ "./authentication" ]
