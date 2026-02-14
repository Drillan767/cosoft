package views

import (
	"cosoft-cli/internal/ui/slack"
	"encoding/json"
	"fmt"
)

type LoginView struct {
	Email    string
	Password string
	Error    *string
}

type LoginCmd struct {
	Email    string
	Password string
}

type Values struct {
	Email struct {
		Email struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"email"`
	} `json:"email"`
	Password struct {
		Password struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"password"`
	} `json:"password"`
}

func (l *LoginView) Update(action Action) (View, Cmd) {
	var values Values

	err := json.Unmarshal(action.Values, &values)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	l.Email = values.Email.Email.Value
	l.Password = values.Password.Password.Value
	l.Error = nil

	if l.Email == "" || l.Password == "" {
		s := ":warning: Tous les champs sont requis"
		l.Error = &s

		return l, nil
	}

	return l, &LoginCmd{
		Email:    l.Email,
		Password: l.Password,
	}
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
