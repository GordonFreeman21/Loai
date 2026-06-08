# Ollama Chat

A lightweight, modern Local AI Chat Client built with Wails v2, Go, and Svelte. This desktop application provides a clean GUI for interacting with local LLMs via Ollama.

## Features

- **Model Selection**: Fetch and switch between available Ollama models.
- **Pull Models**: Interface to pull new models directly from the app.
- **Streaming Responses**: Real-time message generation using Ollama's chat API.
- **Persistent History**: Chat sessions and messages are saved in a local SQLite database.
- **Rich Text Rendering**:
  - Full Markdown support.
  - Syntax highlighting for code blocks.
  - "Copy to Clipboard" button for code snippets.
- **Modern UI**: Dark-mode-first design with a collapsible sidebar and auto-expanding input.

## Tech Stack

- **Backend**: Go (Golang)
- **Desktop Framework**: [Wails v2](https://wails.io/)
- **Frontend**: Svelte + Tailwind CSS
- **Database**: SQLite (CGO-free via `modernc.org/sqlite`)
- **AI Backend**: [Ollama](https://ollama.com/) (REST API)
- **Utilities**: `marked` (Markdown parsing), `highlight.js` (Syntax highlighting)

## Prerequisites

1. **Ollama**: Must be installed and running on `http://localhost:11434`.
2. **Go**: Version 1.21 or higher.
3. **Node.js & npm**: For frontend development.
4. **Wails CLI**: Install via `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.

## Getting Started

### Development Mode

To run the application in development mode with hot-reloading:

```bash
wails dev
```

### Building for Production

To create a production-ready executable:

```bash
wails build
```

The binary will be located in the `build/bin/` directory.

## Project Structure

- `backend/`: Go source code for database, Ollama client, and app logic.
- `frontend/`: Svelte source code, styles, and assets.
- `main.go`: Application entry point and Wails initialization.
- `wails.json`: Wails project configuration.

## License

MIT
