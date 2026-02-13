# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Tomatobot is a Slack bot that provides Pomodoro timer notifications. Users trigger a timer with the `/pomodoro` slash command, and the bot schedules a reminder using Slack's `chat.scheduleMessage` API.

## Build & Run

```bash
# Build
go build -o tomatobot .

# Run (requires environment variables)
SLACK_BOT_TOKEN=xoxb-... SLACK_APP_TOKEN=xapp-... ./tomatobot
```

## Architecture

Single-file Go application (`main.go`) using Socket Mode via `slack-go/slack`. The bot connects over WebSocket (no public URL needed), listens for `/pomodoro` slash commands, and uses `chat.scheduleMessage` to schedule reminders server-side. This means reminders are delivered even if the bot restarts.

## Dependencies

- `github.com/slack-go/slack` — Slack API client with Socket Mode support
- Go modules (`go.mod`) for dependency management

## Configuration

Environment variables:
- `SLACK_BOT_TOKEN` — Bot User OAuth Token (`xoxb-...`)
- `SLACK_APP_TOKEN` — App-Level Token with `connections:write` scope (`xapp-...`)
