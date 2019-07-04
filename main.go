package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/sirupsen/logrus"
)

type handler struct {
	client *client.Config
	config *Config
}

func main() {
	arg := os.Args
	if len(arg) == 1 {
		help()
	}

	config, err := findConfig(configName)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	h := &handler{
		//client: &client.Config{ApiKey: config.ApiKey, LogLevel: logrus.WarnLevel},
		client: &client.Config{ApiKey: config.ApiKey, LogLevel: logrus.DebugLevel},
		config: config,
	}

	if len(arg) < 3 {
		help()
	}

	// attempt to read ID from 3rd parameter
	id := -1
	if len(arg) > 3 {
		if i, err := strconv.Atoi(arg[3]); err == nil {
			id = i
		}
	}

	timeStr, _ := parseArgs(arg[3:]...)

	switch arg[1] {
	case "alert":
		// alert list
		// alert list 5h
		// alert ack 1
		switch arg[2] {
		case "list":
			err = h.alertList(timeStr)
		case "ack":
			// TODO
		default:
			help()
		}
	case "policy":
		// policy list
		// policy test 1
		// policy test 1 1h
		// policy enable 1
		// policy disable 1
		switch arg[2] {
		case "list":
			err = h.policyList()
		case "test":
			err = h.policyTest(id, timeStr)
		case "enable":
			err = h.policyEnable(id, timeStr)
		case "disable":
			err = h.policyDisable(id)
		default:
			help()
		}
	case "filter":
		// filter xxx yyy zzz 1h
		timeStr, restStr := parseArgs(arg[2:]...)
		switch arg[2] {
		case "regex":
			h.filterRegex(restStr, timeStr)
		default:
			h.filterContains(restStr, timeStr)
		}
	default:
		help()
	}

	if err != nil {
		help()
		log.Printf("Error: %s", err)
	}

	/*
		switch arg[1] {
		case "list":
			timeStr, _ := parseArgs(arg[1:]...)
			if timeStr == "" {
				timeStr = "1d"
			}

			switch arg[2] {
			case "filters":
				fallthrough
			case "policies":
				h.listPolicies()
			case "alerts":
				h.listAlerts(timeStr)
			}

		case "help":
			help()

		case "test":
			timeStr, _ := parseArgs(arg[1:]...)
			if timeStr == "" {
				timeStr = "1d"
			}

			if len(arg) < 3 {
				help()
			}
			policyID, err := strconv.Atoi(arg[3])
			if err != nil {
				log.Printf("invalid policy ID: %s expected a number", arg[3])
			}

			switch arg[2] {
			case "filter":
				fallthrough
			case "policy":
				h.testPolicy(policyID, timeStr)
			}

		case "filter":
			timeStr, filter := parseArgs(arg[2:]...)
			log.Printf("filter: '%s' for %s\n", filter, timeStr)

		case "enable":
			timeStr, _ := parseArgs(arg[2:]...)

			if len(arg) < 3 {
				help()
			}
			policyID, err := strconv.Atoi(arg[3])
			if err != nil {
				log.Printf("invalid policy ID: %s expected a number", arg[3])
			}

			switch arg[2] {
			case "filter":
				fallthrough
			case "policy":
				h.enablePolicy(policyID, timeStr)
			}

		case "disable":
			if len(arg) < 3 {
				help()
			}
			policyID, err := strconv.Atoi(arg[3])
			if err != nil {
				log.Printf("invalid policy ID: %s expected a number", arg[3])
			}

			switch arg[2] {
			case "filter":
				fallthrough
			case "policy":
				h.disablePolicy(policyID)
			}

		default:
			log.Printf("unknown parameter: %s", arg[1])
		}
	*/

}

func help() {
	arg := os.Args
	binary := strings.Split(arg[0], "/")[1]

	fmt.Printf("%s alert list               - list all alerts\n", binary)
	fmt.Printf("%s alert list 5h            - list all alerts in the past 5 hours\n", binary)

	fmt.Printf("%s policy list             - list all policies\n", binary)
	fmt.Printf("%s policy test 1             - see what would match policy 1 (use 'ops list policies' to find the number)\n", binary)
	fmt.Printf("%s policy test 1 5m          - see what would match policy 1 in the last 5 minutes\n", binary)

	fmt.Printf("%s policy enable 1           - enable policy 1 \n", binary)
	fmt.Printf("%s policy enable 1 1h        - enable policy 1 for 1 hour\n", binary)
	fmt.Printf("%s policy disable 1          - enable policy 1 for 1 hour\n", binary)

	fmt.Printf("%s filter your filter 1h30m  - create a policy and enable it for 1 hour and 30 minutes\n", binary)

	fmt.Printf("%s help                      - your looking at it\n", binary)

	os.Exit(255)
}
