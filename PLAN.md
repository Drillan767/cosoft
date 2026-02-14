The loop, concretely

Every single interaction (button click, form submission) follows the same 5 steps:

1. Deserialize: Read the user's current View from the DB (keyed by slack_user_id or slack_message_id — you already have the slack_messages table with a payload blob and message_type)
2. Update: Call view.Update(action) — this returns a new View, possibly of a different type (e.g. LoginView → MainMenuView, or MainMenuView → QuickBookView)
3. Serialize: Write the new View back to the DB
4. Render: Call RenderView(newView) to produce a slack.Block
5. Send: POST the blocks to Slack's response_url

What you already have (and it's good)

Your view.go already sketches this out correctly:

- A View interface with Update(state json.RawMessage) View
- LoginView and MainView structs that implement it
- A RenderView(v View) function that dispatches rendering by type
- RenderLoginView that produces slack.Block

The slack_messages table already has the right shape: slack_user_id, slack_message_id, payload (BLOB), and message_type.

What's missing / what needs to change

1. The View needs to be serializable (to/from the DB)

Each View needs to be JSON-marshallable so you can store it in slack_messages.payload. You also need a way to reconstruct the right Go type from the message_type column. Something like:

message_type = "login"    → deserialize payload into LoginView
message_type = "main"     → deserialize payload into MainMenuView
message_type = "quickbook" → deserialize payload into QuickBookView

This is a simple registry/factory — a switch on message_type that returns the right concrete type, then you unmarshal the payload into it.

2. Update needs the action, not just form state

Right now Update(state json.RawMessage) receives raw form values. But it also needs to know which button was pressed (the action_id). For block_actions, the meaningful data is the action ID. For view_submission,
it's the form values. You might want Update to receive something slightly richer — or just the raw Slack payload and let each View parse what it needs.

3. The server handler becomes generic

This is the big win. Instead of handleMenuAction with a manual switch on actionName, and handleQuickbookModal with manual field extraction, your handleInteractions becomes roughly:

func (b *Bot) handleInteractions(w http.ResponseWriter, r *http.Request) {
// 1. Parse the Slack payload, extract user ID and action
// 2. Load current View from DB for this user
// 3. view = view.Update(payload)
// 4. Save new View to DB
// 5. blocks = RenderView(view)
// 6. Send blocks to Slack
}

That's it. One handler. The switch logic moves into each View's Update method, which is where it belongs — each View knows what actions it supports and what View to transition to.

4. Define your Views (matching your TUI pages)

Looking at app.go, your pages are: Landing (main menu), QuickBook, Browse, Reservations, Settings. In Slack, the equivalents would be:
┌──────────────┬──────────────────┬───────────────────────────────────────────────┐
│   TUI Page   │    Slack View    │                State it holds                 │
├──────────────┼──────────────────┼───────────────────────────────────────────────┤
│ (none)       │ LoginView        │ email, password, error                        │
├──────────────┼──────────────────┼───────────────────────────────────────────────┤
│ Landing      │ MainMenuView     │ user info (or just slack_user_id to fetch it) │
├──────────────┼──────────────────┼───────────────────────────────────────────────┤
│ QuickBook    │ QuickBookView    │ duration, nb_people, selected room, step      │
├──────────────┼──────────────────┼───────────────────────────────────────────────┤
│ Browse       │ BrowseView       │ filters, current page                         │
├──────────────┼──────────────────┼───────────────────────────────────────────────┤
│ Reservations │ ReservationsView │ list of reservations                          │
├──────────────┼──────────────────┼───────────────────────────────────────────────┤
│ Settings     │ SettingsView     │ TBD                                           │
└──────────────┴──────────────────┴───────────────────────────────────────────────┘
Each View's Update either mutates its own state and returns self, or returns a different View to navigate (just like LoginView.Update already returns &MainView{} on success).

What to move vs. remove

- server.go lines 140-153 (handleInteractions switch): This entire dispatch gets replaced by the generic loop. The block_actions vs view_submission distinction still matters for parsing the Slack payload, but
after that it all funnels into view.Update(...).
- server.go handleMenuAction (lines 253-281): The switch on actionName ("main-menu", "quick-book") goes away. That routing becomes MainMenuView.Update returning the right next View based on the action.
- server.go handleQuickbookModal / handleLoginModal: The field extraction logic (viewResponse.View.State.Values.Duration...) moves into each View's Update method. handleLoginModal essentially becomes
LoginView.Update + the API call.
- services/auth.go LogInUser / postLogin: These stay, but get called from LoginView.Update (or from a step after Update, depending on whether you want side effects inside Update or outside).
- services/requests.go ShowQuickBook / ShowMainMenu: These get replaced by RenderView. UpdateMessage and DispatchModal stay — they're your "send to Slack" transport layer.

The one question you need to answer

Should Update perform side effects (API calls, DB writes) or should it be pure?

In Bubble Tea, Update returns a Cmd for side effects. You could do the same — Update returns (View, Action) where Action might be "call the login API with these credentials" — and the generic handler executes the
action. This keeps Views testable and pure. But for your project's size, having Update directly call the service might be simpler. Your call.

What the DB actually stores

There's no separate "default states" table. The DB stores one row per user (in slack_messages) representing what that user is currently seeing and the state of that view. That's it.

Think of it like this: if you could freeze your TUI mid-session and save it to disk, then restore it on the next keypress — that's what the DB does for Slack.

Trace through a flow

Let's walk through Browse as if it were in Slack:

Step 1 — User clicks "Browse" on the main menu

MainMenuView.Update sees action_id = "browse" and returns:

return &BrowseView{
Phase:    0,
Date:     time.Now().Format(time.DateOnly),
Hour:     roundToQuarter(time.Now()),
Duration: 0,
NbPeople: 0,
}

This is equivalent to your NewBrowseModel() on line 31 of browse.go. The handler serializes this to JSON and writes it to slack_messages:

slack_user_id = "U12345"
message_type  = "browse"
payload       = {"phase":0,"date":"2026-02-13","hour":"14:30",...}

Then RenderView produces the Slack blocks for phase 0 (the search form with date/hour/duration/nbPeople inputs), and sends them to Slack.

Step 2 — User fills in the form and clicks "Search"

Slack sends a block_actions to your server with user_id = "U12345" and the form values.

The handler:
1. Queries SELECT * FROM slack_messages WHERE slack_user_id = 'U12345'
2. Sees message_type = "browse" → deserializes payload into a BrowseView
3. Calls view.Update(action) — Update reads the submitted values, stores them in the struct, sets Phase = 1, calls the API to get rooms, stores the rooms in the struct, sets Phase = 2
4. Writes the updated BrowseView back to the DB (payload now has rooms, phase=2, etc.)
5. RenderView produces the room selection blocks, sends to Slack

Step 3 — User picks a room

Same cycle. Load from DB → Update (sets Phase = 3, calls booking API) → Save to DB → Render → Send.

The key point you're asking about

I need to keep track of which user is seeing what?

Yes, and it's just one column: message_type. When the handler loads the row, it does:

switch row.MessageType {
case "login":
view = &LoginView{}
case "main":
view = &MainMenuView{}
case "browse":
view = &BrowseView{}
case "quick_book":
view = &QuickBookView{}
}
json.Unmarshal(row.Payload, view)

Now you have the right Go type with the right state. You call view.Update(action), and the returned view might be the same type (user stayed on the same page, just changed a field) or a different type (user
navigated — like Browse phase 0 → phase 2, or MainMenu → Browse).

After Update, you check the type of the returned view to know the new message_type, serialize it, and write it back.

What you don't store

You don't store "default states" or templates. The defaults live in Go code, as constructors — exactly like NewBrowseModel() already does. The DB only holds the current, live state for each active user.

Your slack_messages table

One thing — your current schema has slack_message_id as a separate column, which is good because that's the Slack message timestamp you need to update the right message. But conceptually, there's one "active
session" per slack_user_id. You might want a UNIQUE constraint on slack_user_id so each user has exactly one active view, and you INSERT OR REPLA