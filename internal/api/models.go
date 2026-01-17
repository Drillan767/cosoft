package api

type LoginPayload struct {
	Email    string
	Password string
}

type QuickBookPayload struct {
	Duration int
	NbPeople int
}

type BrowsePayload struct {
	Room      uint
	StarDate  string
	StartHour string
	EndDate   string
	Duration  int
}
