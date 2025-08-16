package k8s

type Resource struct {
	Name, Namespace           string
	RequestCPU, RequestMemory float64
	LimitCPU, LimitMemory     float64
}

type ResourceUsage struct {
	PodName     string
	CPU, Memory float64
}
