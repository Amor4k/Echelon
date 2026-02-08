package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
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

// Filter config struct
type FilterOptions struct {
	CleanMobIds bool
}

// I've noticed while filtering log entries that there is a mob number stated after the ckey ex: (mob_1234). This clutters the output, making it harder to read.
// Going to remove mob numbers from the ckey entries using regex. Might not be the best solution.
var mobPattern = regexp.MustCompile(`\(mob_\d+\)`)

// Filters logs by given ckey. Returns slice of LogEntry and error if any.
func FilterByCkey(filenames []string, ckey string, opts FilterOptions) ([]LogEntry, error) {
	var allResults []LogEntry
	for _, filename := range filenames {
		// Open log files.
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("Failed to open file: %w", err)
		}
		//defer file.Close()

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
				if opts.CleanMobIds {
					entry.Message = mobPattern.ReplaceAllString(entry.Message, "") //REmoves (mob_1234) from log entry for better legilbility.
				}
				results = append(results, entry)
			}

		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("Error reading file: %w", err)
		}
		file.Close()
		// Append results from this file to allResults
		allResults = append(allResults, results...)
	}

	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Timestamp < allResults[j].Timestamp
	})

	// Return result slice and error if any.
	return allResults, nil
}
