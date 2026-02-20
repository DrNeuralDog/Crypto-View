# NotebookLM Setup — CryptoView

## Предварительно

Выполни в терминале:
```bash
nlm login
```
Следуй инструкциям для аутентификации в Google/NotebookLM.

---

## Шаг 1: Создать ноутбук

```bash
nlm notebook create "CryptoView Debug & Docs"
```

Скопируй **Notebook ID** из вывода (формат: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`).

---

## Шаг 2: Добавить источники

Подставь `<NOTEBOOK_ID>` вместо ID из шага 1.

### Документация Go
```bash
nlm source add <NOTEBOOK_ID> --url https://go.dev/doc/ --title "Go Documentation"
nlm source add <NOTEBOOK_ID> --url https://go.dev/doc/effective_go --title "Effective Go"
nlm source add <NOTEBOOK_ID> --url https://pkg.go.dev/std --title "Go Standard Library"
```

### Документация Fyne
```bash
nlm source add <NOTEBOOK_ID> --url https://docs.fyne.io/ --title "Fyne Documentation"
nlm source add <NOTEBOOK_ID> --url https://docs.fyne.io/started/ --title "Fyne Getting Started"
nlm source add <NOTEBOOK_ID> --url https://developer.fyne.io/ --title "Fyne Developer Guide"
```

### Кроссплатформенность и Windows
```bash
nlm source add <NOTEBOOK_ID> --url https://docs.fyne.io/started/cross-compile --title "Fyne Cross-Compilation"
```

### PRD проекта
(выполни из корня проекта CryptoView)
```bash
nlm source add <NOTEBOOK_ID> --file "docs/PRD.md" --title "CryptoView PRD"
```

### (Опционально) Сообщество
```bash
nlm source add <NOTEBOOK_ID> --url https://stackoverflow.com/questions/tagged/golang --title "Go on Stack Overflow"
nlm source add <NOTEBOOK_ID> --url https://github.com/fyne-io/fyne --title "Fyne GitHub"
```

---

## Шаг 3: Обновить ProjectManifest.md

Вставь Notebook ID в `docs/AI_Development/ProjectManifest.md`:
```markdown
- **Docs/Debug Notebook ID:** <твой-notebook-id>
```

---

## Использование

После настройки можно задавать вопросы по проекту:
```bash
nlm query <NOTEBOOK_ID> "Как в Fyne сделать widget.List с обновлением данных?"
```
