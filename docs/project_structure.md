# Project Structure - CryptoView

## Root Directory

```
CryptoView/
├── cmd/
│   └── cryptoview/
│       └── main.go          # Entry point. App init, Window assembly, run.
├── internal/
│   ├── model/               # Pure data structures
│   │   └── coin.go          # Structs with json tags (Coin, PriceResponse)
│   ├── api/                 # Network layer
│   │   ├── client.go        # APIClient interface, http.Client config (timeouts)
│   │   └── coingecko.go     # CoinGecko API implementation
│   ├── service/             # Business logic layer
│   │   └── price_service.go # Caching, currency conversion, auto-refresh
│   └── ui/                  # Presentation layer (Fyne)
│       ├── app.go           # Theme and app initialization
│       ├── main_window.go   # Main window layout assembly
│       ├── components/      # UI components
│       │   ├── coin_list.go # List widget logic (data binding)
│       │   └── toolbar.go   # Currency, theme, language controls
│       └── theme/           # Custom fonts/colors (if needed)
├── pkg/                     # (Optional) Shared utilities
│   └── i18n/
│       └── loc.go           # Localization (ru/en string map)
├── resources/               # Static assets
│   ├── icon.png             # App icon
│   └── coins/               # Coin logos (btc.png, eth.png, etc.)
├── docs/                    # Project documentation
│   ├── PRD.md
│   ├── Implementation.md
│   ├── project_structure.md
│   ├── UI_UX_doc.md
│   ├── CryptoView_Implementation.md  # Task tracking
│   ├── Archive/             # Outdated docs archive
│   ├── OtherHelpfulDocs/    # Reference materials
│   └── WorkflowLogs/        # Development logs
│       ├── DevelopmentLog.md
│       ├── BugLog.md
│       ├── GitLog.md
│       └── UserInteractionLog.md
├── Makefile                 # Build commands (build, run, clean)
├── go.mod
└── go.sum
```

## Detailed Structure

### cmd/cryptoview/
Точка входа приложения. Инициализация Fyne App, создание главного окна, запуск event loop.

### internal/model/
Чистые структуры данных для парсинга JSON API. Без бизнес-логики.

### internal/api/
HTTP-клиент и провайдеры внешних API (CoinGecko). Таймауты, retry, error handling.

### internal/service/
Бизнес-логика: кэширование цен, конвертация валют, автообновление по таймеру.

### internal/ui/
Весь код Fyne: окна, виджеты, темы. Связь с service через каналы/колбэки.

### pkg/
Переиспользуемые утилиты. i18n — простая локализация (map[string]string).

### resources/
Статические файлы: иконки приложения, логотипы монет. Встраиваются через `go:embed`.

### docs/
- **PRD.md** — Product Requirements Document
- **Implementation.md** — План реализации
- **project_structure.md** — Структура проекта (этот файл)
- **UI_UX_doc.md** — Спецификации UI/UX
- **CryptoView_Implementation.md** — Трекинг задач (completed/pending)
- **WorkflowLogs/** — Логи разработки, багов, Git, взаимодействий
- **Archive/** — Архив устаревших документов
- **OtherHelpfulDocs/** — Справочные материалы (не изменять агенту)

## Configuration

- **go.mod** — Go modules, зависимости
- **Makefile** — `make build`, `make run`, `make clean`
- **.gitignore** — Исключения для Git

## Build & Deployment

- `go build ./cmd/cryptoview` — сборка
- `fyne package` — упаковка для Windows/Linux/macOS
- Кросс-компиляция: `GOOS=linux GOARCH=amd64 go build ...`

## Code Style

Проект на **Go**. Используется стандартный `gofmt` и `go vet`. Стиль: [Effective Go](https://go.dev/doc/effective_go).
