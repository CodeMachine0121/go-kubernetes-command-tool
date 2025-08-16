package main

import (
	"context"
	"go-k8s-tools/internal/cli"
	"go-k8s-tools/internal/core"
	"go-k8s-tools/internal/k8s"
)

func main() {

	container := core.BuildContainer()

	err := container.Invoke(func(service k8s.IK8sService) {
		terminalUiService := cli.NewTerminalUIModel(context.Background(), service, "ctbc-csiw", 1000)
		terminalUiService.Run()
	})

	if err != nil {
		panic(err)
	}
}
