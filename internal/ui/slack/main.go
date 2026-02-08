package slack

type BlockElement interface {
	blockElement()
}

type Block struct {
	Blocks []BlockElement `json:"blocks"`
}

type BlockPayload struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Emoji bool   `json:"emoji"`
}
