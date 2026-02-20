CryptoView/
├── cmd/
│   └── cryptoview/
│       └── main.go          # Точка входа. Тут собираем App, Window и запускаем.
├── internal/
│   ├── model/               # Чистые данные (структуры)
│   │   └── coin.go          # Structs с json-тегами (Coin, PriceResponse)
│   ├── api/                 # Слой работы с сетью
│   │   ├── client.go        # Интерфейс APIClient и настройки http.Client (таймауты!)
│   │   └── coingecko.go     # Реализация запросов к CoinGecko
│   ├── service/             # Бизнес-логика (прослойка)
│   │   └── price_service.go # Кэширование, конвертация валют, авто-обновление
│   └── ui/                  # Весь Fyne код (Presentation Layer)
│       ├── app.go           # Инициализация темы и самого приложения
│       ├── main_window.go   # Сборка главного окна (Layout)
│       ├── components/      # Мелкие UI элементы
│       │   ├── coin_list.go # Логика виджета списка (Binding данных)
│       │   └── toolbar.go   # Кнопки валют, темы и языка
│       └── theme/           # Кастомные шрифты или цвета (если нужны)
├── pkg/                     # (Опционально) Общие утилиты
│   └── i18n/                # Локализация (простая мапа строк ru/en)
│       └── loc.go
├── resources/               # Статика
│   ├── icon.png             # Иконка приложения
│   └── coins/               # Папка с логотипами (btc.png, eth.png...)
├── Makefile                 # Команды сборки (build, run, clean)
├── go.mod
└── go.sum
