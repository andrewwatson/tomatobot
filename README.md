# tomatobot

A Slack bot for Pomodoro timer notifications using slash commands. Users trigger timers with `/pomodoro` and receive scheduled reminders via `chat.scheduleMessage`, so reminders persist even if the bot restarts.

## Slack App Setup

### 1. Create a Slack App

Visit [api.slack.com/apps](https://api.slack.com/apps) and select **Create New App** → **From scratch**. Name it "tomatobot" and select your workspace.

### 2. Enable Socket Mode

In your app settings:
- Go to **Settings** → **Socket Mode** → toggle **ON**
- Click **Generate an App-Level Token**
- Name it (e.g., "tomatobot-socket")
- Grant the `connections:write` scope
- Save the token (format: `xapp-...`) — you'll need this to run the bot

### 3. Add the `/pomodoro` Slash Command

In **Features** → **Slash Commands**:
- Click **Create New Command**
- Command: `/pomodoro`
- Short Description: `Start a pomodoro timer`
- Usage Hint: `[minutes] [description]`
- Save

### 4. Add Bot Token Scopes

In **Features** → **OAuth & Permissions**, under **Scopes** → **Bot Token Scopes**, add:
- `chat:write`
- `chat:write.public`
- `commands`

### 5. Install App to Your Workspace

Click **Install to Workspace** at the top of **OAuth & Permissions**, then authorize. Copy the **Bot User OAuth Token** (format: `xoxb-...`) — you'll need this to run the bot.

### 6. Set Environment Variables

```bash
export SLACK_BOT_TOKEN=xoxb-...
export SLACK_APP_TOKEN=xapp-...
```

## Build

```bash
go build -o tomatobot .
```

## Run

```bash
SLACK_BOT_TOKEN=xoxb-... SLACK_APP_TOKEN=xapp-... ./tomatobot
```

The bot connects via WebSocket (no public URL required) and listens for `/pomodoro` slash commands.

## Usage

- **`/pomodoro`** — Starts a 25-minute timer
- **`/pomodoro 15`** — Starts a 15-minute timer
- **`/pomodoro 25 write report`** — Starts a 25-minute timer labeled "write report"
- **`/pomodoro write report`** — Starts a 25-minute timer labeled "write report"

The bot responds immediately with a confirmation and schedules a reminder message to be delivered after the timer expires.

## How It Works

The bot uses Slack's `chat.scheduleMessage` API to schedule reminders server-side. This means:
- Reminders are delivered even if the bot process restarts
- No background goroutines needed
- All state is managed by Slack's servers
