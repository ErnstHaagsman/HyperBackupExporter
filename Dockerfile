FROM golang:1.26 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /hyperbackup-exporter


FROM alpine:3.23.4 AS build-release-stage

WORKDIR /

COPY --from=build-stage /hyperbackup-exporter /hyperbackup-exporter

EXPOSE 6533

USER 1000:1000

ENTRYPOINT ["/hyperbackup-exporter"]