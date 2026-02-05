package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type LogEntry struct {
	Timestamp  string                 `json:"ts"`
	Category   string                 `json:"cat"`
	Message    string                 `json:"msg"`
	Data       interface{}            `json:"data"`
	WorldState map[string]interface{} `json:"w-state"`
	ID         int                    `json:"id"`
	Version    string                 `json:"s-ver"`
}

// Filters logs by given ckey. Returns slice of LogEntry and error if any.
func FilterByCkey(filename string, ckey string) ([]LogEntry, error) {
	// Open log files.
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	//TODO: Go through each line, parse JSON, check if ckey matches. If yes, append to result slice.
	var results []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		//parsing each line as json
		var entry LogEntry
		err := json.Unmarshal([]byte(line), &entry)
		//skip invalid json lines
		if err != nil {
			continue
		}

		if strings.Contains(entry.Message, ckey) { //This could filter false positives, ex: ckey: ben, this also filters ckey: ruben. Need to improve later.
			results = append(results, entry)
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	// Return result slice and error if any.

	return results, nil
}
