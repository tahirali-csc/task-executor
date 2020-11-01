package scm

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	ssh2 "golang.org/x/crypto/ssh"
)

type BasicAuth struct {
	User     string
	Password string
}
type SSHAuth struct {
	PrivateKey string
}

type AuthConfig struct {
	BasicAuth *BasicAuth
	SSHAuth   *SSHAuth
}

type CloneOptions struct {
	RepoURL       string
	Branch        string
	Auth          AuthConfig
	Depth         int
	SingleBranch  bool
	ReferenceName string
}

func Clone(dir string, clone *CloneOptions) error {

	depth := clone.Depth
	if depth <= 0 {
		depth = 1
	}

	cloneOpt := &git.CloneOptions{
		URL:        clone.RepoURL,
		RemoteName: "",
		//ReferenceName:     plumbing.NewBranchReferenceName(clone.ReferenceName),
		SingleBranch:      clone.SingleBranch,
		NoCheckout:        false,
		Depth:             depth,
		RecurseSubmodules: 0,
		Progress:          os.Stdout,
		Tags:              0,
	}

	if clone.Auth.BasicAuth != nil {
		cloneOpt.Auth = newBasicAuth(clone.Auth.BasicAuth)
	} else {
		auth, err := newSSHAuth(clone.Auth.SSHAuth)
		if err != nil {
			return err
		}
		cloneOpt.Auth = auth
	}

	_, err := git.PlainClone(dir, false, cloneOpt)
	if err != nil {
		return err
	}

	return nil

}

func newBasicAuth(basicAuth *BasicAuth) transport.AuthMethod {
	return &http.BasicAuth{
		Username: basicAuth.User,
		Password: basicAuth.Password,
	}
}

func newSSHAuth(sshAuth *SSHAuth) (transport.AuthMethod, error) {
	publicKeys, err := ssh.NewPublicKeys("git", []byte(sshAuth.PrivateKey), "")
	if err != nil {
		return nil, err
	}

	publicKeys.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
	return publicKeys, nil
}
