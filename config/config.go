package config

import (
	"fmt"
	"strings"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	TestResultFile string `env:"TEST_RESULT,required"`
	GitHub
}

type GitHub struct {
	URL        string `env:"GITHUB_API_URL,required"`
	SHA        string `env:"GITHUB_SHA,required"`
	Token      string `env:"GITHUB_TOKEN,required"`
	Workspace  string `env:"GITHUB_WORKSPACE,required"`
	Repository string `env:"GITHUB_REPOSITORY,required"`
}

func Load() (Config, error) {
	var cfg Config
	if err := envdecode.StrictDecode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (self *Config) Owner() (string, error) {
	repoArray, err := self.verifyRepository()
	if err != nil {
		return "", err
	}

	return repoArray[0], nil
}

func (self *Config) Repo() (string, error) {
	repoArray, err := self.verifyRepository()
	if err != nil {
		return "", err
	}

	return repoArray[1], nil
}

func (self *Config) TestResult() string {
	return self.Workspace + self.TestResultFile
}

func (self *Config) verifyRepository() ([]string, error) {
	repoArray := strings.Split(self.GitHub.Repository, "/")
	if len(repoArray) != 2 {
		return nil, fmt.Errorf("Invalid environment variable 'GITHUB_REPOSITORY'. should be 'owner/repository' instead of '%s'",
			self.GitHub.Repository)
	}

	return repoArray, nil
}
