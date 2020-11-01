package main

import (
	"flag"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/scm"
	"github.com/task-executor/pkg/utils"
)

func main() {
	utils.InitLogs(log.DebugLevel)

	repoURL := flag.String("repoURL", "", "Clone URL")
	branch := flag.String("branch", "", "Branch")
	user := flag.String("user", "", "User")
	password := flag.String("password", "", "Password")
	privateKey := flag.String("privateKey", "", "Private Key")
	cloneDir := flag.String("cloneDir", "", "Clone Directory")

	flag.Parse()

	if cloneDir == nil {
		tempDir, err := ioutil.TempDir(*cloneDir, "app")
		if err != nil {
			log.Fatal(err)
		}
		cloneDir = &tempDir
	}
	log.Println("Cloning at:::" + *cloneDir)

	cloneOpts := &scm.CloneOptions{
		RepoURL: *repoURL,
		Branch:  *branch,
	}

	if user != nil {
		cloneOpts.Auth.BasicAuth = &scm.BasicAuth{
			User:     *user,
			Password: *password,
		}
	} else {
		cloneOpts.Auth.SSHAuth = &scm.SSHAuth{
			PrivateKey: *privateKey,
		}
	}

	err := scm.Clone(*cloneDir, cloneOpts)
	if err != nil {
		log.Fatal(err)
	}
}
