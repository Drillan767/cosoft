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

## Previous reservations

**Note** : This requires to know if such information is available.

Would display a list of reservations the user made. When picking one, user will be able to book the related meeting room again, only at date / times where the room is available

# CLI actions

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

## Note
- The element ensuring something's available needs to be its own function, as it needs to be reused in the several instances
- Rooms can be booked retroactively up to 30mn prior to current hour
