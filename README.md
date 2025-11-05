# CoSoft CLI

CLI tool to interact with the CoSoft API, and more easily manage meeting rooms reservations.

## Features

- Interactive mode with calendar view and context-aware menus
- Direct command execution for scripting and quick actions
- List available meeting rooms with details
- Book meeting rooms with comprehensive options
- View your upcoming bookings
- Cancel existing reservations
- Calendar view of meeting room availability
- Authentication status check with logged-in user display
- Interactive login with email and password

## Installation

```bash
# Clone the repository
git clone [repository-url]
cd cosoft

# Install dependencies
npm install

# Link the CLI globally
npm link
```

## Usage

### Interactive Mode

By default, running the CLI without any commands starts it in interactive mode:

```bash
cosoft
```

This opens the calendar view and provides a menu-driven interface to navigate through all features.

### Direct Commands

You can also execute specific commands directly:

```bash
# Show help
cosoft --help

# View calendar (default: today)
cosoft calendar

# View calendar for a specific date
cosoft calendar 2025-09-15

# List available meeting rooms
cosoft list

# Book a meeting room (interactive)
cosoft book

# Book a meeting room (direct with parameters)
cosoft book --room-name "CALL BOX 3" --date 2025-09-03 --start-time 13:30 --end-time 14:30

# Book multiple rooms from JSON file
cosoft book --batch-file examples/batch-bookings.json

# Book multiple rooms from JSON string
cosoft book --batch-json '[{"roomName":"CALL BOX 3","date":"2025-09-03","startTime":"13:30","endTime":"14:30"},{"roomName":"Salle RHUBIK","date":"2025-09-03","startTime":"15:00","endTime":"16:00"}]'

# View your bookings (all dates)
cosoft my-bookings

# View your bookings for a specific date
cosoft my-bookings 2025-09-15

# Cancel a reservation (interactive)
cosoft cancel

# Cancel a specific reservation by ID
cosoft cancel --booking-id "40adf2d0-3a3b-45cd-bb47-b34c011201c1"

# Cancel multiple bookings by IDs
cosoft cancel --batch-ids "id1,id2,id3"

# Cancel multiple bookings from JSON file
cosoft cancel --batch-file examples/batch-cancellations.json

# Cancel multiple bookings from JSON string
cosoft cancel --batch-json '["booking-id-1","booking-id-2"]'

# Check authentication status
cosoft auth

# Login to CoSoft
cosoft login

# Start as MCP server (Model Context Protocol)
cosoft mcp-server
```

### Command Options

```
Usage: cosoft [command] [options]

Commands:
  calendar [date]         Show the calendar view for meeting rooms
  book                    Book a meeting room or batch of rooms
    --room-name <name>    Name of the room to book
    --date <date>         Date in YYYY-MM-DD format
    --start-time <time>   Start time in HH:MM format
    --end-time <time>     End time in HH:MM format
    --batch-file <path>   Path to JSON file containing batch booking data
    --batch-json <json>   JSON string containing batch booking data
  list                    List all available meeting rooms
  my-bookings [date]      View your bookings
  cancel                  Cancel an existing booking or batch of bookings
    --booking-id <id>     ID of the booking to cancel
    --batch-file <path>   Path to JSON file containing batch cancellation data
    --batch-json <json>   JSON string containing batch cancellation data
    --batch-ids <ids>     Comma-separated list of booking IDs to cancel
  auth                    Check authentication status and display logged-in user
  login                   Login to CoSoft and save authentication tokens
  mcp-server              Start as MCP (Model Context Protocol) server

Options:
  -i, --interactive       Run in interactive mode
  -h, --help              Display help information
  -v, --version           Display version information
```

## Authentication

The CLI requires authentication to access the CoSoft API. You can authenticate easily using the login command:

```bash
cosoft login
```

This command will interactively prompt you for your email and password, then save your authentication tokens securely.

The CLI also uses a `.env` file for non-sensitive configuration:

```
COSOFT_API_BASE_URL=https://hub612.cosoft.fr
COSOFT_SPACE_ID=your-space-id
COSOFT_CATEGORY_ID=your-category-id
```

## Batch Operations

### Batch Booking

You can book multiple rooms at once using batch operations:

**JSON File Format for Batch Bookings (`examples/batch-bookings.json`):**
```json
[
  {
    "roomName": "CALL BOX 3",
    "date": "2025-09-06",
    "startTime": "09:00",
    "endTime": "09:30"
  },
  {
    "roomName": "Salle RHUBIK", 
    "date": "2025-09-06",
    "startTime": "10:00",
    "endTime": "11:00"
  }
]
```

### Batch Cancellation

You can cancel multiple bookings at once:

**JSON File Format for Batch Cancellations (`examples/batch-cancellations.json`):**
```json
{
  "bookingIds": [
    "booking-id-1",
    "booking-id-2",
    "booking-id-3" 
  ]
}
```

**Alternative format (array of IDs):**
```json
[
  "booking-id-1",
  "booking-id-2",
  "booking-id-3"
]
```

## MCP Server Mode

The CoSoft CLI can run as a Model Context Protocol (MCP) server, enabling integration with other applications via stdio-based JSON message communication.

### Starting the MCP Server

```bash
cosoft mcp-server
```

The server listens for JSON messages on stdin and responds with JSON on stdout.

### Available MCP Commands

- `list_tools` - Get list of available commands and their parameters
- `list` - List all available meeting rooms
- `book` - Book a meeting room (requires roomName, date, startTime, endTime)
- `cancel` - Cancel a booking by ID (requires bookingId)
- `my-bookings` - List your current bookings (optional date filter)
- `calendar` - View calendar for a date (optional date parameter)

### MCP Message Format

**Request:**
```json
{
  "id": "req-123",
  "type": "request", 
  "command": "book",
  "parameters": {
    "roomName": "CALL BOX 3",
    "date": "2025-10-05",
    "startTime": "14:00",
    "endTime": "15:30"
  }
}
```

**Response:**
```json
{
  "id": "req-123",
  "type": "response",
  "status": "success",
  "result": {
    "bookingId": "generated-id",
    "roomName": "CALL BOX 3",
    "date": "2025-10-05",
    "startTime": "14:00",
    "endTime": "15:30"
  }
}
```

## Development

```bash
# Run the CLI
npm start

# Run linting
npm run lint
```