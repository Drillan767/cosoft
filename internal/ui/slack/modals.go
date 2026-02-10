package slack

type ModalWrapper struct {
	TriggerId string `json:"trigger_id"`
	View      Modal  `json:"view"`
}

type ModalTitle struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ModalAction struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Modal struct {
	Type            string         `json:"type"`
	CallbackId      string         `json:"callback_id"`
	Title           ModalTitle     `json:"title"`
	Blocks          []BlockElement `json:"blocks"`
	Close           ModalAction    `json:"close"`
	Submit          ModalAction    `json:"submit"`
	PrivateMetadata string         `json:"private_metadata"`
}

func NewLogin(responseUrl string) Modal {
	return Modal{
		Type:       "modal",
		CallbackId: "login_modal",
		Title: ModalTitle{
			Type: "plain_text",
			Text: "Connexion",
		},
		Close: ModalAction{
			Type: "plain_text",
			Text: "Fermer",
		},
		Submit: ModalAction{
			Type: "plain_text",
			Text: "Connexion",
		},
		PrivateMetadata: responseUrl,
		Blocks: []BlockElement{
			NewMrkDwn(":information_source:  Pour réserver une salle, il faut d'abord vous identifier."),
			NewInput("Email", "email"),
			NewInput("Mot de passe", "password"),
			NewContext(":warning: Le mot de passe est affiché en clair dans le champ"),
		},
	}
}

func NewQuickbook(responseUrl string) Modal {
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

	return Modal{
		Type:       "modal",
		CallbackId: "quickbook_modal",
		Title: ModalTitle{
			Type: "plain_text",
			Text: "Réservation rapide",
		},
		Close: ModalAction{
			Type: "plain_text",
			Text: "Fermer",
		},
		Submit: ModalAction{
			Type: "plain_text",
			Text: "Réserver une salle",
		},
		PrivateMetadata: responseUrl,
		Blocks: []BlockElement{
			NewRadio("Durée de réservation", "duration", durationChoices),
			NewDivider(),
			NewRadio("Taille de la salle", "nbPeople", nbPeopleChoices),
		},
	}
}
