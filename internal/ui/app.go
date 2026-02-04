package ui

import tea "github.com/charmbracelet/bubbletea"

type NavigateMsg struct {
	Page string
}

type BackToMenuMsg struct{}

type AppModel struct {
	currentPage          PageType
	allowBackNav         bool
	quitting             bool
	landingModel         *LandingModel
	quickBookModel       *QuickBookModel
	browseModel          *BrowseModel
	reservationListModel *ReservationListModel
	settingsModel        *SettingsModel
	// Add others
}

type PageType string

const (
	PageLanding      PageType = "landing"
	PageQuickBook    PageType = "quick-book"
	PageBrowse       PageType = "browse"
	PageReservations PageType = "reservations"
	PageSettings     PageType = "settings"
)

func NewAppModel(startPage string, allowBackNav bool) *AppModel {
	return &AppModel{
		currentPage:          PageType(startPage),
		allowBackNav:         allowBackNav,
		landingModel:         NewLandingModel(),
		quickBookModel:       NewQuickBookModel(),
		browseModel:          NewBrowseModel(),
		reservationListModel: NewReservationListModel(),
		settingsModel:        NewSettingsModel(),
		// Add others
	}
}

func (m *AppModel) Init() tea.Cmd {
	switch m.currentPage {
	case "landing":
		return m.landingModel.Init()
	case "quick-book":
		return m.quickBookModel.Init()
	case "browse":
		return m.browseModel.Init()
	case "reservations":
		return m.reservationListModel.Init()
	case "settings":
		return m.settingsModel.Init()
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
			m.quickBookModel = NewQuickBookModel()
			return m, m.quickBookModel.Init()
		case "browse":
			m.browseModel = NewBrowseModel()
			return m, m.browseModel.Init()
		case "reservations":
			m.reservationListModel = NewReservationListModel()
			return m, m.reservationListModel.Init()
		case "settings":
			m.settingsModel = NewSettingsModel()
			return m, m.settingsModel.Init()
		default:
			return m, nil
		}
	case BackToMenuMsg:
		if m.allowBackNav {
			m.currentPage = PageLanding
			m.landingModel = NewLandingModel()
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
			m.landingModel = NewLandingModel()
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
	case PageBrowse:
		newModel, cmd := m.browseModel.Update(msg)
		m.browseModel = newModel.(*BrowseModel)
		return m, cmd
	case PageReservations:
		newModel, cmd := m.reservationListModel.Update(msg)
		m.reservationListModel = newModel.(*ReservationListModel)
		return m, cmd
	case PageSettings:
		newModel, cmd := m.settingsModel.Update(msg)
		m.settingsModel = newModel.(*SettingsModel)
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
	case PageBrowse:
		return m.browseModel.View()
	case PageReservations:
		return m.reservationListModel.View()
	case PageSettings:
		return m.settingsModel.View()
	// Others
	default:
		return "unknown page"
	}
}
