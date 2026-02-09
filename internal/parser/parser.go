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
// Going to remove mob numbers from the ckey entries using regexp. Might not be the best solution.
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
		//defer file.Close() commented this out because of multiple files.

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

// As we support multiple log files, we should check if they are of the same round and warn user if not.
func GetRoundID(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()

		//As we are only checking one line, anon struct makes more sense
		var metadata struct {
			RoundID string `json:"round-id"`
		}
		//get round id
		if err := json.Unmarshal([]byte(line), &metadata); err == nil {
			return metadata.RoundID, nil
		}
	}

	return "", fmt.Errorf("could not read round ID from %s", filename)
}

// Checks if all log files are of the same round, warns user if not and prompts for continuation
func ValidateRoundIDs(filenames []string) error {

	//Get first file's round id
	firstRoundID, err := GetRoundID(filenames[0])
	if err != nil {
		return fmt.Errorf("Warning: %v", err)
	}

	//iterate through the rest of the log files to check if any of them has a different ID
	for _, logfile := range filenames[1:] {
		roundID, err := GetRoundID(logfile)
		if err != nil {
			return fmt.Errorf("Warning: %v", err)
		}

		if roundID != firstRoundID {
			return fmt.Errorf("Round ID mismatch! %s is from round %s, but %s is from round %s", filenames[0], firstRoundID, logfile, firstRoundID)
		}
	}

	return nil
}
