package k8s

import "context"

type IK8sService interface {
	GetTotalResource(ctx context.Context, namespace string) []Resource
	GetPodResourceUsage(ctx context.Context, namespace string) []ResourceUsage
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

func NewK8sService(proxy IK8sProxy) IK8sService {
	return &K8sService{
		proxy: proxy,
	}
}
