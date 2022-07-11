package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/k8s-suicide-job/pkg/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	fs := flag.NewFlagSet("k8s-suicide-job", flag.ExitOnError)

	loggerConfig := logger.Flags(fs, "logger")
	k8sConfig := k8s.Flags(fs, "k8s")

	namespace := flags.String(fs, "", "job", "namespace", "Namespace of the job", "default", nil)
	name := flags.String(fs, "", "job", "name", "Name of the job", "", nil)
	duration := flags.Duration(fs, "", "job", "duration", "TTL Duration after finished", time.Minute*2, nil)

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	k8sClient, err := k8s.New(k8sConfig)
	logger.Fatal(err)

	job, err := k8sClient.BatchV1().Jobs(*namespace).Get(context.Background(), *name, v1.GetOptions{})
	logger.Fatal(err)

	if job.Status.Succeeded != 1 {
		logger.Warn("Job didn't succeeded, nothing to do")
		return
	}

	ttlSeconds := int32(duration.Seconds())
	job.Spec.TTLSecondsAfterFinished = &ttlSeconds
	_, err = k8sClient.BatchV1().Jobs(*namespace).Update(context.Background(), job, v1.UpdateOptions{})
	logger.Fatal(err)
}
