package main

//I have no idea about GUI's so we are vibecoding this bitch!
//Improvements/optimizations are welcome.

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Amor4k/Echelon/internal/parser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LogAnalyzer struct {
	ckey          string
	inputFiles    []string
	outputDir     string
	cleanMobIDs   bool
	afterMins     *float64
	beforeMins    *float64
	outputFormat  string
	progressBar   *widget.ProgressBar
	statusLabel   *widget.Label
	fileListLabel *widget.Label
}

type DropZone struct {
	box       *fyne.Container
	label     *widget.Label
	rect      *canvas.Rectangle
	isHovered bool
}

func NewDropZone() *DropZone {
	label := widget.NewLabel("Drag & drop log files here")
	label.Alignment = fyne.TextAlignCenter

	rect := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
	rect.CornerRadius = 5
	rect.StrokeWidth = 2
	rect.StrokeColor = theme.Color(theme.ColorNamePrimary)

	box := container.NewVBox(
		widget.NewSeparator(),
		label,
		widget.NewSeparator(),
	)

	box.Resize(fyne.NewSize(400, 60))

	return &DropZone{
		box:   box,
		label: label,
		rect:  rect,
	}
}

func (d *DropZone) SetHovered(hovered bool) {
	d.isHovered = hovered
	if hovered {
		d.rect.FillColor = theme.Color(theme.ColorNameHover)
		d.label.SetText("Drop files here.")
	} else {
		d.rect.FillColor = theme.Color(theme.ColorNameInputBackground)
		d.label.SetText("Drag & drop log files here")
	}
	d.rect.Refresh()
	d.label.Refresh()
}

