# React + Go Calculator

A full-stack calculator with a React/TypeScript frontend and a Go REST API. It supports addition, subtraction, multiplication, division, powers, square roots, and percentages. Every arithmetic operation is performed by the backend.

## Requirements

- Go 1.27.1 or newer, matching `backend/go.mod`.
- Node.js 22.12 or newer and npm (tested locally with Node 22.20.0).
- Alternatively, Docker with Docker Compose v2 to run both services without installing Go or Node.

Clone the repository:

```sh
git clone https://github.com/AndresVelezR/calculator---react-go.git
cd calculator---react-go
```

## Run locally

From the repository root, start the backend in one terminal:

```sh
cd backend
go run ./cmd/server
```

The API listens on `http://localhost:8080`. In another terminal, start the frontend:

```sh
cd frontend
npm ci
cp .env.example .env
npm run dev
```

Open **http://localhost:5173**. The backend allows this exact browser origin by default. If Vite chooses another port because 5173 is occupied, free the port or update `CORS_ORIGIN` when starting the backend.

| Variable | Where it is read | Local default | Docker default |
| --- | --- | --- | --- |
| `VITE_API_URL` | Frontend, at build/dev-server startup | `http://localhost:8080` | `http://localhost:8080` |
| `CORS_ORIGIN` | Backend, at startup | `http://localhost:5173` | `http://localhost:3000` |

The frontend `.env` is ignored by Git. Restart Vite after changing it. `CORS_ORIGIN` is a single exact origin (scheme, host, and port, without a trailing slash), for example:

```sh
cd backend
CORS_ORIGIN=http://localhost:5173 go run ./cmd/server
```

For a production frontend build, run `npm run build` in `frontend/`. Static files are generated in `frontend/dist/`.

## Run with Docker

From the repository root:

```sh
docker compose up --build
```

Open **http://localhost:3000**. The API remains at **http://localhost:8080**. Stop a locally running backend first to free port 8080. Compose waits for the backend health check before starting Nginx.

To run in the background and inspect or stop the services:

```sh
docker compose up --build -d
docker compose ps
docker compose logs
docker compose down
```

Compose accepts `VITE_API_URL` and `CORS_ORIGIN` from the shell or a root `.env` file. The API URL must be reachable from the **browser**; the Compose service name `backend` is only reachable inside Docker. Rebuild the frontend after changing `VITE_API_URL`, since Vite embeds it into the static JavaScript. Opening the UI through a different host requires a matching `CORS_ORIGIN` and a browser-reachable API URL.

Verified with `docker compose up --build -d`: both images built, backend healthy, Nginx returned HTTP 200, `/health` returned `{"status":"ok"}`, and a real API request returned `{"result":42}` for `6 × 7`. A browser calculation from `localhost:3000` also passed.

## Using the calculator

Enter a number, select an operation, enter a second number, and press `=`. `C` resets the calculator; `⌫` deletes the last input digit. Decimal input accepts one decimal point. Errors appear below the display, and keys are disabled while a request is pending.

- **Power:** `2`, `xʸ`, `3`, `=` gives `8`.
- **Square root:** `9`, `√` gives `3`. Finish a pending binary calculation before using square root.
- **Percentage:** `50`, `%`, `20`, `=` gives `10` (20% of 50).
- Operations chain left to right: `2 + 3 × 4 =` gives `20`. This is a keypad calculator, not an expression parser with operator precedence.
- Selecting another operator before entering the second number replaces the pending operator. After a result, entering a digit starts a new number; selecting an operator reuses the result. Pressing `=` without a pending operation reports a validation error.

## REST API

### Health

```sh
curl http://localhost:8080/health
# 200: {"status":"ok"}
```

### Calculate

`POST /api/v1/calculate` accepts one JSON object. `operation` and `a` are required. `b` is required except for `square_root`, where it is ignored if supplied. Numbers must be JSON numbers, not strings; `null` is treated as an absent number. Operation names are case-sensitive.

| Operation | Meaning | Example inputs | Result |
| --- | --- | --- | --- |
| `add` | a + b | a=1, b=2 | 3 |
| `subtract` | a − b | a=1, b=2 | -1 |
| `multiply` | a × b | a=6, b=7 | 42 |
| `divide` | a ÷ b | a=1, b=4 | 0.25 |
| `power` | a raised to b | a=2, b=3 | 8 |
| `square_root` | square root of a | a=9 | 3 |
| `percentage` | a × b / 100 | a=50, b=20 | 10 |

```sh
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":1,"b":2}'
# 200: {"result":3}

curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"square_root","a":9}'
# 200: {"result":3}

curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"percentage","a":50,"b":20}'
# 200: {"result":10}
```

### Validation and domain errors

Calculation errors use `{"error":"readable message"}`. Invalid requests return **400**; valid requests whose arithmetic cannot produce a finite real result return **422**. Unexpected internal response-encoding errors return **500**.

