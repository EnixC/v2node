FROM alpine:3.22 AS assets

ARG GEOIP_URL=https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/geoip.dat
ARG GEOSITE_URL=https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/geosite.dat

RUN apk add --no-cache ca-certificates curl \
    && mkdir -p /assets \
    && curl --fail --location --retry 3 --output /assets/geoip.dat "${GEOIP_URL}" \
    && curl --fail --location --retry 3 --output /assets/geosite.dat "${GEOSITE_URL}"

FROM golang:1.27.1-alpine AS builder

ARG VERSION=dev

WORKDIR /app
COPY go.mod go.sum ./
RUN GOEXPERIMENT=jsonv2 go mod download
COPY . .
RUN CGO_ENABLED=0 GOEXPERIMENT=jsonv2 go build \
    -trimpath \
    -ldflags "-X github.com/wyx2685/v2node/cmd.version=${VERSION} -s -w -buildid=" \
    -o /out/v2node .

FROM alpine:3.22

ARG VERSION=dev
ARG REVISION=unknown
ARG SOURCE=https://github.com/wyx2685/v2node

LABEL org.opencontainers.image.title="v2node" \
      org.opencontainers.image.description="V2Board node backend" \
      org.opencontainers.image.source="${SOURCE}" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${REVISION}"

RUN apk add --no-cache ca-certificates tzdata \
    && mkdir -p /etc/v2node /usr/local/share/xray

COPY --from=builder /out/v2node /usr/local/bin/v2node
COPY --from=assets /assets/geoip.dat /usr/local/share/xray/geoip.dat
COPY --from=assets /assets/geosite.dat /usr/local/share/xray/geosite.dat

ENV TZ=Asia/Shanghai \
    XRAY_LOCATION_ASSET=/usr/local/share/xray \
    V2NODE_LOG_LEVEL=warning \
    V2NODE_TIMEOUT=15 \
    V2NODE_RETRY_COUNT=1

WORKDIR /etc/v2node
VOLUME ["/etc/v2node"]
STOPSIGNAL SIGTERM

ENTRYPOINT ["v2node", "server"]
