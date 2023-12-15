FROM vibioh/scratch

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT [ "/k8s-cleaner" ]

ARG VERSION
ENV VERSION ${VERSION}

ARG GIT_SHA
ENV GIT_SHA ${GIT_SHA}

ARG TARGETOS
ARG TARGETARCH

COPY release/k8s-cleaner_${TARGETOS}_${TARGETARCH} /k8s-cleaner
