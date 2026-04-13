# UI/UX Documentation - CryptoView

## Overview

CryptoView — desktop-приложение для мониторинга курсов криптовалют. Минималистичный интерфейс с фокусом на читаемость и отзывчивость.

## Layout

### Structure

Используется `fyne.Container` с `layout.Border`:

```
┌─────────────────────────────────────────────────────────┐
│  [Logo]              [Spacer]    [USD▼] [Lang] [Theme]  │  ← Header (Top)
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │  [Icon] Bitcoin | BTC     $XX,XXX   +2.5%  12:34  │  │
│  │  [Icon] Ethereum | ETH   $X,XXX    -1.2%  12:34  │  │  ← Body (Center)
│  │  ...                                              │  │     widget.List
│  └─────────────────────────────────────────────────┘  │
│                                                         │
├─────────────────────────────────────────────────────────┤
│  [Status: OK / Loading / Network Error]                  │  ← Footer (Bottom)
└─────────────────────────────────────────────────────────┘
```

### Header (Top)
- **Left:** App logo (small)
- **Center:** Spacer
- **Right:** Fiat Select (USD/EUR/RUB), Language toggle (RU/EN), Theme toggle (Sun/Moon icon)

### Body (Center)
- **widget.List** — обязательно List, не VBox (view recycling, memory efficiency)
- Каждая строка: иконка | название | тикер | цена | 24h % | время обновления

### Footer (Bottom)
- Статус: OK / Loading / Network Error

## Currency Card

Каждая строка списка содержит:
- **Icon/Logo** — локальный файл (resources/coins/btc.png и т.д.)
- **Name | Ticker** — например "Bitcoin | BTC"
- **Current Price** — в выбранной фиатной валюте
- **24h Change** — в процентах, цвет: зелёный (рост), красный (падение)
- **Last Update** — формат HH:MM:SS

## Controls

| Control | Type | Action |
|---------|------|--------|
| Fiat Currency | Dropdown/Select | Мгновенное обновление данных |
| Language | Button/Toggle | RU / EN (default: EN) |
| Theme | Icon Button | Light / Dark |
| Refresh | Button | Ручное обновление данных |

## States

| State | UI Element |
|-------|------------|
| **Loading** | `widget.ProgressBarInfinite` |
| **Error** | Красный текст статуса внизу списка |
| **Normal** | Список с данными |

## Responsive Behavior

Desktop-only. Окно фиксированного или минимального размера. List скроллится при большом количестве элементов.

## Accessibility

- Контрастные цвета для 24h change (зелёный/красный)
- Читаемые шрифты (Fyne default)
- Кнопки с иконками + tooltip при необходимости

## Theme

- **Light:** Светлый фон, тёмный текст
- **Dark:** Тёмный фон, светлый текст
- Переключение через Fyne theme API
