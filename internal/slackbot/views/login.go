package views

import (
	"cosoft-cli/internal/ui/slack"
	"encoding/json"
)

type LoginView struct {
	Email    string
	Password string
	Error    *string
}

func (l *LoginView) Update(state json.RawMessage) View {
	// TODO: parse state
	var tmp struct {
		Password slackInput `json:"password"`
		Email    slackInput `json:"email"`
	}
	err := json.Unmarshal(state, &tmp)
	if err != nil {
		// TODO: log error and start from scratch
		return l
	}

	l.Email = tmp.Email.Value
	l.Password = tmp.Password.Value
	/*
		if l.Email != "" && l.Password != "" {
			// TODO: login
			return &MainView{}
		}
	*/

	s := "Tous les champs sont requis"
	l.Error = &s
	return l
}

func RenderLoginView(l *LoginView) slack.Block {
	loginBlocks := []slack.BlockElement{
		slack.NewMrkDwn(":information_source:  Pour réserver une salle, il faut d'abord vous identifier."),
		slack.NewInput("Email", "email"),
		slack.NewInput("Mot de passe", "password"),
		slack.NewContext(":warning: Le mot de passe est affiché en clair dans le champ"),
		slack.NewButtons([]slack.ChoicePayload{{"Connexion", "login"}}),
	}

	if l.Error != nil {
		loginBlocks[3] = slack.NewContext(*l.Error)
	}

	blocks := slack.Block{
		ResponseType: "ephemeral",
		Blocks:       loginBlocks,
	}

	return blocks
}
