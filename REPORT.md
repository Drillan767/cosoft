# CoSoft CLI - Analysis and Observations

## What This Project Does

A Node.js CLI tool that wraps the CoSoft meeting room booking API to provide a better UX than their web interface. It supports:

- Viewing meeting room availability in a calendar view
- Booking single or multiple rooms (batch operations)
- Viewing and canceling reservations
- Interactive menu-driven navigation
- Direct command execution for scripting
- MCP (Model Context Protocol) server for integration with other tools

## Architecture Overview

### File Structure
```
src/
├── cli.js                    # Entry point, Commander.js setup
├── config/
│   ├── env.js               # Environment & auth token storage
│   └── api.js               # API client & auth logic
├── commands/
│   ├── login.js             # Authentication flow
│   ├── listRooms.js         # List available rooms
│   ├── calendarView.js      # Interactive calendar UI
│   ├── bookRoom.js          # Room booking (450+ lines)
│   ├── myBookings.js        # View user bookings
│   ├── cancelBooking.js     # Cancel bookings
│   └── mcpServer.js         # MCP protocol server
└── utils/
    └── sharedMenus.js       # Reusable navigation menus
```

### Technology Stack
- **CLI Framework:** Commander.js v14
- **HTTP Client:** Axios v1.11
- **UI Components:** @inquirer (prompts), chalk (colors), cli-table3 (tables)
- **Date Handling:** date-fns
- **Bundling:** esbuild + pkg for distribution

## How It Works

### 1. Authentication

**Storage:**
- JWT and refresh tokens stored in `.auth` file (project root)
- API config stored in `.env` (base URL, space ID, category ID)

**Flow:**
1. User runs `cosoft login`
2. Prompts for email/password
3. POST to `/v2/api/api/users/login`
4. Saves `JwtToken` and refresh token (from Set-Cookie header) to `.auth`
5. All subsequent requests include both tokens in cookie header

**Problem:** Refresh tokens are saved but never used to renew expired access tokens. If tokens expire mid-session, user gets generic error and must re-login.

### 2. API Client Pattern

**Request Handler (`makeApiRequest()`):**
```javascript
makeApiRequest({
  method: 'POST',
  path: '/v2/api/api/...',
  body: { ... },
  includeReferer: true  // Mimics browser behavior
})
```

**Defensive Response Parsing:**
The API is inconsistent, so code checks multiple response structures:
- Sometimes `responseData.data` (array)
- Sometimes direct array
- Sometimes nested object paths
- Error fields: `CartHasError`, `Error`, `ErrorMessage`, `DisabledItem`, `HasAlreadyOrdered`

This suggests the CoSoft API is poorly designed/documented.

### 3. Interactive Mode

**Navigation Flow:**
```
Entry → Calendar View (hub)
  ├─ Book → Standard Menu → Calendar
  ├─ Cancel → Standard Menu → Calendar
  ├─ My Bookings → Standard Menu → Calendar
  └─ List → Standard Menu → Calendar
```

**Implementation:**
- Commands accept `{ interactive: true }` option
- After command execution, shows menu via `showStandardMenu()`
- Menu imports commands dynamically (lazy loading)
- Recursive calls maintain session until user exits

**Calendar View:**
- Fixed layout: 8:00-19:00, 15-minute slots, 44 total slots
- 6 rotating colors for rooms (deterministic hash based on room name)
- Legend: `█` = your bookings, `=` = others, `·` = past, space = available
- Fetches all rooms, then gets busy times for each (sequential API calls)

### 4. Booking Flow

**Single Booking:**
1. Fetch all available rooms
2. If interactive, prompt for: room, date, start time, end time
3. POST to items endpoint to get room details
4. POST to busy times endpoint to check availability
5. Client-side validation: check if slot conflicts with existing bookings
6. POST to cart endpoint (add to cart)
7. Check cart for errors (8+ different error conditions)
8. POST to payment endpoint (confirm booking)
9. Display success message

**Batch Booking:**
- JSON file or inline JSON with array of bookings
- Processes sequentially (no parallelization)
- No rollback if one fails mid-batch
- Shows summary at end

