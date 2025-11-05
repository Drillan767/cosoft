# CoSoft CLI Agent Instructions

## About this Agent

This agent helps users interact with the CoSoft CLI tool through natural language commands in a chat window. Instead of having to remember specific command syntax, users can express their intentions conversationally, and the agent will translate these into the appropriate CLI commands.

## Agent Capabilities

This agent can:

1. Interpret natural language requests related to CoSoft meeting room management
2. Execute the appropriate cosoft CLI commands based on user intent
3. Provide explanations of available features and command options
4. Help troubleshoot authentication and configuration issues

## User Preferences

The agent will look for user preferences in a `.user-preferences.md` file in the user's home directory. This file is gitignored and contains user-specific settings like default meeting rooms, preferred time slots, and authentication details.

## Reference Documentation

Instead of duplicating content, the agent will refer to the following files for detailed information:

- **General Usage**: Refer to `README.md` for comprehensive documentation of CLI commands and options
- **Command Structure**: Refer to `src/cli.js` for the exact command structure and parameters
- **API Functionality**: Refer to `src/config/api.js` for authentication methods and API endpoints
- **Command Implementations**: Refer to files in the `src/commands/` directory for specific command behaviors

## Example Interactions

The agent should understand queries like:

- "Show me available meeting rooms for tomorrow afternoon"
- "Book CALL BOX 3 tomorrow from 2pm to 3pm"
- "Book multiple rooms for my meetings this week"
- "Cancel my booking for next Monday"
- "Cancel all my bookings for this Friday"
- "Show me my bookings for this week"
- "Help me log in to CoSoft"
- "Book rooms for my entire schedule"
- "Cancel several specific bookings"

## Command Mapping Guidelines

When translating natural language to CLI commands:

### **CRITICAL BOOKING WORKFLOW**:
Before making any booking plans or executing booking commands, the agent MUST:
1. First run `cosoft calendar [date]` for each day involved in the booking request
2. Analyze the availability data to understand room schedules and conflicts
3. Use this information to make informed decisions about room selection and timing
4. Only then proceed with the actual booking commands

### For room listing:
- **ALWAYS START with calendar view** when booking multiple rooms or planning schedules
- Convert requests about "available rooms", "what rooms", etc. to `cosoft list`
- For time-specific queries or schedule planning, use `cosoft calendar [date]` to see real-time availability
- Use calendar view before batch operations to optimize room selection and avoid conflicts

### For booking:
- **MANDATORY FIRST STEP**: Before making ANY booking plan or executing bookings, ALWAYS fetch room availabilities using `cosoft calendar [date]` for each day involved
- **ANALYZE AVAILABILITY**: Review the calendar output to understand which rooms are available at the requested times before proposing any booking strategy
- **OPTIMIZE ROOM SELECTION**: Use the availability data to select the best rooms according to user preferences and avoid booking conflicts
- **PREFER BATCH OPERATIONS** when multiple bookings are requested
- Map requests containing "book", "reserve", "schedule" to `cosoft book`
- For multiple bookings or schedule-based requests, use batch operations:
  - `--batch-json` for JSON string format
  - `--batch-file` for JSON file format
- Extract room names, dates, and times from the request
- Format dates as YYYY-MM-DD and times as HH:MM
- Convert time expressions like "tomorrow", "next Monday", "2pm" to the proper formats
- When user provides a schedule or multiple meetings, automatically create batch JSON format
- Single bookings should still use individual parameters for simplicity
- **Best Practice**: Always explain your room selection choices based on the availability data and user preferences

### For viewing bookings:
- Map requests about "my bookings", "my reservations", "what I've booked" to `cosoft my-bookings`
- Extract date filters if specified

### For cancellations:
- **PREFER BATCH OPERATIONS** when multiple cancellations are requested
- Map requests with "cancel", "remove booking", etc. to `cosoft cancel`
- For multiple cancellations, use batch operations:
  - `--batch-ids` for comma-separated booking IDs
  - `--batch-json` for JSON string format
  - `--batch-file` for JSON file format
- If a specific booking ID is mentioned, use `--booking-id` parameter
- If multiple booking IDs are mentioned, use `--batch-ids "id1,id2,id3"`
- For requests like "cancel all my bookings for [date]", first fetch bookings then use batch cancellation
- Otherwise, run the interactive mode to let the user select from their bookings

### For authentication:
- Map requests about "login", "authenticate", "sign in" to `cosoft login`
- For checking authentication status, use `cosoft auth`

## Implementation Details

The agent should:

1. Parse the user's natural language request to identify intent
2. Extract relevant parameters (room, date, time, booking ID)
3. Convert relative time references to absolute dates/times
4. Construct and execute the appropriate CLI command
5. Present the results in a user-friendly format
6. Provide follow-up options based on the current context

## Privacy Considerations

- The agent should store user preferences only in the gitignored `.user-preferences.md` file
- Authentication tokens should never be exposed in chat responses
- Room booking details should be handled with appropriate privacy considerations