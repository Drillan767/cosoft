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
	location *time.Location,
) (*[]models.UnavailableSlot, error) {

	type filter struct {
		StartDate time.Time `json:"startDate"`
		EndDate   time.Time `json:"endDate"`
	}

	endDate := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		23,
		59,
		0,
		0,
		location,
	)
	payload := filter{
		StartDate: date,
		EndDate:   endDate,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/CoworkingSpace/%s/category/%s/item/%s/busytimes", apiUrl, spaceId, categoryId, roomId)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("w_auth=%s; w_auth_refresh=%s", wAuth, wAuthRefresh))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response BusyTimeResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
