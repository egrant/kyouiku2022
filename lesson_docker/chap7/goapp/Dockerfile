FROM golang:1.18.3-bullseye AS build
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o /server

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /server /server
EXPOSE 8888
USER nonroot:nonroot
ENTRYPOINT ["/server"]

