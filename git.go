package main

import (
	AppConfig "coding.net/code-watcher/config"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
	"regexp"
)

func FetchRepo(branch string) string {
	r, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		log.Fatalln(err.Error())
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
		log.Fatalf(err.Error())
	}

	if err = r.Fetch(&git.FetchOptions{
		RemoteName: remote.Config().Name,
	}); err != nil {
		log.Fatalf(err.Error())
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
