package main

import (
	"fmt"
	"os/user"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/maintenance"
)

func (h handler) maintenanceCreate(policyID, policyName, timeStr string) error {

	user, err := user.Current()
	if err != nil {
		return err
	}

	// create maintenance
	maintenanceClient, err := maintenance.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating maintenance client")
	}

	startTime := time.Now()
	endTime := addTime(time.Now(), timeStr)

	_, err = maintenanceClient.Create(nil, &maintenance.CreateRequest{
		Description: fmt.Sprintf("ops-cli filter %s by %s ", policyName, user.Username),
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

	return nil
}
