package main

import (
	"fmt"
	"strings"
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
	fmt.Printf("Alert History:\n")
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

	for i, a := range alertResult.Alerts {

		tags := ""
		if a.Acknowledged == false && a.Status == "open" {
			tags += "NEW"
		} else if a.Acknowledged == false && a.Status == "closed" {
			tags += "---"
		} else {
			tags += "ACK"
		}
		opsTags := ""
		if len(tags) > 0 {
			opsTags = fmt.Sprintf("(%s)", strings.Join(a.Tags, ","))
		}
		fmt.Printf("%3d: [%3s] %v %s %s\n", i, tags, a.CreatedAt.In(h.config.timeZone).Format("01-02-06 15:04:05"), a.Message, opsTags)
	}
	if len(alertResult.Alerts) == 0 {
		fmt.Printf("no alerts matched or found\n")
	}
	return nil
}

func (h handler) findAlertByID(id int) (*alert.Alert, error) {
	alertClient, err := alert.NewClient(h.client)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating alert client")
	}

	alertResult, err := alertClient.List(nil, &alert.ListAlertRequest{})
	if err != nil {
		return nil, err
	}

	if len(alertResult.Alerts) < id {
		return nil, fmt.Errorf("alert %d not found\n", id)
	}

	if id < 0 {
		return nil, fmt.Errorf("alert id not specified\n")
	}

	return &alertResult.Alerts[id], nil

}

func (h handler) alertAck(id int) error {
	alertDetails, err := h.findAlertByID(id)
	if err != nil {
		return err
	}

	if alertDetails.Status != "open" || alertDetails.Acknowledged == true {
		fmt.Printf("alert is already ack'd or closed: %s\n", alertDetails.Message)
		return nil
	}

	alertClient, err := alert.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating alert client")
	}

	result, err := alertClient.Acknowledge(nil, &alert.AcknowledgeAlertRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: alertDetails.Id,
		User:            h.config.Prefix,
	})

	if err != nil {
		return err
	}

	fmt.Printf("ack'n alert: %s... %s\n", alertDetails.Message, result.Result)
	return nil
}

func (h handler) alertAll() error {

	alertClient, err := alert.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating alert client")
	}

	alertResult, err := alertClient.List(nil, &alert.ListAlertRequest{Query: "status: open AND acknowledged: false"})
	if err != nil {
		return err
	}

	if len(alertResult.Alerts) == 0 {
		fmt.Printf("no open unacknowledged alerts\n")
		return nil
	}
	for _, a := range alertResult.Alerts {
		result, err := alertClient.Acknowledge(nil, &alert.AcknowledgeAlertRequest{
			IdentifierType:  alert.ALERTID,
			IdentifierValue: a.Id,
			User:            h.config.Prefix,
		})
		if err != nil {
			return err
		}
		fmt.Printf("ack'n alert: %s... %s\n", a.Message, result.Result)
	}

	return nil
}
