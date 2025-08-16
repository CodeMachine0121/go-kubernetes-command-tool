package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-k8s-tools/internal/core"
	"go-k8s-tools/internal/k8s"
)

func main() {

	container := core.BuildContainer()

	err := container.Invoke(func(k8sService k8s.IK8sService) {

		resources := k8sService.GetTotalResource(context.Background(), "ctbc-csiw")
		resourceUsages := k8sService.GetPodResourceUsage(context.Background(), "ctbc-csiw")

		for _, resource := range resources {
			jsonData, _ := json.Marshal(resource)
			fmt.Println(string(jsonData))
		}

		for _, usage := range resourceUsages {
			jsonData, _ := json.Marshal(usage)
			fmt.Println(string(jsonData))
		}
	})

	if err != nil {
		panic(err)
	}
}
