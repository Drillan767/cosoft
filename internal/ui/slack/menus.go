package slack

func MainMenu() Block {
	// Note: will certainly require a slack user id as parameter to retrieve the correct infos
	return Block{
		Blocks: []BlockElement{
			NewMrkDwn("Si vous voyez ceci, c'est que vous vous êtes connecté(e) avec succès."),
		},
	}
}
