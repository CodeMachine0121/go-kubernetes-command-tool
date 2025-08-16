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
		resourceUsagesPercentages := k8sService.GetPercentageOfResourceUsage(context.Background(), "ctbc-csiw")

		for _, usage := range resourceUsagesPercentages {
			jsonData, _ := json.Marshal(usage)
			fmt.Println(string(jsonData))
		}
	})

	if err != nil {
		panic(err)
	}
}
