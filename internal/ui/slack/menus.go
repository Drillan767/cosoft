package slack

import (
	"cosoft-cli/internal/storage"
	"fmt"
)

var (
	durationChoices = []ChoicePayload{
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

	nbPeopleChoices = []ChoicePayload{
		{
			"Une personne",
			"1",
		},
		{
			"Deux personnes ou plus",
			"2",
		},
	}
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
			NewMenuItem(
				"*Mes réservations*\nVoir et annuler vos futures réservations",
				"Accéder",
				"reservations",
			),
		},
	}
}

func QuickBookMenu() Block {
	return Block{
		Blocks: []BlockElement{
			NewHeader("Réservation rapide"),
			NewSelect(
				"Durée de la réservation",
				"Sélectionner",
				"duration",
				durationChoices,
			),
			NewSelect(
				"Capacité",
				"Sélectionner",
				"nbPeople",
				nbPeopleChoices,
			),
			NewButtons([]ChoicePayload{{"Annuler", "cancel"}, {"Réserver", "quick-book"}}),
		},
	}
}

func BrowseMenu() Block {
	return Block{
		Blocks: []BlockElement{
			NewHeader("Réserver une salle"),
			NewDatePicker("Date", "date", "Date"),
			NewTimePicker("*Heure*\nFormat autorisé : 15h00, 15:15", "time", "Heure"),
			NewSelect(
				"Durée de la réservation",
				"Sélectionner",
				"duration",
				durationChoices,
			),
			NewSelect(
				"Capacité",
				"Sélectionner",
				"nbPeople",
				nbPeopleChoices,
			),
			NewButtons([]ChoicePayload{{"Annuler", "cancel"}, {"Voir les salles disponibles", "browse"}}),
		},
	}
}
