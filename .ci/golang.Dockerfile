ARG BASE_IMG=golang:1.19-alpine
FROM ${BASE_IMG} as golang

RUN apk add -q make 
# build-base

COPY devtools.sh devtools.sh
RUN sh devtools.sh
