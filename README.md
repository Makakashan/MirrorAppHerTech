# Mirror App Her Tech

Mirror is a web application in the form of a user's “digital mirror”: it helps track emotions, behavior, and recurring patterns, then reflects them back through dialogue, analytics, and history.

## About the project

The repository contains two main parts:

- `backend` — a Go service with HTTP API, user profile logic, and AI integration.
- `frontend` — an Angular application that shows the chat, analytics, mood calendar, and additional screens.

The app is designed for long-term observation of a user's emotional state and for visualizing recurring patterns in a clear, intuitive way.

## Features

- reflective AI chat;
- message history storage;
- analytics for the user's state;
- metrics and calendar visualization;
- suggestion cards and demo data for presentation scenarios;
- Docker Compose startup.

## Architecture

### Backend

The backend is written in Go and exposes the following API endpoints:

- `GET /api/health` — service health check;
- `GET /api/messages` — chat history;
- `POST /api/chat` — send a message and get a response;
- `GET /api/suggestions` — suggestion cards;
- `GET /api/metrics` — analytics metrics.

The service also:

- loads the user profile from `backend/data/seeds.json`;
- builds the system prompt from `backend/data/system_prompt.md`;
- keeps session context, excerpts, and triggers;
- can run as an API server or serve the frontend as static files.

### Frontend

The frontend is built with Angular 21 and uses:

- routing for application screens;
- HTTP client for API communication;
- animations and view transitions;
- `angular-three` for 3D rendering;
- `lucide-angular` for icons.

Main screens:

- home dashboard;
- chat;
- analytics;
- calendar;
- settings.

## Repository structure

- `backend/` — Go server side
- `frontend/` — Angular client side
- `backend/data/` — seed data, prompts, and system instructions
- `docker-compose.yml` — containerized startup
- `.env.example` — environment variables example

## Requirements

To run the project locally, you need:

- Go 1.26+
- Node.js and npm
- Docker and Docker Compose if you want containerized deployment

## Environment variables

Main backend variables:

- `ADDR` — HTTP server address, default `:8080`
- `DATA_DIR` — data directory, default `data`
- `AI_API_KEY` — API key for the AI provider
- `AI_INVOKE_URL` — API URL for model calls
- `AI_MODEL` — model name
- `MIRROR_DEMO_DATA` — enables demo data unless set to `false`, `0`, or `off`
- `FRONTEND_DIST` — path to the built frontend, if the backend should serve static files

See also `.env.example`.

## Local development

### Option 1: Docker Compose

1. Create `.env` based on `.env.example`.
2. Set `AI_API_KEY` and other values as needed.
3. Start the project:

```bash
docker compose up --build
```

After startup, the app will be available through the frontend container at `http://localhost`.

### Option 2: Run backend and frontend separately

#### Backend

```bash
cd backend
go run ./cmd/server
```

By default, the server listens on `:8080`.

#### Frontend

```bash
cd frontend
npm install
npm start
```

After startup, the frontend will be available at `http://localhost:4200/`.

## Build

### Frontend

```bash
cd frontend
npm run build
```

### Backend

```bash
cd backend
go build ./cmd/server
```

## Testing

### Frontend

```bash
cd frontend
npm test
```

### Backend

Backend tests are located in `backend/internal/server/server_test.go`.

Run them with:

```bash
cd backend
go test ./...
```

## Data and demo scenario

The project uses files from `backend/data/`:

- `seeds.json` — initial user profile and starting data;
- `system_prompt.md` — system prompt template;
- `prompts.md` — additional prompts and hints.

When `MIRROR_DEMO_DATA` is enabled, the backend starts with demo messages, metrics, and suggestion cards for presentation purposes.

## API overview

### `GET /api/health`

Checks whether the service is healthy.

### `GET /api/messages`

Returns the current chat history.

### `POST /api/chat`

Sends a user message and returns:

- the assistant response;
- updated metrics;
- analysis, if one was produced by AI.

Request body:

```json
{
  "message": "I feel anxious before the meeting"
}
```

### `GET /api/suggestions`

Returns an array of suggestion cards.

### `GET /api/metrics`

Returns the current profile metrics and calendar data.

## Technology stack

- Go
- Angular
- TypeScript
- RxJS
- Three.js
- Docker

## License

Not specified yet.
