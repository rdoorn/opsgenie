package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const configName = ".opscli"
const configName2 = "opscli.config"

type Config struct {
	ApiKey          string
	ApiKeyEncrypted string
	TeamID          string
	Timezone        string
	timeZone        *time.Location
	Prefix          string
}

func readConfig(fileName string) (*Config, error) {
	c := &Config{}
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	c.timeZone, err = time.LoadLocation(c.Timezone)
	if err != nil {
		return nil, fmt.Errorf("error reading timezone: %s", err)
	}

	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	c.Prefix = fmt.Sprintf("%s_CLI", strings.ToUpper(user.Username))

	if c.ApiKeyEncrypted != "" {
		decText, err := Decrypt(c.ApiKeyEncrypted, MySecret)
		if err != nil {
			fmt.Println("error decrypting key: ", err)
		}
		c.ApiKey = string(decText)
	}

	return c, nil
}

func findConfig() (*Config, error) {
	// first search in current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	searchPaths := []string{dir, user.HomeDir, "/etc/opscli"}

	for _, search := range searchPaths {
		paths := strings.Split(search, "/")

		for i := len(paths); i > 0; i-- {
			searchPath := strings.Join(paths[:i], "/")
			if _, err := os.Stat(fmt.Sprintf("%s/%s", searchPath, configName)); err == nil {
				return readConfig(fmt.Sprintf("%s/%s", searchPath, configName))
			}
			if _, err := os.Stat(fmt.Sprintf("%s/%s", searchPath, configName2)); err == nil {
				return readConfig(fmt.Sprintf("%s/%s", searchPath, configName2))
			}

		}
	}

	return nil, fmt.Errorf("config file not found in path: %s", configName)
}
