package k8s

import (
	"context"
	"fmt"
	"go-k8s-tools/pkg/utils"
	"path/filepath"

	"github.com/samber/lo"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

type IK8sProxy interface {
	GetTotalResource(ctx context.Context, namespace string) []Resource
	GetPodResourceUsage(ctx context.Context, namespace string) []ResourceUsage
}

type K8sProxy struct {
	clientset     *kubernetes.Clientset
	metricsClient *metricsclient.Clientset
	config        *rest.Config
}

// NewK8sProxy returns a new K8sProxy instance
func NewK8sProxy() IK8sProxy {
	// If no kubeconfig path provided, use default locations
	configPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(fmt.Errorf("未找到 kubeconfig，請確認路徑: %s", configPath))
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("Failed to create kubernet clientSet"))
	}

	metricsClient, err := metricsclient.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("Failed to create metrics clientSet"))
	}

	return &K8sProxy{
		clientset:     clientSet,
		config:        config,
		metricsClient: metricsClient,
	}
}

// GetTotalResource returns the total resource of all deployments in the specified namespace
func (c *K8sProxy) GetTotalResource(ctx context.Context, namespace string) []Resource {

	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	return lo.FlatMap(deployments.Items, func(deployment v1.Deployment, _ int) []Resource {
		return lo.Map(deployment.Spec.Template.Spec.Containers, func(container corev1.Container, _ int) Resource {
			var requestCPU, requestMemory, limitCPU, limitMemory float64

			if cpuReq := container.Resources.Requests[corev1.ResourceCPU]; cpuReq.MilliValue() != 0 {
				requestCPU = utils.CalcCPU(cpuReq)
			}

			if memoryReq := container.Resources.Requests[corev1.ResourceMemory]; memoryReq.MilliValue() != 0 {
				requestMemory = utils.CalcMemory(memoryReq)
			}

			if cpuLimit := container.Resources.Limits[corev1.ResourceCPU]; cpuLimit.MilliValue() != 0 {
				limitCPU = utils.CalcCPU(cpuLimit)
			}
			if memoryLimit := container.Resources.Limits[corev1.ResourceMemory]; memoryLimit.MilliValue() != 0 {
				limitMemory = utils.CalcMemory(memoryLimit)
			}

			return Resource{
				Name:          container.Name,
				Namespace:     namespace,
				RequestCPU:    requestCPU,
				RequestMemory: requestMemory,
				LimitCPU:      limitCPU,
				LimitMemory:   limitMemory,
			}
		})
	})
}

// GetPodResourceUsage returns the resource usage of all pods in the specified namespace
func (c *K8sProxy) GetPodResourceUsage(ctx context.Context, namespace string) []ResourceUsage {
	resources := make([]ResourceUsage, 0)

	if c.metricsClient == nil {
		panic(fmt.Errorf("metrics proxy is nil"))
	}

	metrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, podMetric := range metrics.Items {
		for _, container := range podMetric.Containers {
			cpuUsage := utils.CalcCPU(*container.Usage.Cpu())
			memoryUsage := utils.CalcMemory(*container.Usage.Memory())

			resources = append(resources, ResourceUsage{
				PodName: container.Name,
				CPU:     cpuUsage,
				Memory:  memoryUsage,
			})
		}
	}
	return resources
}
