package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockK8sProxy is a mock implementation of IK8sProxy interface
type MockK8sProxy struct {
	mock.Mock
}

func (m *MockK8sProxy) GetTotalResource(ctx context.Context, namespace string) []Resource {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]Resource)
}

func (m *MockK8sProxy) GetPodResourceUsage(ctx context.Context, namespace string) []ResourceUsage {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]ResourceUsage)
}

func TestK8sService_GetPercentageOfResourceUsage(t *testing.T) {
	// Arrange
	mockProxy := new(MockK8sProxy)
	service := NewK8sService(mockProxy)
	ctx := context.Background()
	namespace := "test-namespace"

	// Mock data setup
	totalResources := []Resource{
		{
			Name:          "container1",
			Namespace:     namespace,
			RequestCPU:    100, // 100m CPU request
			RequestMemory: 256, // 256Mi Memory request
			LimitCPU:      200, // 200m CPU limit
			LimitMemory:   512, // 512Mi Memory limit
		},
		{
			Name:          "container2",
			Namespace:     namespace,
			RequestCPU:    50,  // 50m CPU request
			RequestMemory: 128, // 128Mi Memory request
			LimitCPU:      100, // 100m CPU limit
			LimitMemory:   256, // 256Mi Memory limit
		},
	}

	podResourceUsages := []ResourceUsage{
		{
			PodName: "container1",
			CPU:     80,  // 80m CPU usage
			Memory:  200, // 200Mi Memory usage
		},
		{
			PodName: "container2",
			CPU:     30,  // 30m CPU usage
			Memory:  100, // 100Mi Memory usage
		},
	}

	// Setup mock expectations
	mockProxy.On("GetTotalResource", ctx, namespace).Return(totalResources)
	mockProxy.On("GetPodResourceUsage", ctx, namespace).Return(podResourceUsages)

	// Act
	result := service.GetPercentageOfResourceUsage(ctx, namespace)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// Expected calculations:
	// container1: CPU = 80/200 = 40%, Memory = 200/512 = 39.06%
	// container2: CPU = 30/100 = 30%, Memory = 100/256 = 39.06%
	assert.Equal(t, float64(40), result[0].CPUPercentage)
	assert.Equal(t, 39.06, result[0].MemoryPercentage)
	assert.Equal(t, float64(30), result[1].CPUPercentage)
	assert.Equal(t, 39.06, result[1].MemoryPercentage)

	// Verify mock expectations were met
	mockProxy.AssertExpectations(t)
}
