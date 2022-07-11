# k8s-suicide-job

Update the TTLSecondsAfterFinished of a Kubernetes job if succeeded.

[![Build](https://github.com/ViBiOh/k8s-suicide-job/workflows/Build/badge.svg)](https://github.com/ViBiOh/k8s-suicide-job/actions)
[![codecov](https://codecov.io/gh/ViBiOh/k8s-suicide-job/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/k8s-suicide-job)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_k8s-suicide-job&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_k8s-suicide-job)

## Getting started

Golang binary is built with static link. You can download it directly from the [Github Release page](https://github.com/ViBiOh/k8s-suicide-job/releases) or build it by yourself by cloning this repo and running `make`.

A Docker image is available for `amd64`, `arm` and `arm64` platforms on Docker Hub: [vibioh/k8s-suicide-job](https://hub.docker.com/r/vibioh/k8s-suicide-job/tags).

You can configure app by passing CLI args or environment variables (cf. [Usage](#usage) section). CLI override environment variables.

## Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of k8s-suicide-job:
  -duration duration
        [job] TTL Duration after finished {K8S_SUICIDE_JOB_DURATION} (default 2m0s)
  -k8sConfig string
        [k8s] Path to kubeconfig file {K8S_SUICIDE_JOB_K8S_CONFIG} (default "/Users/macbook/.kube/config")
  -loggerJson
        [logger] Log format as JSON {K8S_SUICIDE_JOB_LOGGER_JSON}
  -loggerLevel string
        [logger] Logger level {K8S_SUICIDE_JOB_LOGGER_LEVEL} (default "INFO")
  -loggerLevelKey string
        [logger] Key for level in JSON {K8S_SUICIDE_JOB_LOGGER_LEVEL_KEY} (default "level")
  -loggerMessageKey string
        [logger] Key for message in JSON {K8S_SUICIDE_JOB_LOGGER_MESSAGE_KEY} (default "message")
  -loggerTimeKey string
        [logger] Key for timestamp in JSON {K8S_SUICIDE_JOB_LOGGER_TIME_KEY} (default "time")
  -name string
        [job] Name of the job {K8S_SUICIDE_JOB_NAME}
  -namespace string
        [job] Namespace of the job {K8S_SUICIDE_JOB_NAMESPACE}
```
