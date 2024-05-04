# k8s-cleaner

Update the TTLSecondsAfterFinished of a Kubernetes job if succeeded.

[![Build](https://github.com/ViBiOh/k8s-cleaner/workflows/Build/badge.svg)](https://github.com/ViBiOh/k8s-cleaner/actions)

## Getting started

Golang binary is built with static link. You can download it directly from the [GitHub Release page](https://github.com/ViBiOh/k8s-cleaner/releases) or build it by yourself by cloning this repo and running `make`.

A Docker image is available for `amd64`, `arm` and `arm64` platforms on Docker Hub: [vibioh/k8s-cleaner](https://hub.docker.com/r/vibioh/k8s-cleaner/tags).

You can configure app by passing CLI args or environment variables (cf. [Usage](#usage) section). CLI override environment variables.

## Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of k8s-cleaner:
  --graceDuration     duration  [http] Grace duration when signal received ${K8S_CLEANER_GRACE_DURATION} (default 30s)
  --jobDuration       duration  [job] TTL Duration after succeeded ${K8S_CLEANER_JOB_DURATION} (default 2m0s)
  --jobLabel          string    [job] Label selector for jobs ${K8S_CLEANER_JOB_LABEL} (default "k8s-cleaner=true")
  --jobNamespace      string    [job] Namespace to watch (blank for all, comma separated otherwise) ${K8S_CLEANER_JOB_NAMESPACE} (default "default")
  --k8sConfig         string    [k8s] Path to kubeconfig file ${K8S_CLEANER_K8S_CONFIG} (default "${HOME}/.kube/config")
  --loggerJson                  [logger] Log format as JSON ${K8S_CLEANER_LOGGER_JSON} (default false)
  --loggerLevel       string    [logger] Logger level ${K8S_CLEANER_LOGGER_LEVEL} (default "INFO")
  --loggerLevelKey    string    [logger] Key for level in JSON ${K8S_CLEANER_LOGGER_LEVEL_KEY} (default "level")
  --loggerMessageKey  string    [logger] Key for message in JSON ${K8S_CLEANER_LOGGER_MESSAGE_KEY} (default "msg")
  --loggerTimeKey     string    [logger] Key for timestamp in JSON ${K8S_CLEANER_LOGGER_TIME_KEY} (default "time")
  --okStatus          int       [http] Healthy HTTP Status code ${K8S_CLEANER_OK_STATUS} (default 204)
  --pprofAgent        string    [pprof] URL of the Datadog Trace Agent (e.g. http://datadog.observability:8126) ${K8S_CLEANER_PPROF_AGENT}
  --pprofPort         int       [pprof] Port of the HTTP server (0 to disable) ${K8S_CLEANER_PPROF_PORT} (default 0)
  --telemetryRate     string    [telemetry] OpenTelemetry sample rate, 'always', 'never' or a float value ${K8S_CLEANER_TELEMETRY_RATE} (default "always")
  --telemetryURL      string    [telemetry] OpenTelemetry gRPC endpoint (e.g. otel-exporter:4317) ${K8S_CLEANER_TELEMETRY_URL}
  --telemetryUint64             [telemetry] Change OpenTelemetry Trace ID format to an unsigned int 64 ${K8S_CLEANER_TELEMETRY_UINT64} (default true)
  --url               string    [alcotest] URL to check ${K8S_CLEANER_URL}
  --userAgent         string    [alcotest] User-Agent for check ${K8S_CLEANER_USER_AGENT} (default "Alcotest")
```
