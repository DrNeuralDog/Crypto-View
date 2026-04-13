# CryptoView Project Structure

## Canonical Layout

```text
CryptoView/
├── cmd/
│   └── cryptoview/
│       └── main.go
├── internal/
│   ├── catalog/
│   ├── marketfeed/
│   ├── model/
│   ├── providers/
│   └── ui/
│       ├── assets/
│       ├── components/
│       ├── i18n/
│       └── theme/
├── resources/
│   ├── Logo/
│   ├── coins/
│   └── resources.go
├── docs/
│   ├── AI_Development/
│   ├── Archive/
│   ├── OtherHelpfulDocs/
│   └── designs/
├── build.ps1
├── build_all_os.ps1
├── Makefile
├── test.ps1
├── go.mod
└── go.sum
```

## Responsibilities

### `cmd/cryptoview`

Application startup only. This is where providers, the FX source, the feed, and the UI are composed together.

### `internal/catalog`

Single source of truth for:

- tracked coin order
- canonical coin IDs
- display names
- tickers
- icon asset keys
- alias normalization across external providers

### `internal/model`

Domain-level data models only. Models in this package should not carry UI formatting details or filesystem paths.

### `internal/providers`

Concrete external integrations:

- market-data providers
- FX provider
- HTTP request helpers
- response sanitation and retry-after parsing

### `internal/marketfeed`

Application orchestration layer:

- polling loops
- cooldown/backoff logic
- fallback chain execution
- cached data behavior
- fiat recalculation
- status emission to UI callbacks

### `internal/ui`

Presentation layer built with Fyne:

- main window assembly
- list/controller widgets
- footer states
- i18n formatting
- custom theme
- asset adapter over embedded resources

### `resources`

Static project assets plus the embedded resource registry used at runtime.
The source image files stay here so the PowerShell build scripts can still use them for Windows icon packaging.

## Notes

- `docs/AI_Development/project_structure.md` is the only canonical structure document.
- Old duplicate structure docs should not be used as sources of truth.
