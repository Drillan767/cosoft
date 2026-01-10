# General workflow

## On startup 

- Check if refresh token available through `zalando/go-keyring`
    - If exist, attempt connect
    - If connect fails or nothing found in keyring, display login screen

- Check if something exists at `~/.cosoft-cli/settings.yml`
    - If exist, load and parse the informations
    - If not, create the file

- Redirect to "main"

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

# Architecture and dependency injections

## Option A: Context-based (Simple)
```golang
  // cmd/root.go
var rootCmd = &cobra.Command{
    Use: "cosoft",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    // Initialize dependencies ONCE
        cfg := config.Load()
        authMgr := auth.NewManager()
        apiClient := api.NewClient(authMgr)
        bookingSvc := booking.NewService(apiClient, cfg)

        // Store in context
        ctx := context.WithValue(cmd.Context(), "config", cfg)
        ctx = context.WithValue(ctx, "bookingService", bookingSvc)
        cmd.SetContext(ctx)

        return nil
    },
}
```

```golang
// cmd/quick.go
var quickCmd = &cobra.Command{
    Use: "quick",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Retrieve from context
        svc := cmd.Context().Value("bookingService").(*booking.Service)

        return svc.QuickBook(duration)
    },
}
```

## Option B: Global App State (Cleaner for CLIs)

```golang
  // cmd/root.go
type App struct {
    Config      *config.Config
    AuthManager *auth.Manager
    API         *api.Client
    BookingSvc  *booking.Service
}

var app *App

var rootCmd = &cobra.Command{
    Use: "cosoft",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Initialize app dependencies
        app = &App{
            Config:      config.Load(),
            AuthManager: auth.NewManager(),
        }
        app.API = api.NewClient(app.AuthManager)
        app.BookingSvc = booking.NewService(app.API, app.Config)

        return nil
    },
}
```

```golang
// cmd/quick.go
var quickCmd = &cobra.Command{
    Use: "quick",
    RunE: func(cmd *cobra.Command, args []string) error {
        return app.BookingSvc.QuickBook(duration)
    },
}

func init() {
    rootCmd.AddCommand(quickCmd)
}
```

Option C: Constructor Injection (Most testable)

```golang
// cmd/quick.go
func NewQuickCmd(svc *booking.Service) *cobra.Command {
    return &cobra.Command{
        Use: "quick",
        RunE: func(cmd *cobra.Command, args []string) error {
            return svc.QuickBook(duration)
        },
    }
}
```

```golang
// cmd/root.go
func Execute() {
    cfg := config.Load()
    authMgr := auth.NewManager()
    apiClient := api.NewClient(authMgr)
    bookingSvc := booking.NewService(apiClient, cfg)

    rootCmd := &cobra.Command{Use: "cosoft"}
    rootCmd.AddCommand(NewQuickCmd(bookingSvc))
    rootCmd.AddCommand(NewBrowseCmd(bookingSvc))
    // ... etc

    rootCmd.Execute()
}
```

  My Recommendation: Option B (Global App State)

  Why?
  - Simple and pragmatic for CLI apps
  - No magic context values (type-safe)
  - Easy to understand flow
  - Cobra's PersistentPreRunE runs before ANY command
  - Testable enough (you can mock app.BookingSvc in tests)

  main.go stays minimal:
  package main

  import "cosoft-cli/cmd"

  func main() {
      cmd.Execute()
  }

  All initialization in cmd/root.go:

```golang
  package cmd

  import (
      "github.com/spf13/cobra"
      // ... your internal packages
  )

  type App struct {
      Config     *config.Config
      API        *api.Client
      BookingSvc *booking.Service
  }

  var app *App

  var rootCmd = &cobra.Command{
      Use:   "cosoft",
      Short: "CoSoft meeting room booking CLI",
      PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
          // Skip init for commands that don't need auth (like login)
          if cmd.Name() == "login" {
              return nil
          }

          // Initialize everything
          cfg, err := config.Load()
          if err != nil {
              return err
          }

          authMgr, err := auth.NewManager()
          if err != nil {
              return err
          }

          apiClient := api.NewClient(authMgr)

          app = &App{
              Config:     cfg,
              API:        apiClient,
              BookingSvc: booking.NewService(apiClient, cfg),
          }

          return nil
      },
  }

  func Execute() {
      if err := rootCmd.Execute(); err != nil {
          os.Exit(1)
      }
  }

  func init() {
      rootCmd.AddCommand(quickCmd)
      rootCmd.AddCommand(browseCmd)
      rootCmd.AddCommand(listCmd)
      // ... etc
  }
```