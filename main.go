package main

import (
	AppConfig "coding.net/code-watcher/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func main() {
	repoBranch := AppConfig.AppConfig.RepoBranch
	repoName := AppConfig.AppConfig.RepoName
	c := cron.New(cron.WithSeconds())
	c.AddFunc("*/2 * * * * *", func() {
		hash := FetchRepo(repoBranch)
		if hash == "" {
			log.Println(fmt.Sprintf("branch: %s not found", repoBranch))
			return
		}
		type Branch struct {
			id   int
			refs string
			hash string
		}

		mysqlDSN := AppConfig.AppConfig.MySQLDSN
		db, err := sql.Open("mysql", mysqlDSN)
		if err != nil {
			log.Fatalf(err.Error())
		}

		//查询是最新
		if stat, err := db.Prepare("select id,refs,hash from branches where refs = ? and repo = ?"); err == nil {
			row := stat.QueryRow(repoBranch, repoName)
			b := Branch{}
			scanErr := row.Scan(&b.id, &b.refs, &b.hash)
			if scanErr == sql.ErrNoRows {
				log.Println(fmt.Sprintf("fetch branch 12121: %s new refs: %s", repoBranch, hash))
				//inert
				if _, err := db.Exec("insert branches values (null, ? , ? , ?  )", repoName, repoBranch, hash); err != nil {
					log.Println(err.Error())
					return
				}
				requestJenkins()
				return
			} else {
				if b.hash == hash {
					//已经是最新，continue
					log.Println(fmt.Sprintf(" branch: %s refs: %s exist , continue ....", repoBranch, b.hash))
				} else {
					log.Println(fmt.Sprintf("fetch branch: %s new refs: %s", repoBranch, hash))
					//不是最新，update
					if _, err := db.Exec("update branches set hash = ? where id = ?", hash, b.id); err != nil {
						log.Println(err.Error())
						return
					}
					requestJenkins()
				}
				return
			}

			stat.Close()

		} else {
			log.Fatalf(err.Error())
		}

		defer db.Close()
	})

	c.Start()
	for {
		time.Sleep(time.Hour)
	}
}

func requestJenkins() {
	jenkinsUrl := AppConfig.AppConfig.JenkinsUrl
	jenkinsUsername := AppConfig.AppConfig.JenkinsName
	jenkinsToken := AppConfig.AppConfig.JenkinsToken
	gorequest.New().
		Post(jenkinsUrl).
		SetBasicAuth(jenkinsUsername, jenkinsToken).
		End()

	log.Println("trigger jenkins success！")
}
