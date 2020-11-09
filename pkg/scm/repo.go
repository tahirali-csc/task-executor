package scm

import "time"

type Repository struct {
	ID        string
	Namespace string
	Name      string
	Perm      *Perm
	Branch    string
	Private   bool
	Clone     string
	CloneSSH  string
	Link      string
	Created   time.Time
	Updated   time.Time
}

// Perm represents a user's repository permissions.
type Perm struct {
	Pull  bool
	Push  bool
	Admin bool
}
