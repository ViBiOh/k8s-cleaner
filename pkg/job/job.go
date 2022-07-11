package job

import (
	"context"
	"flag"
	"strings"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// App of package
type App struct {
	k8s       *kubernetes.Clientset
	namespace string
	label     string
	duration  int32
	done      chan struct{}
}

// Config of package
type Config struct {
	namespace *string
	label     *string
	duration  *time.Duration
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string, overrides ...flags.Override) Config {
	return Config{
		namespace: flags.String(fs, prefix, "job", "Namespace", "Namespace to watch (blank for all)", "", overrides),
		label:     flags.String(fs, prefix, "job", "Label", "Label selector for jobs", "", overrides),
		duration:  flags.Duration(fs, prefix, "job", "Duration", "TTL Duration after succeeded", time.Minute*2, overrides),
	}
}

// New creates new App from Config
func New(config Config, k8s *kubernetes.Clientset) App {
	return App{
		k8s:       k8s,
		namespace: strings.TrimSpace(*config.namespace),
		label:     strings.TrimSpace(*config.label),
		duration:  int32((*config.duration).Seconds()),
		done:      make(chan struct{}),
	}
}

// Done close when work is over
func (a App) Done() <-chan struct{} {
	return a.done
}

// Start listening kubernetes event
func (a App) Start(done <-chan struct{}) {
	defer close(a.done)

	jobs, err := a.k8s.BatchV1().Jobs(a.namespace).Watch(context.Background(), v1.ListOptions{
		LabelSelector: a.label,
		Watch:         true,
	})
	logger.Fatal(err)

	logger.Info("Listening jobs in `%s` namespace with `%s` label selector", a.namespace, a.label)

	for event := range jobs.ResultChan() {
		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			continue
		}

		if job.Spec.TTLSecondsAfterFinished != nil && *job.Spec.TTLSecondsAfterFinished != 0 {
			continue
		}

		if job.Status.Succeeded != 1 {
			continue
		}

		logger.Info("Updating TTLSecondsAfterFinished to %d for %s/%s", a.duration, job.Namespace, job.Name)

		job.Spec.TTLSecondsAfterFinished = &a.duration
		if _, err = a.k8s.BatchV1().Jobs(job.Namespace).Update(context.Background(), job, v1.UpdateOptions{}); err != nil {
			logger.Error("unable to update job `%s/%s`: %s", job.Namespace, job.Name, err)
		}
	}
}
