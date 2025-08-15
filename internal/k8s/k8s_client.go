package k8s

import (
	"context"
	"fmt"
	"go-k8s-tools/pkg/utils"

	"github.com/samber/lo"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	clientset     *kubernetes.Clientset
	metricsClient *metricsclient.Clientset
	config        *rest.Config
}

// Clientset returns the underlying Kubernetes clientset
func (c *Client) Clientset() *kubernetes.Clientset {
	return c.clientset
}

// Config returns the underlying Kubernetes configuration
func (c *Client) Config() *rest.Config {
	return c.config
}

func (c *Client) GetTotalResource(ctx context.Context, namespace string) []Resource {

	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	return lo.FlatMap(deployments.Items, func(deployment v1.Deployment, _ int) []Resource {
		return lo.Map(deployment.Spec.Template.Spec.Containers, func(container corev1.Container, _ int) Resource {
			var requestCPU, requestMemory, limitCPU, limitMemory float32

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

// TestToGetDeploymentName tests the connection to the Kubernetes cluster
func (c *Client) TestToGetDeploymentName(ctx context.Context) []string {
	deployments, err := c.clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	deploymentNames := lo.Map(deployments.Items, func(item v1.Deployment, _ int) string {
		return item.Name
	})

	return deploymentNames
}

func (c *Client) GetPodResourceUsage(ctx context.Context, namespace string) ([]ResourceUsage, error) {
	resources := make([]ResourceUsage, 0)

	if c.metricsClient == nil {
		return nil, fmt.Errorf("metrics client is not initialized")
	}

	metrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
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
	return resources, nil
}
