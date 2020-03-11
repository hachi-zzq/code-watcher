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
		MySQLDSN:     viper.GetString("MySQLDSN"),
		RepoUrl:      viper.GetString("RepoUrl"),
		RepoName:     viper.GetString("RepoName"),
		RepoUserName: viper.GetString("RepoUserName"),
		RepoPassword: viper.GetString("RepoPassword"),
		RepoBranch:   viper.GetString("RepoBranch"),
		JenkinsUrl:   viper.GetString("JenkinsUrl"),
		JenkinsName:  viper.GetString("JenkinsName"),
		JenkinsToken: viper.GetString("JenkinsToken"),
	}
}()

func setDefaultConfig() {
	viper.SetDefault("MySQLDSN", "user:password@localhost:3306/dbname")
	viper.SetDefault("RepoUrl", "http://localhost/test.git")
	viper.SetDefault("RepoName", "coding-dev")
	viper.SetDefault("RepoUserName", "")
	viper.SetDefault("RepoPassword", "")
	viper.SetDefault("RepoBranch", "test")
	viper.SetDefault("JenkinsUrl", "http://localhost")
	viper.SetDefault("JenkinsName", "test")
	viper.SetDefault("JenkinsToken", "test")
}
