package job

import (
	"context"
	"encoding/json"
	"flag"
	"strings"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type jobPatch struct {
	Spec struct {
		TTLSecondsAfterFinished int `json:"ttlSecondsAfterFinished"`
	} `json:"spec"`
}

// App of package
type App struct {
	k8s       *kubernetes.Clientset
	done      chan struct{}
	namespace string
	label     string
	payload   []byte
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
		namespace: flags.String(fs, prefix, "job", "Namespace", "Namespace to watch (blank for all)", "default", overrides),
		label:     flags.String(fs, prefix, "job", "Label", "Label selector for jobs", "k8s-cleaner=true", overrides),
		duration:  flags.Duration(fs, prefix, "job", "Duration", "TTL Duration after succeeded", time.Minute*2, overrides),
	}
}

// New creates new App from Config
func New(config Config, k8s *kubernetes.Clientset) App {
	var patch jobPatch
	patch.Spec.TTLSecondsAfterFinished = int((*config.duration).Seconds())

	payload, err := json.Marshal(patch)
	logger.Fatal(err)

	return App{
		k8s:       k8s,
		namespace: strings.TrimSpace(*config.namespace),
		label:     strings.TrimSpace(*config.label),
		payload:   payload,
		done:      make(chan struct{}),
	}
}

// Done close when work is over
func (a App) Done() <-chan struct{} {
	return a.done
}

// Start listening kubernetes event
func (a App) Start(ctx context.Context) {
	defer close(a.done)

	for {
		if a.watchJobs(ctx) {
			return
		}
	}
}

func (a App) watchJobs(ctx context.Context) bool {
	watcher, err := a.k8s.BatchV1().Jobs(a.namespace).Watch(ctx, v1.ListOptions{
		LabelSelector: a.label,
		Watch:         true,
	})
	logger.Fatal(err)

	logger.Info("Listening jobs in `%s` namespace with `%s` label selector", a.namespace, a.label)

	defer watcher.Stop()

	results := watcher.ResultChan()

	for {
		select {
		case <-ctx.Done():
			return true
		case event, ok := <-results:
			if !ok {
				return false
			}

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

			logger.Info("Updating TTLSecondsAfterFinished for %s/%s", job.Namespace, job.Name)

			if _, err := a.k8s.BatchV1().Jobs(job.Namespace).Patch(ctx, job.Name, types.MergePatchType, a.payload, v1.PatchOptions{}); err != nil {
				logger.Error("patch job `%s/%s`: %s", job.Namespace, job.Name, err)
			}
		}
	}
}
