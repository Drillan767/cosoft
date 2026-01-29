# COSOFT CLI

This TUI tool's goal is to simplify any meeting room's booking without having to use the actual website.

## How it works

1. You login for the 1st time using your Cosoft credentials
2. The JWT is stored locally in `~/.cosoft`, alongside other files such as the list of available rooms (to avoid keeping fetching them) and your favorites among them.
3. Once logged in, you can perform several actions

## Main functions

### Calendar view

Display an ASCII representation of a calendar, displaying when rooms are used and if you're the one booking them.

### Quick book

Will book the 1st available room depending on if the room is available right now, and available long enough for you to book it with the duration you picked.

If one of the rooms is your favorite, will pick it for you.

Once the booking is done, you'll get a fancy table that will summarize the details of the booking:

- Room name
- Start time -> end time

### Browse and book

Will allow you to pick a date, time and duration, that will display all rooms available with these filters.

Once the booking is done, you'll get a fancy table that will summarize the details of the booking:

- Room name
- Start time -> end time

### Upcoming reservations (With dynamic number)

Will list all future bookings you've done. The list contains all details about the reservation such as the meeting room, date and time of booking.

It will also allow you to cancel it.

### Previous reservations

Will allow you to have an history of your past reservations, but selecting one of them will allow you to book the specific room again, but you'll have to pick a new date, time and duration for this.


## CLI

All these functions are available directly from the main command itself.

| Command                 | Function                  |
|-------------------------|---------------------------|
| `quick-book`            | Book first available room |
| `browse`                | Browse & book             |
| `options`               | Settings                  |
| `resa` / `reservations` | My reservations           |
| `history`               | Previous reservations     |
