package main

import (
	"context"
	"flag"
	"os"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/alcotest"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/pprof"
	"github.com/ViBiOh/httputils/v4/pkg/server"
	"github.com/ViBiOh/httputils/v4/pkg/telemetry"
	"github.com/ViBiOh/k8s-cleaner/pkg/job"
	"github.com/ViBiOh/k8s-cleaner/pkg/k8s"
)

func main() {
	fs := flag.NewFlagSet("k8s-cleaner", flag.ExitOnError)
	fs.Usage = flags.Usage(fs)

	healthConfig := health.Flags(fs, "")

	alcotestConfig := alcotest.Flags(fs, "")
	telemetryConfig := telemetry.Flags(fs, "telemetry")
	pprofConfig := pprof.Flags(fs, "pprof")
	loggerConfig := logger.Flags(fs, "logger")

	k8sConfig := k8s.Flags(fs, "k8s")
	jobConfig := job.Flags(fs, "job")

	_ = fs.Parse(os.Args[1:])

	alcotest.DoAndExit(alcotestConfig)

	logger.Init(loggerConfig)

	ctx := context.Background()

	healthApp := health.New(ctx, healthConfig)

	telemetryApp, err := telemetry.New(ctx, telemetryConfig)
	logger.FatalfOnErr(ctx, err, "create telemetry")

	defer telemetryApp.Close(ctx)

	logger.AddOpenTelemetryToDefaultLogger(telemetryApp)

	service, version, env := telemetryApp.GetServiceVersionAndEnv()
	pprofService := pprof.New(pprofConfig, service, version, env)

	go pprofService.Start(healthApp.DoneCtx())

	k8sClient, err := k8s.New(k8sConfig)
	logger.FatalfOnErr(ctx, err, "k8s client")

	jobApp := job.New(jobConfig, k8sClient)

	go jobApp.Start(healthApp.DoneCtx())

	healthApp.WaitForTermination(jobApp.Done())
	server.GracefulWait(jobApp.Done())
}
