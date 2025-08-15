package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-k8s-tools/internal/k8s"
)

func main() {

	var factory k8s.IClientFactory = &k8s.ClientFactory{}

	client := factory.NewClient("")
	resources := client.GetTotalResource(context.Background(), "ctbc-csiw")

	resourceUsages, _ := client.GetPodResourceUsage(context.Background(), "ctbc-csiw")

	for _, resource := range resources {
		jsonData, _ := json.Marshal(resource)
		fmt.Println(string(jsonData))
	}

	for _, usage := range resourceUsages {
		jsonData, _ := json.Marshal(usage)
		fmt.Println(string(jsonData))
	}

}
