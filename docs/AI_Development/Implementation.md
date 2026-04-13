# CryptoView Implementation Status

## Current Architecture

CryptoView is implemented as a layered Go desktop application:

- `cmd/cryptoview` is the composition root.
- `internal/catalog` stores tracked coin metadata and alias mapping.
- `internal/providers` contains concrete HTTP integrations for market and FX sources.
- `internal/marketfeed` orchestrates polling, fallback, caching, cooldowns, and fiat recalculation.
- `internal/ui` contains Fyne windows, widgets, i18n, footer status, and theme handling.
- `resources` embeds runtime assets with `go:embed`.

## Implemented Product Behavior

- Fixed-size desktop window with tracked crypto prices
- Fiat switching for `USD`, `EUR`, `RUB`
- Language switching for `EN`, `RU`
- Light and dark themes
- Automatic market polling
- Provider fallback on API/network failure
- Cached market snapshot usage when all live providers fail
- Cached FX-based fiat recalculation without forcing a new market request
- Embedded app/coin icons resolved by canonical coin ID

## Quality Baseline

- Constructor validation returns errors instead of panicking
- Runtime provider logging uses injected `log/slog`
- Domain models no longer carry UI-only fields
- UI assets are bundled instead of loaded from relative filesystem paths
- Local quality-check flow is standardized through `test.ps1`

## Remaining Maintenance Ideas

- Add hosted CI once forge/tooling decisions are finalized
- Add snapshot/UI regression checks if the project starts changing appearance frequently
- Expand provider tests for more malformed payload scenarios if new external sources are added
