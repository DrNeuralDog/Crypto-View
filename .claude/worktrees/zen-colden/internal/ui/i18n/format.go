package i18n

import (
	"fmt"
	"strings"
	"time"
)

func FormatPrice(value float64, fiat FiatCurrency, lang AppLanguage) string {
	symbol := currencySymbol(fiat)
	if symbol == "" {
		symbol = "$"
	}

	switch lang {
	case LangRU:
		return fmt.Sprintf("%s %s", formatDecimal(value, ' ', ','), symbol)
	default:
		return fmt.Sprintf("%s%s", symbol, formatDecimal(value, ',', '.'))
	}
}

func FormatTime(hhmmss string, _ AppLanguage) string {
	if _, err := time.Parse("15:04:05", hhmmss); err != nil {
		return "--:--:--"
	}
	return hhmmss
}

func currencySymbol(fiat FiatCurrency) string {
	switch fiat {
	case FiatUSD:
		return "$"
	case FiatEUR:
		return "\u20ac"
	case FiatRUB:
		return "\u20bd"
	default:
		return ""
	}
}

func formatDecimal(value float64, thousandSep rune, decimalSep rune) string {
	raw := fmt.Sprintf("%.2f", value)
	parts := strings.SplitN(raw, ".", 2)
	intPart := groupThousands(parts[0], thousandSep)
	if len(parts) < 2 {
		return intPart
	}
	return intPart + string(decimalSep) + parts[1]
}

func groupThousands(intPart string, sep rune) string {
	if len(intPart) <= 3 {
		return intPart
	}

	n := len(intPart)
	lead := n % 3
	if lead == 0 {
		lead = 3
	}

	var b strings.Builder
	b.WriteString(intPart[:lead])
	for i := lead; i < n; i += 3 {
		b.WriteRune(sep)
		b.WriteString(intPart[i : i+3])
	}
	return b.String()
}
