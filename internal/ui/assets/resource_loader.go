package assets

import (
	appresources "cryptoview/resources"
	"fyne.io/fyne/v2"
)

func LoadAppIcon() fyne.Resource {
	return appresources.AppIcon()
}

func LoadCoinIcon(id string) fyne.Resource {
	return appresources.CoinIcon(id)
}
