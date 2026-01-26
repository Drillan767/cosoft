package api

import (
	"fmt"
	"io"
	"net/http"
)

type Api struct{}

var (
	apiUrl     = "https://hub612.cosoft.fr/v2/api/api"
	spaceId    = "a4928a70-38c1-42b9-96f9-b2dd00db5b02"
	categoryId = "7f1e5757-b9b9-4530-84ad-b2dd00db5f0f"
)

func NewApi() *Api {
	return &Api{}
}

func (a *Api) prepareHeaderCookies(
	wAuth, wAuthRefresh, method, endpoint string,
	payload io.Reader,
) (*http.Request, *http.Client, error) {
	req, err := http.NewRequest(method, endpoint, payload)

	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{}

	// A token and a refresh token need to be added in the request's header to make an authenticated request.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("w_auth=%s; w_auth_refresh=%s", wAuth, wAuthRefresh))

	return req, client, nil
}
