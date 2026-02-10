package slack

type BlockElement interface {
	blockElement()
}

type Block struct {
	ReplaceOriginal bool           `json:"replace_original,omitempty"`
	ResponseType    string         `json:"response_type,omitempty"`
	Blocks          []BlockElement `json:"blocks"`
}

type BlockPayload struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Emoji bool   `json:"emoji"`
}