func (d *DropZone) Highlight() {
	d.rect.FillColor = theme.Color(theme.ColorNameSuccess)
	d.label.SetText("Files added!")
	d.rect.Refresh()
	d.label.Refresh()

	time.AfterFunc(800*time.Millisecond, func() {
		fyne.Do(func() {
			d.SetHovered(false)
		})
	})
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ECHELON - SS13 Log Analyzer")
	dropZone := NewDropZone()

	//Output selection
	outputFormatSelect := widget.NewSelect([]string{"Readable Log (.log)", "JSON (.json)", "HTML (.html)"}, nil)
	outputFormatSelect.SetSelected("Readable Log (.log)") //Default

	analyzer := &LogAnalyzer{
		cleanMobIDs: true,
		inputFiles:  []string{},
	}

	iconSource, err := fyne.LoadResourceFromPath("cmd/gui/icon.png")
	if err == nil {
		myApp.SetIcon((iconSource))
	}

	//Drag & Drop handling

	myWindow.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		hasNewFiles := false

		for _, uri := range uris {
			path := uri.Path()
			ext := strings.ToLower(filepath.Ext(path))

			if ext != ".log" && ext != ".txt" && ext != ".json" {
				continue
			}

			// dedupe
			exists := false
			for _, f := range analyzer.inputFiles {
				if f == path {
					exists = true
					break
				}
			}

			if !exists {
				analyzer.inputFiles = append(analyzer.inputFiles, path)
				hasNewFiles = true
			}
		}

		if hasNewFiles {
			dropZone.Highlight()
		}

		analyzer.updateFileList()
	})

	// Input fields
	ckeyEntry := widget.NewEntry()
	ckeyEntry.SetPlaceHolder("Enter ckey(s)")

	// File list
	analyzer.fileListLabel = widget.NewLabel("No files selected")
	analyzer.fileListLabel.Wrapping = fyne.TextWrapWord

	// Select files button
	selectFilesBtn := widget.NewButton("Select Log Files", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				return
			}
			reader.Close()

			analyzer.inputFiles = append(analyzer.inputFiles, reader.URI().Path())
			analyzer.updateFileList()
		}, myWindow)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".log", ".json"}))
		fd.Show()
	})

	// Clear files button
	clearFilesBtn := widget.NewButton("Clear Files", func() {
		analyzer.inputFiles = []string{}
		analyzer.updateFileList()
	})

	// Output directory selection
	outputDirLabel := widget.NewLabel("Output: Same as first input file")
	selectOutputBtn := widget.NewButton("Select Output Directory", func() {
		fd := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if dir == nil {
				return
			}

			analyzer.outputDir = dir.Path()
			outputDirLabel.SetText("Output: " + analyzer.outputDir)
		}, myWindow)
		fd.Show()
	})

	// Options
	cleanMobIDs := widget.NewCheck("Clean Mob IDs", func(checked bool) {
		analyzer.cleanMobIDs = checked
	})
	cleanMobIDs.Checked = true

	// Time filtering options
	enableTimeFilter := widget.NewCheck("Enable Time Filtering", nil)

	afterMinsEntry := widget.NewEntry()
	afterMinsEntry.SetPlaceHolder("e.g., 10.5")
	afterMinsEntry.Disable()

	beforeMinsEntry := widget.NewEntry()
	beforeMinsEntry.SetPlaceHolder("e.g., 45.0")
	beforeMinsEntry.Disable()

	enableTimeFilter.OnChanged = func(checked bool) {
		if checked {
			afterMinsEntry.Enable()
			beforeMinsEntry.Enable()
		} else {
			afterMinsEntry.Disable()
			beforeMinsEntry.Disable()
		}
	}

	// Progress bar and status
	analyzer.progressBar = widget.NewProgressBar()
	analyzer.progressBar.Hide()
	analyzer.statusLabel = widget.NewLabel("")

	// Run button
	runButton := widget.NewButton("Filter Logs", func() {
		analyzer.ckey = strings.TrimSpace(ckeyEntry.Text)
		analyzer.outputFormat = outputFormatSelect.Selected

		if analyzer.ckey == "" {
			dialog.ShowInformation("Error", "Please enter a ckey", myWindow)
			return
		}

		if len(analyzer.inputFiles) == 0 {
			dialog.ShowInformation("Error", "Please select at least one log file", myWindow)
			return
		}

		// Parse time filters if enabled
		analyzer.afterMins = nil
		analyzer.beforeMins = nil

		if enableTimeFilter.Checked {
			if afterMinsEntry.Text != "" {
				var val float64
				_, err := fmt.Sscanf(afterMinsEntry.Text, "%f", &val)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Invalid 'After' minutes: %v", err), myWindow)
					return
				}
				analyzer.afterMins = &val
			}

			if beforeMinsEntry.Text != "" {
				var val float64
				_, err := fmt.Sscanf(beforeMinsEntry.Text, "%f", &val)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Invalid 'Before' minutes: %v", err), myWindow)
					return
				}
				analyzer.beforeMins = &val
			}
		}

		// Process files
		go analyzer.processFiles(myWindow)
	})

	// Layout
	fileButtonsBox := container.NewHBox(
		selectFilesBtn,
		clearFilesBtn,
	)

	timeFilterForm := container.NewVBox(
		enableTimeFilter,
		container.NewGridWithColumns(2,
			widget.NewLabel("After (minutes):"),
			afterMinsEntry,
			widget.NewLabel("Before (minutes):"),
			beforeMinsEntry,
		),
	)

	content := container.NewVBox(
		widget.NewLabelWithStyle("ECHELON - SS13 Log Analyzer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),

		widget.NewLabel("Player Ckey:"),
		ckeyEntry,

		widget.NewSeparator(),
		widget.NewLabel("Log Files:"),
		dropZone.box,
		fileButtonsBox,
		analyzer.fileListLabel,

		widget.NewSeparator(),
		widget.NewLabel("Output Directory:"),
		selectOutputBtn,
		outputDirLabel,

		widget.NewSeparator(),
		widget.NewLabel("Options:"),
		cleanMobIDs,
		widget.NewLabel("Output Format:"),
		outputFormatSelect,

		widget.NewSeparator(),
		widget.NewLabel("Time Filtering (relative to round start):"),
		timeFilterForm,

		widget.NewSeparator(),
		runButton,
		analyzer.progressBar,
		analyzer.statusLabel,
	)

	scrollContainer := container.NewScroll(content)
	myWindow.SetContent(scrollContainer)
	myWindow.Resize(fyne.NewSize(600, 800))
	myWindow.ShowAndRun()
}

