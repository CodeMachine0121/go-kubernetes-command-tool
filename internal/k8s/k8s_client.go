package k8s

import (
	"context"

	"github.com/samber/lo"
	v1 "k8s.io/api/apps/v1"
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

// TestToGetDeploymentName tests the connection to the Kubernetes cluster
func (c *Client) TestToGetDeploymentName(ctx context.Context) []string {
	_, err := c.clientset.Discovery().ServerVersion()
	deployments, _ := c.clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	deploymentNames := lo.Map(deployments.Items, func(item v1.Deployment, _ int) string {
		return item.Name
	})

	return deploymentNames
}
