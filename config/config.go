package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charleswan/grab-gitlab-group/utils"
	"github.com/spf13/viper"
)

func Get() *ConfigFields {
	configOnce.Do(func() {
		viper.SetConfigName("config")                    // name of config file (without extension)
		viper.SetConfigType("toml")                      // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(utils.GetCurrentDirectory()) // path to look for the config file in
		viper.AddConfigPath(".")                         // optionally look for config in the working directory
		if err := viper.ReadInConfig(); err != nil {     // Find and read the config file
			panic(err)
		}
		configInstance = newConfig()
		configInstance.limiterMax = viper.Get("sys.limiter_max").(int64)
		configInstance.cookie = viper.Get("auth.gitlab_cookie").(string)

		groupURL := viper.Get("url.gitlab_group").(string)
		if strings.HasSuffix(groupURL, "/") {
			configInstance.groupURL = fmt.Sprintf("%s-/children.json", groupURL)
		} else {
			configInstance.groupURL = fmt.Sprintf("%s/-/children.json", groupURL)
		}
	})

	return configInstance
}

var configOnce sync.Once
var configInstance *ConfigFields

type ConfigFields struct {
	limiterMax   int64
	cookie       string
	groupURL     string
	gitPrefixURL string
	clonePath    string
}

func newConfig() *ConfigFields {
	return &ConfigFields{}
}

func (cfg *ConfigFields) GetLimiterMax() int64 {
	return cfg.limiterMax
}

func (cfg *ConfigFields) GetCookie() string {
	return cfg.cookie
}

func (cfg *ConfigFields) GetGroupURL() string {
	return cfg.groupURL
}

func (cfg *ConfigFields) GetGitPrefixURL(projectName string) string {
	gitPrefixURL := viper.Get("url.gitlab_git_ssh_prefix").(string)
	if strings.HasSuffix(gitPrefixURL, "/") {
		return fmt.Sprintf("%s%s.git", cfg.gitPrefixURL, projectName)
	}

	return fmt.Sprintf("%s/%s.git", cfg.gitPrefixURL, projectName)
}

func (cfg *ConfigFields) GetClonePath(projectName string) string {
	clonePath := viper.Get("path.clone_dest_path").(string)
	if strings.HasSuffix(clonePath, "/") {
		return fmt.Sprintf("%s%s", cfg.clonePath, projectName)
	}

	return fmt.Sprintf("%s/%s", cfg.clonePath, projectName)
}
