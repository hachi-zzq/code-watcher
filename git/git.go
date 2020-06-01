package git

import (
	AppConfig "coding.net/code-watcher/config"
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"log"
	"regexp"
)

func FetchRepo(repoName, branch string) string {
	fs := osfs.New("git_storage/" + repoName + "_" + branch)
	r, err := git.Init(filesystem.NewStorage(fs, cache.NewObjectLRUDefault()), nil)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			r, err = git.Open(filesystem.NewStorage(fs, cache.NewObjectLRUDefault()), nil)
		} else {
			log.Panic(err.Error())
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

	if refers,e := remote.List(&git.ListOptions{}); e != nil{
		log.Panic(e)
	}else {
		for _,refer := range refers {
			if  refer.Name() == plumbing.NewBranchReferenceName(branch) {
				return refer.Hash().String()
			}
		}
	}

	return ""
}
