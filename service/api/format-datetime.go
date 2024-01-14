package api

import "time"

// formatDatetime allows to format a string in RFC3339 format to a string in the following format:
// "2006-01-02 15:04:05"
// Function will return the new format
func formatDatetime(datetime string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return "", err
	}
	return parsedTime.Format("2006-01-02 15:04:05"), nil
}
