package common

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/shared/models"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func BuildCalendar(
	maxLabelLength, displayedHours int,
	rooms []models.RoomUsage,
	userBookings []api.Reservation,
) []string {
	for _, room := range rooms {
		if len(room.Name)+1 > maxLabelLength {
			maxLabelLength = len(room.Name) + 1
		}
	}

	rows := make([]string, len(rooms)+1)

	rows[0] = createCalendarHeader(maxLabelLength, displayedHours)

	for i, room := range rooms {
		rows[i+1] = createCalendarRow(room, maxLabelLength, userBookings)
	}

	return rows
}

func createCalendarHeader(labelLength, displayedHours int) string {
	spacing := 2
	result := ""

	for i := 0; i < displayedHours; i++ {
		if i+8 < 10 {
			result += "0"
		}
		result += fmt.Sprintf("%dh%s", i+8, strings.Repeat(" ", spacing))
	}

	return strings.Repeat(" ", labelLength-1) + result
}

func createCalendarRow(
	row models.RoomUsage,
	labelLength int,
	userBookings []api.Reservation,
) string {
	type parsedSlot struct {
		Start time.Time
		End   time.Time
	}

	location, _ := LoadLocalTime()

	var slots []parsedSlot
	spacing := labelLength - len(row.Name)
	columns := ""

	for _, slot := range row.UsedSlots {
		start, _ := time.ParseInLocation("2006-01-02T15:04:05", slot.Start, location)
		end, _ := time.ParseInLocation("2006-01-02T15:04:05", slot.End, location)

		slots = append(slots, parsedSlot{
			Start: start,
			End:   end,
		})
	}

	var userSlots []parsedSlot
	for _, slot := range userBookings {
		if slot.ItemName != row.Name {
			// User booking not matching current row, skipping.
			continue
		}

		start, _ := time.ParseInLocation("2006-01-02T15:04:05", slot.Start, location)
		end, _ := time.ParseInLocation("2006-01-02T15:04:05", slot.End, location)
		userSlots = append(userSlots, parsedSlot{
			Start: start,
			End:   end,
		})
	}

	now := GetClosestQuarterHour()

	year, month, day := slots[0].Start.Date()
	baseDate := time.Date(year, month, day, 0, 0, 0, 0, location)
	startTime := baseDate.Add(8 * time.Hour)
	endTime := baseDate.Add(23 * time.Hour)
	counter := 0

	current := startTime

	for !current.After(endTime) {
		occupied := false
		ownReservation := false
		slotEnd := current.Add(15 * time.Minute)

		for _, slot := range slots {
			if current.Before(slot.End) && slotEnd.After(slot.Start) {
				occupied = true
				break
			}
		}

		for _, uSlot := range userSlots {
			if current.Before(uSlot.End) && slotEnd.After(uSlot.Start) {
				ownReservation = true
				break
			}
		}

		symbol := " "

		if occupied {
			symbol = "░"
		}

		if ownReservation {
			symbol = "█"
		}

		isNow := current.Equal(now)
		nextSlot := current.Add(15 * time.Minute)
		nextIsNow := nextSlot.Equal(now)

		if isNow {
			// If current time, color the cell's background in red,
			symbol = lipgloss.NewStyle().Background(lipgloss.Color("#f45656")).Render(symbol)
		}

		if counter%4 == 3 && !nextIsNow {
			// If not current time, simply display a normal pipe.
			symbol += "│"
		}

		columns += symbol
		counter++

		current = current.Add(15 * time.Minute)
	}

	return row.Name + strings.Repeat(" ", spacing) + "│" + columns
}
