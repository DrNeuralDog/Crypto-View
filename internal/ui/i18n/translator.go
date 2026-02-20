package i18n

import "sync"

type Translator struct {
	mu   sync.RWMutex
	lang AppLanguage
}

func NewTranslator(defaultLanguage AppLanguage) *Translator {
	lang := defaultLanguage
	if _, ok := translations[lang]; !ok {
		lang = LangEN
	}
	return &Translator{lang: lang}
}

func (t *Translator) SetLanguage(lang AppLanguage) {
	if _, ok := translations[lang]; !ok {
		return
	}
	t.mu.Lock()
	t.lang = lang
	t.mu.Unlock()
}

func (t *Translator) Language() AppLanguage {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.lang
}

func (t *Translator) T(key string) string {
	t.mu.RLock()
	lang := t.lang
	t.mu.RUnlock()

	if message, ok := translations[lang][key]; ok {
		return message
	}
	if message, ok := translations[LangEN][key]; ok {
		return message
	}
	return key
}

var translations = map[AppLanguage]map[string]string{
	LangEN: {
		"app.title":               "CryptoView",
		"status.label":            "Status:",
		"status.ok":               "OK",
		"status.loading":          "Loading...",
		"status.error.network":    "Network error",
		"toolbar.refresh.tooltip": "Refresh",
		"toolbar.lang.en":         "EN",
		"toolbar.lang.ru":         "RU",
	},
	LangRU: {
		"app.title":               "CryptoView",
		"status.label":            "Статус:",
		"status.ok":               "OK",
		"status.loading":          "Загрузка...",
		"status.error.network":    "Ошибка сети",
		"toolbar.refresh.tooltip": "Обновить",
		"toolbar.lang.en":         "EN",
		"toolbar.lang.ru":         "RU",
	},
}
