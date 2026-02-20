# CryptoView - Task Tracking (Implementation)

> **Назначение:** Отслеживание выполненных и текущих задач. Не путать с DevelopmentLog (логирование действий).

## Текущий этап: Stage 2 — Core Features (MVP)

## Задачи

### Stage 1: Foundation & Setup
- [x] Initialize Go module (`go mod init`)
- [x] Create directory structure (cmd, internal, resources, docs)
- [x] Add Fyne dependency
- [x] Create minimal `main.go` with empty window
- [x] Configure Makefile (build, run, clean)

### Stage 2: Core Features (MVP)
- [x] Define JSON structs for API response (model/coin.go)
- [x] Implement HTTP client with timeout (api/client.go)
- [x] Implement CoinGecko provider (api/coingecko.go)
- [ ] Create `GetPrices(currency string)` — console output first
- [x] Build basic Fyne window with `layout.Border`
- [x] Create list with mock data (widget.List)
- [x] Connect real API data to list
- [x] Implement fiat currency selector (USD/EUR/RUB)

### Stage 3: Advanced Features
- [ ] Implement theme switch (Light/Dark)
- [ ] Add localization (i18n, RU/EN)
- [ ] Add refresh button
- [ ] Implement 60s auto-refresh (goroutine + channel)
- [ ] Add network error handling (banner/toast, no crash)
- [ ] Add coin icons to resources

### Stage 4: Polish & Optimization
- [ ] Loading state (ProgressBarInfinite)
- [ ] Error state (red status text)
- [ ] Build binaries for Windows, Linux, macOS
- [ ] Final testing and bug fixes

---

**Следующая задача к выполнению:** Create `GetPrices(currency string)` � console output first