func (a *LogAnalyzer) updateFileList() {
	if len(a.inputFiles) == 0 {
		a.fileListLabel.SetText("No files selected")
		return
	}

	fileList := fmt.Sprintf("%d file(s) selected:\n", len(a.inputFiles))
	for i, file := range a.inputFiles {
		fileName := filepath.Base(file)
		fileList += fmt.Sprintf("%d. %s\n", i+1, fileName)
	}
	a.fileListLabel.SetText(fileList)
}

func (a *LogAnalyzer) processFiles(window fyne.Window) {
	fyne.Do(func() {
		a.progressBar.Show()
		a.statusLabel.SetText("Validating round IDs...")
	})

	// Validate round IDs
	err := parser.ValidateRoundIDs(a.inputFiles)
	if err != nil {
		fyne.Do(func() {
			a.progressBar.Hide()
			a.statusLabel.SetText("Ready")
		})

		// Show warning and ask to continue
		dialog.ShowConfirm("Round ID Mismatch",
			fmt.Sprintf("%v\n\nDo you want to continue anyway?", err),
			func(continue_ bool) {
				if continue_ {
					a.doFiltering(window)
				} else {
					a.statusLabel.SetText("Cancelled")
				}
			}, window)
		return
	}

	a.doFiltering(window)
}

func (a *LogAnalyzer) doFiltering(window fyne.Window) {
	fyne.Do(func() {
		a.progressBar.SetValue(0.3)
		a.statusLabel.SetText("Filtering logs...")
	})

	// Create filter options
	opts := parser.FilterOptions{
		CleanMobIds: a.cleanMobIDs,
		AfterMins:   a.afterMins,
		BeforeMins:  a.beforeMins,
	}

	// Filter using the parser
	results, err := parser.FilterByCkey(a.inputFiles, a.ckey, opts)
	if err != nil {
		fyne.Do(func() {
			a.progressBar.Hide()
			a.statusLabel.SetText("Error")
			dialog.ShowError(fmt.Errorf("Filtering failed: %v", err), window)
		})

		return
	}

	fyne.Do(func() {
		a.progressBar.SetValue(0.7)
		a.statusLabel.SetText(fmt.Sprintf("Writing %d results to file...", len(results)))
	})

	// Determine output file path
	outputFile := a.getOutputPath()

	// Write results to JSON file
	err = writeResultsToFile(results, outputFile, a.outputFormat, a.ckey)
	if err != nil {
		fyne.Do(func() {
			a.progressBar.Hide()
			a.statusLabel.SetText("Error")
			dialog.ShowError(fmt.Errorf("Failed to write output: %v", err), window)
		})

		return
	}

	fyne.Do(func() {
		a.progressBar.SetValue(1.0)
		a.progressBar.Hide()
		a.statusLabel.SetText("Complete!")
		// Show success dialog
		successMsg := fmt.Sprintf("Successfully filtered %d log entries!\n\nOutput saved to:\n%s",
			len(results), outputFile)
		dialog.ShowInformation("Success", successMsg, window)
	})

}

func (a *LogAnalyzer) getOutputPath() string {
	dir := filepath.Dir(a.inputFiles[0])
	if a.outputDir != "" {
		dir = a.outputDir
	}

	timestamp := time.Now().Format("20060102_150405")

	var ext string
	switch a.outputFormat {
	case "JSON (.json)":
		ext = "json"
	case "HTML (.html)":
		ext = "html"
	default:
		ext = "log"
	}

	outputName := fmt.Sprintf("%s_%s_filtered.%s", a.ckey, timestamp, ext)

	return filepath.Join(dir, outputName)
}
