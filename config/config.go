package config

import (
	"os"
	"strconv"
	"strings"
)

const Version = "1.0.0"

var allowedUsers []int64

func init() {
	usersString := strings.Split(os.Getenv("ALLOWED_USERS"), ",")
	allowedUsers = make([]int64, 0, len(usersString))
	for _, idString := range usersString {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err == nil {
			allowedUsers = append(allowedUsers, id)
		}
	}
}

// IsUserAllowed checks if a user is allowed to use this bot or not
func IsUserAllowed(userID int64) bool {
	if len(allowedUsers) == 0 {
		return true
	}
	for _, id := range allowedUsers {
		if id == userID {
			return true
		}
	}
	return false
}
