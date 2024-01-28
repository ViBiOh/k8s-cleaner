package job

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/concurrent"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
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

	flags.New("Namespace", "Namespace to watch (blank for all, comma separated otherwise)").Prefix(prefix).DocPrefix("job").StringVar(fs, &config.Namespace, "default", overrides)
	flags.New("Label", "Label selector for jobs").Prefix(prefix).DocPrefix("job").StringVar(fs, &config.Label, "k8s-cleaner=true", overrides)
	flags.New("Duration", "TTL Duration after succeeded").Prefix(prefix).DocPrefix("job").DurationVar(fs, &config.Duration, time.Minute*2, overrides)

	return &config
}

type Service struct {
	k8s        *kubernetes.Clientset
	done       chan struct{}
	namespaces []string
	label      string
	payload    []byte
}

func New(config *Config, k8s *kubernetes.Clientset) Service {
	var patch jobPatch
	patch.Spec.TTLSecondsAfterFinished = int((config.Duration).Seconds())

	payload, err := json.Marshal(patch)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "marshal json", slog.Any("error", err))
		os.Exit(1)
	}

	return Service{
		k8s:        k8s,
		namespaces: strings.Split(config.Namespace, ","),
		label:      config.Label,
		payload:    payload,
		done:       make(chan struct{}),
	}
}

func (s Service) Done() <-chan struct{} {
	return s.done
}

func (s Service) Start(ctx context.Context) {
	defer close(s.done)

	var wg sync.WaitGroup

	for _, namespace := range s.namespaces {
		wg.Add(1)

		go s.watchJobs(ctx, &wg, namespace)
	}

	wg.Wait()
}

func (s Service) watchJobs(ctx context.Context, wg *sync.WaitGroup, namespace string) {
	defer wg.Done()

	for {
		if s.watchNamespace(ctx, namespace) {
			return
		}
	}
}

func (s Service) watchNamespace(ctx context.Context, namespace string) bool {
	watcher, err := s.k8s.BatchV1().Jobs(namespace).Watch(ctx, v1.ListOptions{
		LabelSelector: s.label,
		Watch:         true,
	})
	logger.FatalfOnErr(ctx, err, "watch jobs")

	slog.LogAttrs(ctx, slog.LevelInfo, "Listening jobs", slog.String("namespace", namespace), slog.String("label", s.label))

	var done bool

	concurrent.ChanUntilDone(ctx, watcher.ResultChan(), func(event watch.Event) {
		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			return
		}

		if job.Spec.TTLSecondsAfterFinished != nil && *job.Spec.TTLSecondsAfterFinished != 0 {
			return
		}

		if job.Status.Succeeded != 1 {
			return
		}

		slog.LogAttrs(ctx, slog.LevelInfo, "Updating TTLSecondsAfterFinished", slog.String("namespace", job.Namespace), slog.String("name", job.Name))

		if _, err := s.k8s.BatchV1().Jobs(job.Namespace).Patch(ctx, job.Name, types.MergePatchType, s.payload, v1.PatchOptions{}); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "patch job", slog.String("namespace", job.Namespace), slog.String("name", job.Name), slog.Any("error", err))
		}
	}, func() {
		select {
		case <-ctx.Done():
			done = true
		default:
		}

		watcher.Stop()
	})

	return done
}
