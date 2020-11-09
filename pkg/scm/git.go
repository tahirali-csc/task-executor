package scm

import "time"

// GitService provides access to git resources.
type GitService interface {
}

type Commit struct {
	Sha       string
	Message   string
	Author    Signature
	Committer Signature
	Link      string
}

// Signature identifies a git commit creator.
type Signature struct {
	Name  string
	Email string
	Date  time.Time

	// Fields are optional. The provider may choose to
	// include account information in the response.
	Login  string
	Avatar string
}
