package main

import (
	"go-k8s-tools/internal/cli"
	"go-k8s-tools/internal/core"
)

func main() {

	container := core.BuildContainer()

	err := container.Invoke(func(terminalUiService cli.ITerminalUIService) {

		terminalUiService.ShowUsagePercentage("ctbc-csiw")
	})

	if err != nil {
		panic(err)
	}
}
