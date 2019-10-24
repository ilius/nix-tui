package nixtui

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	ui "github.com/ilius/clui"
)

func createView() {
	winWidth := 50
	winHeight := 10
	view := ui.AddWindow(0, 0, winWidth, winHeight, "Nix-TUI")
	view.SetTitleButtons(ui.ButtonMaximize | ui.ButtonClose)
	view.SetMaximized(true)
	winWidth, winHeight = view.Size()

	frame := ui.CreateFrame(view, ui.AutoSize, ui.AutoSize, ui.BorderNone, ui.AutoSize)
	frame.SetPack(ui.Vertical)
	frame.SetPaddings(1, 1) // border size
	frame.SetGaps(0, 0)     // gap size between children

	resizeHandlers := []func(event *ui.Event){}

	var searchField *ui.EditField
	{
		hframe := ui.CreateFrame(frame, winWidth-2, 1, ui.BorderNone, ui.Fixed)
		hframe.SetPack(ui.Horizontal)
		label := "Search:"
		width := len(label) + 1
		ui.CreateLabel(hframe, width, 1, label, ui.Fixed)
		searchField = ui.CreateEditField(hframe, winWidth-width-4, "", ui.Fixed)
		searchField.SetConstraints(width+10, 1)
		onResize := func(event *ui.Event) {
			hframe.SetSize(winWidth-2, 1)
			searchField.SetSize(winWidth-width-4, 1)
		}
		resizeHandlers = append(resizeHandlers, onResize)
	}

	inDescCheck := ui.CreateCheckBox(frame, ui.AutoSize, "In package descriptions as well", ui.AutoSize)

	searchButton := ui.CreateButtonSmall(frame, 8, "[Search]", ui.Fixed)
	searchButton.SetAlign(ui.AlignLeft)

	outputView := ui.CreateTextView(frame, ui.AutoSize, 5, ui.AutoSize)
	outputView.SetAutoScroll(true)
	outputPrint := func(format string, args ...interface{}) {
		outputView.AddText([]string{fmt.Sprintf(format, args...) + "\n"})
	}
	outputLog := func(format string, args ...interface{}) {
		outputPrint(time.Now().Format("15:04:05")+": "+format, args...)
	}

	view.OnScreenResize(func(event ui.Event) {
		// outputLog("event=%v", event)
		winWidth = event.Width
		winHeight = event.Height
		outputLog("winWidth=%v", winWidth)
		for _, handler := range resizeHandlers {
			handler(&event)
		}
	})

	searchButton.OnClick(func(event ui.Event) {
		inDesc := inDescCheck.State() == 1
		phrase := searchField.Title()
		outputLog("phrase = %#v, inDesc = %v", phrase, inDesc)
		searchButton.SetActive(false)
	})

	ui.CreateLabel(frame, 1, 1, "", ui.Fixed)

	// ui.CreateLabel(frame, winWidth-2, 1, "Press Ctrl+Q twice to exit!", 1)
	exitButton := ui.CreateButtonSmall(frame, 8, "[Exit]", ui.Fixed)
	exitButton.SetAlign(ui.AlignRight)
	exitButton.OnClick(func(event ui.Event) {
		ui.DeinitLibrary()
		os.Exit(0)
	})

	for _, handler := range resizeHandlers {
		handler(nil)
	}

	ui.ActivateControl(view, searchField)

}

func exitOnInterrupt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for sig := range signalChan {
			if sig.String() == "interrupt" {
				fmt.Println("\nYou pressed Ctrl+C, Goodbye :)")
				os.Exit(0)
			}
		}
	}()
}

func Main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	createView()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}
