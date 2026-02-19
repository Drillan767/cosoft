package views

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/ui/slack"
	"fmt"
)

type ReservationView struct {
	Phase         int
	Reservations  *[]api.Reservation
	ReservationId *string
	Error         *string
}

type ReservationCmd struct {
	Reservations  *[]api.Reservation
	ReservationId *string
}

func (r *ReservationView) Update(action Action) (View, Cmd) {
	fmt.Println(action.ActionID)
	if action.ActionID == "cancel" {
		return r, &LandingCmd{}
	}

	return r, nil
}

func RenderReservationsView(r *ReservationView) slack.Block {
	return slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewHeader("Réservations à venir"),
		},
	}
}
