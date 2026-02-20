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

type BookCmd struct {
	NbPeople   int
	Duration   int
	Datetime   time.Time
	PickedRoom models.Room
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
				Value string `json:"value"`
			} `json:"selected_option"`
		} `json:"duration"`
	} `json:"duration"`
	NbPeople struct {
		NbPeople struct {
			Type           string `json:"type"`
			SelectedOption struct {
				Value string `json:"value"`
			} `json:"selected_option"`
		} `json:"nbPeople"`
	} `json:"nbPeople"`
}

type PickedRoomPayload struct {
	PickRoom struct {
		PickRoom struct {
			Type           string `json:"type"`
			SelectedOption struct {
				Value string `json:"value"`
			} `json:"selected_option"`
		} `json:"pick-room"`
	} `json:"pick-room"`
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

		parsedDt, err := b.criteriaToTime()

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

		nbPeople, duration, err := b.filtersToNumber()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return b, &BrowseCmd{
			NbPeople: nbPeople,
			Duration: duration,
			Datetime: *parsedDt,
		}
	} else if action.ActionID == "pick-room" {
		var pickedRoom PickedRoomPayload

		err := json.Unmarshal(action.Values, &pickedRoom)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		roomId := pickedRoom.PickRoom.PickRoom.SelectedOption.Value

		for _, r := range *b.Rooms {
			if r.Id == roomId {
				b.PickedRoom = &r
				break
			}
		}

		if b.PickedRoom == nil {
			fmt.Println("PickedRoom not found")
			return b, nil
		}
	} else if action.ActionID == "book" {
		nbPeople, duration, _ := b.filtersToNumber()
		t, _ := b.criteriaToTime()

		return b, &BookCmd{
			NbPeople:   nbPeople,
			Duration:   duration,
			PickedRoom: *b.PickedRoom,
			Datetime:   *t,
		}
	} else if action.ActionID == "back" {
		b.Phase = 0
		return b, nil
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
						"*Aucune salle disponible*\nVeuillez changer vos filtres",
						"Retour",
						"back",
					),
				},
			}
		}

		nbPeople, duration, _ := b.filtersToNumber()
		t, _ := b.criteriaToTime()
		end := t.Add(time.Duration(duration) * time.Minute)
		choices := make([]slack.ChoicePayload, len(*b.Rooms))

		for i, room := range *b.Rooms {
			choices[i] = slack.ChoicePayload{
				Text:  room.Name,
				Value: room.Id,
			}
		}

		blocks := []slack.BlockElement{
			slack.NewHeader(fmt.Sprintf("%d salles ont été trouvées", len(*b.Rooms))),
			slack.NewMrkDwn(fmt.Sprintf(
				"%s → %s — %d personnes",
				t.Format("02/01/2006 15:04"),
				end.Format("02/01/2006 15:04"),
				nbPeople,
			)),
			slack.NewDivider(),
			slack.NewSelect(
				"Sélectionez une salle",
				"Salle",
				"pick-room",
				choices,
			),
		}

		if b.PickedRoom != nil {
			blocks = append(
				blocks,
				slack.NewPreview(
					fmt.Sprintf("*%s*\n%.2f crédits", b.PickedRoom.Name, b.PickedRoom.Price),
					b.PickedRoom.Image,
					b.PickedRoom.Name,
				),
				slack.NewButtons([]slack.ChoicePayload{{"Réserver", "book"}}),
			)
		}

		blocks = append(blocks,
			slack.NewDivider(),
			slack.NewButtons([]slack.ChoicePayload{
				{"Retour à l'accueil", "cancel"},
				{"Modifier les filtres", "back"},
			}),
		)

		return slack.Block{
			Blocks: blocks,
		}

	case 2:
		duration, _ := strconv.Atoi(b.Duration)
		startTime, _ := b.criteriaToTime()
		endTime := startTime.Add(time.Duration(duration) * time.Minute)
		dateFormat := "02/01/2006 15:04"
		paidPrice := b.PickedRoom.Price * (float64(duration) / 60)

		return slack.Block{
			Blocks: []slack.BlockElement{
				slack.BlockElement(slack.NewMrkDwn(":white_check_mark: *Réservation réussie !*")),

				slack.BlockElement(slack.NewMultiMarkdown([]string{
					fmt.Sprintf("*Salle de réunion :*\n%s", b.PickedRoom.Name),
					fmt.Sprintf("*Durée :*\n%s → %s", startTime.Format(dateFormat), endTime.Format(dateFormat)),
					fmt.Sprintf("*Coût :*\n%.2f credits", paidPrice),
				})),
				slack.BlockElement(slack.NewMenuItem(
					"Vous pouvez maintenant revenir à l'accueil",
					"Retour",
					"cancel",
				)),
			},
		}
	}
	return slack.Block{}
}

func (b *BrowseView) criteriaToTime() (*time.Time, error) {
	dt := fmt.Sprintf("%s %s", b.Date, b.Time)

	parsedDt, err := time.Parse("2006-01-02 15:04", dt)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &parsedDt, nil
}

func (b *BrowseView) filtersToNumber() (int, int, error) {
	nbPeople, err := strconv.Atoi(b.NbPeople)

	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}

	duration, err := strconv.Atoi(b.Duration)

	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}

	return nbPeople, duration, nil
}
