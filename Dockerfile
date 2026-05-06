ARG BUILDER_IMAGE="docker.io/golang:1.26.2-alpine3.23@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1"
ARG BASE_IMAGE="gcr.io/distroless/static-debian13:nonroot@sha256:e3f945647ffb95b5839c07038d64f9811adf17308b9121d8a2b87b6a22a80a39"

FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} AS builder


WORKDIR /usr/src

COPY . .

ARG RELEASE="unknown"
ARG TARGETARCH
ARG TARGETOS

ENV CGO_ENABLED="0"
ENV GOARCH="$TARGETARCH"
ENV GOOS="$TARGETOS"

RUN --mount=type=cache,target=/go/pkg/mod \
	go build -ldflags="-s -w -X 'go.yunus-emre.dev/url-shortaner/pkg/version.Version=${RELEASE}'" -trimpath -o ./api-server ./cmd/api-server/main.go

FROM ${BASE_IMAGE}

COPY --chown=65532:65532 --from=builder --link /usr/src/api-server /usr/local/bin/api-server

USER 65532:65532

STOPSIGNAL SIGINT

ENTRYPOINT [ "/usr/local/bin/api-server" ]

CMD [ "/usr/local/bin/api-server" ]
