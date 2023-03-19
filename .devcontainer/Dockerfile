FROM mcr.microsoft.com/devcontainers/go:0-1.20

# RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
#     make grep

COPY .ci/devtools.sh devtools.sh
RUN sh devtools.sh
