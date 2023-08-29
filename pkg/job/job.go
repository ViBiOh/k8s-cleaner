package job

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/ViBiOh/flags"
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

type Config struct {
	Namespace string
	Label     string
	Duration  time.Duration
}

func Flags(fs *flag.FlagSet, prefix string, overrides ...flags.Override) *Config {
	var config Config

	flags.New("Namespace", "Namespace to watch (blank for all)").Prefix(prefix).DocPrefix("job").StringVar(fs, &config.Namespace, "default", overrides)
	flags.New("Label", "Label selector for jobs").Prefix(prefix).DocPrefix("job").StringVar(fs, &config.Label, "k8s-cleaner=true", overrides)
	flags.New("Duration", "TTL Duration after succeeded").Prefix(prefix).DocPrefix("job").DurationVar(fs, &config.Duration, time.Minute*2, overrides)

	return &config
}

type Service struct {
	k8s       *kubernetes.Clientset
	done      chan struct{}
	namespace string
	label     string
	payload   []byte
}

func New(config *Config, k8s *kubernetes.Clientset) Service {
	var patch jobPatch
	patch.Spec.TTLSecondsAfterFinished = int((config.Duration).Seconds())

	payload, err := json.Marshal(patch)
	if err != nil {
		slog.Error("marshal json", "err", err)
		os.Exit(1)
	}

	return Service{
		k8s:       k8s,
		namespace: config.Namespace,
		label:     config.Label,
		payload:   payload,
		done:      make(chan struct{}),
	}
}

func (a Service) Done() <-chan struct{} {
	return a.done
}

func (a Service) Start(ctx context.Context) {
	defer close(a.done)

	for {
		if a.watchJobs(ctx) {
			return
		}
	}
}

func (a Service) watchJobs(ctx context.Context) bool {
	watcher, err := a.k8s.BatchV1().Jobs(a.namespace).Watch(ctx, v1.ListOptions{
		LabelSelector: a.label,
		Watch:         true,
	})
	if err != nil {
		slog.Error("watch jobs", "err", err)
		os.Exit(1)
	}

	slog.Info("Listening jobs", "namespace", a.namespace, "label", a.label)

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

			slog.Info("Updating TTLSecondsAfterFinished", "namespace", job.Namespace, "name", job.Name)

			if _, err := a.k8s.BatchV1().Jobs(job.Namespace).Patch(ctx, job.Name, types.MergePatchType, a.payload, v1.PatchOptions{}); err != nil {
				slog.Error("patch job", "err", err, "namespace", job.Namespace, "name", job.Name)
			}
		}
	}
}
