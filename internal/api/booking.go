package api

import (
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

func (a *Api) GetAvailableRooms(wAuth, wAuthRefresh string, payload QuickBookPayload) error {

	type DateTimePayload struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	dtp := DateTimePayload{
		Start: payload.DateTime.Format(time.RFC3339),
		End:   payload.DateTime.Add(time.Duration(payload.Duration)).Format(time.RFC3339),
	}

	jsonValue, err := json.Marshal(dtp)

	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/CoworkingSpace/%s/category/%s/items", apiUrl, spaceId, categoryId)

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return err
	}

	q := req.URL.Query()

	q.Add("capacity", strconv.Itoa(payload.NbPeople))
	q.Add("datewithhours", string(jsonValue))

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("w_auth=%s; w_auth_refresh=%s", wAuth, wAuthRefresh))

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	response := AvailableRoomsResponse{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	// rooms := []models.Room{}

	return nil
}
