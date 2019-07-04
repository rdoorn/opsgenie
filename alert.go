package main

import (
	"fmt"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
)

func (h handler) alertList(timeframe string) error {
	// only add time filter, if a time has been specified
	var query string
	if timeframe != "" {
		history := subTime(time.Now(), timeframe)
		query = fmt.Sprintf("createdAt > %d", history.Unix())
	}
	return h.listAlertsQuery(query)
}

func (h handler) listAlertsQuery(query string) error {
	alertClient, err := alert.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating alert client")
	}

	alertResult, err := alertClient.List(nil, &alert.ListAlertRequest{Query: query})
	if err != nil {
		return err
	}

	fmt.Printf("Alerts:\n")
	for i, a := range alertResult.Alerts {

		tags := ""
		if a.Acknowledged == false && a.Status == "open" {
			tags += "NEW"
		} else if a.Acknowledged == false && a.Status == "closed" {
			tags += "---"
		} else {
			tags += "ACK"
		}
		/*
			if contains(a.Tags, "filtered") {
				tags += "filtered"
			} else {
				if a.Acknowledged == false && a.Status == "open" {
					tags += "new page"
				} else {
					tags += "ack'd page"
				}
			}*/
		fmt.Printf("%.2d: [%3s] %v %s (tags: %+v)\n", i, tags, a.CreatedAt.In(h.config.timeZone).Format("01-02-06 15:04:05"), a.Message, a.Tags)
	}
	return nil
}
