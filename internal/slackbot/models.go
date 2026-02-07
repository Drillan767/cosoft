package slackbot

type Request struct {
	UserId      string
	Command     string
	Text        string
	ResponseUrl string
}
