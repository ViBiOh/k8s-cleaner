package k8s

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/ViBiOh/flags"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	File string
}

func Flags(fs *flag.FlagSet, prefix string, overrides ...flags.Override) *Config {
	var defaultConfig string
	if home := homedir.HomeDir(); home != "" {
		defaultConfig = filepath.Join(home, ".kube", "config")
	}

	var config Config

	flags.New("config", "Path to kubeconfig file").Prefix(prefix).DocPrefix("k8s").StringVar(fs, &config.File, defaultConfig, overrides)

	return &config
}

func New(config *Config) (*kubernetes.Clientset, error) {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		if len(config.File) == 0 {
			return nil, fmt.Errorf("get in-cluster config: %w", err)
		}

		k8sConfig, err = clientcmd.BuildConfigFromFlags("", config.File)
		if err != nil {
			return nil, fmt.Errorf("get cluster config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return clientset, nil
}
