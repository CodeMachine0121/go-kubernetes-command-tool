package utils

import "k8s.io/apimachinery/pkg/api/resource"

func CalcCPU(cpuValue resource.Quantity) float32 {
	return float32(cpuValue.MilliValue()) / 1000
}

func CalcMemory(memoryValue resource.Quantity) float32 {
	return float32(memoryValue.MilliValue()) / (1024 * 1024 * 1024)
}
