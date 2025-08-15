package k8s

type Resource struct {
	Name, Namespace           string
	RequestCPU, RequestMemory float32
	LimitCPU, LimitMemory     float32
}

type ResourceUsage struct {
	PodName     string
	CPU, Memory float32
}
