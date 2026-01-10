# COSOFT CLI

This TUI tool's goal is to simplify any meeting room's booking without having to use the actual website.

## Main functions

### Calendar view

## On startup 

- Check if refresh token available through `zalando/go-keyring`
    - If exist, attempt connect
    - If connect fails or nothing found in keyring, display login screen

- Check if the directory `~/.cosoft-cli` exists and have both `settings.yml` and `rooms.json`
    - If exist, load and parse the informations
    - If not, create the files

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