package api

type LoginPayload struct {
	Email    string
	Password string
}

type QuickBookPayload struct {
	Duration int
}
