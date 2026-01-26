package api

import (
	"cosoft-cli/shared/models"
	"time"

	"github.com/google/uuid"
)

type LoginPayload struct {
	Email    string
	Password string
}

type CosoftAvailabilityPayload struct {
	DateTime time.Time
	Duration int
	NbPeople int
}

type CosoftBookingPayload struct {
	CosoftAvailabilityPayload
	Room        models.Room
	UserCredits float64
}

type BrowsePayload struct {
	Room      uint
	StarDate  string
	StartHour string
	EndDate   string
	Duration  int
	NbPeople  int
}

type PriceResponse struct {
	Credits float64 `json:"Credits"`
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

type DateTimePayload struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type AvailabilityBodyPayload struct {
	Capacity         int             `json:"capacity"`
	CategoryId       string          `json:"categoryId"`
	CoworkingSpaceId string          `json:"coworkingSpaceId"`
	DateTime         DateTimePayload `json:"datewithhours"`
}

type DateTimeAlt struct {
	Date  string            `json:"date"`
	Times []DateTimePayload `json:"times"`
}

type DateTime struct {
	Start      string    `json:"start"`
	End        string    `json:"end"`
	Type       string    `json:"type"`
	TimeSlotId any       `json:"timeSlotId"` // Will always be null
	Id         uuid.UUID `json:"id"`
}

type RoomBookingCartPayload struct {
	CoworkingSpaceId string      `json:"coworkingSpaceId"`
	CategoryId       string      `json:"categoryId"`
	ItemId           string      `json:"itemId"`
	CartId           string      `json:"cartId"`
	DateTimeAlt      DateTimeAlt `json:"startenddate_"`
	DateTime         []DateTime  `json:"startenddate"`
}

type RoomBookingPayload struct {
	IsUser           bool                     `json:"isUser"`
	IsPerson         bool                     `json:"isPerson"`
	IsVatRequired    bool                     `json:"isVatRequired"`
	IsStatusRequired bool                     `json:"isStatusRequired"`
	CGV              bool                     `json:"cgv"`
	SocietyName      string                   `json:"societyname"`
	SocietyVat       string                   `json:"societyvat"`
	SocietySiret     string                   `json:"societysiret"`
	SocietyStatus    string                   `json:"societystatus"`
	FirstName        string                   `json:"firstname"`
	LastName         string                   `json:"lastname"`
	Address          string                   `json:"address"`
	City             string                   `json:"city"`
	ZipCode          string                   `json:"zipCode"`
	Phone            string                   `json:"phone"`
	Email            string                   `json:"email"`
	Cart             []RoomBookingCartPayload `json:"cart"`
	PaymentType      string                   `json:"paymentType"`
}

type Reservation struct {
	OrderResourceRentId string  `json:"OrderResourceRentId"`
	ItemName            string  `json:"ItemName"`
	Start               string  `json:"Start"`
	End                 string  `json:"End"`
	Credits             float64 `json:"Credits"`
}

type FutureBookingsResponse struct {
	Total int           `json:"total"`
	Data  []Reservation `json:"data"`
}

type CancellationPayload struct {
	Id string `json:"Id"`
}

type BusyTimeResponse struct {
	Data    []models.UnavailableSlot `json:"data"`
	Error   string                   `json:"Error"`
	Message string                   `json:"Message"`
	Fields  string                   `json:"Fields"`
}
