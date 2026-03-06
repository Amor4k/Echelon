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
	pinnedScroll    *container.Scroll
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

	e.pinnedScroll = container.NewScroll(pinnedEntries)
	e.pinnedScroll.SetMinSize(fyne.NewSize(0, 150))

	// Toggle button
	toggleBtn := widget.NewButton("▼ Collapse", nil)
	toggleBtn.OnTapped = func() {
		if e.pinnedScroll.Visible() {
			e.pinnedScroll.Hide()
			toggleBtn.SetText("▶ Expand Pinned")
		} else {
			e.pinnedScroll.Show()
			toggleBtn.SetText("▼ Collapse")
		}
	}

	headerLabel := widget.NewLabelWithStyle("📌 Pinned Logs:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	header := container.NewBorder(nil, nil, headerLabel, toggleBtn, container.NewMax())

	return container.NewBorder(header, nil, nil, nil, e.pinnedScroll)
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
	// Get non-deleted entries
	var displayIndices []int
	for i := range e.results {
		if !e.deletedIndices[i] {
			displayIndices = append(displayIndices, i)
		}
	}

	list := widget.NewList(
		func() int {
			return len(displayIndices)
		},
		func() fyne.CanvasObject {
			return e.createLogEntryTemplate()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			e.updateLogEntry(displayIndices[id], item)
		},
	)

	return container.NewBorder(
		widget.NewLabel(fmt.Sprintf("All logs (%d entries):", len(displayIndices))),
		nil,
		nil,
		nil,
		list,
	)
}

func (e *LogEditor) createLogEntryTemplate() fyne.CanvasObject {
	timestamp := widget.NewLabel("")
	timestamp.TextStyle.Bold = true

	category := widget.NewLabel("")

	message := widget.NewLabel("")
	message.Wrapping = fyne.TextWrapWord

	pinBtn := widget.NewButton("📌", nil)
	pinBtn.Importance = widget.LowImportance

	deleteBtn := widget.NewButton("🗑️", nil)
	deleteBtn.Importance = widget.DangerImportance

	// Store widgets in a simple VBox for easy access
	return container.NewVBox(
		container.NewHBox(
			timestamp,
			category,
			widget.NewLabel(""), // Spacer
			pinBtn,
			deleteBtn,
		),
		message,
		widget.NewSeparator(),
	)
}

func (e *LogEditor) updateLogEntry(index int, item fyne.CanvasObject) {
	entry := e.results[index]

	// Access the VBox container
	vbox := item.(*fyne.Container)
	headerRow := vbox.Objects[0].(*fyne.Container)
	messageLabel := vbox.Objects[1].(*widget.Label)

	// Access header widgets
	timestamp := headerRow.Objects[0].(*widget.Label)
	category := headerRow.Objects[1].(*widget.Label)
	pinBtn := headerRow.Objects[3].(*widget.Button)
	deleteBtn := headerRow.Objects[4].(*widget.Button)

	// Update content
	timestamp.SetText("[" + entry.Timestamp + "]")
	category.SetText("[" + entry.Category + "]")
	messageLabel.SetText(entry.Message)

	// Update button callbacks
	pinBtn.OnTapped = func() {
		e.togglePin(index)
	}

	deleteBtn.OnTapped = func() {
		e.deleteLog(index)
	}

	// Update pin button text if already pinned
	if e.pinnedIndices[index] {
		pinBtn.SetText("✓ Pinned")
	} else {
		pinBtn.SetText("📌 Pin")
	}
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
