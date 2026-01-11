package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LayoutConfig holds configuration for the layout appearance
type LayoutConfig struct {
	Header      string
	Footer      string
	ShowHeader  bool
	ShowFooter  bool
	ShowBorder  bool
	BorderColor string
	HeaderColor string
	FooterColor string
}

// DefaultLayoutConfig returns a sensible default configuration
func DefaultLayoutConfig() LayoutConfig {
	return LayoutConfig{
		ShowHeader:  true,
		ShowFooter:  true,
		ShowBorder:  true,
		BorderColor: "#fd4b4bff",
		HeaderColor: "#f45656ff",
		FooterColor: "#888888ff",
	}
}

// Layout wraps any tea.Model and provides consistent styling
type Layout struct {
	config  LayoutConfig
	content tea.Model
	width   int
	height  int

	// Styles
	headerStyle    lipgloss.Style
	footerStyle    lipgloss.Style
	containerStyle lipgloss.Style
}

// NewLayout creates a new layout wrapper around any tea.Model
func NewLayout(content tea.Model, config LayoutConfig) *Layout {
	l := &Layout{
		config:  config,
		content: content,
	}
	l.updateStyles()
	return l
}

// NewLayoutWithDefaults creates a layout with default configuration
func NewLayoutWithDefaults(content tea.Model, header, footer string) *Layout {
	config := DefaultLayoutConfig()
	config.Header = header
	config.Footer = footer
	return NewLayout(content, config)
}

func (l *Layout) updateStyles() {
	// Header style
	l.headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color(l.config.HeaderColor)).
		Padding(0, 1).
		Width(l.width)

	// Footer style
	l.footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(l.config.FooterColor)).
		Padding(0, 1).
		Width(l.width)

	// Container style
	if l.config.ShowBorder {
		l.containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(l.config.BorderColor)).
			Padding(1, 2)
	} else {
		l.containerStyle = lipgloss.NewStyle().
			Padding(1, 2)
	}
}

func (l *Layout) Init() tea.Cmd {
	return l.content.Init()
}

func (l *Layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width
		l.height = msg.Height
		l.updateStyles()
	}

	// Update the wrapped content
	var cmd tea.Cmd
	l.content, cmd = l.content.Update(msg)
	return l, cmd
}

func (l *Layout) View() string {
	contentView := l.content.View()

	// If content is empty (quitting), return empty string
	if contentView == "" {
		return ""
	}

	var sections []string

	// Add header if enabled
	if l.config.ShowHeader && l.config.Header != "" {
		header := l.headerStyle.Render(l.config.Header)
		sections = append(sections, header)
	}

	// Add content
	content := l.containerStyle.Render(contentView)
	sections = append(sections, content)

	// Add footer if enabled
	if l.config.ShowFooter && l.config.Footer != "" {
		footer := l.footerStyle.Render(l.config.Footer)
		sections = append(sections, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// SetHeader updates the header text
func (l *Layout) SetHeader(header string) {
	l.config.Header = header
}

// SetFooter updates the footer text
func (l *Layout) SetFooter(footer string) {
	l.config.Footer = footer
}

// GetContent returns the wrapped content model
func (l *Layout) GetContent() tea.Model {
	return l.content
}

// Helper function to create centered text
func CenterText(text string, width int) string {
	if width <= 0 {
		return text
	}
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}
