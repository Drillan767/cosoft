package ui

import tea "github.com/charmbracelet/bubbletea"

type NavigateMsg struct {
	Page string
}

type BackToMenuMsg struct{}

type AppModel struct {
	currentPage    PageType
	allowBackNav   bool
	quitting       bool
	landingModel   *LandingModel
	quickBookModel *QuickBookModel
	// Add others
}

type PageType string

const (
	PageLanding      PageType = "landing"
	PageQuickBook    PageType = "quick-book"
	PageBrowse       PageType = "browse"
	PageReservations PageType = "resa"
	PageSettings     PageType = "settings"
)

func NewAppModel(startPage string, allowBackNav bool) *AppModel {
	return &AppModel{
		currentPage:    PageType(startPage),
		allowBackNav:   allowBackNav,
		landingModel:   NewLandingModel(),
		quickBookModel: NewQuickBookModel(),
		// Add others
	}
}

func (m *AppModel) Init() tea.Cmd {
	switch m.currentPage {
	case "landing":
		return m.landingModel.Init()
	case "quick-book":
		return m.quickBookModel.Init()
	default:
		return m.landingModel.Init()
	}
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NavigateMsg:
		// Handle "quit" selection
		if msg.Page == "quit" {
			m.quitting = true
			return m, tea.Quit
		}

		m.currentPage = PageType(msg.Page)
		switch msg.Page {
		case "quick-book":
			return m, m.quickBookModel.Init()
		default:
			return m, nil
		}
	case BackToMenuMsg:
		if m.allowBackNav {
			m.currentPage = PageLanding
			return m, m.landingModel.Init()
		}
	case tea.KeyMsg:
		// Handle ctrl+c to quit immediately
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		// Handle ESC BEFORE forwarding to child models
		if msg.String() == "esc" && m.allowBackNav && m.currentPage != PageLanding {
			m.currentPage = PageLanding
			return m, m.landingModel.Init()
		}
	}

	// Forward all messages to the current page's model
	switch m.currentPage {
	case PageLanding:
		newModel, cmd := m.landingModel.Update(msg)
		m.landingModel = newModel.(*LandingModel)
		return m, cmd
	case PageQuickBook:
		newModel, cmd := m.quickBookModel.Update(msg)
		m.quickBookModel = newModel.(*QuickBookModel)
		return m, cmd
	// Add other pages here as you create them
	}

	return m, nil
}

func (m *AppModel) View() string {
	// Return empty view when quitting to avoid showing borders
	if m.quitting {
		return ""
	}

	switch m.currentPage {
	case PageLanding:
		return m.landingModel.View()
	case PageQuickBook:
		return m.quickBookModel.View()
	// Others
	default:
		return "unknown page"
	}
}
