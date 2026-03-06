package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

// Creates pinned log entries at the top of the editor.
func (e *LogEditor) createPinnedSection() *fyne.Container {
	pinnedEntries := container.NewVBox()

	hasPinned := false
	for i, entry := range e.results {
		if e.pinnedIndices[i] && !e.deletedIndices[i] {
			hasPinned = true
			pinnedEntries.Add(e.createPinnedLogEntry(i, entry))
		}
	}

	if !hasPinned {
		pinnedEntries.Add(widget.NewLabel("No pinned logs yet"))
	}

	scroll := container.NewScroll(pinnedEntries)
	scroll.SetMinSize(fyne.NewSize(0, 150)) //Fixed height for pinned section, might change later

	return container.NewBorder(
		widget.NewLabelWithStyle("📌 Pinned Logs:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		scroll,
	)
}

// Pins Log entry, adds it to the pinned section
func (e *LogEditor) createPinnedLogEntry(index int, entry parser.LogEntry) *fyne.Container {

	timestamp := widget.NewRichTextFromMarkdown("**[" + entry.Timestamp + "]**")
	timestamp.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNamePrimary

	category := widget.NewLabel("[" + entry.Category + "]")
	category.Importance = widget.HighImportance

	message := widget.NewLabel(entry.Message)
	message.Wrapping = fyne.TextWrapWord

	logContent := container.NewHBox(
		timestamp,
		category,
		message,
	)

	//Unpin
	unpinBtn := widget.NewButton("📌 Unpin", func() {
		e.togglePin(index)
	})
	unpinBtn.Importance = widget.MediumImportance

	actions := container.NewHBox(unpinBtn)

	return container.NewVBox(
		container.NewBorder(nil, nil, nil, actions, logContent),
		widget.NewSeparator(),
	)
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
	//First, filter out deleted entries

	var filteredResults []parser.LogEntry
	for i, entry := range e.results {
		if !e.deletedIndices[i] {
			filteredResults = append(filteredResults, entry)
		}
	}

	// Save filtered results
	err := writeResultsToFile(filteredResults, e.outputPath, e.outputFormat, e.ckey)
	if err != nil {
		fyne.Do(func() {
			dialog.ShowError(fmt.Errorf("Failed to save: %v", err), editorWindow)
		})
		return
	}

	// Success msg
	fyne.Do(func() {
		successMsg := fmt.Sprintf("Successfully saved %d log entries!\n\nOutput saved to:\n%s", len(filteredResults), e.outputPath)
		dialog.ShowInformation("Success", successMsg, parentWindow)
		editorWindow.Close()
	})
}
