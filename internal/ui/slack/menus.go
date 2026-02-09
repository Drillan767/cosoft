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
				"Réserver immédiatement une salle de réunion",
				"Accéder",
				"quick-book",
				"quick-book",
			),
		},
	}
}
