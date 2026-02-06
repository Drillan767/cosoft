package slackbot

type Request struct {
	UserId      string `json:"user_id"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseUrl string `json:"response_url"`
}
