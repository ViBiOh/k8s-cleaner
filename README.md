# k8s-cleaner

Update the TTLSecondsAfterFinished of a Kubernetes job if succeeded.

[![Build](https://github.com/ViBiOh/k8s-cleaner/workflows/Build/badge.svg)](https://github.com/ViBiOh/k8s-cleaner/actions)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_k8s-cleaner&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_k8s-cleaner)

## Getting started

Golang binary is built with static link. You can download it directly from the [GitHub Release page](https://github.com/ViBiOh/k8s-cleaner/releases) or build it by yourself by cloning this repo and running `make`.

A Docker image is available for `amd64`, `arm` and `arm64` platforms on Docker Hub: [vibioh/k8s-cleaner](https://hub.docker.com/r/vibioh/k8s-cleaner/tags).

You can configure app by passing CLI args or environment variables (cf. [Usage](#usage) section). CLI override environment variables.

## Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of k8s-cleaner:
  -graceDuration duration
        [http] Grace duration when SIGTERM received {K8S_CLEANER_GRACE_DURATION} (default 30s)
  -jobDuration duration
        [job] TTL Duration after succeeded {K8S_CLEANER_JOB_DURATION} (default 2m0s)
  -jobLabel string
        [job] Label selector for jobs {K8S_CLEANER_JOB_LABEL} (default "k8s-cleaner=true")
  -jobNamespace string
        [job] Namespace to watch (blank for all) {K8S_CLEANER_JOB_NAMESPACE} (default "default")
  -k8sConfig string
        [k8s] Path to kubeconfig file {K8S_CLEANER_K8S_CONFIG} (default "${HOME}/.kube/config")
  -loggerJson
        [logger] Log format as JSON {K8S_CLEANER_LOGGER_JSON}
  -loggerLevel string
        [logger] Logger level {K8S_CLEANER_LOGGER_LEVEL} (default "INFO")
  -loggerLevelKey string
        [logger] Key for level in JSON {K8S_CLEANER_LOGGER_LEVEL_KEY} (default "level")
  -loggerMessageKey string
        [logger] Key for message in JSON {K8S_CLEANER_LOGGER_MESSAGE_KEY} (default "message")
  -loggerTimeKey string
        [logger] Key for timestamp in JSON {K8S_CLEANER_LOGGER_TIME_KEY} (default "time")
  -okStatus int
        [http] Healthy HTTP Status code {K8S_CLEANER_OK_STATUS} (default 204)
  -prometheusAddress string
        [prometheus] Listen address {K8S_CLEANER_PROMETHEUS_ADDRESS}
  -prometheusCert string
        [prometheus] Certificate file {K8S_CLEANER_PROMETHEUS_CERT}
  -prometheusGzip
        [prometheus] Enable gzip compression of metrics output {K8S_CLEANER_PROMETHEUS_GZIP}
  -prometheusIdleTimeout duration
        [prometheus] Idle Timeout {K8S_CLEANER_PROMETHEUS_IDLE_TIMEOUT} (default 10s)
  -prometheusIgnore string
        [prometheus] Ignored path prefixes for metrics, comma separated {K8S_CLEANER_PROMETHEUS_IGNORE}
  -prometheusKey string
        [prometheus] Key file {K8S_CLEANER_PROMETHEUS_KEY}
  -prometheusPort uint
        [prometheus] Listen port (0 to disable) {K8S_CLEANER_PROMETHEUS_PORT} (default 9090)
  -prometheusReadTimeout duration
        [prometheus] Read Timeout {K8S_CLEANER_PROMETHEUS_READ_TIMEOUT} (default 5s)
  -prometheusShutdownTimeout duration
        [prometheus] Shutdown Timeout {K8S_CLEANER_PROMETHEUS_SHUTDOWN_TIMEOUT} (default 5s)
  -prometheusWriteTimeout duration
        [prometheus] Write Timeout {K8S_CLEANER_PROMETHEUS_WRITE_TIMEOUT} (default 10s)
  -url string
        [alcotest] URL to check {K8S_CLEANER_URL}
  -userAgent string
        [alcotest] User-Agent for check {K8S_CLEANER_USER_AGENT} (default "Alcotest")
```
