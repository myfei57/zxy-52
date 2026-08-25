FROM golang:1.23.12 AS build
WORKDIR /src
ENV GOPROXY=off GOSUMDB=off
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN go build -mod=vendor -o /out/wastegen ./cmd/wastegen

FROM golang:1.23.12
WORKDIR /app
ENV GOPROXY=off GOSUMDB=off
COPY --from=build /out/wastegen /app/wastegen
EXPOSE 8901
CMD ["/app/wastegen", "-addr", "0.0.0.0:8901", "-root", "/app/data"]
