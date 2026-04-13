package appresources

import (
	"embed"
	"path"
	"sync"

	"cryptoview/internal/catalog"
	"fyne.io/fyne/v2"
)

//go:embed Logo/* coins/*
var embeddedFiles embed.FS

var (
	cacheMu sync.RWMutex
	cache   = make(map[string]fyne.Resource)
)

func AppIcon() fyne.Resource {
	return load("Logo/CryptoView Icon.png")
}

func CoinIcon(id string) fyne.Resource {
	coin, ok := catalog.Lookup(id)
	if !ok {
		return nil
	}
	return load(coin.IconAsset)
}

func load(name string) fyne.Resource {
	cacheMu.RLock()
	resource, ok := cache[name]
	cacheMu.RUnlock()
	if ok {
		return resource
	}

	data, err := embeddedFiles.ReadFile(name)
	if err != nil {
		return nil
	}

	resource = fyne.NewStaticResource(path.Base(name), data)

	cacheMu.Lock()
	cache[name] = resource
	cacheMu.Unlock()
	return resource
}
