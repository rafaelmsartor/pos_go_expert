from golang:1.22 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather_api

FROM scratch
WORKDIR /app
COPY --from=build /app/weather_api .
ENTRYPOINT ["./weather_api"]