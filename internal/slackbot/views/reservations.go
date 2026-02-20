package views

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/ui/slack"
	"fmt"
	"time"
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
	if r.Error != nil {
		return slack.Block{
			Blocks: []slack.BlockElement{
				slack.NewContext(*r.Error),
			},
		}
	}

	blocks := slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewHeader("Mes réservations"),
			slack.NewMrkDwn(fmt.Sprintf("*%d* réservation(s) à venir", len(*r.Reservations))),
		},
	}

	var list []slack.BlockElement

	for _, r := range *r.Reservations {
		parsedStart, err := time.Parse("2006-01-02T15:04:05", r.Start)

		if err != nil {
			fmt.Println(err)
			return slack.Block{}
		}

		parsedEnd, err := time.Parse("2006-01-02T15:04:05", r.End)

		if err != nil {
			fmt.Println(err)
			return slack.Block{}
		}

		duration := parsedEnd.Sub(parsedStart).Minutes()

		paidPrice := r.Credits * (float64(duration) / 60)

		dateFormat := "02/01/2006 15:04"
		list = append(list, slack.BlockElement(slack.NewMenuItem(
			fmt.Sprintf(
				"*%s*\n%s → %s · %.02f crédits",
				r.ItemName,
				parsedStart.Format(dateFormat),
				parsedEnd.Format(dateFormat),
				paidPrice,
			),
			"Annuler",
			r.OrderResourceRentId,
		)))
	}

	blocks.Blocks = list

	return blocks
}
