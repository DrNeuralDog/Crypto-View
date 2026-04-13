# CryptoView UI / UX Notes

## Product Shape

CryptoView is a compact desktop utility for glanceable crypto price checks.
The UX goal is speed and clarity, not dense analytics.

## Main Window

- Fixed-size desktop window
- Header with logo, fiat selector, language selector, and theme toggle
- Scrollable center list with one row per tracked coin
- Footer with loading / ok / warning / error status

## Coin Row

Each row presents:

- embedded coin icon
- coin name and ticker
- current price in the selected fiat
- 24h percentage change with color coding
- last update time formatted for the active locale

## Interaction Rules

- Fiat switching is local and immediate when cached FX data is available
- Language switching updates header, footer, and formatted text without rebuilding the application
- Theme switching toggles between light and dark palettes
- Closing the window must stop the feed cleanly

## Status UX

- `Loading`: visible progress indicator
- `OK`: healthy provider status
- `Warning`: fallback provider active, rate-limited cache usage, or offline cached mode
- `Error`: no market data available

## Visual Constraints

- Prioritize readability over density
- Keep the window compact
- Preserve high contrast for price-change colors
- Use embedded assets so the UI does not depend on relative filesystem paths at runtime

## Accessibility / Maintainability

- Text should remain legible in both themes
- Controls use stable internal values (`EN`, `RU`, `USD`, `EUR`, `RUB`)
- UI formatting lives in the i18n layer, not in domain or provider code
