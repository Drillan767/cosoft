# General workflow

## Main display

- Calendar view of upcoming reservations in grid
- Menu with quick actions
    - Quick book
    - Browse & book
    - Upcoming reservations (3)
    - Previous reservations
    - Settings
    - Quit

### Calendar view

1. Fetch all rooms (from cache)
2. Fetch YOUR bookings for today
3. Fetch busy times for each room (parallel API calls)
4. Build availability matrix [room][timeslot]
5. Render grid

```
CoSoft - Today (Friday, Jan 10)

           08:00   09:00   10:00   11:00   12:00   13:00   14:00   15:00   16:00   17:00   18:00   19:00
  CALL BOX 1  │     │█████│█████│     │·····│·····│     │     │     │     │     │     │
  CALL BOX 2  │     │     │=====│=====│·····│·····│=====│=====│     │     │     │     │
  CALL BOX 3  │█████│█████│     │     │·····│·····│     │     │█████│     │     │     │
  CALL BOX 5  │     │     │     │=====│·····│·····│     │     │     │     │     │     │
```

Color/Pattern scheme:
```golang
// Terminal-friendly options
const (
    YourBooking   = "█"  // Solid block (with color)
    OtherBooking  = "▓"  // Medium shade (different color)
    PastTime      = "░"  // Light shade (gray)
    Available     = " "  // Empty space
)

// internal/ui/calendar.go
func RenderCalendar(view CalendarView) string {
    var buf strings.Builder

    // Header with time labels
    buf.WriteString("         ")
    for hour := 8; hour <= 19; hour++ {
        buf.WriteString(fmt.Sprintf("%-8s", fmt.Sprintf("%02d:00", hour)))
    }
    buf.WriteString("\n")

    // Room rows
    for _, room := range view.Rooms {
        buf.WriteString(fmt.Sprintf("%-13s", room.Name))
        for _, slot := range room.TimeSlots {
            buf.WriteString(renderSlot(slot))
        }
        buf.WriteString("\n")
    }

    // Legend
    buf.WriteString("\nLegend: █ Your bookings  │ ▓ Others  │ ░ Past  │   Available\n")

    return buf.String()
}

func renderSlot(slot TimeSlot) string {
    switch slot.Status {
    case YourBooking:
        return color.GreenString("█████")
    case OtherBooking:
        return color.BlueString("▓▓▓▓▓")
    case Past:
        return color.New(color.FgHiBlack).Sprint("░░░░░")
    default:
        return "     "
    }
}

// Parallel fetch
func fetchAllBusyTimes(rooms []Room, date time.Time) (map[string][]BusyTime, error) {
    var wg sync.WaitGroup
    results := make(map[string][]BusyTime)
    var mu sync.Mutex

    for _, room := range rooms {
        wg.Add(1)
        go func(r Room) {
            defer wg.Done()
            busyTimes, err := api.GetBusyTimes(r.ID, date)
            if err != nil {
                log.Warn("Failed to fetch busy times for", r.Name)
                return
            }
            mu.Lock()
            results[r.ID] = busyTimes
            mu.Unlock()
        }(room)
    }

    wg.Wait()
    return results, nil
}
```

## Book first available room

- User selects the duration (30mn / 1h)
- Load all rooms informations / bookings
- Check if any is free for the next span of time previously picked, right now
- If any of the user's favorite rooms are available, pick those first
- Proceed to do the booking / payment on the fly
- Display a table with booking's informations

## Browse & book

User will go through a multistep form

- Step 1
    - User picks the date and time they want to make a booking (default is today, rounded to the closest quarter (14:04 => show rooms available from 14:00 as well))
- Step 2
    - User picks a room from the filtered result and will be able to set one of them as favorite (this will update the yml file)
- Step 3
    - User picks a duration time
- Result
    - Proceed to do the booking / payment on the fly
    - Display a table with booking's informations

## Default options

User will be able to define the default duration for his booking. This will simply select the option by default but the user will still be able to pick something else

- Favorite rooms
- Prefered duration: 30mn / 1h / 2h

## Upcoming reservations

Will display all reservations as a list. The list itself will contain all needed informations, such as date, time (09:00 -> 10:00) and meeting room. Selecting a date (then confirming) will cancel the reservation

## Previous reservations

**Note** : This requires to know if such information is available.

Would display a list of reservations the user made. When picking one, user will be able to book the related meeting room again, only at date / times where the room is available

## Quit

Take a wild fucking guess


# CLI actions

All the options from the main menu will also be available as a CLI command to skip a step if you already know what you need

| Command                 | Function                  |
|-------------------------|---------------------------|
| `quick-book`            | Book first available room |
| `browse`                | Browse & book             |
| `options`               | Settings                  |
| `resa` / `reservations` | My reservations           |
| `history`               | Previous reservations     |
| `calendar [--date?]`    | Displays the calendar     |

# Slack bot

HUB612 staff members have expressed their desire to be able to book a room through a Slack bot, which makes sense since the CLI format might be hard for them to use. This project might be able to concile everything, but some considerations need to made.

## Note
- The element ensuring something's available needs to be its own function, as it needs to be reused in the several instances
- Rooms can be booked retroactively up to 30mn prior to current hour

## Structure

```golang
type User struct {
    Id string
    FirstName string
    LastName string
    Email string
    Jwt string
    Credits int
    ExpiresAt time.Time
    SlackUserID *string
    CreatedAt time.Time
}

type Room struct {
    Id string
    Name string
    MaxUsers int
    Price int
    CreatedAt time.Time
}

type Reservation strut {
    Id uint
    Date time.Time
    Room Room
    User User
    Duration int
    Cost int
    CreatedAt time.Time
}

```

---

# Authentication flow:

1. User types `./cosoft` in terminal
2. `root.go` is triggered through `main.go` and will try to display "landing" (root.go:28)
3. However, a PreRunE hook will ensure that the user is logged in first (root.go:25)
4. The hook consists of loading the json file and check that the timestamp next to the jwt token is valid (auth/service.go:39-44) and will return true if everything went smoothly
5. If not valid / not existing, then display the login form (root:61)
6. When / if the form is complete and valid (ui/login:145), display the spinner and start trying to log the user (ui/login:44)
7. If auth fails, display the error, remove spinner and reset the form (ui/login:105)
8. If auth succeeds (ui/login:99), remove the loader and stop the form.
9. Back at `root.go`, we ensure that everything went ok, then try to get the user, then actually store the user once it's retrieve (root:77)
10. Since the requirement no longer blocks the flow, `root.go` can now run the "landing" script (root:28)
11. By the way, the `StartApp()` function is the one responsible for display data in the layout (ui/main:24)

# Cancellation

You **cannot** cancel a reservation that has already started

URL:

```
https://hub612.cosoft.fr/v2/api/api/Reservation/cancel-order
```

Payload:

```
{
    Id: d378b30e-a30b-4267-aee7-b3d601030b1b // OrderResourceRentId
}
```
