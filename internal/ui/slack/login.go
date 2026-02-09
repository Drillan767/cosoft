package slack

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
