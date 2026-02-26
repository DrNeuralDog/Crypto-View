package ui

import (
	"image/color"

	"cryptoview/internal/ui/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FooterState string

const (
	FooterStateLoading FooterState = "loading"
	FooterStateOK      FooterState = "ok"
	FooterStateWarning FooterState = "warning"
	FooterStateError   FooterState = "error"
)

type FooterController struct {
	root       fyne.CanvasObject
	translator *i18n.Translator

	indicator    *canvas.Circle
	statusLabel  *widget.Label
	statusValue  *canvas.Text
	progress     *widget.ProgressBarInfinite
	progressWrap *fyne.Container

	state      FooterState
	customText string
}

func NewFooterController(translator *i18n.Translator) *FooterController {
	if translator == nil {
		translator = i18n.NewTranslator(i18n.LangEN)
	}

	indicator := canvas.NewCircle(color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF})
	indicatorWrap := container.NewGridWrap(fyne.NewSize(8, 8), indicator)

	statusLabel := widget.NewLabel("")
	statusValue := canvas.NewText("", color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF})
	statusValue.TextSize = 16

	progress := widget.NewProgressBarInfinite()
	progress.Hide()
	progressWrap := container.NewPadded(progress)
	progressWrap.Hide()

	row := container.NewHBox(
		container.NewCenter(indicatorWrap),
		container.NewCenter(statusLabel),
		container.NewCenter(statusValue),
	)
	root := container.NewVBox(container.NewPadded(row), progressWrap)

	controller := &FooterController{
		root:         root,
		translator:   translator,
		indicator:    indicator,
		statusLabel:  statusLabel,
		statusValue:  statusValue,
		progress:     progress,
		progressWrap: progressWrap,
	}
	controller.SetOK()
	return controller
}

func (f *FooterController) CanvasObject() fyne.CanvasObject {
	return f.root
}

func (f *FooterController) SetLanguage(language i18n.AppLanguage) {
	f.translator.SetLanguage(language)
	f.applyState()
}

func (f *FooterController) SetLoading() {
	f.state = FooterStateLoading
	f.customText = ""
	f.applyState()
}

func (f *FooterController) SetOK() {
	f.state = FooterStateOK
	f.customText = ""
	f.applyState()
}

func (f *FooterController) SetOKWithMessage(msg string) {
	f.state = FooterStateOK
	f.customText = msg
	f.applyState()
}

func (f *FooterController) SetError(msg string) {
	f.state = FooterStateError
	f.customText = msg
	f.applyState()
}

func (f *FooterController) SetWarning(msg string) {
	f.state = FooterStateWarning
	f.customText = msg
	f.applyState()
}

func (f *FooterController) applyState() {
	f.statusLabel.SetText(f.translator.T("status.label"))

	switch f.state {
	case FooterStateLoading:
		f.indicator.FillColor = color.NRGBA{R: 0x90, G: 0xA4, B: 0xAE, A: 0xFF}
		f.indicator.Refresh()
		f.statusValue.Text = f.translator.T("status.loading")
		f.statusValue.Color = color.NRGBA{R: 0x90, G: 0xA4, B: 0xAE, A: 0xFF}
		f.statusValue.Refresh()
		f.progress.Show()
		f.progressWrap.Show()
		f.progress.Start()
	case FooterStateError:
		f.indicator.FillColor = color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}
		f.indicator.Refresh()
		message := f.customText
		if message == "" {
			message = f.translator.T("status.error.network")
		}
		f.statusValue.Text = message
		f.statusValue.Color = color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}
		f.statusValue.Refresh()
		f.progress.Stop()
		f.progress.Hide()
		f.progressWrap.Hide()
	case FooterStateWarning:
		f.indicator.FillColor = color.NRGBA{R: 0xFF, G: 0x98, B: 0x00, A: 0xFF}
		f.indicator.Refresh()
		message := f.customText
		if message == "" {
			message = f.translator.T("status.warning.cached")
		}
		f.statusValue.Text = message
		f.statusValue.Color = color.NRGBA{R: 0xFF, G: 0x98, B: 0x00, A: 0xFF}
		f.statusValue.Refresh()
		f.progress.Stop()
		f.progress.Hide()
		f.progressWrap.Hide()
	default:
		f.indicator.FillColor = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
		f.indicator.Refresh()
		message := f.customText
		if message == "" {
			message = f.translator.T("status.ok")
		}
		f.statusValue.Text = message
		f.statusValue.Color = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
		f.statusValue.Refresh()
		f.progress.Stop()
		f.progress.Hide()
		f.progressWrap.Hide()
	}
}
