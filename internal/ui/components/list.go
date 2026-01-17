package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Label    string
	Subtitle string
	Value    interface{}
}

type ListItemDelegate struct{}

type ListModel struct {
	list      list.Model
	selected  Item
	confirmed bool
}

func (i Item) FilterValue() string { return i.Label }

func (d ListItemDelegate) Height() int                             { return 2 }
func (d ListItemDelegate) Spacing() int                            { return 1 }
func (d ListItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ListItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(Item)
	isSelected := index == m.Index()

	// Styles
	labelStyle := lipgloss.NewStyle().
		Width(80).
		Padding(0, 2)

	subtitleStyle := lipgloss.NewStyle().
		Width(80).
		Padding(0, 2).
		Foreground(lipgloss.Color("8"))

	if isSelected {
		labelStyle = labelStyle.
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("15")).
			Bold(true)
		subtitleStyle = subtitleStyle.
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("7"))
	}

	output := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(i.Label),
		subtitleStyle.Render(i.Subtitle),
	)

	fmt.Fprint(w, output)
}

func NewListModel(items []Item, title string) ListModel {
	listItems := make([]list.Item, len(items))

	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, ListItemDelegate{}, 80, 20)

	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		Margin(1, 0, 1, 2).
		Foreground(lipgloss.Color("62")).
		Bold(true)

	return ListModel{
		list:      l,
		confirmed: false,
	}
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.selected = m.list.SelectedItem().(Item)
			m.confirmed = true
			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	if m.confirmed {
		return fmt.Sprintf("Selected: %s (Value: %v)\n", m.selected.Label, m.selected.Value)
	}
	return m.list.View()
}

func (m ListModel) GetSelection() *Item {
	if m.confirmed {
		return &m.selected
	}

	return nil
}
