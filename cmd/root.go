package cmd

import (
	"context"
	"fmt"
	"go-k8s-tools/internal/cli"
	"go-k8s-tools/internal/core"
	"go-k8s-tools/internal/k8s"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	namespace string
	interval  int
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "啟動 Kubernetes 資源監控終端介面",
	Long: `resource 命令會啟動一個即時的終端介面，
顯示指定命名空間中所有 Pod 的 CPU 和記憶體使用率。

使用範例：
  gk resource                    # 監控 default 命名空間
  gk resource -n kube-system     # 監控 kube-system 命名空間  
  gk resource -n default -i 2000 # 監控 default 命名空間，每2秒更新一次`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 建立依賴注入容器
		container := core.BuildContainer()

		// 執行資源監控
		return container.Invoke(func(service k8s.IK8sService) error {
			fmt.Printf("🚀 啟動 Kubernetes 資源監控...\n")
			fmt.Printf("📊 命名空間: %s\n", namespace)
			fmt.Printf("⏱️  更新間隔: %dms\n", interval)
			fmt.Printf("💡 按 'q' 或 'Ctrl+C' 退出\n\n")

			terminalUiService := cli.NewTerminalUIModel(
				context.Background(),
				service,
				namespace,
				time.Duration(interval)*time.Millisecond,
			)
			terminalUiService.Run()
			return nil
		})
	},
}

var rootCmd = &cobra.Command{
	Use:   "gk",
	Short: "Kubernetes 資源監控工具",
	Long: `gk 是一個強大的 Kubernetes 資源監控 CLI 工具，
可以幫助您即時監控集群中的資源使用情況。

使用範例：
  gk resource        # 啟動資源監控終端介面`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "執行命令時發生錯誤: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(resourceCmd)

	// 添加命令標誌
	resourceCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "指定要監控的 Kubernetes 命名空間")
	resourceCmd.Flags().IntVarP(&interval, "interval", "i", 1000, "資源監控更新間隔（毫秒）")
}
