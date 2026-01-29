package common

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Common struct{}

func NewCommon() *Common {
	return &Common{}
}

func (c *Common) CreateTable(header []string, rows [][]string) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#fd4b4b"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#fd4b4b")).Bold(true).Align(lipgloss.Center)
			case col == 1:
				return lipgloss.NewStyle().Padding(0, 1).Width(20).Foreground(lipgloss.Color("245"))
			default:
				return lipgloss.NewStyle().Padding(0, 1).Width(14).Foreground(lipgloss.Color("245"))
			}
		}).
		Headers(header...)

	for _, row := range rows {
		t.Row(row...)
	}

	return "\n\n" + t.String() + "\n\n"
}
