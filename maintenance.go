package main

import (
	"fmt"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/maintenance"
)

func (h handler) maintenanceCreate(policyID, policyName, timeStr string) error {

	// create maintenance
	maintenanceClient, err := maintenance.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating maintenance client")
	}

	startTime := time.Now().Add(-1 * time.Minute) // 1 minute earlier to ensure a quick start
	endTime := addTime(time.Now(), timeStr)

	_, err = maintenanceClient.Create(nil, &maintenance.CreateRequest{
		Description: fmt.Sprintf("policy %s enabled by %s ", policyName, h.config.Prefix),
		Time: maintenance.Time{
			Type:      maintenance.Schedule,
			StartDate: &startTime,
			EndDate:   &endTime,
		},
		Rules: []maintenance.Rule{
			maintenance.Rule{
				State: maintenance.Enabled,
				Entity: maintenance.Entity{
					Id:   policyID,
					Type: maintenance.Policy,
				},
				//TeamId: h.config.TeamID,
				//Type:   policy.AlertPolicy,
			},
		},
	})
	fmt.Printf("Policy: %s is enabled for the next %s\n", policyName, timeStr)
	return nil
}

func (h handler) maintananceList(stype maintenance.StatusType) (*maintenance.ListResult, error) {
	maintenanceClient, err := maintenance.NewClient(h.client)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating policy client")
	}
	//maintenance, err := maintenance.ListRequest(nil, &maintenance.ListRequest{TeamId: h.config.TeamID})
	return maintenanceClient.List(nil, &maintenance.ListRequest{Type: stype})
}

func (h handler) maintenanceGet(id string) (*maintenance.GetResult, error) {
	maintenanceClient, err := maintenance.NewClient(h.client)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating policy client")
	}
	//maintenance, err := maintenance.ListRequest(nil, &maintenance.ListRequest{TeamId: h.config.TeamID})
	return maintenanceClient.Get(nil, &maintenance.GetRequest{Id: id})
}
