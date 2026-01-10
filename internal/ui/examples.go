package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Example 1: Simple List wrapped in Layout

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type ListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewListModel(items []list.Item, title string) *ListModel {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &ListModel{list: l}
}

func (m *ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 6) // Account for layout header/footer
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
			}
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ListModel) View() string {
	if m.quitting {
		return ""
	}
	return m.list.View()
}

func (m *ListModel) GetChoice() string {
	return m.choice
}

// ExampleListWithLayout shows how to wrap a list in the layout
func (ui *UI) ExampleListWithLayout() (string, error) {
	items := []list.Item{
		item{title: "Deploy to Production", desc: "Push latest changes to prod"},
		item{title: "Run Tests", desc: "Execute test suite"},
		item{title: "View Logs", desc: "Check application logs"},
		item{title: "Database Backup", desc: "Create DB snapshot"},
	}

	listModel := NewListModel(items, "Select an Action")

	// Wrap in layout
	layout := NewLayoutWithDefaults(
		listModel,
		"COSOFT CLI - Main Menu",
		"Use ↑/↓ to navigate • Enter to select • q to quit",
	)

	p := tea.NewProgram(layout)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if layoutModel, ok := finalModel.(*Layout); ok {
		if listModel, ok := layoutModel.GetContent().(*ListModel); ok {
			return listModel.GetChoice(), nil
		}
	}

	return "", fmt.Errorf("failed to retrieve selection")
}

// Example 2: Table wrapped in Layout

type TableModel struct {
	table    table.Model
	quitting bool
}

func NewTableModel() *TableModel {
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 15},
		{Title: "Region", Width: 15},
	}

	rows := []table.Row{
		{"1", "Server Alpha", "Running", "us-east-1"},
		{"2", "Server Beta", "Stopped", "us-west-2"},
		{"3", "Server Gamma", "Running", "eu-west-1"},
		{"4", "Server Delta", "Maintenance", "ap-south-1"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &TableModel{table: t}
}

func (m *TableModel) Init() tea.Cmd {
	return nil
}

func (m *TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			// Handle selection
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *TableModel) View() string {
	if m.quitting {
		return ""
	}
	return m.table.View()
}

func (m *TableModel) GetSelected() table.Row {
	return m.table.SelectedRow()
}

// ExampleTableWithLayout shows how to wrap a table in the layout
func (ui *UI) ExampleTableWithLayout() (table.Row, error) {
	tableModel := NewTableModel()

	// Wrap in layout with custom config
	config := DefaultLayoutConfig()
	config.Header = "COSOFT CLI - Server Status"
	config.Footer = "Use ↑/↓ to navigate • Enter to select • q to quit"
	config.BorderColor = "#00FF00"
	config.HeaderColor = "#0066CC"

	layout := NewLayout(tableModel, config)

	p := tea.NewProgram(layout)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if layoutModel, ok := finalModel.(*Layout); ok {
		if tableModel, ok := layoutModel.GetContent().(*TableModel); ok {
			return tableModel.GetSelected(), nil
		}
	}

	return nil, fmt.Errorf("failed to retrieve selection")
}

// Example 3: Simple text content wrapped in Layout

type SimpleTextModel struct {
	content  string
	quitting bool
}

func NewSimpleTextModel(content string) *SimpleTextModel {
	return &SimpleTextModel{content: content}
}

func (m *SimpleTextModel) Init() tea.Cmd {
	return nil
}

func (m *SimpleTextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *SimpleTextModel) View() string {
	if m.quitting {
		return ""
	}
	return m.content
}

// ExampleTextWithLayout shows how to wrap simple text content
func (ui *UI) ExampleTextWithLayout(title, content, footer string) error {
	textModel := NewSimpleTextModel(content)

	layout := NewLayoutWithDefaults(textModel, title, footer)

	p := tea.NewProgram(layout)
	_, err := p.Run()
	return err
}

// Example 4: Dashboard with multiple sections
func (ui *UI) ExampleDashboard() error {
	var sections []string

	sections = append(sections, lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00")).
		Render("System Status: ✓ Healthy"))

	sections = append(sections, "\n")

	stats := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render("CPU\n65%"),
		"  ",
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render("Memory\n4.2 GB"),
		"  ",
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render("Disk\n120 GB"),
	)
	sections = append(sections, stats)

	sections = append(sections, "\n\n")
	sections = append(sections, "Recent Activity:")
	sections = append(sections, "  • User logged in from 192.168.1.1")
	sections = append(sections, "  • Database backup completed")
	sections = append(sections, "  • New deployment to production")

	content := strings.Join(sections, "\n")

	dashboardModel := NewSimpleTextModel(content)

	config := DefaultLayoutConfig()
	config.Header = "COSOFT CLI - Dashboard"
	config.Footer = "Press q or Ctrl+C to exit"
	config.BorderColor = "#FFD700"

	layout := NewLayout(dashboardModel, config)

	p := tea.NewProgram(layout)
	_, err := p.Run()
	return err
}
