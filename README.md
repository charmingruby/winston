# winston

Overwatch assistant with analyst, strategist, and reviewer agents.

## Configuration

Set `GEMINI_API_KEY` in the environment or in a `.env` file in the working directory.
The HTTP entrypoint also uses `PORT`.

## MCP server

Build the stdio server:

```sh
go build -o bin/winston-mcp ./cmd/mcp
```

Configure your MCP client to launch the binary. For clients using the
`mcpServers` configuration format:

```json
{
  "mcpServers": {
    "winston": {
      "command": "/absolute/path/to/winston/bin/winston-mcp",
      "env": {
        "GEMINI_API_KEY": "your-api-key"
      }
    }
  }
}
```

Available tools:

- `get_team_roster`: returns configured team players without calling a model.
- `ask_winston`: accepts `{"message": "Suggest a composition for our team"}` and
  returns `{"answer": "..."}`. winston can consult the roster and delegate to
  specialist agents. Each call starts an independent conversation; include any
  relevant previous context in the message.

Agent calls have a three-minute timeout and incur model-provider usage.
Logs go to stderr; stdout is reserved for MCP messages.

For local development, run `go run ./cmd/mcp` from the repository directory.
The server waits for MCP JSON-RPC messages on stdin, not interactive chat input.

## HTTP server

```sh
go run ./cmd/agent
```

The agent is served at `POST /agents/winston`.

## Tests

```sh
go test ./...
```
