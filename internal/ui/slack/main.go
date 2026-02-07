package slack

type BlockElement interface {
	blockElement()
}

type Block struct {
	Blocks []BlockElement `json:"blocks"`
}
