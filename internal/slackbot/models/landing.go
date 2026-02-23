package models

import (
	"cosoft-cli/internal/storage"
	"fmt"
)

type LandingState struct {
	User *storage.User

	nextCount int
}

func NewLandingState(store *storage.Store, userID string) (State, error) {
	_, err := store.UpdateCredits(&userID)
	if err != nil {
		return &LoginState{
			Error: ptr(errInternalError),
		}, fmt.Errorf("update credits: %v", err)
	}

	user, err := store.GetUserData(&userID)
	if err != nil {
		return &LoginState{
			Error: ptr(errInternalError),
		}, fmt.Errorf("get user data: %v", err)
	}

	return &LandingState{
		User: user,
	}, nil
}

func (s *LandingState) Type() string { return landingStateType }

func (s *LandingState) Update(store *storage.Store, params UpdateParams) (State, error) {
	switch params.ActionID {
	case "browse":
		return &BrowseState{}, nil
	case "calendar":
		return newCalendarState(), nil
	case "quick-book":
		return &QuickBookState{}, nil
	case "reservations":
		return newReservationState(store, params.UserID)
	default:
		return s, nil
	}
}

func (s *LandingState) Next() bool {
	s.nextCount++
	return s.nextCount < 2
}
