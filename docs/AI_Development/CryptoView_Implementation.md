# CryptoView Task Tracking

## Current Stage

Stage 4 - Maintenance, Refactoring, and Quality

## Completed

- [x] Establish explicit composition root in `cmd/cryptoview`
- [x] Separate domain catalog, provider layer, feed orchestration, and UI layer
- [x] Replace panic-style feed construction with error-return construction
- [x] Centralize tracked coin metadata and alias normalization
- [x] Remove UI-only fields from the domain coin model
- [x] Inject feed into UI instead of constructing it inside the window layer
- [x] Embed runtime icons and branding with `go:embed`
- [x] Decouple language selector behavior from display-label parsing
- [x] Replace legacy `internal/api` and `internal/service/marketfeed` layout with the new package topology
- [x] Update unit/integration tests to the new architecture
- [x] Standardize local quality checks with `gofmt`, `go vet`, and `go test`
- [x] Sync README and AI-development docs with the current codebase

## Open

- [ ] Add hosted CI if the repository workflow requires remote enforcement
- [ ] Re-authenticate NotebookLM / reinstall `nlm` CLI if second-brain sync is needed again

## Next Recommended Focus

If product work resumes, the next practical area is feature growth on top of the new architecture rather than more structural cleanup.
