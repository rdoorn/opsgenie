package main

import (
	"fmt"
	"log"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
)

// policyList shows all the policies
func (h handler) findPolicyByName(policyName string) (*policy.PolicyProps, error) {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)

	if err != nil {
		return nil, fmt.Errorf("error occured while creating policy client")
	}

	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return nil, err
	}

	for _, r := range result.Policies {
		if r.Name == policyName {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("policy not found")
}

// policyList shows all the policies
func (h handler) policyList() error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)

	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	fmt.Printf("Policies:\n")
	for id, r := range result.Policies {
		enabled := ""
		if r.Enabled {
			enabled = "enabled"
		} else {
			enabled = "disabled"
		}
		fmt.Printf("%2d: [%8s] %s \n", id, enabled, r.Name)

		/*
			details, err := policyClient.GetAlertPolicy(nil, &policy.GetAlertPolicyRequest{Id: r.Id, TeamId: h.config.TeamID})
			if err != nil {
				return err
			}

			fmt.Printf("policy detail: tags%s\n", details.Tags)
		*/

	}

	return nil
}

// policyTest can test a policy against a filter
func (h handler) policyTest(policyID int, timeframe string) error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)

	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	if len(result.Policies)-1 < policyID {
		return fmt.Errorf("policy %d not found\n", policyID)
	}

	if policyID < 0 {
		return fmt.Errorf("policy not specified\n")
	}

	p := result.Policies[policyID]
	details, err := policyClient.GetAlertPolicy(nil, &policy.GetAlertPolicyRequest{Id: p.Id, TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	log.Printf("type:%+v conditions:%+v\n", details.Filter.ConditionMatchType, details.Filter.Conditions)
	query, err := filter2query(details.Filter)
	if err != nil {
		return err
	}
	log.Printf("string filter: %s", query)

	if timeframe != "" {
		history := subTime(time.Now(), timeframe)
		query = fmt.Sprintf("(%s) AND createdAt > %d", query, history.Unix())
	}

	count, err := h.listAlertsQuery(query)
	if count == 0 {
		fmt.Printf("\nIf your filter does not do what you expect, you may test it manually at: https://schubergphilis.app.opsgenie.com/alert/list\n")
	}
	return err
}

func (h handler) policyDisable(policyID int) error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	// get policy details by the id provided
	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	if len(result.Policies) < policyID {
		return fmt.Errorf("policy %d not found\n", policyID)
	}

	if policyID < 0 {
		return fmt.Errorf("policy not specified\n")
	}

	p := result.Policies[policyID]

	//disable policy
	_, err = policyClient.DisablePolicy(nil, &policy.DisablePolicyRequest{TeamId: h.config.TeamID, Id: p.Id, Type: policy.AlertPolicy})
	if err != nil {
		return err
	}
	fmt.Printf("Policy disabled successfuly\n")
	return nil
}

func (h handler) policyEnable(policyID int, timeStr string) error {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)
	if err != nil {
		return fmt.Errorf("error occured while creating policy client")
	}

	result, err := policyClient.ListAlertPolicies(nil, &policy.ListAlertPoliciesRequest{TeamId: h.config.TeamID})
	if err != nil {
		return err
	}

	if len(result.Policies) < policyID {
		return fmt.Errorf("policy %d not found\n", policyID)
	}

	if policyID < 0 {
		return fmt.Errorf("policy not specified\n")
	}

	p := result.Policies[policyID]

	if timeStr == "" {

		//enable policy
		_, err = policyClient.EnablePolicy(nil, &policy.EnablePolicyRequest{TeamId: h.config.TeamID, Id: p.Id, Type: policy.AlertPolicy})
		if err != nil {
			return err
		}
		fmt.Printf("Policy enabled successfuly\n")
	} else {
		return h.maintenanceCreate(p.Id, p.Name, timeStr)
	}
	return nil
}

// input: search key + search value
// output: policyID, policyName, error
func (h handler) createPolicy(key, value string) (string, string, error) {
	//create a policy client
	policyClient, err := policy.NewClient(h.client)
	if err != nil {
		return "", "", fmt.Errorf("error occured while creating policy client")
	}

	condition := og.Contains
	searchKey := ""
	switch key {
	case "contains":
		searchKey = "description"
		condition = og.Contains
	case "regex":
		searchKey = "description"
		condition = og.Matches
	default:
		return "", "", fmt.Errorf("unknown key: %s", key)
	}

	policyName := replaceNonAlphanum(fmt.Sprintf("%s %s %s: %s", h.config.Prefix, searchKey, key, value), "_")
	if p, err := h.findPolicyByName(policyName); err == nil {
		// policy already exists
		return p.Id, p.Name, nil
	}

	req := policy.CreateAlertPolicyRequest{
		MainFields: policy.MainFields{
			PolicyType:        "alert",
			Name:              policyName,
			Enabled:           false,
			PolicyDescription: "created by ops-cli",
			TeamId:            h.config.TeamID,
			Filter: &og.Filter{
				ConditionMatchType: og.MatchAllConditions,
				Conditions: []og.Condition{
					og.Condition{
						Operation:     condition,
						Field:         og.Description,
						ExpectedValue: value,
					}},
			},
		},
		Message: "{{message}}",
		Tags:    []string{"filtered"},
	}
	req.TeamId = h.config.TeamID

	res, err := policyClient.CreateAlertPolicy(nil, &req)

	return res.Id, res.Name, err

}
