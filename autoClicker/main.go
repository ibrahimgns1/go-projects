package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	hook "github.com/robotn/gohook"
)

func main() {

	engine := NewAutoClicker()

	myApp := app.New()
	w := myApp.NewWindow("Go Auto Clicker v2")
	w.Resize(fyne.NewSize(300, 350))

	lblStatus := widget.NewLabel("Durum: BEKLİYOR")
	lblStatus.Alignment = fyne.TextAlignCenter
	entrySpeed := widget.NewEntry()
	entrySpeed.SetText("100")
	entrySpeed.PlaceHolder = "Milisaniye (örn: 100)"

	radioGroup := widget.NewRadioGroup([]string{"Sol Tık", "Sağ Tık"}, func(value string) {
		if value == "Sol Tık" {
			engine.SetButton("left")
		} else {
			engine.SetButton("right")
		}
	})
	radioGroup.SetSelected("Sol Tık")

	btnToggle := widget.NewButton("BAŞLAT (F8)", nil)

	toggleFunc := func() {
		if engine.IsRunning() {
			engine.Stop()
			lblStatus.SetText("Durum: DURDU")
			btnToggle.SetText("BAŞLAT (F8)")
			entrySpeed.Enable()
			radioGroup.Enable()
		} else {

			ms, err := strconv.Atoi(entrySpeed.Text)
			if err != nil || ms < 1 {
				lblStatus.SetText("Hata: Geçersiz Hız!")
				return
			}

			engine.SetInterval(ms)
			engine.Start()

			lblStatus.SetText(fmt.Sprintf("ÇALIŞIYOR (%d ms)", ms))
			btnToggle.SetText("DURDUR (F8)")
			entrySpeed.Disable()
			radioGroup.Disable()
		}
	}

	btnToggle.OnTapped = toggleFunc

	go func() {

		hook.Register(hook.KeyDown, []string{"f8"}, func(e hook.Event) {
			toggleFunc()
		})

		s := hook.Start()
		<-hook.Process(s)
	}()

	content := container.NewVBox(
		widget.NewLabel("Ayarlar"),
		widget.NewSeparator(),
		widget.NewLabel("Tıklama Tipi:"),
		radioGroup,
		widget.NewLabel("Hız (ms):"),
		entrySpeed,
		widget.NewSeparator(),
		lblStatus,
		btnToggle,
	)

	w.SetContent(content)
	w.Show()
	myApp.Run()
}
