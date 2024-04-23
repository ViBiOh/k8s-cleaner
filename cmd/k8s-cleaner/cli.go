package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/alcotest"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
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
	loggerConfig := logger.Flags(fs, "logger")

	k8sConfig := k8s.Flags(fs, "k8s")
	jobConfig := job.Flags(fs, "job")

	_ = fs.Parse(os.Args[1:])

	alcotest.DoAndExit(alcotestConfig)

	logger.Init(loggerConfig)

	ctx := context.Background()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:9999", http.DefaultServeMux))
	}()

	healthApp := health.New(ctx, healthConfig)

	telemetryApp, err := telemetry.New(ctx, telemetryConfig)
	logger.FatalfOnErr(ctx, err, "create telemetry")

	defer telemetryApp.Close(ctx)

	logger.AddOpenTelemetryToDefaultLogger(telemetryApp)

	k8sClient, err := k8s.New(k8sConfig)
	logger.FatalfOnErr(ctx, err, "k8s client")

	jobApp := job.New(jobConfig, k8sClient)

	go jobApp.Start(healthApp.DoneCtx())

	healthApp.WaitForTermination(jobApp.Done())
	server.GracefulWait(jobApp.Done())
}
