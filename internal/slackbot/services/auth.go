package services

import (
	"bytes"
	"cosoft-cli/internal/api"
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

func (s *SlackService) LogInUser(email, password, responseUrl string) (*api.UserResponse, error) {
	apiClient := api.NewApi()

	loginPayload := api.LoginPayload{
		Email:    email,
		Password: password,
	}

	response, err := apiClient.Login(&loginPayload)

	if err != nil {
		return nil, err
	}

	s.postLogin(responseUrl)

	return response, nil
}

func (s *SlackService) postLogin(responseUrl string) {
	/*
		1. Send blocks payload to display the main menu to the response url
		2. Create / update the user in the database based on the response
	*/
}
