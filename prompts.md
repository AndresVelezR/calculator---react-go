# Engineering approach and AI-assisted workflow

## Technical direction

I used AI assistance within a detailed engineering plan. My instructions set
implementation boundaries, specified observable behavior, and defined the
checks required before changes could be committed. I delegated implementation,
test execution, debugging, and documentation work within that scope.

The priorities were correctness, clarity, and maintainability. I explicitly
limited the scope to the calculator requirements and asked for straightforward,
idiomatic code without unnecessary abstractions or additional features.

## Architecture and contracts I specified

The following choices were explicit in my implementation instructions:

| Area | Direction | Engineering rationale |
| --- | --- | --- |
| Backend dependencies | Use Go's standard library, including its HTTP router. | Keep the service small and avoid framework overhead for a single calculation endpoint. |
| Domain boundary | Put arithmetic operations and dispatch in a domain package independent of HTTP. | Test calculation behavior without a server or request context. |
| API contract | Use `POST /api/v1/calculate` with `{operation, a, b}` and consistent JSON success and error responses. | Give the frontend one predictable interface for all supported operations. |
| Required values | Represent request operands as `*float64`; require `b` except for square root. | Distinguish a missing operand from the valid number zero. |
| Error handling | Use sentinel errors and `errors.Is`; map invalid requests to 400 and domain failures to 422. | Keep error classification stable without coupling behavior to message text. |
| Frontend responsibilities | Separate the API client, the `useCalculator` state hook, and display/keypad components. | Keep transport, interaction state, and rendering independently testable. |
| Numeric presentation | Address floating-point display artifacts in the frontend after reproducing them. | Keep display formatting separate from backend arithmetic. |
| Runtime configuration | Configure the API URL and permitted browser origin through environment variables. | Make local and container deployments explicit. |

These rationales explain the design choices in the plan; they are not a claim
that every implementation detail was prescribed in advance.

## How I constrained the work

I divided the implementation into bounded blocks: domain behavior, HTTP API,
frontend client, calculator UI, frontend tests, containers, and documentation.
Each block had a branch and a sequence of logical commits. I required pull
requests to use merge commits and prohibited rewriting the existing history.

I also specified acceptance checks:

- Backend changes: clean `gofmt` output, `go vet ./...`, and `go test ./...`.
- Frontend changes: lint and build, plus tests once the test tooling existed.
- Tests: domain edge cases, handler success and error responses, the HTTP
  client, hook behavior, and user interactions.
- Integration: run the Docker stack and verify health and an actual calculation.
- Documentation: provide executable setup instructions, API examples, design
  assumptions, and measured coverage rather than estimated numbers.

For the anticipated infinity, CORS, and floating-point display issues, I
required reproduction before a fix commit. If a problem did not occur, the
instruction was to skip that fix and report it. This made the proposed fixes
conditional on evidence.

## Representative prompts

The excerpts below condense instructions I supplied during the project. They
preserve the technical requirements but are edited for readability, not verbatim
conversation records.

### Backend: implement a defined contract

> Implement `POST /api/v1/calculate` in `backend/internal/httpapi` using the Go
> standard library. Accept `{"operation":"add","a":1,"b":2}` and return
> `{"result":3}` on success or `{"error":"readable message"}` on failure.
> Use `*float64` for both operands so missing fields are distinguishable from
> zero. Require `a` for every operation and `b` except for `square_root`.
> Return 400 for malformed JSON, missing required fields, and unknown
> operations; return 422 for domain failures. Use `errors.Is`, never message
> comparisons. Keep DTOs and HTTP error mapping in the HTTP package, preserve
> `GET /health`, and test the handler with `net/http/httptest`.

### Frontend: separate state, transport, and presentation

> Implement a typed client that calls
> `POST {VITE_API_URL}/api/v1/calculate` with async/await. Default the base URL
> to `http://localhost:8080` and propagate the backend's error message for
> unsuccessful responses. Put display, previous operand, pending operation,
> error, and loading state in `useCalculator`. Expose actions for digits,
> decimal point, operation selection, equals, clear, and delete. Give Display
> and Keypad explicit props; keep state in the hook. Prevent duplicate decimal
> points and empty input, show backend errors, and support the optional
> operations as well as the four required ones.

### Debugging: demonstrate the failure before changing behavior

> Reproduce the non-finite result problem with curl before fixing it, using
> a power calculation such as `10^1000`. Add a domain error such as
> `ErrNonFiniteResult` and have `Calculate` reject Inf or NaN results, with a
> regression test. Do not record a fix for an issue that was not reproduced.
> Apply the same evidence requirement to browser CORS failures and the
> `0.1 + 0.2` display artifact.

### Delivery: verify the deployment and document the implementation

> Add a multi-stage backend Dockerfile and serve the built frontend through
> Nginx. Use Compose to run both services, with an explicit frontend API URL
> and permitted CORS origin. Verify `docker compose up --build`, `/health`,
> and a real calculation. Document local setup, Docker usage, API success and
> error examples, test commands, actual coverage, and design assumptions.
> Include component and request-sequence Mermaid diagrams.

## Scope of the assistance

AI contributed implementation code, tests, debugging, repository operations,
and substantial documentation. It also helped review and organize supplied
architecture diagrams and correct rendering issues. It exercised technical
judgment while carrying out the plan.

My role in the recorded instructions was to establish the technical direction,
define contracts and constraints, limit scope, and require verification. I
remain responsible for the submitted solution and for explaining its behavior
and tradeoffs. The value of the assistance was accelerating execution within
those requirements, with test and integration results as evidence of the outcome.