```sh
# Invalid JSON: 400
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{'
# {"error":"request must be a valid JSON object with operation, a, and b"}

# Missing operation: 400
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"a":1,"b":2}'
# {"error":"operation is required"}

# Unknown operation: 400
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"modulo","a":1,"b":2}'
# {"error":"unknown operation"}

# Missing first number: 400
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"add","b":2}'
# {"error":"a is required"}

# Missing second number: 400
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"add","a":1}'
# {"error":"b is required for this operation"}

# Invalid numeric type: 400
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"add","a":"one","b":2}'
# {"error":"request must be a valid JSON object with operation, a, and b"}

# Multiple JSON values: 400
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{} {}'
# {"error":"request must contain a single JSON object"}

# Division by zero: 422
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"divide","a":1,"b":0}'
# {"error":"division by zero"}

# Negative square root: 422
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"square_root","a":-1}'
# {"error":"square root of a negative number"}

# Overflow: 422
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"power","a":10,"b":1000}'
# {"error":"result is not a finite number"}

# No real-valued power result: 422
curl -i http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' -d '{"operation":"power","a":-1,"b":0.5}'
# {"error":"result is not a finite number"}
```

Unknown JSON fields, inputs outside the `float64` range, and bodies larger than 4096 bytes also return 400. Routes and HTTP methods are handled by Go's standard router; unknown routes return 404 and unsupported methods return 405.

## Tests and coverage

Backend, from the repository root:

```sh
cd backend
gofmt -l .             # expected: no output
go vet ./...
go test ./...
go test -race ./...
mkdir -p coverage
go test -coverprofile=coverage/coverage.out ./...
go tool cover -func=coverage/coverage.out
go tool cover -html=coverage/coverage.out -o coverage/index.html
```

Frontend, in another terminal from the repository root:

```sh
cd frontend
npm ci
npm run lint
npm run build
npm test -- --run
npm run test:coverage
```

Use `npm test` for watch mode. Open `backend/coverage/index.html` and `frontend/coverage/index.html` for detailed HTML reports. Generated reports are ignored by Git and can be recreated with the commands above.

Measured coverage on the completed implementation:

| Scope | Statements | Branches | Functions | Lines |
| --- | --- | --- | --- | --- |
| Backend, all packages | 86.7% | — | — | — |
| Backend domain | 100.0% | — | — | — |
| Backend HTTP API | 94.1% | — | — | — |
| Frontend application | 100% | 100% | 100% | 100% |

Go coverage measures statements. The backend total includes `cmd/server/main.go` (0% unit coverage); process startup is exercised through the local server and Docker checks. HTTP coverage includes all handler validation paths, routing, CORS, and domain-error mapping; response-writing/encoding fallback paths are not fully unit-covered.

The frontend has **61 passing tests in four files**. Its coverage includes App, Calculator, Display, Keypad, the HTTP client, the hook, and formatting. Only the React bootstrap (`main.tsx`), type-only files, test setup, and tests themselves are excluded. Tests mock `fetch` for client/UI checks and mock the API function for isolated hook checks. Chromium was also used to verify real API calls, errors, and layouts at 320px and 1280px.

## Design decisions and assumptions

- **Go standard library:** `net/http` handles method-aware routes and `encoding/json` handles the wire format. The backend has no third-party dependencies.
- **Separate domain and transport:** `internal/domain` contains arithmetic and sentinel errors, with no HTTP dependencies. `internal/httpapi` owns decoding, validation, JSON responses, CORS, and HTTP status mapping.
- **One calculation endpoint:** an explicit operation dispatcher keeps the contract small and makes all seven operations use the same validation and response path.
- **Pointer input fields:** `*float64` distinguishes a missing or null operand from the valid number zero. HTTP error mapping uses `errors.Is`, including wrapped errors.
- **Finite floating-point arithmetic:** Go `float64` and JavaScript numbers are used. `Calculate` rejects Inf/NaN results before JSON encoding. This is not arbitrary-precision or decimal financial arithmetic. Power follows Go `math.Pow` semantics, including `0^0 = 1`; underflow can produce zero.
- **Presentation rounding:** results display up to ten decimal places without trailing zeroes. Nonzero magnitudes below `1e-10` use up to ten significant digits to avoid displaying them as zero. Original result values remain available for subsequent calculations. Typed decimal input is preserved.
- **Simple React state:** `useCalculator` owns input, pending operations, loading, and errors; display and keypad components receive props. One request runs at a time. No state-management framework or client-side arithmetic engine is needed.
- **Explicit CORS origin:** the backend permits the configured frontend origin and answers preflight requests. This keeps development and Docker origins predictable.
- **Stateless service:** no authentication, database, persisted history, or extra operations are included. Those are outside this calculator's scope.

See [Architecture](docs/ARCHITECTURE.md) for the editable Mermaid diagrams and the [diagram gallery](docs/diagrams/README.md) for the complete set of SVG architecture views.
