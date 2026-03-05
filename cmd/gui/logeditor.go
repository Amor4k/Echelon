package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Amor4k/Echelon/internal/parser"
)

type LogEditor struct {
	results        []parser.LogEntry
	pinnedIndices  map[int]bool
	deletedIndices map[int]bool
	ckey           string
	roundID        string
	outputPath     string
	outputFormat   string
}

func showLogEditor(results []parser.LogEntry, ckey string, roundID string, outputPath string, outputFormat string, parentWindow fyne.Window) {
	editor := &LogEditor{
		results:        results,
		pinnedIndices:  make(map[int]bool),
		deletedIndices: make(map[int]bool),
		ckey:           ckey,
		outputPath:     outputPath,
		outputFormat:   outputFormat,
	}

	editorWindow := fyne.CurrentApp().NewWindow("Log Editor - " + ckey + " - Round " + roundID)
	editorWindow.Resize(fyne.NewSize(1000, 700))

	//Build the Editor UI
	content := editor.buildUI(editorWindow, parentWindow)
	editorWindow.SetContent(content)
	editorWindow.Show()
}

func (e *LogEditor) buildUI(editorWindow fyne.Window, parentWindow fyne.Window) fyne.CanvasObject {

	//Pinned logs (at the top)
	pinnedSection := e.createPinnedSection()

	//All logs section (scrollable)
	logsSection := e.createLogsSection()

	//Action buttons
	saveBtn := widget.NewButton("Save Filtered Logs", func() {
		e.saveFilteredLogs(editorWindow, parentWindow)
	})

	closeBtn := widget.NewButton("Close Without Saving", func() {
		editorWindow.Close()
	})

	buttons := container.NewHBox(saveBtn, closeBtn)

	return container.NewBorder(
		widget.NewLabel(fmt.Sprintf("Viewing %d log entries for %s", len(e.results), e.ckey)),
		buttons,
		nil,
		nil,
		container.NewVSplit(pinnedSection, logsSection),
	)
}

//Sub methods

func (e *LogEditor) createPinnedSection() *fyne.Container {
	// TODO: Implement pinned logs display.
	return container.NewVBox(widget.NewLabel("Pinned logs will appear here"))
}

func (e *LogEditor) createLogsSection() *fyne.Container {
	//TODO: Implement scrollable log list
	return container.NewVBox(widget.NewLabel("All logs will appear here"))
}

func (e *LogEditor) saveFilteredLogs(editorWindow fyne.Window, parentWindow fyne.Window) {
	//TODO: Implement save
}
