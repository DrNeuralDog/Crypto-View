# Implementation Plan for CryptoView

## Feature Analysis

### Identified Features

1. **Cryptocurrency List** — Scrollable list of 7 coins (BTC, ETH, TON, SOL, DOGE, XRP, LTC)
2. **Currency Card** — Icon, name/ticker, price, 24h change (color-coded), last update time
3. **Fiat Selector** — USD, EUR, RUB dropdown, triggers data refresh
4. **Language Switch** — RU/EN toggle (default EN)
5. **Theme Switch** — Light/Dark mode
6. **Refresh Button** — Manual data update
7. **Network Layer** — CoinGecko/CryptoCompare API, 60s auto-refresh, error handling (no crash on network failure)

### Feature Categorization

- **Must-Have:** Cryptocurrency list, currency card data, fiat selector, network interaction, error handling
- **Should-Have:** Theme switch, language switch, refresh button
- **Nice-to-Have:** Custom icons, advanced animations

## Recommended Tech Stack

### Frontend (GUI)
- **Framework:** Fyne v2 — Cross-platform Go GUI, native look, good for desktop utilities
- **Documentation:** https://docs.fyne.io/

### Backend / Core
- **Language:** Go 1.22+
- **Documentation:** https://go.dev/doc/

### API
- **Provider:** CoinGecko API (free tier, no auth for basic endpoints)
- **Documentation:** https://www.coingecko.com/en/api/documentation

### Additional Tools
- **Build:** Makefile, `go build`
- **Packaging:** Fyne `fyne package` for distributable binaries

## Implementation Stages

### Stage 1: Foundation & Setup
**Duration:** 0.5–1 day  
**Dependencies:** None

#### Sub-steps:
- [ ] Initialize Go module (`go mod init`)
- [ ] Create directory structure (cmd, internal, resources, docs)
- [ ] Add Fyne dependency
- [ ] Create minimal `main.go` with empty window
- [ ] Configure Makefile (build, run, clean)

### Stage 2: Core Features (MVP)
**Duration:** 2–3 days  
**Dependencies:** Stage 1 completion

#### Sub-steps:
- [ ] Define JSON structs for API response (model/coin.go)
- [ ] Implement HTTP client with timeout (api/client.go)
- [ ] Implement CoinGecko provider (api/coingecko.go)
- [ ] Create `GetPrices(currency string)` — console output first
- [ ] Build basic Fyne window with `layout.Border`
- [ ] Create list with mock data (widget.List)
- [ ] Connect real API data to list
- [ ] Implement fiat currency selector (USD/EUR/RUB)

### Stage 3: Advanced Features
**Duration:** 1–2 days  
**Dependencies:** Stage 2 completion

#### Sub-steps:
- [ ] Implement theme switch (Light/Dark)
- [ ] Add localization (i18n, RU/EN)
- [ ] Add refresh button
- [ ] Implement 60s auto-refresh (goroutine + channel)
- [ ] Add network error handling (banner/toast, no crash)
- [ ] Add coin icons to resources

### Stage 4: Polish & Optimization
**Duration:** 1 day  
**Dependencies:** Stage 3 completion

#### Sub-steps:
- [ ] Loading state (ProgressBarInfinite)
- [ ] Error state (red status text)
- [ ] Build binaries for Windows, Linux, macOS
- [ ] Final testing and bug fixes

## Resource Links

- [Go Documentation](https://go.dev/doc/)
- [Fyne Documentation](https://docs.fyne.io/)
- [Fyne Getting Started](https://docs.fyne.io/started/hello.html)
- [CoinGecko API](https://www.coingecko.com/en/api/documentation)
- [Effective Go](https://go.dev/doc/effective_go)
