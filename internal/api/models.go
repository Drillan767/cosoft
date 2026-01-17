package api

import "time"

type LoginPayload struct {
	Email    string
	Password string
}

type QuickBookPayload struct {
	DateTime time.Time
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

type CreditResponse struct {
	ParsedValue float64 `json:"parsedValue"`
}

type PriceResponse struct {
	Credits CreditResponse `json:"Credits"`
}

type RoomResponse struct {
	Id      string          `json:"Id"`
	Name    string          `json:"Name"`
	NbUsers int             `json:"NbUsers"`
	Prices  []PriceResponse `json:"Prices"`
}

type AvailableRoomsResponse struct {
	VisitedItems   []RoomResponse `json:"VisitedItems"`
	UnvisitedItems []RoomResponse `json:"UnvisitedItems"`
}
