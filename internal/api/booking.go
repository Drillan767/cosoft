package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *Api) GetFutureBookings(wAuth, wAuthRefresh string) (int, error) {

	type FutureBookingsResponse struct {
		Total int `json:"total"`
	}

	req, err := http.NewRequest(
		"GET",
		apiUrl+"/Reservations/get-current-and-incoming?PerPage=5&Page=1",
		nil,
	)

	if err != nil {
		return 0, err
	}

	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("w_auth=%s; w_auth_refresh=%s", wAuth, wAuthRefresh))

	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	response := FutureBookingsResponse{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	return response.Total, nil
}
