package k8s

import (
	"context"
	"math"
)

type IK8sService interface {
	GetTotalResource(ctx context.Context, namespace string) []Resource
	GetPodResourceUsage(ctx context.Context, namespace string) []ResourceUsage
	GetPercentageOfResourceUsage(ctx context.Context, namespace string) []ResourceUsagePercentage
}

type ResourceUsagePercentage struct {
	Name             string  `json:"name"`
	CPUPercentage    float64 `json:"cpu_percentage"`
	MemoryPercentage float64 `json:"memory_percentage"`
}

type K8sService struct {
	proxy IK8sProxy
}

func (s *K8sService) GetTotalResource(ctx context.Context, namespace string) []Resource {
	return s.proxy.GetTotalResource(ctx, namespace)
}

func (s *K8sService) GetPodResourceUsage(ctx context.Context, namespace string) []ResourceUsage {
	return s.proxy.GetPodResourceUsage(ctx, namespace)
}

func (s *K8sService) GetPercentageOfResourceUsage(ctx context.Context, namespace string) []ResourceUsagePercentage {
	totalResources := s.proxy.GetTotalResource(ctx, namespace)
	podUsages := s.proxy.GetPodResourceUsage(ctx, namespace)

	// Create a map for quick lookup of resource limits by container name
	resourceMap := make(map[string]Resource)
	for _, resource := range totalResources {
		resourceMap[resource.Name] = resource
	}

	var results []ResourceUsagePercentage

	for _, usage := range podUsages {
		if resource, exists := resourceMap[usage.PodName]; exists {
			var cpuPercentage, memoryPercentage float64

			// Calculate CPU percentage (usage / limit * 100)
			if resource.LimitCPU > 0 {
				cpuPercentage = math.Floor(usage.CPU/resource.LimitCPU*10000) / 100
			}

			// Calculate Memory percentage (usage / limit * 100)
			if resource.LimitMemory > 0 {
				memoryPercentage = math.Floor(usage.Memory/resource.LimitMemory*10000) / 100
			}

			results = append(results, ResourceUsagePercentage{
				Name:             usage.PodName,
				CPUPercentage:    cpuPercentage,
				MemoryPercentage: memoryPercentage,
			})
		}
	}

	return results
}

func NewK8sService(proxy IK8sProxy) IK8sService {
	return &K8sService{
		proxy: proxy,
	}
}
