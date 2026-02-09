package services

import (
	"bytes"
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (s *SlackService) IsSlackAuthenticated(request models.Request) bool {
	cookies, err := s.store.HasActiveToken(&request.UserId)

	if err != nil || cookies == nil {
		return false
	}

	apiClient := api.NewApi()
	err = apiClient.GetAuth(cookies.WAuth, cookies.WAuthRefresh)

	return err == nil
}

func (s *SlackService) DisplayLogin(request models.Request) {
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
		return
	}

	type LoginWrapper struct {
		TriggerId string      `json:"trigger_id"`
		View      slack.Modal `json:"view"`
	}

	loginForm := slack.NewLogin(request.ResponseUrl)

	loginWrapper := LoginWrapper{
		TriggerId: request.TriggerId,
		View:      loginForm,
	}

	jsonBlocks, err := json.Marshal(loginWrapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/views.open", bytes.NewBuffer(jsonBlocks))

	if err != nil {
		fmt.Println(err)
		return
	}

	accessToken := os.Getenv("SLACK_ACCESS_TOKEN")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	_, err = http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *SlackService) LogInUser(email, password, slackUserId, responseUrl string) error {
	apiClient := api.NewApi()

	loginPayload := api.LoginPayload{
		Email:    email,
		Password: password,
	}

	response, err := apiClient.Login(&loginPayload)

	if err != nil {
		return err
	}

	return s.postLogin(*response, slackUserId, responseUrl)
}

func (s *SlackService) postLogin(user api.UserResponse, slackUserId, responseUrl string) error {

	// Save / update user in database
	err := s.store.SetUser(
		&user,
		user.JwtToken,
		user.RefreshToken,
		&slackUserId,
	)

	if err != nil {
		return err
	}

	// Display main menu with user's info
	storedUser := storage.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Credits:   user.Credits,
	}

	mainMenu := slack.MainMenu(storedUser)

	jsonMenu, err := json.Marshal(mainMenu)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", responseUrl, bytes.NewBuffer(jsonMenu))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)

	return err
}
