package cli

import (
	"context"
	"fmt"
	"go-k8s-tools/internal/k8s"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type resourceUsageMsg []k8s.ResourceUsagePercentage

type ResourceTerminalUi struct {
	ctx        context.Context
	k8sService k8s.IK8sService
	namespace  string
	interval   time.Duration
	usages     []k8s.ResourceUsagePercentage
	error      error
}

func NewResourceTerminalUi(ctx context.Context, svc k8s.IK8sService, namespace string, interval time.Duration) *ResourceTerminalUi {
	return &ResourceTerminalUi{
		ctx:        ctx,
		k8sService: svc,
		namespace:  namespace,
		interval:   interval,
	}
}

func (m *ResourceTerminalUi) Init() tea.Cmd {
	return tea.Tick(m.interval, func(t time.Time) tea.Msg {
		usages := m.k8sService.GetPercentageOfResourceUsage(m.ctx, m.namespace)
		return resourceUsageMsg(usages)
	})
}

func (m *ResourceTerminalUi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case resourceUsageMsg:
		m.usages = msg
		return m, tea.Tick(m.interval, func(t time.Time) tea.Msg {
			usages := m.k8sService.GetPercentageOfResourceUsage(m.ctx, m.namespace)
			return resourceUsageMsg(usages)
		})
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ResourceTerminalUi) View() string {
	var view string

	// Header section with border
	view += "################################################################\n"
	view += "#                🚀 Kubernetes Resource Monitor                #\n"
	view += "#                     Resource Usage (%)                       #\n"
	view += "################################################################\n"

	// Table header - fixed width columns
	view += "# Pod Name               #  CPU %  # Memory % # Status #\n"
	view += "################################################################\n"

	// Handle empty data
	if len(m.usages) == 0 {
		view += "#                    No data available                        #\n"
		view += "#                Loading resource usage...                    #\n"
	} else {
		// Data rows with proper alignment
		for _, usage := range m.usages {
			status := getStatusIndicator(usage.CPUPercentage, usage.MemoryPercentage)

			view += fmt.Sprintf("# %-22s # %6.2f%% # %7.2f%% # %-6s #\n",
				truncateName(usage.Name, 22),
				usage.CPUPercentage,
				usage.MemoryPercentage,
				status)
		}
	}

	view += "################################################################\n"

	// Bottom status and instructions
	view += "\n"
	view += "📊 Status Legend:\n"
	view += "   🟢 Normal (<70%)   🟡 Warning (70-90%)   🔴 High (>90%)\n"
	view += "\n"
	view += "⌨️  Controls: Press 'q' or 'Ctrl+C' to exit\n"

	// Display error if any
	if m.error != nil {
		view += fmt.Sprintf("\n⚠️  Error: %v\n", m.error)
	}

	return view
}

// Get status indicator
func getStatusIndicator(cpu, memory float64) string {
	maxUsage := cpu
	if memory > cpu {
		maxUsage = memory
	}

	switch {
	case maxUsage >= 90:
		return "🔴 High"
	case maxUsage >= 70:
		return "🟡 Warn"
	default:
		return "🟢 OK"
	}
}

// Truncate long names
func truncateName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	if maxLen <= 3 {
		return name[:maxLen]
	}
	return name[:maxLen-3] + "..."
}

func (m *ResourceTerminalUi) Run() {
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
