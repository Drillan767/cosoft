package api

import (
	"bytes"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DateTimePayload struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type AvailabilityBodyPayload struct {
	Capacity         int             `json:"capacity"`
	CategoryId       string          `json:"categoryId"`
	CoworkingSpaceId string          `json:"coworkingSpaceId"`
	DateTime         DateTimePayload `json:"datewithhours"`
}

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

func (a *Api) GetAvailableRooms(wAuth, wAuthRefresh string, payload QuickBookPayload) ([]models.Room, error) {

	req, err := a.prepareRoomAvailabilityRequest(payload)
	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("w_auth=%s; w_auth_refresh=%s", wAuth, wAuthRefresh))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := AvailableRoomsResponse{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	rooms := make([]models.Room, 0, len(response.UnvisitedItems)+len(response.VisitedItems))
	for _, room := range response.VisitedItems {

		// Filter out Hubmit and Clhub
		if room.NbUsers > 10 {
			continue
		}

		mr := models.Room{
			Id:      room.Id,
			Name:    room.Name,
			NbUsers: room.NbUsers,
			Price:   room.Prices[0].Credits,
		}

		rooms = append(rooms, mr)
	}

	for _, room := range response.UnvisitedItems {

		// Filter out Hubmit and Clhub
		if room.NbUsers > 10 {
			continue
		}

		mr := models.Room{
			Id:      room.Id,
			Name:    room.Name,
			NbUsers: room.NbUsers,
			Price:   room.Prices[0].Credits,
		}

		rooms = append(rooms, mr)
	}

	return rooms, nil
}

func (a *Api) prepareRoomAvailabilityRequest(payload QuickBookPayload) (*http.Request, error) {
	dtp := DateTimePayload{
		Start: payload.DateTime.Format(time.RFC3339),
		End:   payload.DateTime.Add(time.Duration(payload.Duration) * time.Minute).Format(time.RFC3339),
	}

	abp := AvailabilityBodyPayload{
		Capacity:         payload.NbPeople,
		CategoryId:       categoryId,
		CoworkingSpaceId: spaceId,
		DateTime:         dtp,
	}

	jsonDtp, err := json.Marshal(dtp)

	if err != nil {
		return nil, err
	}

	jsonAbp, err := json.Marshal(abp)

	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/CoworkingSpace/%s/category/%s/items", apiUrl, spaceId, categoryId)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonAbp))

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("capacity", strconv.Itoa(payload.NbPeople))
	q.Add("datewithhours", string(jsonDtp))
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (a *Api) debug(text string) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	file.WriteString(fmt.Sprintf("%s \n\n", text))
}
