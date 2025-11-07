# 🏒 Stinky Socks Hockey Schedule Downloader

Go CLI tool that fetches Stinky Socks Hockey games from the Stinky Socks website (https://www.stinkysocks.net/) and syncs them to your Google Calendar.

---

## Workflow

```
+------------+      +---------------------+      +-----------------+
| Fetch API  | ---> | Filter & Transform  | ---> | Google Calendar |
| Schedule   |      | Schedule Items      |      | Events          |
+------------+      +---------------------+      +-----------------+
```

1. Fetch schedule data from the Stinky Socks API.
2. Filter for your team and convert raw data into Event structs.
3. Sync events to your Google Calendar (create/update as needed).

---

## Features

- Fetch schedules from Stinky Socks Hockey website.
- Filter games by your favorite team.
- Automatically create or update events in Google Calendar.
- Unit-tested, modular, and interview-ready.

---

## Prerequisites

Install dependencies:

```bash
go get github.com/PuerkitoBio/goquery
go get -u google.golang.org/api/calendar/v3
go get -u golang.org/x/oauth2/google
```

---

## Google Calendar API Setup

1. Follow https://developers.google.com/workspace/calendar/api/quickstart/go
2. Enable the API.
3. Configure the OAuth consent screen.
4. Authorize credentials for a **desktop application**.

> ⚠️ Important: If you are not a test user, add your Gmail to the OAuth consent screen:
>
> - Go to Google Cloud Console → APIs & Services → OAuth consent screen
> - Scroll to Test users → + Add Users → add your Gmail
> - Save changes

Place the downloaded credentials.json in the project root.

---

## Setup Environment Variable

```bash
export KREEZEE_COOKIE="YOUR_SESSION_COOKIE_HERE"
```

---

## Running the Application

```bash
go run cmd/main.go
```

**Example output:**

```
Fetching schedule for team containing 'Puck Luck' from 2025-11-06 to 2026-11-06...
Filtered Events:
Puck Luck vs Other Team @ Rink 1 — Fri, 07 Nov 2025 14:00:00 EST to Fri, 07 Nov 2025 15:00:00 EST
Another Team vs Puck Luck @ Rink 2 — Sat, 08 Nov 2025 16:00:00 EST to Sat, 08 Nov 2025 17:00:00 EST
✅ Event created: Puck Luck vs Other Team (link-to-google-calendar)
🔄 Event updated: Another Team vs Puck Luck (link-to-google-calendar)
```

---

## Running Tests

```bash
go test ./...
```

Expected output:

```
ok      schedule-downloader/auth        0.176s
ok      schedule-downloader/google-calendar    0.238s
ok      schedule-downloader/schedule    0.034s
?       schedule-downloader/cmd [no test files]
?       schedule-downloader/models      [no test files]
```

---

## Project Structure

```
schedule-downloader/
├── auth/              # OAuth helper
├── google-calendar/   # Google Calendar sync logic
├── schedule/          # Fetching and transforming schedule data
├── cmd/main.go        # CLI entrypoint
├── go.mod             # Module dependencies
├── go.sum
└── credentials.json   # Google API credentials
```

---
