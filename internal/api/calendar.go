package api

import (
	"bytes"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (a *Api) GetRoomBusyTime(
	wAuth, wAuthRefresh string,
	roomId string,
	date time.Time,
) (*[]models.UnavailableSlot, error) {

	type filter struct {
		StartDate time.Time `json:"startDate"`
		EndDate   time.Time `json:"endDate"`
	}

	a.debug(date.Format(time.RFC3339))
	a.debug(roomId)

	endDate := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		23,
		59,
		0,
		0,
		time.UTC,
	)
	payload := filter{
		StartDate: date,
		EndDate:   endDate,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		a.debug("failed to marshal payload")
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/CoworkingSpace/%s/category/%s/item/%s/busytimes", apiUrl, spaceId, categoryId, roomId)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))

	if err != nil {
		a.debug("failed to create request")
		return nil, err
	}

	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("w_auth=%s; w_auth_refresh=%s", wAuth, wAuthRefresh))

	resp, err := client.Do(req)

	if err != nil {
		a.debug("failed to do request")
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	a.debug(string(data))

	var slots []models.UnavailableSlot
	if err := json.Unmarshal(data, &slots); err != nil {
		a.debug("failed to unmarshal response: " + err.Error())
		return nil, err
	}

	return &slots, nil
}
