package k8s

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

func NewClient(kubeConfigPath string) *Client {

	config, err := BuildConfig(kubeConfigPath)

	if err != nil {
		fmt.Errorf("Failed to found kubenets config")
		return nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Failed to create kubernet clientset")
		return nil
	}

	metricsClient, err := metricsclient.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Failed to create metrics clientset")
		return nil
	}

	return &Client{
		clientset:     clientset,
		config:        config,
		metricsClient: metricsClient,
	}
}

func BuildConfig(kubeconfigPath string) (*rest.Config, error) {
	// If no kubeconfig path provided, use default locations
	if kubeconfigPath == "" {
		kubeconfigPath = GetDefaultKubeconfigPath()
	}

	// Try to use kubeconfig file first
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		// Fall back to in-cluster config if kubeconfig fails
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubeconfig (%s) and in-cluster: %w", kubeconfigPath, err)
		}
	}

	return config, nil
}

func GetDefaultKubeconfigPath() string {
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}
	return ""
}
