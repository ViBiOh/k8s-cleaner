package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/alcotest"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/prometheus"
	"github.com/ViBiOh/httputils/v4/pkg/server"
	"github.com/ViBiOh/k8s-cleaner/pkg/job"
	"github.com/ViBiOh/k8s-cleaner/pkg/k8s"
)

func main() {
	fs := flag.NewFlagSet("k8s-cleaner", flag.ExitOnError)
	fs.Usage = flags.Usage(fs)

	promServerConfig := server.Flags(fs, "prometheus", flags.NewOverride("Port", uint(9090)), flags.NewOverride("IdleTimeout", 10*time.Second), flags.NewOverride("ShutdownTimeout", 5*time.Second))
	healthConfig := health.Flags(fs, "")

	alcotestConfig := alcotest.Flags(fs, "")
	prometheusConfig := prometheus.Flags(fs, "prometheus", flags.NewOverride("Gzip", false))
	loggerConfig := logger.Flags(fs, "logger")

	k8sConfig := k8s.Flags(fs, "k8s")
	jobConfig := job.Flags(fs, "job")

	if err := fs.Parse(os.Args[1:]); err != nil {
		slog.Error("parse flags", "err", err)
		os.Exit(1)
	}

	alcotest.DoAndExit(alcotestConfig)
	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	ctx := context.Background()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:9999", http.DefaultServeMux))
	}()

	promServer := server.New(promServerConfig)
	prometheusApp := prometheus.New(prometheusConfig)
	healthApp := health.New(healthConfig)

	k8sClient, err := k8s.New(k8sConfig)
	if err != nil {
		slog.Error("k8s client", "err", err)
		os.Exit(1)
	}

	jobApp := job.New(jobConfig, k8sClient)

	go jobApp.Start(healthApp.Done(ctx))
	endCtx := healthApp.End(ctx)

	go promServer.Start(endCtx, "prometheus", prometheusApp.Handler())

	healthApp.WaitForTermination(jobApp.Done())
	server.GracefulWait(promServer.Done(), jobApp.Done())
}
