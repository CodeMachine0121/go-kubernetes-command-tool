package main

import (
	"context"
	"go-k8s-tools/internal/k8s"

	"github.com/samber/lo"
)

func main() {

	var factory k8s.IClientFactory = &k8s.ClientFactory{}

	client := factory.NewClient("")
	deploymentNames := client.TestToGetDeploymentName(context.Background())

	lo.ForEach(deploymentNames, func(item string, _ int) {
		println(item)
	})
}
