package k8s

import (
	"context"

	"github.com/samber/lo"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

// Clientset returns the underlying Kubernetes clientset
func (c *Client) Clientset() *kubernetes.Clientset {
	return c.clientset
}

// Config returns the underlying Kubernetes configuration
func (c *Client) Config() *rest.Config {
	return c.config
}

type Resource struct {
	Name, Namespace           string
	RequestCPU, RequestMemory float32
	LimitCPU, LimitMemory     float32
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
				requestCPU = float32(cpuReq.MilliValue()) / 1000
			}

			if memoryReq := container.Resources.Requests[corev1.ResourceMemory]; memoryReq.MilliValue() != 0 {
				requestMemory = float32(memoryReq.MilliValue()) / (1024 * 1024 * 1024)
			}

			if cpuLimit := container.Resources.Limits[corev1.ResourceCPU]; cpuLimit.MilliValue() != 0 {
				limitCPU = float32(cpuLimit.MilliValue()) / 1000
			}
			if memoryLimit := container.Resources.Limits[corev1.ResourceMemory]; memoryLimit.MilliValue() != 0 {
				limitMemory = float32(memoryLimit.MilliValue()) / (1024 * 1024 * 1024)
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
