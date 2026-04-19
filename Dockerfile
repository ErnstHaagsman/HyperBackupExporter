FROM golang:1.26 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /hyperbackup-exporter


FROM gcr.io/distroless/base-debian13 AS build-release-stage

WORKDIR /

COPY --from=build-stage /hyperbackup-exporter /hyperbackup-exporter

EXPOSE 6533

USER nonroot:nonroot

ENTRYPOINT ["/hyperbackup-exporter"]