**Issues:**
- bookRoom() function is 450+ lines
- Booking logic duplicated for batch operations
- API response doesn't return booking ID (can't immediately cancel)
- No optimistic concurrency (must fetch busy times each time)

### 5. Data Flow Oddities

**Room IDs:**
- Rooms have different ID formats for different API calls
- Some endpoints use ItemId, others use different identifiers
- Code maps between these formats

**Field Name Chaos:**
- Mix of camelCase and PascalCase (ItemId, FirstName, LastName)
- Price fields: EuroTTC, EuroHT, Credits (code checks all)
- Availability: IsLocked field (unclear what this means)
- Busy times filtering done client-side (API doesn't support date range)

**Response Structures:**
```javascript
// List rooms endpoint
{ VisitedItems: [...], UnvisitedItems: [...] }

// Cart endpoint
{ CartHasError: bool, Error: string, ErrorMessage: string, ... }

// Busy times endpoint
[{ StartDate, EndDate, ... }]  // simple array

// User auth endpoint
{ data: { User: { IsAuth: bool, ... } } }  // nested object
```

No consistency = lots of defensive code.

### 6. MCP Server Implementation

**Protocol:** JSON Lines over stdio

**Available Commands:**
- list_tools, list, book, cancel, my-bookings, calendar

**Implementation Hack:**
- Hijacks console.log to capture booking results
- Returns placeholder "generated-booking-id" instead of real ID
- Reuses existing command functions (not ideal architecture)

**Missing:**
- Authentication flow (assumes pre-authenticated)
- Proper result extraction from API responses
- Full calendar data (simplified for MCP)

## Major Issues

### 1. No Token Refresh Logic
**Problem:** Refresh tokens saved but never used. When access token expires, user must re-login.

**Fix:** Implement token refresh interceptor in API client:
```javascript
if (response.status === 401) {
  newAccessToken = await refreshAccessToken(refreshToken)
  retryRequest()
}
```

### 2. API Fragility
**Problem:** Different endpoints return wildly different structures. Code has defensive checks everywhere.

**Fix:** Create response normalizers for each endpoint type. Centralize error parsing.

### 3. Code Duplication
**Problem:** booking logic appears twice (single + batch), 450+ line functions

**Fix:** Extract common logic into shared functions. Break large functions into smaller units.

### 4. Calendar Performance
**Problem:** Fetches busy times sequentially for each room (10+ rooms = 10+ API calls)

**Fix:** Parallelize with `Promise.all()` (already used elsewhere, just not here)

### 5. Batch Operations Reliability
**Problem:**
- Sequential processing (slow)
- No rollback on partial failure
- No retry logic

**Fix:**
- Parallel processing with concurrency limit
- Option to cancel previous bookings if batch fails
- Add retry with exponential backoff

### 6. Hard-coded Values
**Problem:**
- Calendar: 8:00-19:00, 15-min slots (magic numbers)
- Only 6 colors for rooms
- Space ID and Category ID in env but still somewhat hard-coded

**Fix:** Extract constants to config file. Support multiple CoSoft instances.

### 7. MCP Server Implementation
**Problem:** console.log hijacking, placeholder IDs, incomplete data

**Fix:** Refactor commands to return structured data instead of printing. MCP should call internal functions, not command wrappers.

### 8. No Booking Confirmation ID
**Problem:** After booking, no way to reference it immediately. Must query my-bookings to find it.

**Fix:** Parse booking ID from API response, display to user, store in memory for quick cancel.

### 9. Input Validation
**Problem:** Some validation happens after user input, some before API call, inconsistent

**Fix:** Validate all inputs before making API calls. Centralize validation logic.

### 10. Error Messages
**Problem:** Mix of generic ("API error") and specific ("CALL BOX 3 is already booked")

**Fix:** Map API error codes to user-friendly messages. Provide actionable guidance.

## Recommendations for Golang Rewrite

### Architecture Improvements

1. **Clean Separation:**
   ```
   cmd/          # CLI commands (Cobra)
   internal/
     ├── api/    # API client
     ├── auth/   # Authentication & token management
     ├── ui/     # Terminal UI (bubbletea)
     └── models/ # Data structures
   pkg/          # Public packages (if library)
   ```

2. **API Client Design:**
   - Single HTTP client with interceptors
   - Automatic token refresh
   - Response normalizers per endpoint
   - Structured error types
   - Context-based cancellation
   - Connection pooling

3. **Data Models:**
   - Define structs for all API responses
   - Use json tags for marshaling
   - Validation with struct tags
   - Type safety throughout

4. **Concurrency:**
   - Worker pools for batch operations
   - Goroutines + channels for parallel API calls
   - Context for timeouts and cancellation
   - Rate limiting to avoid API throttling

5. **Testing:**
   - Mock HTTP responses
   - Table-driven tests
   - Integration tests with test server
   - Benchmark batch operations

6. **Configuration:**
   - Use viper for config management
   - Support multiple environments
   - Secure credential storage (keychain/keyring)
   - Environment-specific settings

### Technical Choices

**CLI Framework:** Cobra (standard in Go ecosystem)
**HTTP Client:** net/http with custom RoundTripper
**Terminal UI:** bubbletea (for interactive mode)
**Tables:** pterm or tablewriter
**Colors:** fatih/color or lipgloss
**Date/Time:** time package (standard library)
**JSON:** encoding/json (standard library)
**Testing:** testify for assertions
**Mocking:** gomock or testify/mock

### Optimization Opportunities

1. **Caching Layer:**
   - In-memory cache for rooms (TTL: 5 minutes)
   - Cache busy times (TTL: 30 seconds)
   - Invalidate on booking/cancellation

2. **Connection Pooling:**
   - Single HTTP client instance
   - Keepalive connections
   - Configurable timeouts

3. **Parallel Operations:**
   - Concurrent room fetching
   - Parallel busy times queries
   - Batch operations with worker pools

4. **Request Batching:**
   - If API supports it, batch requests
   - Otherwise, parallelize with rate limiting

5. **Smart Validation:**
   - Validate before network calls
   - Early exit on invalid input
   - Client-side business logic when possible

### Security Considerations

1. **Token Storage:**
   - Use OS keychain/credential manager
   - Encrypt at rest
   - Never log tokens
   - Clear from memory after use

2. **Input Sanitization:**
   - Validate all user inputs
   - Prevent injection attacks
   - Use parameterized queries if applicable

3. **Error Handling:**
   - Don't leak sensitive info in errors
   - Log security events
   - Rate limit login attempts

### Code Quality

1. **Linting:** golangci-lint with strict config
2. **Formatting:** gofmt, goimports
3. **Documentation:** godoc comments
4. **Error Wrapping:** fmt.Errorf with %w
5. **Logging:** structured logging (zap or zerolog)
6. **Metrics:** Track API latency, error rates

## Current Code Quality Assessment

**Strengths:**
- Modular command structure
- Comprehensive feature set
- Good use of existing libraries
- Batch operations support
- Both interactive and non-interactive modes
- MCP integration (innovative)

**Weaknesses:**
- Large functions (bookRoom.js is 450+ lines)
- Code duplication (booking logic)
- Defensive programming required everywhere (API fragility)
- No token refresh
- Magic numbers throughout
- No tests
- Inconsistent error handling
- console.log hijacking in MCP
- Sequential API calls in hot paths

**Technical Debt:**
- Refactor large functions
- Extract common logic
- Add comprehensive tests
- Implement token refresh
- Normalize API responses
- Add input validation layer
- Improve MCP implementation

## Bottom Line

This is a functional CLI that successfully wraps a terrible API. The "vibe coding" approach worked to get something running, but the lack of structure shows in large functions, duplicated code, and defensive error handling everywhere.

For a Golang rewrite, focus on:
1. **Type safety** - Structs for all API responses
2. **Concurrency** - Parallel operations where possible
3. **Clean architecture** - Separate concerns, small functions
4. **Robust error handling** - Don't rely on defensive checks everywhere
5. **Token management** - Proper refresh logic
6. **Testing** - Unit and integration tests from the start

The CoSoft API is clearly a mess (inconsistent responses, poor error handling, weird field names), so your client needs to normalize this chaos into a clean interface. Golang's strong typing and explicit error handling will help significantly.
