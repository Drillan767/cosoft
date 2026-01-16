package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HeaderParts holds the different sections of the header
type HeaderParts struct {
	Left    string
	Center  string
	Right   string
	Credits string
}

// UpdateHeaderMsg is sent by child models to update header parts
type UpdateHeaderMsg struct {
	Center  *string // nil means don't update
	Right   *string
	Credits *string
}

// LayoutConfig holds configuration for the layout appearance
type LayoutConfig struct {
	Header      HeaderParts
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
func NewLayoutWithDefaults(content tea.Model, headerLeft, footer string) *Layout {
	config := DefaultLayoutConfig()
	config.Header.Left = headerLeft
	config.Footer = footer
	return NewLayout(content, config)
}

func (l *Layout) updateStyles() {
	// Header style - no Width set, we control width via content string length
	l.headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color(l.config.HeaderColor)).
		Padding(0, 1)

	// Footer style
	l.footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(l.config.FooterColor)).
		Padding(0, 1).
		Width(l.width)

	contentWidth := l.width
	if l.config.ShowBorder && contentWidth > 0 {
		contentWidth = contentWidth - 2 // Account for border (1 left + 1 right)
		if contentWidth < 0 {
			contentWidth = 0
		}
		l.containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(l.config.BorderColor)).
			Padding(1, 2).
			Width(contentWidth)
	} else if contentWidth > 0 {
		if contentWidth < 0 {
			contentWidth = 0
		}
		l.containerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Width(contentWidth)
	} else {
		// Fallback when width is not yet set
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

	case UpdateHeaderMsg:
		// Handle header updates from child models
		if msg.Center != nil {
			l.config.Header.Center = *msg.Center
		}
		if msg.Right != nil {
			l.config.Header.Right = *msg.Right
		}
		if msg.Credits != nil {
			l.config.Header.Credits = *msg.Credits
		}
		return l, nil
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
	if l.config.ShowHeader && l.hasHeaderContent() {
		header := l.renderHeader()
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

// hasHeaderContent checks if any header part has content
func (l *Layout) hasHeaderContent() bool {
	h := l.config.Header
	return h.Left != "" || h.Center != "" || h.Right != "" || h.Credits != ""
}

// renderHeader renders the multi-part header with proper alignment
func (l *Layout) renderHeader() string {
	h := l.config.Header
	width := l.width
	if width <= 0 {
		width = 80 // fallback width
	}

	// Match the container width
	innerWidth := width - 2
	if innerWidth < 10 {
		innerWidth = 10
	}

	// Build left part
	leftPart := h.Left

	// Build center part
	centerPart := h.Center

	// Build right part (user info + credits)
	var rightParts []string
	if h.Right != "" {
		rightParts = append(rightParts, h.Right)
	}
	if h.Credits != "" {
		rightParts = append(rightParts, h.Credits)
	}
	rightPart := strings.Join(rightParts, " | ")

	// Calculate spacing
	leftLen := len(leftPart)
	centerLen := len(centerPart)
	rightLen := len(rightPart)

	// Calculate available space for gaps
	totalContentLen := leftLen + centerLen + rightLen
	availableSpace := innerWidth - totalContentLen

	if availableSpace < 2 {
		// Not enough space, just concatenate with minimal spacing
		return l.headerStyle.Render(leftPart + " " + centerPart + " " + rightPart)
	}

	var leftGap, rightGap int

	if centerLen > 0 {
		// With center content: distribute space on both sides
		leftGap = availableSpace / 2
		rightGap = availableSpace - leftGap
	} else {
		// No center: all space between left and right
		leftGap = availableSpace
		rightGap = 0
	}

	// Build the header line
	var headerLine string
	if centerLen > 0 {
		headerLine = leftPart + strings.Repeat(" ", leftGap) + centerPart + strings.Repeat(" ", rightGap) + rightPart
	} else {
		headerLine = leftPart + strings.Repeat(" ", leftGap) + rightPart
	}

	return l.headerStyle.Render(headerLine)
}

func (l *Layout) SetHeaderLeft(text string) {
	l.config.Header.Left = text
}

func (l *Layout) SetHeaderCenter(location string) {
	l.config.Header.Center = location
}

func (l *Layout) SetHeaderRight(userInfo string) {
	l.config.Header.Right = userInfo
}

func (l *Layout) SetHeaderCredits(credits string) {
	l.config.Header.Credits = credits
}

func (l *Layout) SetFooter(footer string) {
	l.config.Footer = footer
}

func (l *Layout) GetContent() tea.Model {
	return l.content
}

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
