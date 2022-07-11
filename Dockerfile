FROM vibioh/scratch

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT [ "/k8s-suicide-job" ]

ARG VERSION
ENV VERSION=${VERSION}

ARG TARGETOS
ARG TARGETARCH

COPY release/k8s-suicide-job_${TARGETOS}_${TARGETARCH} /k8s-suicide-job
