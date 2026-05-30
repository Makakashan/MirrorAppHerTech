#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_PID=""
FRONTEND_PID=""
BACKEND_PORT="${PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-4200}"

ensure_port_free() {
  local port="$1"
  local label="$2"

  if command -v lsof >/dev/null 2>&1 && lsof -ti ":${port}" >/dev/null 2>&1; then
    echo "${label} port ${port} is already in use. Stop that process or set a different port." >&2
    exit 1
  fi
}

cleanup() {
  if [[ -n "${FRONTEND_PID}" ]] && kill -0 "${FRONTEND_PID}" 2>/dev/null; then
    kill "${FRONTEND_PID}" 2>/dev/null || true
  fi
  if [[ -n "${BACKEND_PID}" ]] && kill -0 "${BACKEND_PID}" 2>/dev/null; then
    kill "${BACKEND_PID}" 2>/dev/null || true
  fi
}

trap cleanup EXIT INT TERM

ensure_port_free "${BACKEND_PORT}" "Backend"
ensure_port_free "${FRONTEND_PORT}" "Frontend"

cd "${ROOT_DIR}/backend"
GOCACHE="${GOCACHE:-/tmp/go-cache}" ADDR=":${BACKEND_PORT}" go run ./cmd/server &
BACKEND_PID=$!

cd "${ROOT_DIR}/frontend"
./node_modules/.bin/ng serve --host 127.0.0.1 --port "${FRONTEND_PORT}" --proxy-config proxy.conf.json &
FRONTEND_PID=$!

wait -n "${BACKEND_PID}" "${FRONTEND_PID}"
