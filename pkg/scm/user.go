package scm

import "time"

// User represents a user account.
type User struct {
	Login   string
	Name    string
	Email   string
	Avatar  string
	Created time.Time
	Updated time.Time
}
