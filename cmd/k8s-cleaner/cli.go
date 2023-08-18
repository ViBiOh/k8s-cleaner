package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

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

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	alcotest.DoAndExit(alcotestConfig)

	logger.Init(loggerConfig)

	ctx := context.Background()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:9999", http.DefaultServeMux))
	}()

	healthApp := health.New(healthConfig)

	telemetryApp, err := telemetry.New(ctx, telemetryConfig)
	if err != nil {
		slog.Error("create telemetry", "err", err)
		os.Exit(1)
	}

	defer telemetryApp.Close(ctx)

	k8sClient, err := k8s.New(k8sConfig)
	if err != nil {
		slog.Error("k8s client", "err", err)
		os.Exit(1)
	}

	jobApp := job.New(jobConfig, k8sClient)

	go jobApp.Start(healthApp.Done(ctx))

	healthApp.WaitForTermination(jobApp.Done())
	server.GracefulWait(jobApp.Done())
}
