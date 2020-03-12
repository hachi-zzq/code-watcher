package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	//mysql
	MySQLDSN string

	//repo url
	RepoName     string
	RepoUrl      string
	RepoUserName string
	RepoPassword string
	RepoBranch   string

	//jenkins
	JenkinsUrl   string
	JenkinsName  string
	JenkinsToken string
}

var AppConfig *Config = func() *Config {
	setDefaultConfig()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Println("can not find config file")
	}
	viper.AutomaticEnv()
	return &Config{
		MySQLDSN:     viper.GetString("MySQL_DSN"),
		RepoUrl:      viper.GetString("REPO_URL"),
		RepoName:     viper.GetString("REPO_NAME"),
		RepoUserName: viper.GetString("REPO_USERNAME"),
		RepoPassword: viper.GetString("REPO_PASSWORD"),
		RepoBranch:   viper.GetString("REPO_BRANCH"),
		JenkinsUrl:   viper.GetString("JENKINS_URL"),
		JenkinsName:  viper.GetString("JENKINS_NAME"),
		JenkinsToken: viper.GetString("JENKINS_TOKEN"),
	}
}()

func setDefaultConfig() {
	viper.SetDefault("MySQL_DSN", "user:password@localhost:3306/dbname")
	viper.SetDefault("REPO_URL", "http://localhost/test.git")
	viper.SetDefault("REPO_NAME", "coding-dev")
	viper.SetDefault("REPO_USERNAME", "")
	viper.SetDefault("REPO_PASSWORD", "")
	viper.SetDefault("REPO_BRANCH", "test")
	viper.SetDefault("JENKINS_URL", "http://localhost")
	viper.SetDefault("JENKINS_NAME", "test")
	viper.SetDefault("JENKINS_TOKEN", "test")
}
