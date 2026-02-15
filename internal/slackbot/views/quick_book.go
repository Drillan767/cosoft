package views

import (
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
)

type QuickBookView struct {
	Phase      int
	NbPeople   string
	Duration   string
	Rooms      *[]models.Room
	PickedRoom *models.Room
	Error      *string
}

type QuickBookCmd struct {
	NbPeople int
	Duration int
	Rooms    []models.Room
}

type QuickBookValues struct {
	Duration struct {
		Duration struct {
			SelectedOption struct {
				Value string `json:"value"`
			} `json:"selected_option"`
		} `json:"duration"`
	} `json:"duration"`
	NbPeople struct {
		NbPeople struct {
			SelectedOption struct {
				Value string `json:"value"`
			} `json:"selected_option"`
		} `json:"nbPeople"`
	} `json:"nbPeople"`
}

func (qb *QuickBookView) Update(action Action) (View, Cmd) {

	switch action.ActionID {
	case "cancel":
		return qb, &LandingCmd{}
	case "quick-book":
		var values QuickBookValues

		err := json.Unmarshal(action.Values, &values)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		qb.Duration = values.Duration.Duration.SelectedOption.Value
		qb.NbPeople = values.NbPeople.NbPeople.SelectedOption.Value
		qb.Error = nil

		if qb.NbPeople == "" || qb.Duration == "" {
			s := ":warning: Tous les champs sont requis"
			qb.Error = &s

			return qb, nil
		}

		nbPeople, err := strconv.Atoi(qb.NbPeople)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		duration, err := strconv.Atoi(qb.Duration)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Everything looks clear, submitting.
		return qb, &QuickBookCmd{
			NbPeople: nbPeople,
			Duration: duration,
		}

	default:
		return qb, nil

	}
}

func RenderQuickBookView(qb *QuickBookView) slack.Block {
	blocks := slack.QuickBookMenu()

	switch qb.Phase {
	case 0:
		if qb.Error != nil {
			blocks.Blocks = slices.Insert(
				blocks.Blocks,
				3,
				slack.BlockElement(slack.NewContext(*qb.Error)),
			)
		}

		return blocks

	case 2:
		blocks.Blocks = slices.Insert(
			blocks.Blocks,
			len(blocks.Blocks),
			slack.BlockElement(slack.NewMrkDwn(":large_green_circle: Une salle a été trouvée !")),
			slack.BlockElement(slack.NewMrkDwn("Réservation en cours...")),
		)

		return blocks
	case 3:
		// duration, _ := strconv.Atoi(qb.Duration)
		// startTime := common.GetClosestQuarterHour()
		// endTime := startTime.Add(time.Duration(duration) * time.Minute)
		// dateFormat := "02/01/2006 15:04"
		// paidPrice := qb.PickedRoom.Price * (float64(duration) / 60)

		blocks.Blocks = slices.Insert(
			blocks.Blocks,
			len(blocks.Blocks),
			slack.BlockElement(slack.NewDivider()),
			slack.BlockElement(slack.NewMrkDwn(":white_check_mark: *Réservation réussie !*")),
			/*
				slack.BlockElement(slack.NewMultiMarkdown([]string{
					fmt.Sprintf("*Salle de réunion :*\n%s", qb.PickedRoom.Name),
					fmt.Sprintf("*Durée :\n%s → %s", startTime.Format(dateFormat), endTime.Format(dateFormat)),
					fmt.Sprintf("*Coût :*\n%.2f credits", paidPrice),
				})),
				slack.BlockElement(slack.NewMenuItem(
					"Vous pouvez maintenant revenir à l'accueil",
					"Retour",
					"cancel",
				)),
			*/
		)

		return blocks
	default:
		return slack.QuickBookMenu()
	}
}
