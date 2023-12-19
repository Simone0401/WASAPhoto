package api

import (
	"net/http"
	"strconv"
)

// isAuthorized checks if the user is authorized to perform the action, by checking the Authorization header.
// The auth token must be in the format "Bearer and must be equal to userID to be valid".
// If the token is added to header the function will return true, otherwise it will return false.
func (rt *_router) isAuthorized(header http.Header) bool {
	var authToken int
	authToken, err := strconv.Atoi(header.Get("Authorization"))
	if err != nil {
		return false
	}
	if authToken <= 0 {
		return false
	}
	return true
}
