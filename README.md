# COSOFT CLI

This TUI tool's goal is to simplify any meeting room's booking without having to use the actual website.

## How it works

1. You log in for the 1st time using your Cosoft credentials
2. Your authentication information are stored in a SQLite database in your user config directory. This location changes
based on your OS.
3. Once logged in, you can perform several actions:

## Main functions

### Calendar view

Display an ASCII representation of a calendar, displaying when rooms are used and if you're the one booking them.

### Quick book

Will book the 1st available room depending on if the room is available right now, and available long enough for you to
book it with the duration you picked.

Once the booking is done, you'll get a fancy table that will summarize the details of the booking:

- Room name
- Start time → end time
- Credits paid for the booking

### Browse and book

Will allow you to pick a date, time and duration, and a room size that will display all rooms available with these filters.

Once the booking is done, you'll get a fancy table that will summarize the details of the booking:

- Room name
- Start time → end time
- Credits spent

### Upcoming reservations (With dynamic number)

Will list all future bookings you've done. The list contains all details about the reservation such as the meeting room, date and time of booking.

It will also allow you to cancel it.

### Previous reservations (Coming soon)

Will allow you to have a history of your past reservations, but selecting one of them will allow you to book the specific room again, but you'll have to pick a new date, time and duration for this.

## Non-interactive booking

Will allow you to book a meeting room with various options:

| Parameter | Shortcut | default | Description                                                                          |
|-----------|----------|---------|--------------------------------------------------------------------------------------|
| capacity  | c        | 1       | Will filter rooms by size                                                            |
| name      | n        |         | Will book a specific room if available. Run `./cosoft rooms` to see what's available |
| time      | t        |         | If provided, will book your room at the desired time.                                |
| duration  | d        | 30      | Indicates the booking's duration. Must be between 30 and 120.                        |

## CLI

All these functions are available directly from the main command itself.

| Command   | Function                                            |
|-----------|-----------------------------------------------------|
| *(empty)* | Displays the interactive menu                       |
| `book`    | Non interactive booking with parameters (see above) |
| `rooms`   | List all available rooms                            |


# Installation

1. Download the compiled binary of your choice at the latest release available [here](https://github.com/Drillan767/cosoft/releases).
2. In the directory where the binary has been downloaded run the following commands:

Linux / MacOS:

```bash
chmod +x  meeting-room-booker-{os}-{processor} # Allows the program to be run
sudo mv meeting-room-booker-{os}-{processor} /usr/local/bin/cosoft
```

You can now run `cosoft` from anywhere in your terminal.

Windows:

Feel free to rename the program as you like and place it anywhere you'll remember, then:

1. Search for "cmd" in the start menu
2. Once the terminal opened, change the directory to get to where the program is
3. Type `meeting-room-booker-windows-amd64.exe` and hit the "enter" key.