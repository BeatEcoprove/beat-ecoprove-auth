ARG GO_VERSION=1.23.2
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

ARG TARGETARCH=amd64

COPY . /src

RUN go generate ./...

RUN --mount=type=cache,target=/go/pkg/mod/ \
  CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server cmd/identity-service/main.go

FROM alpine:latest AS final

ARG UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  appuser

RUN mkdir -p /app/www && chown -R appuser:appuser /app && chmod -R 755 /app

USER appuser

COPY --from=build /bin/server /app

WORKDIR /app

EXPOSE 3000

ENTRYPOINT [ "./server" ]
