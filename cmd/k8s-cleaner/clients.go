package main

import (
	"context"
	"fmt"

	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/pprof"
	"github.com/ViBiOh/httputils/v4/pkg/request"
	"github.com/ViBiOh/httputils/v4/pkg/telemetry"
	"github.com/ViBiOh/k8s-cleaner/pkg/k8s"
	"k8s.io/client-go/kubernetes"
)

type clients struct {
	telemetry *telemetry.Service
	pprof     *pprof.Service
	health    *health.Service

	k8s *kubernetes.Clientset
}

func newClients(ctx context.Context, config configuration) (clients, error) {
	var output clients
	var err error

	logger.Init(ctx, config.logger)

	output.telemetry, err = telemetry.New(ctx, config.telemetry)
	if err != nil {
		return output, fmt.Errorf("telemetry: %w", err)
	}

	logger.AddOpenTelemetryToDefaultLogger(output.telemetry)
	request.AddOpenTelemetryToDefaultClient(output.telemetry.MeterProvider(), output.telemetry.TracerProvider())

	service, version, env := output.telemetry.GetServiceVersionAndEnv()
	output.pprof = pprof.New(config.pprof, service, version, env)

	output.k8s, err = k8s.New(config.k8s)
	if err != nil {
		return output, fmt.Errorf("k8s: %w", err)
	}

	output.health = health.New(ctx, config.health)

	return output, nil
}

func (c clients) Start() {
	go c.pprof.Start(c.health.DoneCtx())
}

func (c clients) Close(ctx context.Context) {
	c.telemetry.Close(ctx)
}
