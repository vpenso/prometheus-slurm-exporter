ARG ARCH="amd64"
ARG GOARCH="amd64"
ARG GOOS="linux"
FROM docker.io/golang:1.17-bullseye AS build

# In the future, this should be refactored using HEREDOC
# However, not yet fully supported by some tools
RUN set -ex; \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends -y \
    \
    slurm-wlm

COPY . /go/src
WORKDIR src

RUN set -ex; \
    go build -v -o /go/bin/prometheus-slurm-exporter

FROM docker.io/${ARCH}/debian:bullseye-slim

# In the future, this should be refactored using HEREDOC
# However, not yet fully supported by some tools
RUN set -ex; \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends -y \
    \
    slurm-wlm \
    \
    && apt-get -y autoclean; apt-get -y autoremove; \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /go/bin/prometheus-slurm-exporter /bin/slurm_exporter

EXPOSE 8080
USER   nobody
# The Slurm and Munge files that need to be present at runtime
VOLUME /etc/slurm/slurm.conf /etc/munge/munge.key /run/munge/munge.socket.2
ENTRYPOINT  ["/bin/slurm_exporter"]
