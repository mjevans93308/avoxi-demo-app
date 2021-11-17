# for building (from Dockerfile directory)
# docker build -f Dockerfile -t geolocate-ip-demo-app:0.0.1 .
# for running
# docker run --env-file ./geolocate-ip-demo-app.env --publish 8080:8080 geolocate-ip-demo-app:0.0.1
# ssh'ing to running container:
# docker ps (to grab container id)
# docker exec -it <container_id> /bin/sh

# multistage dockerfile to minimise our final image size
FROM golang:alpine AS builder

ENV GO111MODULE=on

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o geolocate-ip-demo-app main/main.go

FROM alpine:latest
COPY --from=builder /build/geolocate-ip-demo-app .

# copy in env file so we avoid exposing sensitive env vars during docker build/run process
COPY --from=builder /build/geolocate-ip-demo-app.env .
RUN apk --no-cache add ca-certificates

EXPOSE 8080
# CMD corresponds to the cobra commands created in main/main.go
CMD ["./geolocate-ip-demo-app", "serve"]