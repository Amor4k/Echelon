package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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

	pinnedContainer *fyne.Container
	logsContainer   *fyne.Container
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
	e.pinnedContainer = e.createPinnedSection()

	//All logs section (scrollable)
	e.logsContainer = e.createLogsSection()

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
		container.NewVSplit(e.pinnedContainer, e.logsContainer),
	)
}

//Sub methods

func (e *LogEditor) createPinnedSection() *fyne.Container {
	// TODO: Implement pinned logs display.
	return container.NewVBox(widget.NewLabel("Pinned logs will appear here"))
}

func (e *LogEditor) createLogsSection() *fyne.Container {
	logEntries := container.NewVBox()

	for i, entry := range e.results {
		//skip delted entries
		if e.deletedIndices[i] {
			continue
		}

		logEntries.Add(e.createLogEntry(i, entry))
	}

	//Wrap in a scroll container
	scroll := container.NewScroll(logEntries)

	return container.NewBorder(
		widget.NewLabel("All logs:"),
		nil,
		nil,
		nil,
		scroll,
	)
}

func (e *LogEditor) createLogEntry(index int, entry parser.LogEntry) *fyne.Container {
	//Timestamp (blue)
	timestamp := widget.NewRichTextFromMarkdown("**[" + entry.Timestamp + "]**")
	timestamp.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNamePrimary

	// var categoryColor theme.ColorName
	// if entry.Category == "attack" {
	// 	categoryColor = theme.ColorNameError //Red
	// } else {
	// 	categoryColor = theme.ColorNameSuccess //Green
	// }

	category := widget.NewLabel("[" + entry.Category + "]")
	category.Importance = widget.HighImportance

	//Message
	message := widget.NewLabel(entry.Message)
	message.Wrapping = fyne.TextWrapWord

	// Create a colored box for category

	categoryBox := container.NewHBox(category)

	//Log content (timestamp + category + message)

	logContent := container.NewHBox(
		timestamp,
		categoryBox,
		message,
	)

	// Pinning
	//Might need refactoring!

	pinBtn := widget.NewButton("📌 Pin", func() {
		e.togglePin(index)
	})
	pinBtn.Importance = widget.LowImportance

	// Delete btn
	deleteBtn := widget.NewButton("🗑️ Delete", func() {
		e.deleteLog(index)
	})
	deleteBtn.Importance = widget.DangerImportance

	// Action buttons
	actions := container.NewHBox(pinBtn, deleteBtn)

	// Full entry w/ seperator

	entryContainer := container.NewVBox(
		container.NewBorder(nil, nil, nil, actions, logContent),
		widget.NewSeparator(),
	)

	return entryContainer
}

func (e *LogEditor) togglePin(index int) {
	if e.pinnedIndices[index] {
		delete(e.pinnedIndices, index)
	} else {
		e.pinnedIndices[index] = true
	}
	e.refreshUI()
}

func (e *LogEditor) deleteLog(index int) {
	e.deletedIndices[index] = true
	e.refreshUI()
}

func (e *LogEditor) refreshUI() {
	// Rebuild pinned section
	e.pinnedContainer.Objects = e.createPinnedSection().Objects
	e.pinnedContainer.Refresh()

	// Rebuild logs section
	e.logsContainer.Objects = e.createLogsSection().Objects
	e.logsContainer.Refresh()

}

func (e *LogEditor) saveFilteredLogs(editorWindow fyne.Window, parentWindow fyne.Window) {
	//TODO: Implement save
}
