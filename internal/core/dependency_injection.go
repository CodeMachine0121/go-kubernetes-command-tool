package core

import (
	"go-k8s-tools/internal/k8s"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	if err := container.Provide(k8s.NewK8sProxy); err != nil {
		panic(err)
	}

	if err := container.Provide(k8s.NewK8sService); err != nil {
		panic(err)
	}

	return container
}
