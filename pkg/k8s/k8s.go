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

// Config of package
type Config struct {
	file *string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string, overrides ...flags.Override) Config {
	var defaultConfig string
	if home := homedir.HomeDir(); home != "" {
		defaultConfig = filepath.Join(home, ".kube", "config")
	}

	return Config{
		file: flags.String(fs, prefix, "k8s", "config", "Path to kubeconfig file", defaultConfig, overrides),
	}
}

// New creates new App from Config
func New(config Config) (*kubernetes.Clientset, error) {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		if len(*config.file) == 0 {
			return nil, fmt.Errorf("unable to get in-cluster config: %s", err)
		}

		k8sConfig, err = clientcmd.BuildConfigFromFlags("", *config.file)
		if err != nil {
			return nil, fmt.Errorf("unable to get cluster config: %s", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %s", err)
	}

	return clientset, nil
}
