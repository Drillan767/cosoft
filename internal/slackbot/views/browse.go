package views

import (
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

type BrowseView struct {
	Phase      int
	NbPeople   string
	Duration   string
	Date       string
	Time       string
	Rooms      *[]models.Room
	PickedRoom *models.Room
	Error      *string
}

type BrowseCmd struct {
	NbPeople int
	Duration int
	Datetime time.Time
	Rooms    []models.Room
}

type BrowsePayload struct {
	Time struct {
		Time struct {
			Type         string `json:"type"`
			SelectedTime string `json:"selected_time"`
		} `json:"time"`
	} `json:"time"`
	Date struct {
		Date struct {
			Type         string `json:"type"`
			SelectedDate string `json:"selected_date"`
		} `json:"date"`
	} `json:"date"`
	Duration struct {
		Duration struct {
			Type           string `json:"type"`
			SelectedOption struct {
				Text struct {
					Type  string `json:"type"`
					Text  string `json:"text"`
					Emoji bool   `json:"emoji"`
				} `json:"text"`
				Value string `json:"value"`
			} `json:"selected_option"`
		} `json:"duration"`
	} `json:"duration"`
	NbPeople struct {
		NbPeople struct {
			Type           string `json:"type"`
			SelectedOption struct {
				Text struct {
					Type  string `json:"type"`
					Text  string `json:"text"`
					Emoji bool   `json:"emoji"`
				} `json:"text"`
				Value string `json:"value"`
			} `json:"selected_option"`
		} `json:"nbPeople"`
	} `json:"nbPeople"`
}

func (b *BrowseView) Update(action Action) (View, Cmd) {
	if action.ActionID == "cancel" {
		return b, &LandingCmd{}
	} else if action.ActionID == "browse" {
		var values BrowsePayload

		err := json.Unmarshal(action.Values, &values)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		b.Duration = values.Duration.Duration.SelectedOption.Value
		b.NbPeople = values.NbPeople.NbPeople.SelectedOption.Value

		b.Date = values.Date.Date.SelectedDate
		b.Time = values.Time.Time.SelectedTime

		if b.NbPeople == "" || b.Duration == "" {
			s := ":warning: Tous les champs sont requis"
			b.Error = &s

			fmt.Println(s)

			return b, nil
		}

		// If untouched from Slack, the fields won't be filled even if valid.
		// So we set (again) their default value.
		if b.Date == "" {
			b.Date = time.Now().Format(time.DateOnly)
		}
		if b.Time == "" {
			b.Time = common.GetClosestQuarterHour().Format("15:04")
		}

		dt := fmt.Sprintf("%s %s", b.Date, b.Time)

		parsedDt, err := time.Parse("2006-01-02 15:04", dt)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if parsedDt.Before(time.Now()) {
			s := ":warning: Veuillez choisir une date dans le futur"
			b.Error = &s

			return b, nil
		}

		if parsedDt.Minute()%15 != 0 {
			s := ":warning: Veuillez choisir un quart d'heure (14h00, 15h15, 16h30, 17h45....)"
			b.Error = &s

			return b, nil
		}

		nbPeople, err := strconv.Atoi(b.NbPeople)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		duration, err := strconv.Atoi(b.Duration)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return b, &BrowseCmd{
			NbPeople: nbPeople,
			Duration: duration,
			Datetime: parsedDt,
		}

	} else if strings.HasPrefix(action.ActionID, "book-") {
		// A room has been picked
		// Return this for now.
		return b, nil
	}

	return b, nil
}

func RenderBrowseView(b *BrowseView) slack.Block {
	switch b.Phase {
	case 0:
		blocks := slack.BrowseMenu()
		if b.Error != nil {
			blocks.Blocks = slices.Insert(
				blocks.Blocks,
				5,
				slack.BlockElement(slack.NewContext(*b.Error)),
			)
		}

		return blocks
	case 1:
		if len(*b.Rooms) == 0 {
			return slack.Block{
				Blocks: []slack.BlockElement{
					slack.NewMenuItem(
						"*Aucune salle disponible\nVeuillez changer vos filtres*",
						"Retour",
						"cancel",
					),
				},
			}
		}

		return slack.Block{
			Blocks: []slack.BlockElement{
				slack.NewMrkDwn(fmt.Sprintf("%d salles ont été trouvées", len(*b.Rooms))),
			},
		}
	}
	return slack.Block{}
}
