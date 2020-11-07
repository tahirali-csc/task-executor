package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/scm"
	"github.com/task-executor/pkg/utils"
)

func readSecrets() {
	log.Println("Reading secrets")
	dat, _ := ioutil.ReadFile("/etc/secret-volume/username")
	log.Println(string(dat))

	dat, _ = ioutil.ReadFile("/etc/secret-volume/password")
	log.Println(string(dat))
}

const BasicAuth = "basic-auth"
const SSHAuth = "ssh-auth"

func main() {
	utils.InitLogs(log.DebugLevel)

	repoURL := flag.String("repo", "", "Clone URL")
	branch := flag.String("branch", "", "Branch")
	secretType := flag.String("secret-type", "", "Type of secret")
	cloneDir := flag.String("clone-dir", "", "Clone Directory")
	flag.Parse()

	if repoURL == nil {
		log.Println("Repo URL is missing")
		return
	}
	if branch == nil {
		log.Println("Branch is missing")
		return
	}
	if secretType == nil {
		log.Println("Secret Type is missing")
		return
	}

	cloneOpts := &scm.CloneOptions{
		RepoURL: *repoURL,
		Branch:  *branch,
	}

	switch *secretType {
	case BasicAuth:

		user := os.Getenv("USER")
		if len(strings.TrimSpace(user)) == 0 {
			log.Println("User ID is missing")
			return
		}

		password := os.Getenv("PASSWORD")
		if len(strings.TrimSpace(user)) == 0 {
			log.Println("Password is missing")
			return
		}

		cloneOpts.Auth.BasicAuth = &scm.BasicAuth{
			User:     user,
			Password: password,
		}

	case SSHAuth:
		sshKey := os.Getenv("SSHKey")
		if len(strings.TrimSpace(sshKey)) == 0 {
			log.Println("SSH Key is missing")
			return
		}

		cloneOpts.Auth.SSHAuth = &scm.SSHAuth{
			PrivateKey: sshKey,
		}

	default:
		log.Error("Unsupported secret credentials")
		return
	}
	//readSecrets()

	if cloneDir == nil || len(*cloneDir) == 0{
		tempDir, err := ioutil.TempDir(*cloneDir, "app")
		if err != nil {
			log.Fatal(err)
		}
		cloneDir = &tempDir
	}
	log.Println("Cloning at:::" + *cloneDir)

	err := scm.Clone(*cloneDir, cloneOpts)
	if err != nil {
		log.Fatal(err)
	}
}
