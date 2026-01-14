package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserResponse struct {
	JwtToken  string  `json:"JwtToken"`
	Id        string  `json:"Id"`
	FirstName string  `json:"FirstName"`
	LastName  string  `json:"LastName"`
	Email     string  `json:"Email"`
	Credits   float64 `json:"Credits"`
}

type AuthPayload struct {
	IsAuth  bool          `json:"isAuth"`
	Message string        `json:"Message"`
	User    *UserResponse `json:"User"`
}

func (a *Api) Login(payload *LoginPayload) (*UserResponse, error) {
	values := map[string]string{
		"email":    payload.Email,
		"password": payload.Password,
	}

	jsonValues, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(
		baseUrl+"/v2/api/api/users/login",
		"application/json",
		bytes.NewBuffer(jsonValues),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	response := AuthPayload{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if !response.IsAuth || response.User == nil {
		return nil, fmt.Errorf("Wrong username / password")
	}

	return response.User, nil
}
