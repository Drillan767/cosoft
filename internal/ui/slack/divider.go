package slack

type Divider struct {
	Type string `json:"type"`
}

func (Divider) blockElement() {}

func NewDivider() Divider {
	return Divider{"divider"}
}
