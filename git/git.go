package git

import (
	AppConfig "coding.net/code-watcher/config"
	"fmt"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"log"
	"regexp"
)

func FetchRepo(repoName, branch string) string {
	fmt.Printf(branch)
	fs := osfs.New("git_storage/" + repoName + "_" + branch)
	r, err := git.Init(filesystem.NewStorage(fs, cache.NewObjectLRUDefault()), nil)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			r, err = git.Open(filesystem.NewStorage(fs, cache.NewObjectLRUDefault()), nil)
		} else {
			log.Fatalln(err.Error())
		}
	}

	repoUrl := AppConfig.AppConfig.RepoUrl
	repoUsername := AppConfig.AppConfig.RepoUserName
	repoPass := AppConfig.AppConfig.RepoPassword

	if repoUsername != "" {
		if r, err := regexp.Compile("(http://|https://)"); err == nil {
			repoUrl = r.ReplaceAllString(repoUrl, fmt.Sprintf("%s%s:%s@", "${1}", repoUsername, repoPass))
		}
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoUrl},
	})

	if err != nil {
		if err == git.ErrRemoteExists {
			remote, _ = r.Remote("origin")
		} else {
			log.Fatalf(err.Error())
		}
	}

	if err = r.Fetch(&git.FetchOptions{
		RemoteName: remote.Config().Name,
	}); err != nil {

		if err == git.NoErrAlreadyUpToDate {
			log.Printf("cc")
		} else {
			log.Fatalf(err.Error())
		}

	}

	branchHash := ""

	if ref, err := r.Reference(plumbing.NewRemoteReferenceName(remote.Config().Name, branch), true); err == nil {
		if ref != nil {
			branchHash = ref.Hash().String()
		}
	} else {
		log.Fatalf(err.Error())
	}
	return branchHash
}
