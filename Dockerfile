ARG ARCH="amd64"
ARG GOARCH="amd64"
ARG GOOS="linux"
FROM docker.io/golang:1.21-bookworm AS build

COPY . /go/src
WORKDIR /go/src

RUN set -ex; \
    GOARCH=${GOARCH} GOOS=${GOOS} go build -v -o /go/bin/prometheus-slurm-exporter

FROM docker.io/${ARCH}/debian:bookworm-slim

# The exporter shells out to sinfo/squeue/sacct/sdiag/sshare at scrape time,
# so the Slurm client tools (and munge for authenticated calls) must be
# present in the runtime image.
RUN set -ex; \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends -y \
    \
    slurm-wlm \
    munge \
    \
    && apt-get -y autoclean; apt-get -y autoremove; \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /go/bin/prometheus-slurm-exporter /bin/slurm_exporter

EXPOSE 8080
USER   nobody
# The Slurm and Munge files that need to be present at runtime are
# bind-mounted at run time (VOLUME cannot be used for individual files):
#   docker run -d \
#     -v /etc/slurm/slurm.conf:/etc/slurm/slurm.conf:ro \
#     -v /etc/munge/munge.key:/etc/munge/munge.key:ro \
#     -v /run/munge/munge.socket.2:/run/munge/munge.socket.2 \
#     -p 8080:8080 prometheus-slurm-exporter
ENTRYPOINT  ["/bin/slurm_exporter"]
