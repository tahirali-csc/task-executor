package scm

import "net/http"

// Webhook defines a webhook for repository events.
type Webhook interface {
	Repository() Repository
}

type WebhookService interface {
	// Parse returns the parsed the repository webhook payload.
	Parse(req *http.Request, fn interface{}) (Webhook, error)
}

// PushHook represents a push hook, eg push events.
type PushHook struct {
	Ref     string
	BaseRef string
	Repo    Repository
	Before  string
	After   string
	Commit  Commit
	Sender  User
}

func (h *PushHook) Repository() Repository               { return h.Repo }