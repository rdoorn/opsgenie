package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const configName = ".opsconfig"

type Config struct {
	ApiKey   string
	TeamID   string
	Timezone string
	timeZone *time.Location
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

	return c, nil
}

func findConfig(fileName string) (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	paths := strings.Split(dir, "/")

	for i := len(paths); i > 0; i-- {
		searchPath := strings.Join(paths[:i], "/")
		if _, err := os.Stat(fmt.Sprintf("%s/%s", searchPath, fileName)); err == nil {
			return readConfig(fmt.Sprintf("%s/%s", searchPath, fileName))
		}

	}
	return nil, fmt.Errorf("config file not found in path: %s", configName)
}
