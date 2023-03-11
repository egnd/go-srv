ARG BASE_IMG=alpine
FROM ${BASE_IMG} as base
RUN apk add -q tzdata ca-certificates

FROM scratch
ARG TARGETOS
ARG TARGETARCH
WORKDIR /
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY configs configs
COPY bin/${TARGETOS}-${TARGETARCH} /usr/local/bin/
ENTRYPOINT ["go-srv"]
