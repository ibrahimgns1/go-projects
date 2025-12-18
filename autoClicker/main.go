package main

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type modernTheme struct{}

func (m modernTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{R: 139, G: 92, B: 246, A: 255} // Violet 500
	case theme.ColorNameBackground:
		return color.RGBA{R: 15, G: 23, B: 42, A: 255} // Slate 900
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 30, G: 41, B: 59, A: 255} // Slate 800
	case theme.ColorNameForeground:
		return color.RGBA{R: 248, G: 250, B: 252, A: 255} // Slate 50
	case theme.ColorNameButton:
		return color.RGBA{R: 51, G: 65, B: 85, A: 255}
	case theme.ColorNameSeparator:
		return color.RGBA{R: 71, G: 85, B: 105, A: 255}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m modernTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m modernTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m modernTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNamePadding {
		return 12
	}
	if name == theme.SizeNameInputRadius {
		return 8
	}
	if name == theme.SizeNameSelectionRadius {
		return 8
	}
	return theme.DefaultTheme().Size(name)
}

func main() {
	engine := NewAutoClicker()

	myApp := app.New()
	myApp.Settings().SetTheme(&modernTheme{})

	w := myApp.NewWindow("Go Auto Clicker Pro")
	w.Resize(fyne.NewSize(360, 500))

	title := canvas.NewText("Auto Clicker", color.RGBA{R: 139, G: 92, B: 246, A: 255})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	lblStatus := widget.NewLabelWithStyle("DURUM: BEKLİYOR", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	entrySpeed := widget.NewEntry()
	entrySpeed.SetText("100")

	entryX := widget.NewEntry()
	entryX.SetText("0")
	entryY := widget.NewEntry()
	entryY.SetText("0")

	lblCurrentPos := widget.NewLabelWithStyle("0 , 0", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})

	radioGroup := widget.NewRadioGroup([]string{"Sol Tık", "Sağ Tık"}, func(value string) {
		if value == "Sol Tık" {
			engine.SetButton("left")
		} else {
			engine.SetButton("right")
		}
	})
	radioGroup.Horizontal = true
	radioGroup.SetSelected("Sol Tık")

	locationGroup := widget.NewRadioGroup([]string{"Mevcut", "Sabit"}, func(value string) {
		if value == "Mevcut" {
			engine.SetUseLocation(false)
			entryX.Disable()
			entryY.Disable()
		} else {
			engine.SetUseLocation(true)
			entryX.Enable()
			entryY.Enable()
		}
	})
	locationGroup.Horizontal = true
	locationGroup.SetSelected("Mevcut")

	btnGetPosition := widget.NewButtonWithIcon("Koordinat Al", theme.SearchIcon(), func() {
		x, y := robotgo.Location()
		entryX.SetText(fmt.Sprintf("%d", x))
		entryY.SetText(fmt.Sprintf("%d", y))
	})

	entryX.Disable()
	entryY.Disable()

	btnToggle := widget.NewButton("BAŞLAT (F8)", nil)
	btnToggle.Importance = widget.HighImportance

	toggleFunc := func() {
		if engine.IsRunning() {
			engine.Stop()
			lblStatus.SetText("DURUM: DURDURULDU")
			btnToggle.SetText("BAŞLAT (F8)")
			btnToggle.Importance = widget.HighImportance
			btnToggle.Refresh()
			entrySpeed.Enable()
			radioGroup.Enable()
			locationGroup.Enable()
			btnGetPosition.Enable()
			if locationGroup.Selected == "Sabit" {
				entryX.Enable()
				entryY.Enable()
			}
		} else {
			ms, err := strconv.Atoi(entrySpeed.Text)
			if err != nil || ms < 1 {
				lblStatus.SetText("HATA: Geçersiz Hız!")
				return
			}

			engine.SetInterval(ms)

			if locationGroup.Selected == "Sabit" {
				x, errX := strconv.Atoi(entryX.Text)
				y, errY := strconv.Atoi(entryY.Text)
				if errX != nil || errY != nil || x < 0 || y < 0 {
					lblStatus.SetText("HATA: Geçersiz Koordinat!")
					return
				}
				engine.SetLocation(x, y)
			}

			engine.Start()
			lblStatus.SetText(fmt.Sprintf("ÇALIŞIYOR (%d ms)", ms))
			btnToggle.SetText("DURDUR (F8)")
			btnToggle.Importance = widget.DangerImportance
			btnToggle.Refresh()

			entrySpeed.Disable()
			radioGroup.Disable()
			locationGroup.Disable()
			entryX.Disable()
			entryY.Disable()
			btnGetPosition.Disable()
		}
	}

	btnToggle.OnTapped = toggleFunc

	go func() {
		hook.Register(hook.KeyDown, []string{"f8"}, func(e hook.Event) {
			fyne.DoAndWait(func() {
				toggleFunc()
			})
		})
		s := hook.Start()
		<-hook.Process(s)
	}()

	go func() {
		for {
			x, y := robotgo.Location()
			fyne.DoAndWait(func() {
				lblCurrentPos.SetText(fmt.Sprintf("%d , %d", x, y))
			})
			time.Sleep(time.Millisecond * 100)
		}
	}()

	settingsCard := widget.NewCard("Genel Ayarlar", "", container.NewVBox(
		widget.NewLabel("Tıklama Tipi"),
		radioGroup,
		widget.NewLabel("Tıklama Hızı (ms)"),
		entrySpeed,
	))

	locationCard := widget.NewCard("Konum Ayarları", "", container.NewVBox(
		locationGroup,
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabel("X"), entryX),
			container.NewVBox(widget.NewLabel("Y"), entryY),
		),
		btnGetPosition,
		container.NewHBox(layout.NewSpacer(), widget.NewLabel("Mevcut:"), lblCurrentPos, layout.NewSpacer()),
	))

	content := container.NewPadded(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			settingsCard,
			locationCard,
			layout.NewSpacer(),
			lblStatus,
			btnToggle,
		),
	)

	w.SetContent(content)
	w.SetFixedSize(true)
	w.ShowAndRun()
}
