package cli

import (
	"context"
	"fmt"
	"go-k8s-tools/internal/k8s"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type resourceUsageMsg []k8s.ResourceUsagePercentage

type TerminalUIModel struct {
	ctx        context.Context
	k8sService k8s.IK8sService
	namespace  string
	interval   time.Duration
	usages     []k8s.ResourceUsagePercentage
	error      error
}

func NewTerminalUIModel(ctx context.Context, svc k8s.IK8sService, namespace string, interval time.Duration) *TerminalUIModel {
	return &TerminalUIModel{
		ctx:        ctx,
		k8sService: svc,
		namespace:  namespace,
		interval:   interval,
	}
}

func (m *TerminalUIModel) Init() tea.Cmd {
	return tea.Tick(m.interval, func(t time.Time) tea.Msg {
		usages := m.k8sService.GetPercentageOfResourceUsage(m.ctx, m.namespace)
		return resourceUsageMsg(usages)
	})
}

func (m *TerminalUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case resourceUsageMsg:
		m.usages = msg
		return m, tea.Tick(m.interval, func(t time.Time) tea.Msg {
			usages := m.k8sService.GetPercentageOfResourceUsage(m.ctx, m.namespace)
			return resourceUsageMsg(usages)
		})
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *TerminalUIModel) View() string {
	view := "Resource Usage (%)\n"
	view += "Pod Name           CPU %    Memory %\n"
	view += "-----------------------------------\n"
	for _, usage := range m.usages {
		view += fmt.Sprintf("%-18s %-8.2f %-8.2f\n", usage.Name, usage.CPUPercentage, usage.MemoryPercentage)
	}
	view += "\n按 q 離開"
	return view
}

func (m *TerminalUIModel) Run() {
	model := NewTerminalUIModel(m.ctx, m.k8sService, m.namespace, m.interval)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("錯誤: %v\n", err)
	}
}
