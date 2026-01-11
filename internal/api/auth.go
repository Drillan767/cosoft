package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (a *Api) Login(payload *LoginPayload) error {
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

	fmt.Println(resp.Body)

	return nil
}

func (a *Api) LoginFromRefreshToken() {
	// ...
}
