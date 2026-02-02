package parser

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
	//TODO: Open log files.
	//TODO: Go through each line, parse JSON, check if ckey matches. If yes, append to result slice.
	// Return result slice and error if any.

	return nil, nil
}
