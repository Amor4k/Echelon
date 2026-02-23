package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Amor4k/Echelon/internal/parser"
)

func writeResultsToFile(results []parser.LogEntry, outputPath string, format string, ckey string) error {
	switch format {
	case "JSON (.json)":
		return writeJSON(results, outputPath)
	case "HTML (.html)":
		return writeHTML(results, outputPath, ckey)
	default:
		return writeLog(results, outputPath)
	}
}

func writeLog(results []parser.LogEntry, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range results {
		fmt.Fprintf(file, "[%s] [%s] %s\n", entry.Timestamp, entry.Category, entry.Message)
	}
	return nil
}

func writeJSON(results []parser.LogEntry, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")
	return encoder.Encode(results)
}

func writeHTML(results []parser.LogEntry, outputPath string, ckey string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	//HTML header
	fmt.Fprintf(file, `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>ECHELON - Log Analysis for %s</title>
    <style>
        body {
            font-family: 'Courier New', monospace;
            background-color: #1e1e1e;
            color: #d4d4d4;
            padding: 20px;
        }
        h1 {
            color: #4ec9b0;
            border-bottom: 2px solid #4ec9b0;
        }
        .log-entry {
            margin: 5px 0;
            padding: 5px;
            border-left: 3px solid #666;
        }
        .timestamp {
            color: #569cd6;
            font-weight: bold;
        }
        .category {
            display: inline-block;
            padding: 2px 8px;
            border-radius: 3px;
            font-size: 0.9em;
            margin: 0 5px;
        }
        .attack {
            background-color: #f48771;
            color: #000;
        }
        .game {
            background-color: #4ec9b0;
            color: #000;
        }
        .message {
            color: #ce9178;
        }
    </style>
</head>
<body>
    <h1>ECHELON Log Analysis</h1>
    <p><strong>Player:</strong> %s</p>
    <p><strong>Total Entries:</strong> %d</p>
    <hr>
`, ckey, ckey, len(results))

	//Write each log entry

	for _, entry := range results {
		categoryClass := "game"
		if entry.Category == "attack" {
			categoryClass = "attack"
		}

		fmt.Fprintf(file, `    <div class="log-entry">
        <span class="timestamp">%s</span>
        <span class="category %s">%s</span>
        <span class="message">%s</span>
    </div>
`, entry.Timestamp, categoryClass, entry.Category, entry.Message)
	}

	// HTML footer
	fmt.Fprintf(file, `</body>
</html>`)

	return nil
}
