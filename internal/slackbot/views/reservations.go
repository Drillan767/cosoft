package views

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/ui/slack"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ReservationView struct {
	Phase             int
	Reservations      *[]api.Reservation
	PickedReservation *api.Reservation
	ReservationId     *string
	BookingStarted    bool
	Error             *string
}

type ReservationCmd struct {
	Reservations  *[]api.Reservation
	ReservationId *string
}

type CancelReservationCmd struct {
	ReservationId *string
}

func (r *ReservationView) Update(action Action) (View, Cmd) {
	if action.ActionID == "back" {
		return r, &LandingCmd{}
	}

	if action.ActionID == "cancel" {
		return r, &CancelReservationCmd{
			ReservationId: r.ReservationId,
		}
	}

	// Action id is the uuid of a selected booking.
	if err := uuid.Validate(action.ActionID); err == nil {
		r.ReservationId = &action.ActionID

		for _, reservation := range *r.Reservations {
			if reservation.OrderResourceRentId == action.ActionID {
				r.PickedReservation = &reservation
				break
			}
		}

		if r.PickedReservation == nil {
			fmt.Println("Could not find a picked reservation")
			return r, nil
		}

		// Ensure the reservation hasn't already started
		bookinStartsAt, err := time.ParseInLocation("2006-01-02T15:04:05", r.PickedReservation.Start, time.Local)

		if err != nil {
			fmt.Println(err)
			return r, nil
		}

		fmt.Println("started:", bookinStartsAt.Before(time.Now()), bookinStartsAt, time.Now())

		if bookinStartsAt.Before(time.Now()) {
			r.BookingStarted = true
		}
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

	switch r.Phase {
	case 0:
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
				"Sélectionner",
				r.OrderResourceRentId,
			)))
		}

		if r.PickedReservation != nil {

			if r.BookingStarted {
				list = append(
					list,
					slack.BlockElement(slack.NewContext(":warning: Impossible d'annuler une réservation déjà commencée")),
				)
			} else {
				list = append(
					list,
					slack.BlockElement(slack.NewDivider()),
					slack.BlockElement(slack.NewMenuItem(
						fmt.Sprintf(
							"Confirmer l'annulation de \"%s\" ?",
							r.PickedReservation.ItemName,
						),
						"Annuler réservation",
						"cancel",
					)),
				)
			}
		}

		list = append(
			list,
			slack.BlockElement(slack.NewDivider()),
			slack.NewButtons([]slack.ChoicePayload{{"Retour", "back"}}),
		)

		blocks.Blocks = list

		return blocks
	case 1:
		return slack.Block{
			Blocks: []slack.BlockElement{
				slack.NewHeader(":white_check_mark: Annulation réussie réussie !"),
				slack.NewMrkDwn("Vous pouvez maintenant retourner à l'accueil"),
				slack.NewDivider(),
				slack.NewButtons([]slack.ChoicePayload{{"Retour", "back"}}),
			},
		}
	default:
		return slack.Block{}
	}

}
