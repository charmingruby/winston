# Winston

An Overwatch assistant for my personal squad of friends.

## What is this?

A coach for our squad, helping us pick heroes, build compositions, and plan strategies around how we like to play.

## Agents

- **Winston:** coordinates the conversation, consults the roster, checks recommendations, and synthesizes the final answer.
- **Analyst:** recommends heroes, compositions, synergies, and swaps.
- **Strategist:** turns the composition into an actionable gameplay plan.

Composition advice follows Winston → Analyst → Strategist → Winston. Simple questions can be answered directly.

## Discord bot: `@Winston`

Run a dedicated Gateway bot with `go run ./cmd/discord` (or `task run:discord`).
No public URL, HTTPS endpoint, public key, or command registration is needed.
The existing HTTP and MCP commands remain independent.

### Bot setup

1. Create an application in the [Discord Developer Portal](https://discord.com/developers/applications).
2. Under **Bot**, obtain the bot token. Set `DISCORD_BOT_TOKEN` in `.env`, along
   with `GEMINI_API_KEY` and `GEMINI_MODEL`. Never commit or share the token.
3. Install the bot in your server using the `bot` scope, with **View Channels**,
   **Send Messages**, and **Read Message History** permissions in the channels it
   should use. For threads, also grant **Send Messages in Threads** and ensure
   the bot has access to the thread.
4. Start `go run ./cmd/discord` and send:
   `@Winston G sup e S DPS em Colosseo contra Hazard e Ashe: o que jogar?`
   Use an actual Discord mention, not just the text `@Winston`.

The bot listens to guild messages and responds only when mentioned. Bot and
webhook messages are ignored to prevent loops. The privileged **Message Content**
intent is not needed because Discord exposes the content of messages mentioning
the bot. DMs are not handled.

Each mention starts an independent Winston conversation. Include context again
for follow-up questions. Answers reply to the original message without pinging
users or roles, and long answers are split into multiple messages.

The bot accepts up to eight simultaneous requests (extra mentions are ignored
while busy) and limits generation to three minutes. Run one instance;
duplicate events are not deduplicated. DiscordGo handles reconnects and REST rate
limits. Ctrl+C closes the connection and cancels in-flight work without waiting
for requests to finish.

## Tests

```sh
go test -race ./...
```

Tests cover context propagation, canceled and busy handlers, and Unicode-safe
message splitting without credentials, a Gateway connection, or model calls.
