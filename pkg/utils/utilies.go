package utils

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

func CalcCPU(cpuValue resource.Quantity) float64 {
	value := float64(cpuValue.MilliValue()) / 1000
	return value
}

func CalcMemory(memoryValue resource.Quantity) float64 {
	value := float64(memoryValue.MilliValue()) / (1024 * 1024 * 1024)
	return value
}
