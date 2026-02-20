## 1. Introduction

**Project Name:** CryptoView
**Type:** Desktop Application (Utility)
**Platforms:** Windows, Linux, macOS
**Stack:** Go 1.22+, Fyne (GUI), Public HTTP API
**Goal:** Create a lightweight, cross-platform application for monitoring cryptocurrency exchange rates.
**Educational Goal:** To master network operations (`net/http`), JSON parsing, concurrency patterns (goroutines/channels), and binding data to a UI in Go.

## 2. Functional Requirements

### 2.1. Cryptocurrency List

The application must display a **Scrollable List** of the following assets:

1. **Bitcoin (BTC)**
2. **Ethereum (ETH)**
3. **TON Coin (TON)**
4. **Solana (SOL)**
5. **Dogecoin (DOGE)**
6. **Ripple (XRP)**
7. **Litecoin (LTC)**

### 2.2. Currency Card Data

Each row in the list must contain:

* **Icon/Logo** of the coin (locally stored in assets).
* **Name and Ticker** (e.g., Bitcoin | BTC).
* **Current Price** in the selected fiat currency.
* **24h Change** (in percentage, color-coded red/green).
* **Last Update Time** (HH:MM:SS format).

### 2.3. Controls & Settings

The Toolbar or Header must contain:

1. **Fiat Currency Selector:** Dropdown/Select -> USD ($), EUR (€), RUB (₽). *Action:* Instantly triggers a data refresh/recalculation.
2. **Language Switch:** Button or Toggle -> RU / EN. (Default:  **EN** ).
3. **Theme Switch:** Icon Button (Sun/Moon) -> Light / Dark.
4. **Refresh Button:** Manually force a data update (in addition to auto-refresh).

### 2.4. Network Interaction

* **API:** Use a public API (Recommended: **CoinGecko API** or **CryptoCompare** — free tier, no complex auth required for basic endpoints).
* **Auto-Refresh:** Data must update automatically every  **60 seconds** .
* **Error Handling:** If there is no internet connection, show a "Network Error" banner or toast notification. The app  **must not crash** .

---

## 3. UI/UX (User Interface)

### 3.1. Layout

Use `fyne.Container` with a `layout.Border`.

* **Top (Header):**
  * *Left:* App Logo (small).
  * *Center:* Empty/Spacer.
  * *Right:* Select (USD/RUB/EUR), Icon (Lang), Icon (Theme).
* **Center (Body):**
  * `widget.List`: An optimized list widget. *Critical:* Use `List` instead of a simple `VBox` to ensure memory efficiency (view recycling).

### 3.2. States

* **Loading:** While fetching data for the first time -> `widget.ProgressBarInfinite`.
* **Error:** If the API request fails -> Red status text at the bottom of the list.

---

## 4. Technical Architecture

**Directory Structure (Simplified Clean Architecture):**

**Plaintext**

```
CryptoView/
├── cmd/
│   └── cryptoview/
│       └── main.go        # App initialization & entry point
├── internal/
│   ├── api/               # NETWORK LAYER
│   │   ├── client.go      # HTTP Client configuration (timeouts)
│   │   └── provider.go    # Logic for CoinGecko requests (GetRates)
│   ├── model/             # DATA STRUCTURES
│   │   └── currency.go    # JSON structs for parsing API responses
│   └── ui/                # PRESENTATION LAYER
│       ├── main_window.go # Window assembly
│       ├── crypto_list.go # List rendering logic
│       └── toolbar.go     # Control buttons
├── resources/             # STATIC ASSETS
│   ├── icons/             # btc.png, eth.png, etc.
│   └── i18n/              # en.json, ru.json (or a simple string map)
└── go.mod
```

### Key Technical Challenges (For Learning):

1. **JSON Unmarshalling:** Mapping Go structs to API responses (e.g., tags like `json:"current_price"`).
2. **Concurrency:** Network requests **must not block** the UI thread.
   * *Bad:* Click button -> Interface freezes -> Data arrives.
   * *Good:* Click button -> Launch goroutine `go fetchPrices()` -> UI remains responsive -> Data arrives via Channel -> Update `List`.
3. **HTTP Client:** Configuring `http.Client` with a timeout (e.g., 10 seconds) to prevent the app from hanging indefinitely if the API is down.

---

## 5. Implementation Plan (Roadmap)

1. **Day 1: Core & Network**
   * Define `structs` for the JSON response.
   * Write a `GetPrices(currency string)` function that prints data to the console.
2. **Day 2: Basic UI**
   * Initialize the Fyne window.
   * Create a basic list with Mock Data (static placeholders).
3. **Day 3: Integration**
   * Connect real data to the UI list.
   * Implement the currency switcher (re-fetching data with the new `vs_currency` parameter).
4. **Day 4: Polish**
   * Implement Theme switching (Dark/Light).
   * Add Localization (simple string map).
   * Build binaries for Windows and Linux.
