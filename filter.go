package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
)

func (h handler) filterContains(query, timeStr string) error {
	result, err := h.createPolicy("contains", query)
	if err != nil {
		return fmt.Errorf("create policy failed: %s", err)
	}
	return h.maintenanceCreate(result.Id, result.Name, timeStr)
}

func (h handler) filterRegex(query, timeStr string) {
}

func (h handler) createPolicy(key, value string) (*policy.CreateResult, error) {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating policy client")
	}

	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	req := policy.CreateAlertPolicyRequest{
		MainFields: policy.MainFields{
			PolicyType:        "alert",
			Name:              fmt.Sprintf("ops-cli (%s) %s: %s", user.Username, key, value),
			Enabled:           false,
			PolicyDescription: "created by ops-cli",
			TeamId:            h.config.TeamID,
			Filter: &og.Filter{
				ConditionMatchType: og.MatchAllConditions,
				Conditions: []og.Condition{
					og.Condition{
						Operation:     og.Contains,
						Field:         og.Description,
						ExpectedValue: value,
					}},
			},
		},
		Message: "{{message}}",
		Tags:    []string{"filtered"},
	}
	req.TeamId = h.config.TeamID

	log.Printf("req: %+v", req)
	//policyClient.
	res, err := policyClient.CreateAlertPolicy(nil, &req)

	log.Printf("res: %+v, err: %s", res, err)
	return res, err

}
