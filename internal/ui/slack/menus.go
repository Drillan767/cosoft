package slack

import (
	"cosoft-cli/internal/storage"
	"fmt"
)

func MainMenu(user storage.User) Block {
	welcomeMessage := fmt.Sprintf(
		"Vous êtes connecté(e) en tant que *%s %s* (%s)",
		user.FirstName,
		user.LastName,
		user.Email,
	)

	creditsMessage := fmt.Sprintf("Il vous reste *%.2f* credits", user.Credits)

	return Block{
		Blocks: []BlockElement{
			NewMrkDwn(welcomeMessage),
			NewMrkDwn(creditsMessage),
			NewDivider(),
			NewHeader("Menu principal"),
			NewMenuItem(
				"*Réservation rapide*\nRéserver immédiatement une salle de réunion",
				"Accéder",
				"quick-book",
			),
			NewMenuItem(
				"*Parcourir les salles*\nRéserver une salle pour une date ultérieure",
				"Accéder",
				"browse",
			),
		},
	}
}

func QuickBookMenu() Block {
	durationChoices := []ChoicePayload{
		{
			"30 minutes",
			"30",
		},
		{
			"1 heure",
			"60",
		},
		{
			"1 heure 30",
			"90",
		},
		{
			"2 heures",
			"120",
		},
	}

	nbPeopleChoices := []ChoicePayload{
		{
			"Une personne",
			"1",
		},
		{
			"Deux personnes ou plus",
			"2",
		},
	}

	return Block{
		Blocks: []BlockElement{
			NewHeader("Réservation rapide"),
			NewRadio("Durée de réservation", "duration", durationChoices),
			NewDivider(),
			NewRadio("Taille de la salle", "nbPeople", nbPeopleChoices),
			NewButtons([]ChoicePayload{{"Annuler", "cancel"}, {"Réserver", "quick-book"}}),
		},
	}
}
