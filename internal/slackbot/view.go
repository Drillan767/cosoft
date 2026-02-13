package slackbot

import (
	"cosoft-cli/internal/ui/slack"
	"encoding/json"
)

type View interface {
	Update(state json.RawMessage) View
}

type slackInput struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

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
	if l.Email != "" && l.Password != "" {
		// TODO: login
		return &MainView{}
	}

	s := "Tous les champs sont requis"
	l.Error = &s
	return l
}

type MainView struct {
}

func (m *MainView) Update(state json.RawMessage) View {
	var tmp struct{}

	/*
		Suivant l'action qui a été cliquée, retourner la vue correspondante
		- Si quickbook, renvoyer un future QuickBookView{}
	*/
}

func RenderView(v View) slack.Block {
	switch v := v.(type) {
	case *LoginView:
		return RenderLoginView(v)
	case *MainView:
		return slack.Block{}
	default:
		return slack.Block{}
	}
}

func RenderLoginView(l *LoginView) slack.Block {
	// TODO: display error
	return slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewMrkDwn(":information_source:  Pour réserver une salle, il faut d'abord vous identifier."),
			slack.NewInput("Email", "email"),
			slack.NewInput("Mot de passe", "password"),
			slack.NewContext(":warning: Le mot de passe est affiché en clair dans le champ"),
			slack.NewButtons([]slack.ChoicePayload{{"Connexion", "login"}}),
		},
	}
}
