package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Item represents a selectable option with label, subtitle, and value.
type Item[T comparable] struct {
	Label    string
	Subtitle string
	Value    T
}

// ListField is a huh-compatible field that displays a list of items
// with labels and subtitles.
type ListField[T comparable] struct {
	value    *T
	items    []Item[T]
	title    string
	cursor   int
	focused  bool
	selected bool

	// huh configuration
	key        string
	width      int
	height     int
	accessible bool
	theme      *huh.Theme
	keymap     *huh.KeyMap
	position   huh.FieldPosition

	// validation
	validate func(T) error
	err      error
}

// NewListField creates a new ListField with the given items and title.
func NewListField[T comparable](items []Item[T], title string) *ListField[T] {
	return &ListField[T]{
		items:    items,
		title:    title,
		cursor:   0,
		width:    80,
		height:   10,
		validate: func(T) error { return nil },
	}
}

// Value sets the pointer to store the selected value.
func (f *ListField[T]) Value(value *T) *ListField[T] {
	f.value = value
	return f
}

// Key sets the field's key for form access.
func (f *ListField[T]) Key(key string) *ListField[T] {
	f.key = key
	return f
}

// Title sets the field's title.
func (f *ListField[T]) Title(title string) *ListField[T] {
	f.title = title
	return f
}

// Validate sets the validation function.
func (f *ListField[T]) Validate(validate func(T) error) *ListField[T] {
	f.validate = validate
	return f
}

// --- huh.Field interface implementation ---

func (f *ListField[T]) Init() tea.Cmd {
	return nil
}

func (f *ListField[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if f.cursor > 0 {
				f.cursor--
			}
		case "down", "j":
			if f.cursor < len(f.items)-1 {
				f.cursor++
			}
		case "enter":
			if len(f.items) > 0 {
				f.selected = true
				if f.value != nil {
					*f.value = f.items[f.cursor].Value
				}
				f.err = f.validate(f.items[f.cursor].Value)
				return f, func() tea.Msg { return huh.NextField() }
			}
		}
	}
	return f, nil
}

func (f *ListField[T]) View() string {
	if len(f.items) == 0 {
		return "No items available"
	}

	var b strings.Builder

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true)

	if f.theme != nil {
		titleStyle = f.theme.Focused.Title
	}

	// Not focused: show condensed view (title + selected value)
	if !f.focused {
		if f.selected && f.title != "" {
			selectedItem := f.items[f.cursor]
			b.WriteString(titleStyle.Render(f.title))
			b.WriteString(" ")
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(selectedItem.Label))
			b.WriteString("\n")
		}
		return b.String()
	}

	// Focused: show full list
	if f.title != "" {
		b.WriteString(titleStyle.Render(f.title))
		b.WriteString("\n\n")
	}

	for i, item := range f.items {
		isCursor := i == f.cursor

		labelStyle := lipgloss.NewStyle().Padding(0, 2).Width(f.width / 2)
		subtitleStyle := lipgloss.NewStyle().
			Padding(0, 2).
			Width(f.width / 2).
			Foreground(lipgloss.Color("8"))

		if isCursor {
			labelStyle = labelStyle.
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("15")).
				Bold(true)
			subtitleStyle = subtitleStyle.
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("7"))
		}

		b.WriteString(labelStyle.Render(item.Label))
		b.WriteString("\n")
		if item.Subtitle != "" {
			b.WriteString(subtitleStyle.Render(item.Subtitle))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (f *ListField[T]) Blur() tea.Cmd {
	f.focused = false
	return nil
}

func (f *ListField[T]) Focus() tea.Cmd {
	f.focused = true
	return nil
}

func (f *ListField[T]) Error() error {
	return f.err
}

func (f *ListField[T]) Run() error {
	return huh.Run(f)
}

func (f *ListField[T]) RunAccessible(w io.Writer, r io.Reader) error {
	// Simple accessible mode: list options and read selection
	fmt.Fprintln(w, f.title)
	for i, item := range f.items {
		fmt.Fprintf(w, "%d. %s", i+1, item.Label)
		if item.Subtitle != "" {
			fmt.Fprintf(w, " - %s", item.Subtitle)
		}
		fmt.Fprintln(w)
	}

	var selection int
	fmt.Fprint(w, "Enter selection number: ")
	_, err := fmt.Fscanln(r, &selection)
	if err != nil {
		return err
	}

	if selection < 1 || selection > len(f.items) {
		return fmt.Errorf("invalid selection: %d", selection)
	}

	f.cursor = selection - 1
	f.selected = true
	if f.value != nil {
		*f.value = f.items[f.cursor].Value
	}

	return nil
}

func (f *ListField[T]) Skip() bool {
	return len(f.items) == 0
}

func (f *ListField[T]) Zoom() bool {
	return false
}

func (f *ListField[T]) KeyBinds() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("up/k", "move up")),
		key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("down/j", "move down")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}
}

func (f *ListField[T]) WithTheme(theme *huh.Theme) huh.Field {
	f.theme = theme
	return f
}

func (f *ListField[T]) WithAccessible(accessible bool) huh.Field {
	f.accessible = accessible
	return f
}

func (f *ListField[T]) WithKeyMap(keymap *huh.KeyMap) huh.Field {
	f.keymap = keymap
	return f
}

func (f *ListField[T]) WithWidth(width int) huh.Field {
	f.width = width
	return f
}

func (f *ListField[T]) WithHeight(height int) huh.Field {
	f.height = height
	return f
}

func (f *ListField[T]) WithPosition(position huh.FieldPosition) huh.Field {
	f.position = position
	return f
}

func (f *ListField[T]) GetKey() string {
	return f.key
}

func (f *ListField[T]) GetValue() any {
	if f.value != nil {
		return *f.value
	}
	if len(f.items) > 0 {
		return f.items[f.cursor].Value
	}
	return nil
}

// IsSelected returns true if the user has made a selection.
func (f *ListField[T]) IsSelected() bool {
	return f.selected
}

// SelectedItem returns the currently selected item.
func (f *ListField[T]) SelectedItem() *Item[T] {
	if len(f.items) == 0 {
		return nil
	}
	return &f.items[f.cursor]
}
