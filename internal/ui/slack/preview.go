package slack

type Preview struct {
	Type      string           `json:"type"`
	Text      BlockPayload     `json:"text"`
	Accessory PreviewAccessory `json:"accessory"`
}

type PreviewAccessory struct {
	Type     string `json:"type"`
	ImageUrl string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

func (Preview) blockElement() {}

func NewPreview(text, imageUrl, altText string) Preview {
	return Preview{
		Type: "section",
		Text: BlockPayload{
			Type: "mrkdwn",
			Text: text,
		},
		Accessory: PreviewAccessory{
			Type:     "image",
			ImageUrl: imageUrl,
			AltText:  altText,
		},
	}
}
