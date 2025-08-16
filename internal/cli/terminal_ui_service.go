package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-k8s-tools/internal/k8s"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ITerminalUIService interface {
	ShowUsagePercentage(namespace string)
}

type TerminalUIService struct {
	k8sService k8s.IK8sService
}

// Bubbletea model for the resource usage display
type resourceModel struct {
	resources  []k8s.ResourceUsagePercentage
	namespace  string
	k8sService k8s.IK8sService
	width      int
	height     int
	quitting   bool
}

// Message types for Bubbletea
type tickMsg time.Time
type updateResourcesMsg []k8s.ResourceUsagePercentage

// Initialize the model
func (m resourceModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(),
	)
}

// Tick command for periodic updates
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages and state changes
func (m resourceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "r":
			// Refresh trigger (could be connected to actual data fetching)
			return m, tickCmd()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		// Refresh data from k8s service
		if m.k8sService != nil {
			newData := m.k8sService.GetPercentageOfResourceUsage(context.Background(), m.namespace)
			m.resources = newData
		}
		return m, tickCmd()

	case updateResourcesMsg:
		m.resources = []k8s.ResourceUsagePercentage(msg)
		return m, nil
	}

	return m, nil
}

// View renders the terminal UI
func (m resourceModel) View() string {
	if m.quitting {
		return "\n感謝使用 Kubernetes 資源監控工具！\n"
	}

	// Define styles
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#F25D94")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	progressBarStyle := lipgloss.NewStyle().
		Width(30)

	podNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true).
		Width(25)

	percentageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true)

	lowUsageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575"))

	mediumUsageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFAA00"))

	highUsageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	// Header
	var s strings.Builder
	s.WriteString(headerStyle.Render("🚀 Kubernetes 資源使用率監控"))
	s.WriteString("\n\n")
	s.WriteString(titleStyle.Render("📊 Pod 資源使用率"))
	s.WriteString("\n\n")

	if len(m.resources) == 0 {
		s.WriteString("⏳ 載入資源資料中...\n\n")
	} else {
		// Table header
		s.WriteString(fmt.Sprintf("%-25s %-15s %-30s %-15s %-30s\n",
			lipgloss.NewStyle().Bold(true).Render("Pod 名稱"),
			lipgloss.NewStyle().Bold(true).Render("CPU (%)"),
			lipgloss.NewStyle().Bold(true).Render("CPU 使用率"),
			lipgloss.NewStyle().Bold(true).Render("Memory (%)"),
			lipgloss.NewStyle().Bold(true).Render("Memory 使用率"),
		))
		s.WriteString(strings.Repeat("─", 120))
		s.WriteString("\n")

		// Resource data
		for _, resource := range m.resources {
			// Format pod name
			podName := podNameStyle.Render(resource.Name)

			// CPU usage
			cpuPercentageStr := percentageStyle.Render(fmt.Sprintf("%.2f%%", resource.CPUPercentage))
			cpuProgress := createProgressBar(resource.CPUPercentage, progressBarStyle)

			// Memory usage
			memoryPercentageStr := percentageStyle.Render(fmt.Sprintf("%.2f%%", resource.MemoryPercentage))
			memoryProgress := createProgressBar(resource.MemoryPercentage, progressBarStyle)

			s.WriteString(fmt.Sprintf("%-25s %-15s %-30s %-15s %-30s\n",
				podName,
				getColoredPercentage(resource.CPUPercentage, cpuPercentageStr, lowUsageStyle, mediumUsageStyle, highUsageStyle),
				cpuProgress,
				getColoredPercentage(resource.MemoryPercentage, memoryPercentageStr, lowUsageStyle, mediumUsageStyle, highUsageStyle),
				memoryProgress,
			))
		}
	}

	// Footer with instructions
	s.WriteString("\n" + strings.Repeat("─", 120) + "\n")
	s.WriteString("⌨️  操作說明: ")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render("q"))
	s.WriteString(" 退出 | ")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render("r"))
	s.WriteString(" 重新整理 | ")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render("Ctrl+C"))
	s.WriteString(" 強制退出")
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("🔄 自動更新: 每 2 秒 | ⏰ 更新時間: %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(time.Now().Format("15:04:05"))))

	return s.String()
}

// Create a visual progress bar
func createProgressBar(percentage float64, barStyle lipgloss.Style) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}

	// Get width from the bar style
	width := barStyle.GetWidth()
	if width == 0 {
		width = 30 // fallback width
	}

	filled := int(percentage / 100 * float64(width))
	empty := width - filled

	var bar strings.Builder
	bar.WriteString("[")

	// Colored progress based on usage level
	var fillChar string
	var fillColor lipgloss.Color
	if percentage < 50 {
		fillColor = lipgloss.Color("#04B575") // Green
		fillChar = "█"
	} else if percentage < 80 {
		fillColor = lipgloss.Color("#FFAA00") // Orange
		fillChar = "█"
	} else {
		fillColor = lipgloss.Color("#FF0000") // Red
		fillChar = "█"
	}

	fillStyle := lipgloss.NewStyle().Foreground(fillColor)
	bar.WriteString(fillStyle.Render(strings.Repeat(fillChar, filled)))
	bar.WriteString(strings.Repeat("░", empty))
	bar.WriteString("]")

	// Apply the progressBarStyle to the entire progress bar
	return barStyle.Render(bar.String())
}

// Get colored percentage text based on usage level
func getColoredPercentage(percentage float64, text string, lowStyle, mediumStyle, highStyle lipgloss.Style) string {
	if percentage < 50 {
		return lowStyle.Render(text)
	} else if percentage < 80 {
		return mediumStyle.Render(text)
	} else {
		return highStyle.Render(text)
	}
}

func (t *TerminalUIService) ShowUsagePercentage(namespace string) {
	percentage := t.k8sService.GetPercentageOfResourceUsage(context.Background(), namespace)

	// Create the initial model with data
	model := resourceModel{
		resources:  percentage,
		namespace:  namespace,
		k8sService: t.k8sService,
		width:      80,
		height:     24,
	}

	// Create and start the Bubbletea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("錯誤: %v\n", err)
	}
}

func NewTerminalUiService(k8sService k8s.IK8sService) ITerminalUIService {
	return &TerminalUIService{
		k8sService: k8sService,
	}
}
