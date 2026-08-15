# AGENTS.md

> Comprehensive agent and developer guide for **`goensemblrest`**.

---

## 1. Project Overview

`goensemblrest` is an idiomatic, high-performance Go client library for the [Ensembl REST API](https://rest.ensembl.org/). It serves as the official Go counterpart to [`pyEnsemblRest`](https://github.com/gawbul/pyEnsemblRest), providing 100% endpoint parity, thread safety, client-side sliding-window rate limiting, exponential backoff retries for transient errors, deep `context.Context` awareness, and idiomatic error handling with standard Go `errors.Is` and `errors.As`.

### Key Characteristics

- **Module Path**: `github.com/gawbul/goensemblrest`
- **Go Version**: Go `1.26.x` (target), compatible with Go 1.24+
- **Zero External Dependencies**: Uses Go standard library exclusively (`net/http`, `context`, `encoding/json`, `sync`, `time`, `net/url`, `bytes`, `errors`, `fmt`, `io`, `regexp`, `strconv`, `strings`).
- **Concurrency**: Safe for concurrent use across multiple goroutines.
- **License**: MIT License (Copyright (c) 2020-2026 Steve Moss).

---

## 2. Repository Layout & File Map

```
.
├── .github/
│   └── workflows/
│       ├── nightly.yaml        # Nightly live API drift check (03:00 UTC)
│       ├── pull_request.yaml   # PR validation: lint, vet, gofmt, unit tests, coverage, live smoke tests
│       └── push_tag.yaml       # Release automation on git tag push (v*)
├── examples/
│   └── main.go                 # Executable demonstration script showcasing library capabilities
├── Makefile                    # Standard developer automation targets
├── go.mod                      # Module definition (zero external dependencies)
├── LICENSE                     # MIT License
├── README.md                   # Public documentation, quickstart, and endpoint catalog
│
├── client.go                   # Client struct, configuration options (WithBaseURL, WithTimeout, etc.)
├── request.go                  # HTTP execution engine, path resolution, header handling, retry loop
├── ratelimit.go                # Thread-safe sliding-window rate limiter and header telemetry parser
├── errors.go                   # Sentinel errors, APIError struct, HTTP status mappings, Unwrap() logic
├── endpoints.go                # EndpointSpec definitions, EndpointsTable (100+ endpoints), Call() dynamic dispatcher
├── types.go                    # RequestOption definitions, request/response structs, domain data models
│
├── archive.go                  # Archive domain methods (GetArchiveByID, GetArchiveByMultipleIDs)
├── comparative.go              # Comparative genomics methods (Gene trees, Homology, CAFE, Alignments)
├── ga4gh.go                    # Global Alliance for Genomics and Health (GA4GH) & Beacon endpoints
├── info.go                     # Information & metadata methods (Ping, Species, Assembly, Biotypes, etc.)
├── ld.go                       # Linkage Disequilibrium methods (Pairwise, Region, Variant)
├── lookup.go                   # Genomic feature lookup methods (Symbol, ID, Multiple IDs/Symbols)
├── mapping.go                  # Coordinate and feature mapping methods (Assembly, cDNA, CDS, Translation)
├── ontology.go                 # Ontology and taxonomy methods (Ancestors, Descendants, Taxonomy)
├── overlap.go                  # Genomic region overlap methods (ID, Region, Translation)
├── phenotype.go                # Phenotype annotation methods (Accession, Gene, Region, Term)
├── regulation.go               # Regulation methods (Binding matrices)
├── sequence.go                 # Biological sequence retrieval methods (ID, Region, Multiple IDs/Regions)
├── transcript.go               # Transcript haplotype methods
├── variation.go                # Genetic variation & recoder methods (ID, PMCID, PMID, Multiple IDs)
├── vep.go                      # Variant Effect Predictor methods (HGVS, ID, Region)
├── xrefs.go                    # External cross-reference methods (ID, Symbol, Name)
│
├── client_test.go              # Unit tests for Client initialization, options, rate limiter, and error retry logic
├── endpoints_test.go           # Offline mock unit tests covering all 100+ domain methods
├── example_test.go             # Runnable GoDoc example tests (ExampleNewClient, ExampleClient_Call, etc.)
└── live_test.go                # Live integration smoke tests against rest.ensembl.org (guarded by //go:build live)
```

---

## 3. Core Architecture & Design Patterns

### 3.1 Client Configuration (Functional Options)

The `Client` struct in `client.go` is initialized using functional options (`ClientOption`):

```go
client, err := ensembl.NewClient(
    ensembl.WithBaseURL("https://rest.ensembl.org"),  // Trims trailing slashes
    ensembl.WithTimeout(60 * time.Second),            // Request timeout
    ensembl.WithRateLimit(15, time.Second),           // Sliding window limit (reqs, window)
    ensembl.WithMaxAttempts(5),                       // Max retry attempts
    ensembl.WithUserAgent("custom-agent/1.0"),        // User-Agent header
    ensembl.WithHeader("X-Custom", "value"),          // Default persistent header
    ensembl.WithHTTPClient(customHTTPClient),         // Custom *http.Client
)
```

### 3.2 Sliding-Window Rate Limiting (`ratelimit.go`)

- **Algorithm**: In-memory timestamp sliding-window queue (`timestamps []time.Time`).
- **Default Limit**: 15 requests per 1 second (matching Ensembl's server-side limits).
- **Concurrency**: Protected by `sync.Mutex`.
- **Context Cancellation**: `Wait(ctx context.Context)` aborts immediately if `ctx.Done()` fires while sleeping.
- **Telemetry Parsing**: Automatically updates rate limit metrics from response headers (`X-RateLimit-Reset`, `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Period`, `Retry-After`).
- **Inspection**: Retrieve current telemetry at any time via `client.RateLimit()`.

### 3.3 Request Execution & Exponential Backoff Retries (`request.go`)

- **Path Template Interpolation**: `resolvePath` resolves `{{param}}` templates using `paramRegex`.
  - **Colon Preservation**: Escapes path characters via `url.PathEscape` but specifically preserves colons (`:`) without percent-encoding (`%3A` -> `:`). This is mandatory for genomic coordinate formats (e.g., `X:1000..2000:1` or `homo_sapiens:BRCA2`).
- **Retry Strategy**: Automatically retries up to `maxAttempts` (default: 5) for transient conditions:
  - HTTP Statuses: `408 Request Timeout`, `500 Internal Server Error`, `502 Bad Gateway`, `503 Service Unavailable`, `504 Gateway Timeout`.
  - HTTP Status `429 Too Many Requests` (uses `Retry-After` header duration if present).
  - Ensembl Transient Body Strings: `"something bad has happened"`, `"Something went wrong while fetching from LDFeatureContainerAdaptor"`, `"timeout"`.
- **Backoff Formula**: Sleep duration is `attempt * (wallTime * 2)` or the server-provided `Retry-After` seconds.
- **Response Decoding**: `unmarshalResponse` decodes into:
  - `*string`: raw string content.
  - `*[]byte`: raw bytes.
  - Any struct / map pointer: JSON unmarshaled.

### 3.4 Error Handling & Sentinel Errors (`errors.go`)

`*ensembl.APIError` encapsulates the HTTP status code, server message, rate limit state, and raw response body. It implements `Unwrap()` to support idiomatic Go standard library error inspections:

```go
var record ensembl.LookupRecord
err := client.GetLookupByID(ctx, "UNKNOWN_ID", &record)
if err != nil {
    // 1. Sentinel inspection with errors.Is
    if errors.Is(err, ensembl.ErrNotFound) {
        // HTTP 404
    } else if errors.Is(err, ensembl.ErrBadRequest) {
        // HTTP 400
    } else if errors.Is(err, ensembl.ErrRateLimit) {
        // HTTP 429
    } else if errors.Is(err, ensembl.ErrServiceUnavailable) {
        // HTTP 503
    } else if errors.Is(err, ensembl.ErrMaxRetriesReached) {
        // Retry attempts exhausted
    }

    // 2. Structured error inspection with errors.As
    var apiErr *ensembl.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Message)
    }
}
```

### 3.5 Dual Dispatch Model (`endpoints.go`)

1. **Strongly-Typed Domain Methods**: Clean, exported methods on `*Client` across 16 domain files (e.g. `client.GetLookupByID(ctx, id, target, opts...)`).
2. **Dynamic Dispatch**: `client.Call(ctx, apiCallName, pathParams, bodyData, target, opts...)` allows invoking any endpoint by string key against `EndpointsTable`.
3. **Endpoint Catalog Introspection**: `client.Endpoints()` returns a clone of all available `EndpointSpec` metadata (name, doc, method, URL, content type, POST parameters).

### 3.6 Request Options (`types.go`)

Per-call customization using `RequestOption`:
- `WithQuery(key, value string)`: Add a single URL query param.
- `WithQueryParams(map[string]string)`: Set multiple query params.
- `WithURLValues(url.Values)`: Set query params from `url.Values`.
- `WithRequestContentType(contentType string)`: Override Content-Type / Accept (e.g. `"text/x-fasta"`, `"text/x-gff3"`).
- `WithRequestHeader(key, value string)`: Add a custom header to a specific request.

---

## 4. Development Workflow & Commands

All standard development tasks are orchestrated via the `Makefile` or standard Go toolchain.

### Make Targets

| Target | Command | Purpose |
|---|---|---|
| `make all` | `lint test-race build` | **Default**: Runs linter, race-detected tests, and builds all binaries |
| `make build` | `go build -v ./... && go build -v -o bin/example ./examples` | Compiles the package and example executable |
| `make test` | `go test -v ./...` | Runs offline unit tests |
| `make test-race` | `go test -v -race -cover ./...` | Runs unit tests with Go data race detector and coverage report |
| `make test-live` | `ENSEMBL_LIVE_TESTS=1 go test -v -tags=live -race -run TestLiveEnsemblAPI ./...` | Runs integration smoke tests against live `rest.ensembl.org` |
| `make test-coverage` | `go test -v -race -coverprofile=coverage.out ... && go tool cover -html=coverage.out` | Generates HTML test coverage report (`coverage.html`) |
| `make lint` | `go vet ./...` | Runs Go static analysis / vetting |
| `make format` | `go fmt ./...` | Formats all Go source files according to standard Go styling |
| `make example` | `go run ./examples/main.go` | Executes the interactive example application |
| `make clean` | `rm -rf bin/ coverage.out coverage.html` | Removes build artifacts and coverage files |

### Direct Go Toolchain Equivalents

- **Run all offline tests**:
  ```bash
  go test -v -race ./...
  ```
- **Check formatting**:
  ```bash
  gofmt -s -l .
  ```
- **Verify go.mod cleanliness**:
  ```bash
  go mod tidy
  git diff --exit-code go.mod
  ```

---

## 5. Coding Standards & Conventions

When modifying or extending `goensemblrest`, agents MUST adhere strictly to the following guidelines:

1. **Zero External Dependencies**:
   - Never add third-party dependencies to `go.mod`. Use only the Go standard library.
2. **Context Propagation**:
   - All network and API-dispatching methods MUST accept `ctx context.Context` as their first parameter.
3. **Target Output Parameter**:
   - Domain methods MUST accept `target any` as the destination pointer for decoding the response.
   - If `target` is `nil`, the response body is discarded safely without error.
4. **Variadic Request Options**:
   - Domain methods MUST accept `opts ...RequestOption` at the end of their signature.
5. **Path Parameter Safety**:
   - In `resolvePath`, never let colons (`:`) in genomic region identifiers (such as chromosome coordinates `13:32889611..32973805:1`) be converted to `%3A`.
6. **Thread Safety**:
   - Any mutable state in `Client` or `rateLimiter` must be guarded with appropriate mutex locks (`sync.RWMutex` / `sync.Mutex`).
7. **Error Wrapping & Mapping**:
   - Always return wrapped errors using `fmt.Errorf("%w", ...)` where appropriate.
   - Ensure HTTP status codes map to the correct sentinel error in `(*APIError).Unwrap()`.
8. **Formatting & Static Checks**:
   - All code must pass `gofmt -s`, `go vet ./...`, and `go test -race ./...` with zero failures or warnings.

---

## 6. How to Add or Modify an Endpoint

To add a new Ensembl REST API endpoint or update an existing one, follow this step-by-step procedure:

1. **Update `EndpointsTable` in `endpoints.go`**:
   - Add the specification to `EndpointsTable` with `Name`, `Doc`, `URL` (using `{{param}}` for path variables), `Method` (e.g. `http.MethodGet` or `http.MethodPost`), `ContentType`, and `PostParameters` (if POST).
2. **Add Domain Wrapper Method**:
   - Locate the corresponding domain file (e.g., `lookup.go`, `variation.go`, `comparative.go`) or create a new `[domain].go` file if adding a new domain category.
   - Define the typed method on `*Client`:
     ```go
     // GetLookupByExample retrieves ...
     func (c *Client) GetLookupByExample(ctx context.Context, param string, target any, opts ...RequestOption) error {
         params := map[string]string{"param": param}
         return c.Call(ctx, "getLookupByExample", params, nil, target, opts...)
     }
     ```
3. **Add Data Models in `types.go` (if applicable)**:
   - Define clean, exported Go structs with appropriate `json:"field_name,omitempty"` tags matching the Ensembl JSON schema.
4. **Add Unit Test in `endpoints_test.go`**:
   - Add a test case in `TestAllDomainEndpoints` under the relevant domain subtest using the mock `httptest.Server`.
5. **Add Live Smoke Test in `live_test.go` (optional/if stable)**:
   - If the endpoint operates on stable reference data, add a smoke test in `TestLiveEnsemblAPI`.
6. **Validate Changes**:
   - Run `make format && make lint && make test-race && make build`.

---

## 7. CI/CD & Testing Pipelines

The repository uses GitHub Actions for continuous integration and automated releases:

1. **Pull Request Workflow (`.github/workflows/pull_request.yaml`)**:
   - Triggers on PRs and pushes to `main` modifying Go files or workflow configs.
   - Runs `go mod tidy` check, `go vet ./...`, `gofmt` verification.
   - Executes offline unit tests with race detection (`-race`) and atomic coverage profiling (`-coverprofile=coverage.out`).
   - Compiles the `./examples` binary.
   - Executes live smoke tests against `rest.ensembl.org` (`continue-on-error: true` to prevent third-party network outages from blocking PRs).
2. **Nightly API Drift Check (`.github/workflows/nightly.yaml`)**:
   - Runs on schedule at `03:00 UTC` daily.
   - Executes the live test suite against the production Ensembl REST API to detect upstream contract changes, deprecations, or schema drift.
3. **Release Workflow (`.github/workflows/push_tag.yaml`)**:
   - Triggers when pushing git tags matching `v*` (e.g. `v0.1.0`).
   - Runs full test suite, builds multi-platform example artifacts, and publishes a GitHub Release with auto-generated release notes.
