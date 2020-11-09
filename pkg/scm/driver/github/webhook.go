package github

import (
	"encoding/json"
	"fmt"
	"github.com/guregu/null"
	"github.com/task-executor/pkg/scm"
	"io/ioutil"
	"net/http"
	"time"
)

type webhookService struct {
}

func (s *webhookService) Parse(req *http.Request, fn interface{}) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var hook scm.Webhook
	switch req.Header.Get("X-GitHub-Event") {
	case "push":
		hook, err = s.parsePushHook(data)
	default:
	}
	if err != nil {
		return nil, err
	}

	return hook, nil

}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	return convertPushHook(dst), err
}

func convertPushHook(src *pushHook) *scm.PushHook {
	dst := &scm.PushHook{
		Ref:     src.Ref,
		BaseRef: src.BaseRef,
		Before:  src.Before,
		After:   src.After,
		Commit: scm.Commit{
			Sha:     src.After,
			Message: src.Head.Message,
			Link:    src.Compare,
			Author: scm.Signature{
				Login: src.Head.Author.Username,
				Email: src.Head.Author.Email,
				Name:  src.Head.Author.Name,
				Date:  src.Head.Timestamp.ValueOrZero(),
			},
			Committer: scm.Signature{
				Login: src.Head.Committer.Username,
				Email: src.Head.Committer.Email,
				Name:  src.Head.Committer.Name,
				Date:  src.Head.Timestamp.ValueOrZero(),
			},
		},
		Repo: scm.Repository{
			ID:        fmt.Sprint(src.Repository.ID),
			Namespace: src.Repository.Owner.Login,
			Name:      src.Repository.Name,
			Branch:    src.Repository.DefaultBranch,
			Private:   src.Repository.Private,
			Clone:     src.Repository.CloneURL,
			CloneSSH:  src.Repository.SSHURL,
			Link:      src.Repository.HTMLURL,
		},
		Sender: *convertUser(&src.Sender),
	}
	// fix https://github.com/drone/go-scm/issues/8
	//if scm.IsTag(dst.Ref) && src.Head.ID != "" {
	//	dst.Commit.Sha = src.Head.ID
	//	dst.After = src.Head.ID
	//}
	return dst
}

type (
	// github push webhook payload
	pushHook struct {
		Ref     string `json:"ref"`
		BaseRef string `json:"base_ref"`
		Before  string `json:"before"`
		After   string `json:"after"`
		Compare string `json:"compare"`
		Head    struct {
			ID        string    `json:"id"`
			TreeID    string    `json:"tree_id"`
			Distinct  bool      `json:"distinct"`
			Message   string    `json:"message"`
			Timestamp null.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"committer"`
			Added    []interface{} `json:"added"`
			Removed  []interface{} `json:"removed"`
			Modified []string      `json:"modified"`
		} `json:"head_commit"`
		Commits []struct {
			ID        string    `json:"id"`
			TreeID    string    `json:"tree_id"`
			Distinct  bool      `json:"distinct"`
			Message   string    `json:"message"`
			Timestamp null.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"committer"`
			Added    []interface{} `json:"added"`
			Removed  []interface{} `json:"removed"`
			Modified []string      `json:"modified"`
		} `json:"commits"`
		Repository struct {
			ID    int64 `json:"id"`
			Owner struct {
				Login     string `json:"login"`
				AvatarURL string `json:"avatar_url"`
			} `json:"owner"`
			Name          string `json:"name"`
			FullName      string `json:"full_name"`
			Private       bool   `json:"private"`
			Fork          bool   `json:"fork"`
			HTMLURL       string `json:"html_url"`
			SSHURL        string `json:"ssh_url"`
			CloneURL      string `json:"clone_url"`
			DefaultBranch string `json:"default_branch"`
		} `json:"repository"`
		Pusher user `json:"pusher"`
		Sender user `json:"sender"`
	}
)

type user struct {
	ID      int         `json:"id"`
	Login   string      `json:"login"`
	Name    string      `json:"name"`
	Email   null.String `json:"email"`
	Avatar  string      `json:"avatar_url"`
	Created time.Time   `json:"created_at"`
	Updated time.Time   `json:"updated_at"`
}

func convertUser(from *user) *scm.User {
	return &scm.User{
		Avatar:  from.Avatar,
		Email:   from.Email.String,
		Login:   from.Login,
		Name:    from.Name,
		Created: from.Created,
		Updated: from.Updated,
	}
}
