package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
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
	CleanMobIds       bool
	AfterMins         *float64 //nil ptr == no filter
	BeforeMins        *float64 //nil ptr == no filter
	SecondCkey        string
	InteractionWindow float64 //Default 5.0 (minutes)
}

// I've noticed while filtering log entries that there is a mob number stated after the ckey ex: (mob_1234). This clutters the output, making it harder to read.
// Going to remove mob numbers from the ckey entries using regex. Might not be the best solution.
var mobPattern = regexp.MustCompile(`\(mob_\d+\)`)

// isNearAnyTime checks if t is within windowMinutes of any time in the list
func isNearAnyTime(t time.Time, times []time.Time, windowMinutes float64) bool {
	for _, otherTime := range times {
		diff := t.Sub(otherTime).Minutes()
		if diff < 0 {
			diff = -diff //abs
		}
		if diff <= windowMinutes {
			return true
		}
	}
	return false
}

func filterInteractions(entries []LogEntry, ckey1 string, ckey2 string, windowMinutes float64) ([]LogEntry, error) {
	if windowMinutes == 0 {
		windowMinutes = 5.0 //Default time window is 5 minutes
	}

	var player1Times []time.Time
	var player2Times []time.Time

	//First pass: collect timestamps for each player
	for _, entry := range entries {
		if strings.Contains(entry.Message, ckey1) {
			t, err := ParseTimestamp(entry.Timestamp)
			if err == nil {
				player1Times = append(player1Times, t)
			}
		}
		if strings.Contains(entry.Message, ckey2) {
			t, err := ParseTimestamp(entry.Timestamp)
			if err == nil {
				player2Times = append(player2Times, t)
			}
		}
	}

	//Second pass: include entries near interaction times
	var results []LogEntry
	for _, entry := range entries {
		entryTime, err := ParseTimestamp(entry.Timestamp)
		if err != nil {
			continue
		}

		//Check if entry involves either player
		hasPlayer1 := strings.Contains(entry.Message, ckey1)
		hasPlayer2 := strings.Contains(entry.Message, ckey2)

		if !hasPlayer1 && !hasPlayer2 {
			//neither player involved, skip
			continue
		}
		//Check if the other player was active nearby in time
		if hasPlayer1 && isNearAnyTime(entryTime, player2Times, windowMinutes) {
			results = append(results, entry)
		} else if hasPlayer2 && isNearAnyTime(entryTime, player1Times, windowMinutes) {
			results = append(results, entry)
		}
	}
	return results, nil
}

// Filters logs by given ckey. Returns slice of LogEntry and error if any.
func FilterByCkey(filenames []string, ckey string, opts FilterOptions) ([]LogEntry, error) {
	var allResults []LogEntry
	var roundStartTime time.Time

	// Get round start time if we need time filtering
	if opts.AfterMins != nil || opts.BeforeMins != nil {
		var err error
		roundStartTime, err = GetRoundstartTime(filenames[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get round start time: %w", err)
		}
	}

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
			//If interaction filtering is enabled, check for second ckey.
			if opts.SecondCkey != "" {
				if strings.Contains(entry.Message, ckey) || strings.Contains(entry.Message, opts.SecondCkey) {
					results = append(results, entry)
				}
			} else {
				if strings.Contains(entry.Message, ckey) {
					results = append(results, entry)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			file.Close()
			return nil, fmt.Errorf("Error reading file: %w", err)
		}

		file.Close()
		allResults = append(allResults, results...)
	}

	if opts.SecondCkey != "" {
		var err error
		allResults, err = filterInteractions(allResults, ckey, opts.SecondCkey, opts.InteractionWindow)
		if err != nil {
			return nil, err
		}
	}

	//Apply time filtering
	if opts.AfterMins != nil || opts.BeforeMins != nil {
		var filtered []LogEntry
		for _, entry := range allResults {
			entryTime, err := ParseTimestamp(entry.Timestamp)
			if err != nil {
				continue
			}

			minutesElapsed := entryTime.Sub(roundStartTime).Minutes()

			if opts.AfterMins != nil && minutesElapsed < *opts.AfterMins {
				continue
			}

			if opts.BeforeMins != nil && minutesElapsed > *opts.BeforeMins {
				filtered = append(filtered, entry)
			}
		}
		allResults = filtered
	}

	//Clean mob IDs
	if opts.CleanMobIds {
		for i := range allResults {
			allResults[i].Message = mobPattern.ReplaceAllString(allResults[i].Message, "")
		}
	}

	//Sort chronologically
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Timestamp < allResults[j].Timestamp
	})
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

//Helper functions to filter with given timespamps

// Parses roundstart time from second line of log as relative start timestamp.
func ParseTimestamp(ts string) (time.Time, error) {

	//Passed value "2006-01-02 15:04:05.000" is apparently the standard (day 1 - month 2 - hour 3 - minute 4 - second 5 - year 6)
	t, err := time.Parse("2006-01-02 15:04:05.000", ts)
	if err == nil {
		return t, err
	}

	//sometimes this doesn't work, we try parsing without miliseconds.
	t, err = time.Parse("2006-01-02 15:04:05", ts)
	if err == nil {
		return t, err
	}

	//If both does not work, return error

	return time.Time{}, fmt.Errorf("Unable to parse timestamp %s", ts)

}

// Gets the roundstart time from the first time stamp of the passed log file.
func GetRoundstartTime(filename string) (time.Time, error) {

	file, err := os.Open(filename)
	if err != nil {
		return time.Time{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		//skip to the second line for the first timestamp
		if scanner.Scan() {
			var entry LogEntry
			if err := json.Unmarshal([]byte(scanner.Text()), &entry); err == nil {
				return ParseTimestamp(entry.Timestamp)
			}
		}
	}

	//If we reach here, we f'd up
	return time.Time{}, fmt.Errorf("Could not parse roundstart time from %s", filename)
}
