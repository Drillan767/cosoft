package api

type LoginPayload struct {
	Email    string
	Password string
}

type QuickBookPayload struct {
	Duration int
}

type BrowePayload struct {
	Room      uint
	StarDate  string
	StartHour string
	EndDate   string
	Duration  int
}
